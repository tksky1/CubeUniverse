package main

import (
	"CubeUniverse/universalFuncs"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"log"
	"time"
)

// UniverseVersion CubeUniverse版本号
const UniverseVersion = "dev0.1"

var clientSet *kubernetes.Clientset
var dynamicClient *dynamic.DynamicClient

func main() {
	//TODO:删掉测试内容
	test()

	log.SetPrefix("[ControlBackend]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("正在加载ControlBackend", UniverseVersion)

	clientSet = universalFuncs.GetClientSet()
	dynamicClient = universalFuncs.GetDynamicClient()

	for { //等待osd加载完成
		_, _, _, _, osdStat := universalFuncs.CheckCephComponent(clientSet)
		if osdStat {
			break
		}
		time.Sleep(5 * time.Second)
	}

	// 初始化ceph-api账号
	err := SetCubeUniverseAccount()
	if err != nil {
		log.Println(err)
	}

	// 获取ceph-token为后期数据做准备
	err = GetCephToken()
	if err != nil {
		log.Println(err)
	}

	// 后端内容...

}
