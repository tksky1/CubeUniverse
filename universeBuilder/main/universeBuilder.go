// @Description  UniverseBuilder 自动检测和部署CubeUniverse缺失的组件
package main

import (
	"CubeUniverse/universalFuncs"
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"log"
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
		time.Sleep(7 * time.Second)
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
		err := universalFuncs.CreateYaml(universalFuncs.GetParentDir()+"/deployment/UniverseOperator.yml", "cubeuniverse")
		if err != nil {
			log.Println("启动UniverseOperator失败，请检查CubeUniverse项目文件是否完好！\n", err)
		}
		ret = false
	}

	if !dashboard {
		log.Println("启动CubeUniverse-DashBoard..")
		err := universalFuncs.CreateYaml(universalFuncs.GetParentDir()+"/deployment/UniverseDashBoard.yml", "cubeuniverse")
		if err != nil {
			log.Println("启动UniverseDashBoard失败，请检查CubeUniverse项目文件是否完好！\n", err)
		}
		ret = false
	}

	if !controlBackend {
		log.Println("启动CubeUniverse-controlBackend..")
		err := universalFuncs.CreateYaml(universalFuncs.GetParentDir()+"/deployment/ControlBackend.yml", "cubeuniverse")
		if err != nil {
			log.Println("启动controlBackend失败，请检查CubeUniverse项目文件是否完好！\n", err)
		}
		ret = false
	}
	return true //TODO：调试完改成ret
}

// 启动ceph组件
func buildCeph() (ret bool) {

	operator, rbdplugin, mon, mgr, osd := universalFuncs.CheckCephComponent(clientSet)
	if !operator {
		log.Println("启动ceph-operator..")
		err := universalFuncs.PatchCrdFromYaml(universalFuncs.GetParentDir()+"/deployment/storage/crds.yaml", "", clientSet, dynamicClient)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)
		err = universalFuncs.PatchCrdFromYaml(universalFuncs.GetParentDir()+"/deployment/storage/common.yaml", "rook-ceph", clientSet, dynamicClient)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)
		err = universalFuncs.PatchCrdFromYaml(universalFuncs.GetParentDir()+"/deployment/storage/operator.yaml", "rook-ceph", clientSet, dynamicClient)
		if err != nil {
			log.Fatal(err)
		}
		return false
	}

	if !rbdplugin {
		log.Println("启动ceph-cluster..")
		err := universalFuncs.PatchYaml(universalFuncs.GetParentDir()+"/deployment/storage/cluster.yaml", "rook-ceph")
		if err != nil {
			log.Println("启动ceph-cluster失败，请检查CubeUniverse项目文件是否完好！\n", err)
		}

		err = universalFuncs.PatchYaml(universalFuncs.GetParentDir()+"/deployment/storage/toolbox.yaml", "rook-ceph")
		if err != nil {
			log.Println("启动ceph-toolbox失败，请检查CubeUniverse项目文件是否完好！\n", err)
		}

		err = universalFuncs.PatchYaml(universalFuncs.GetParentDir()+"/deployment/storage/dashboard-external-https.yaml", "rook-ceph")
		if err != nil {
			log.Println("启动ceph-dashboard失败，请检查CubeUniverse项目文件是否完好！\n", err)
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

	if !universalFuncs.CheckMysqlStat(clientSet) {
		storage, err := clientSet.StorageV1().StorageClasses().Get(context.Background(), "cubeuniverse-fs-storage", v1.GetOptions{})
		if storage == nil || err != nil {
			err = universalFuncs.PatchCrdFromYaml(universalFuncs.GetParentDir()+"/deployment/storage/filesystem-storageclass.yaml",
				"", clientSet, dynamicClient)
			if err != nil {
				log.Println("启动文件存储失败，请检查CubeUniverse项目文件是否完好！\n", err)
			}
			time.Sleep(1 * time.Second)
		}
		err = universalFuncs.PatchCrdFromYaml(universalFuncs.GetParentDir()+"/deployment/storage/mysql-pre.yaml",
			"cubeuniverse", clientSet, dynamicClient)
		if err != nil {
			log.Println("启动mysql-pre失败，请检查CubeUniverse项目文件是否完好！\n", err)
		}
		err = universalFuncs.PatchCrdFromYaml(universalFuncs.GetParentDir()+"/deployment/storage/mysql.yaml",
			"cubeuniverse", clientSet, dynamicClient)
		if err != nil {
			log.Println("启动mysql失败，请检查CubeUniverse项目文件是否完好！\n", err)
		}
		return false
	}

	return true
}
