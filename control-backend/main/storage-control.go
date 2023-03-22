package main

import (
	"CubeUniverse/universalFuncs"
	"bytes"
	"context"
	"errors"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"regexp"
	"strconv"
)

// 块存储/文件存储/对象存储相关控制

//<-------块存储-------->

func CheckBlockStorage() bool {
	storage, err := clientSet.StorageV1().StorageClasses().Get(context.Background(), "cubeuniverse-block-storage", v1.GetOptions{})
	if storage == nil || err != nil {
		return false
	}
	return true
}

func CreateBlockStorage() error {
	if CheckBlockStorage() {
		return errors.New("块存储已存在！")
	}
	err := universalFuncs.ApplyCrdFromYaml(universalFuncs.GetParentDir()+"/deployment/storage/block-storageclass.yaml", "", clientSet, dynamicClient)
	return err
}

// ApplyBlockPVC 为客户创建/更新块存储PVC，指定名字（必须全小写、只能用-.隔开）、命名空间、申请容量（整数，GB）
// 很可能会出现命名空间不存在等err，要正确处理告知前端
func ApplyBlockPVC(name string, namespace string, volume int) error {
	match, _ := regexp.MatchString("[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*", name)
	if !match {
		return errors.New("输入的pvc名字不合法！请使用全英文小写，用-或.隔开")
	}
	pvcBytes, err := os.ReadFile(universalFuncs.GetParentDir() + "/deployment/consumeTemplate/block-pvc.yaml")
	if err != nil {
		return err
	}
	pvcBytes = bytes.Replace(pvcBytes, []byte("sample-pvc"), []byte(name), 1)
	pvcBytes = bytes.Replace(pvcBytes, []byte("sample-namespace"), []byte(namespace), 1)
	pvcBytes = bytes.Replace(pvcBytes, []byte("1Gi"), append([]byte(strconv.Itoa(volume)), 'G', 'i'), 1)
	err = universalFuncs.ApplyBytes(pvcBytes, namespace)
	return err
}

// ListBlockSystemPVC 返回所有通过CubeUniverse创建的块存储PVC列表，内含名字、命名空间、容量等
func ListBlockSystemPVC() ([]corev1.PersistentVolumeClaim, error) {
	selector := labels.SelectorFromSet(map[string]string{"pvc-provider": "cubeuniverse", "pvc-type": "block"})
	listPVC, err := clientSet.CoreV1().PersistentVolumeClaims("").List(context.TODO(), v1.ListOptions{LabelSelector: selector.String()})
	pvcs := listPVC.Items
	return pvcs, err
}

//<--------文件存储--------->

func CheckFileSystemStorage() bool {
	storage, err := clientSet.StorageV1().StorageClasses().Get(context.Background(), "cubeuniverse-fs-storage", v1.GetOptions{})
	if storage == nil || err != nil {
		return false
	}
	return true
}

func CreateFileSystemStorage() error {
	if CheckFileSystemStorage() {
		return errors.New("文件系统存储已存在！")
	}
	err := universalFuncs.ApplyCrdFromYaml(universalFuncs.GetParentDir()+"/deployment/storage/filesystem-storageclass.yaml", "", clientSet, dynamicClient)
	return err
}

// ApplyFileSystemPVC 为客户创建/更新文件存储PVC，指定名字（必须全小写、只能用-.隔开）、命名空间、申请容量（整数，GB）
// 很可能会出现名称已存在/命名空间不存在等err，要正确处理告知前端
func ApplyFileSystemPVC(name string, namespace string, volume int) error {
	match, _ := regexp.MatchString("[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*", name)
	if !match {
		return errors.New("输入的pvc名字不合法！请使用全英文小写，用-或.隔开")
	}
	pvcBytes, err := os.ReadFile(universalFuncs.GetParentDir() + "/deployment/consumeTemplate/fs-pvc.yaml")
	if err != nil {
		return err
	}
	pvcBytes = bytes.Replace(pvcBytes, []byte("sample-pvc"), []byte(name), 1)
	pvcBytes = bytes.Replace(pvcBytes, []byte("sample-namespace"), []byte(namespace), 1)
	pvcBytes = bytes.Replace(pvcBytes, []byte("1Gi"), append([]byte(strconv.Itoa(volume)), 'G', 'i'), 1)
	err = universalFuncs.ApplyBytes(pvcBytes, namespace)
	return err
}

// ListFileSystemPVC 返回所有通过CubeUniverse创建的文件存储PVC列表，内含名字、命名空间、容量等
func ListFileSystemPVC() ([]corev1.PersistentVolumeClaim, error) {
	selector := labels.SelectorFromSet(map[string]string{"pvc-provider": "cubeuniverse", "pvc-type": "filesystem"})
	listPVC, err := clientSet.CoreV1().PersistentVolumeClaims("").List(context.TODO(), v1.ListOptions{LabelSelector: selector.String()})
	pvcs := listPVC.Items
	return pvcs, err
}

//<--------对象存储-------->

func CheckObjectStorage() bool {
	storage, err := clientSet.StorageV1().StorageClasses().Get(context.Background(), "cubeuniverse-obj-storage", v1.GetOptions{})
	if storage == nil || err != nil {
		return false
	}
	return true
}

func CreateObjectStorage() error {
	if CheckObjectStorage() {
		return errors.New("对象存储已存在！")
	}
	err := universalFuncs.ApplyCrdFromYaml(universalFuncs.GetParentDir()+"/deployment/storage/object-storageclass.yaml", "", clientSet, dynamicClient)
	return err
}

// ApplyObjectBucket 为客户创建/更新对象存储的bucket，指定名字（必须全小写、只能用-.隔开）、命名空间
// 很可能会出现名称已存在/命名空间不存在等err，要正确处理告知前端
func ApplyObjectBucket(name string, namespace string) error {
	match, _ := regexp.MatchString("[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*", name)
	if !match {
		return errors.New("输入的对象桶名字不合法！请使用全英文小写，用-或.隔开")
	}
	pvcBytes, err := os.ReadFile(universalFuncs.GetParentDir() + "/deployment/consumeTemplate/obj-bucket-claim.yaml")
	if err != nil {
		return err
	}
	pvcBytes = bytes.ReplaceAll(pvcBytes, []byte("sample-bucket"), []byte(name))
	pvcBytes = bytes.Replace(pvcBytes, []byte("sample-namespace"), []byte(namespace), 1)
	err = universalFuncs.ApplyCrdFromBytes(pvcBytes, namespace, clientSet, dynamicClient)
	return err
}

// DeleteObjectBucket 删除对象桶声明
func DeleteObjectBucket(name string, namespace string) error {
	return DeleteCRD("objectbucket.io", "v1alpha1", "objectbucketclaims", name, namespace)
}

// ListBlockPVC 返回所有通过CubeUniverse创建的块存储bucket-Claim列表，内含名字、命名空间等
func ListBlockPVC() ([]corev1.PersistentVolumeClaim, error) {
	selector := labels.SelectorFromSet(map[string]string{"pvc-provider": "cubeuniverse", "pvc-type": "object"})
	listPVC, err := clientSet.CoreV1().PersistentVolumeClaims("").List(context.TODO(), v1.ListOptions{LabelSelector: selector.String()})
	pvcs := listPVC.Items
	return pvcs, err
}

//<------其他/通用------>

// DeletePVC 删除块存储/文件存储PVC，删除后用户使用pvc的数据将被删除，应提示用户
func DeletePVC(name string, namespace string) error {
	err := clientSet.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), name, v1.DeleteOptions{})
	return err
}

// DeleteCRD 用于删除CRD，比如对象存储对应的bucket-claim
func DeleteCRD(group string, version string, resource string, name string, namespace string) error {
	crdMeta := schema.GroupVersionResource{Group: group, Version: version, Resource: resource}
	err := dynamicClient.Resource(crdMeta).Namespace(namespace).Delete(context.TODO(), name, v1.DeleteOptions{})
	return err
}