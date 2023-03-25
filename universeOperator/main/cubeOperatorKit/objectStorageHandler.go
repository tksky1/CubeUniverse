package cubeOperatorKit

import (
	"bytes"
	"context"
	"crypto/md5"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
	ObjectStorageHandler
	这里处理直接针对ceph对象存储的操作
*/

// SessionAndBucketName 存返回的session和对应的bucketName（由rook随机生成）
type SessionAndBucketName struct {
	sess       *session.Session
	bucketName string
}

// 缓存用，依赖md5
var sessionCacheMap map[[16]byte]*SessionAndBucketName

// <----------CRUD功能，供外部调用----------->

// GetObject 访问指定对象，返回对象的Value
func GetObject(namespace, bucketClaimName, key string) (objectValue []byte, errors error) {
	sessWithBucketName, err := GetObjectStorageSession(namespace, bucketClaimName)
	if err != nil {
		return nil, err
	}
	sess := sessWithBucketName.sess
	downloader := s3manager.NewDownloader(sess)
	buf := aws.NewWriteAtBuffer([]byte{})
	_, err = downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(sessWithBucketName.bucketName),
			Key:    aws.String(key),
		})
	return buf.Bytes(), err
}

// PutObject 发送对象Put请求到ceph
func PutObject(namespace, bucketClaimName, key string, value []byte) error {
	sessWithBucketName, err := GetObjectStorageSession(namespace, bucketClaimName)
	if err != nil {
		return err
	}
	sess := sessWithBucketName.sess
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(sessWithBucketName.bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(value),
	})
	return err
}

// DeleteObject 删除指定对象
func DeleteObject(namespace, bucketClaimName, key string) error {
	sessWithBucketName, err := GetObjectStorageSession(namespace, bucketClaimName)
	if err != nil {
		return err
	}
	sess := sessWithBucketName.sess
	svc := s3.New(sess)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	_, err = svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(sessWithBucketName.bucketName),
		Key:    aws.String(key),
	})
	return err
}

// ListObjectFromBucket 列出某bucket的全部Object的key（List方法只能查key，不能查value）
func ListObjectFromBucket(namespace, bucketClaimName string) (keys []string, err error) {
	sessWithBucketName, err := GetObjectStorageSession(namespace, bucketClaimName)
	if err != nil {
		return nil, err
	}
	sess := sessWithBucketName.sess
	svc := s3.New(sess)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	output, err := svc.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(sessWithBucketName.bucketName),
	})
	if err != nil {
		return nil, err
	}
	var key []string
	for _, object := range output.Contents {
		key = append(key, *object.Key)
	}
	return key, nil
}

//	<----------工具函数---------->

// GetObjectStorageSession 获取指定bucketClaim的session
func GetObjectStorageSession(namespace, bucketClaimName string) (*SessionAndBucketName, error) {
	bucketClaimMD5 := md5.Sum([]byte(namespace + bucketClaimName))
	sessionCached, contains := sessionCacheMap[bucketClaimMD5]
	if contains {
		return sessionCached, nil
	}

	cm, err := clientSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), bucketClaimName, v1.GetOptions{})
	if err != nil {
		panic("configMap获取失败：" + err.Error())
	}
	bucketName := cm.Data["BUCKET_NAME"]
	bucketHost := cm.Data["BUCKET_HOST"]
	secret, err := clientSet.CoreV1().Secrets(namespace).Get(context.TODO(), bucketClaimName, v1.GetOptions{})
	if err != nil {
		panic("secret获取失败：" + err.Error())
	}

	codedAccessId := secret.Data["AWS_ACCESS_KEY_ID"]
	accessID := string(codedAccessId)
	if err != nil {
		panic("secret获取失败：" + err.Error())
	}

	codedAccessKey := secret.Data["AWS_SECRET_ACCESS_KEY"]
	accessKey := string(codedAccessKey)

	sess := session.Must(session.NewSession(&aws.Config{
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(accessID, accessKey, ""),
		Region:           aws.String(endpoints.CnNorth1RegionID),
		Endpoint:         aws.String("http://" + bucketHost + ":80")}))
	ret := &SessionAndBucketName{bucketName: bucketName, sess: sess}
	sessionCacheMap[bucketClaimMD5] = ret
	return ret, nil
}
