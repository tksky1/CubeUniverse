package cubeControl

import (
	"CubeUniverse/universalFuncs"
	"log"
	"time"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// UniverseVersion CubeUniverse版本号
const UniverseVersion = "0.1alpha"

var ClientSet *kubernetes.Clientset
var DynamicClient *dynamic.DynamicClient

func Init() {

	log.SetPrefix("[ControlBackend]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("正在加载ControlBackend", UniverseVersion)

	log.Println("等待文件存储和sql服务启动..")
	for { //等待postgresql加载完成
		mysql := universalFuncs.CheckMysqlStat(ClientSet)
		if mysql {
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

}
