package cubeOperatorKit

/*	-------------
	供外部调用的接口，接入缓存
	-------------
*/

// #cgo CXXFLAGS: -std=c++11
// #cgo LDFLAGS: -lstdc++
// #include "cache.h"
import "C"
import "encoding/json"

// GetObject 访问指定对象，返回对象的Value
func GetObject(namespace, bucketClaimName, key string) (objectValue []byte, errors error) {
	cacheKey := C.CString(namespace + bucketClaimName + key)
	cacheOut := C.ask(cacheKey)
	outString := C.GoString(cacheOut)
	if outString == "" {
		objectValue, err := GetObjectS3(namespace, bucketClaimName, key)
		if err != nil {
			return nil, err
		}
		C.insr(cacheKey, C.CString(string(objectValue)))
		return objectValue, nil
	}
	return []byte(outString), nil
}

// PutObject 发送对象Put请求到ceph
func PutObject(namespace, bucketClaimName, key string, value []byte) error {
	err := PutObjectS3(namespace, bucketClaimName, key, value)
	if err != nil {
		return err
	}
	cacheKey := C.CString(namespace + bucketClaimName + key)
	cacheKey2 := C.CString("list:" + namespace + bucketClaimName)
	C.insr(cacheKey, C.CString(string(value)))
	C.del(cacheKey2)
	return nil
}

// DeleteObject 删除指定对象
func DeleteObject(namespace, bucketClaimName, key string) error {
	err := DeleteObjectS3(namespace, bucketClaimName, key)
	if err != nil {
		return err
	}
	cacheKey := C.CString(namespace + bucketClaimName + key)
	cacheKey2 := C.CString("list:" + namespace + bucketClaimName)
	C.insr(cacheKey, C.CString(""))
	C.insr(cacheKey2, C.CString(""))
	return nil
}

// ListObjectFromBucket 列出某bucket的全部Object的key（List方法只能查key，不能查value）
func ListObjectFromBucket(namespace, bucketClaimName string) (keys []string, errors error) {
	cacheKey := C.CString("list:" + namespace + bucketClaimName)
	cacheOut := C.ask(cacheKey)
	outString := C.GoString(cacheOut)
	if outString == "" {
		objectValue, err := ListObjectFromBucketS3(namespace, bucketClaimName)
		if err != nil {
			return nil, err
		}
		byteList, err := json.Marshal(objectValue)
		if err != nil {
			return nil, err
		}
		C.insr(cacheKey, C.CString(string(byteList)))
		return objectValue, nil
	}
	err := json.Unmarshal([]byte(outString), &keys)
	return keys, err
}
