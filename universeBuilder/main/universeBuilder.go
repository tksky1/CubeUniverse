// @Description  UniverseBuilder 自动检测和部署CubeUniverse缺失的组件
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
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

// UniverseVersion CubeUniverse版本号
const UniverseVersion = "dev0.1"

var clientSet *kubernetes.Clientset
var dynamicClient *dynamic.DynamicClient
var startTime time.Time

func main() {
	log.SetPrefix("[UniverseBuilder]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("正在加载UniverseBuilder", UniverseVersion)
	clientSet = universalFuncs.GetClientSet()
	dynamicClient = universalFuncs.GetDynamicClient()
	startTime = time.Now()
	log.Println("启动CubeUniverse组件..")
	for !buildCube() {
		if time.Now().Sub(startTime) > time.Second*120 {
			log.Println("搭建CubeUniverse组件已超过120秒，请检查集群网络、组件健康情况！")
		}
		time.Sleep(3 * time.Second)
	}
	log.Println("CubeUniverse组件已正常运行，启动ceph组件..")
	startTime = time.Now()
	for !buildCeph() {
		time.Sleep(5 * time.Second)
		if time.Now().Sub(startTime) > time.Minute*30 {
			log.Println("搭建ceph集群已超过30分钟，请检查本log和集群健康情况！")
		}
	}
	log.Println("Ceph已正常运行.")
}

// TODO: 调试完需把buildCube的println改成panic

// 启动CubeUniverse组件
func buildCube() (ret bool) {

	operator, dashboard, controlBackend, _ := universalFuncs.CheckCubeUniverseComponent(clientSet)
	ret = true
	if !operator {
		log.Println("启动CubeUniverse-operator..")
		err := universalFuncs.ApplyYaml(universalFuncs.GetParentDir()+"/deployment/UniverseOperator.yml", "cubeuniverse")
		if err != nil {
			log.Println("启动UniverseOperator失败，请检查CubeUniverse项目文件是否完好！\n", err)
		}
		ret = false
	}

	if !dashboard {
		log.Println("启动CubeUniverse-DashBoard..")
		err := universalFuncs.ApplyYaml(universalFuncs.GetParentDir()+"/deployment/UniverseDashBoard.yml", "cubeuniverse")
		if err != nil {
			log.Println("启动UniverseDashBoard失败，请检查CubeUniverse项目文件是否完好！\n", err)
		}
		ret = false
	}

	if !controlBackend {
		log.Println("启动CubeUniverse-controlBackend..")
		err := universalFuncs.ApplyYaml(universalFuncs.GetParentDir()+"/deployment/ControlBackend.yml", "cubeuniverse")
		if err != nil {
			log.Println("启动controlBackend失败，请检查CubeUniverse项目文件是否完好！\n", err)
		}
		ret = false
	}
	return true //TODO：调试完改成False
}

// 启动ceph组件
func buildCeph() (ret bool) {

	operator, rbdplugin, mon, mgr, osd := universalFuncs.CheckCephComponent(clientSet)
	if !operator {
		log.Println("启动ceph-operator..")
		err := universalFuncs.ApplyCrdFromYaml(universalFuncs.GetParentDir()+"/deployment/storage/crds.yaml", "", clientSet, dynamicClient)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)
		err = universalFuncs.ApplyCrdFromYaml(universalFuncs.GetParentDir()+"/deployment/storage/common.yaml", "rook-ceph", clientSet, dynamicClient)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)
		err = universalFuncs.ApplyCrdFromYaml(universalFuncs.GetParentDir()+"/deployment/storage/operator.yaml", "rook-ceph", clientSet, dynamicClient)
		if err != nil {
			log.Fatal(err)
		}
		return false
	}

	if !rbdplugin {
		log.Println("启动ceph-cluster..")
		err := universalFuncs.ApplyYaml(universalFuncs.GetParentDir()+"/deployment/storage/cluster.yaml", "rook-ceph")
		if err != nil {
			log.Println("启动ceph-cluster失败，请检查CubeUniverse项目文件是否完好！\n", err)
		}

		err = universalFuncs.ApplyYaml(universalFuncs.GetParentDir()+"/deployment/storage/toolbox.yaml", "rook-ceph")
		if err != nil {
			log.Println("启动ceph-toolbox失败，请检查CubeUniverse项目文件是否完好！\n", err)
		}

		return false
	}

	if !mon {
		log.Println("ceph-monitor未启动，等待..")
		if time.Now().Sub(startTime) > time.Minute*20 {
			log.Println("ceph已开始构建超过20分钟，monitor仍未启动，请删除所有节点上的/var/lib/rook文件夹并重新安装集群！")
		}
		return false
	}

	if !mgr {
		log.Println("ceph-mgr未启动，等待..")
		if time.Now().Sub(startTime) > time.Minute*20 {
			log.Println("ceph已开始构建超过20分钟，mgr仍未启动，请删除所有节点上的/var/lib/rook文件夹并重新安装集群！")
		}
		return false
	}

	if !osd {
		log.Println("ceph-osd未启动，等待..")
		if time.Now().Sub(startTime) > time.Minute*30 {
			log.Println("ceph已开始构建超过30分钟，osd仍未启动，请确保节点都已安装一个没有文件系统的空磁盘，并重新安装集群！")
		}
		return false
	}

	//TODO: 删除这些调试用代码
	err := CreateFileSystemStorage()
	//if err != nil {
	//	panic(err)
	//}
	err = createObjectStorage()
	//if err != nil {
	//	panic(err)
	//}
	err = applyFileSystemPVC("test-fspvc", "cubeuniverse", 10)
	if err != nil {
		panic(err)
	}
	err = applyObjectBucket("test-obj-bucket", "cubeuniverse")
	if err != nil {
		panic(err)
	}
	err = DeleteCRD("objectbucket.io", "v1alpha1", "objectbucketclaims", "test-obj-bucket-claim", "cubeuniverse")
	if err != nil {
		panic(err)
	}
	pvclist, err := ListFileSystemPVC()
	for _, pvc := range pvclist {
		print(pvc.Name, " ", pvc.Spec.Resources.Requests.Storage().String())
	}
	if err != nil {
		panic(err)
	}

	return true
}

// TODO：调试完删掉
func applyFileSystemPVC(name string, namespace string, volume int) error {
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

func applyObjectBucket(name string, namespace string) error {
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

func createObjectStorage() error {
	err := universalFuncs.ApplyCrdFromYaml(universalFuncs.GetParentDir()+"/deployment/storage/object-storageclass.yaml", "", clientSet, dynamicClient)
	return err
}

func CreateFileSystemStorage() error {
	err := universalFuncs.ApplyCrdFromYaml(universalFuncs.GetParentDir()+"/deployment/storage/filesystem-storageclass.yaml", "", clientSet, dynamicClient)
	return err
}

func DeleteCRD(group string, version string, resource string, name string, namespace string) error {
	crdMeta := schema.GroupVersionResource{Group: group, Version: version, Resource: resource}
	err := dynamicClient.Resource(crdMeta).Namespace(namespace).Delete(context.TODO(), name, v1.DeleteOptions{})
	return err
}

func ListFileSystemPVC() ([]corev1.PersistentVolumeClaim, error) {
	selector := labels.SelectorFromSet(map[string]string{"pvc-provider": "cubeuniverse", "pvc-type": "filesystem"})
	listPVC, err := clientSet.CoreV1().PersistentVolumeClaims("").List(context.TODO(), v1.ListOptions{LabelSelector: selector.String()})
	pvcs := listPVC.Items
	return pvcs, err
}
