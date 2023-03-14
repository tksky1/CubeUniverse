package main

import (
	"CubeUniverse/universalFuncs"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"log"
)

// UniverseVersion CubeUniverse版本号
const UniverseVersion = "dev0.1"

var clientSet *kubernetes.Clientset
var dynamicClient *dynamic.DynamicClient

func main() {
	log.SetPrefix("[ControlBackend]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("正在加载ControlBackend", UniverseVersion)

	clientSet = universalFuncs.GetClientSet()
	dynamicClient = universalFuncs.GetDynamicClient()

	// 后端内容...

}
