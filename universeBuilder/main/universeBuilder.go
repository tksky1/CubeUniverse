// @Description  UniverseBuilder 自动部署CubeUniverse组件
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
	log.SetPrefix("[UniverseBuilder]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("正在加载UniverseBuilder", UniverseVersion)
	clientSet = universalFuncs.GetClientSet()
	dynamicClient = universalFuncs.GetDynamicClient()
	buildCube()
	buildCeph()
}

// TO DO: 调试完需把buildCube的println改成panic

// 启动CubeUniverse组件
func buildCube() {
	log.Println("启动CubeUniverse组件..")
	operator, dashboard, controlBackend := universalFuncs.CheckCubeUniverseComponent(clientSet)
	if !operator {
		log.Println("启动CubeUniverse-operator..")
		err := universalFuncs.ApplyYaml(universalFuncs.GetParentDir() + "/deployment/UniverseOperator.yml")
		if err != nil {
			log.Println("启动UniverseOperator失败，请检查CubeUniverse项目文件是否完好！\n", err)
		}
	}

	if !dashboard {
		log.Println("启动CubeUniverse-DashBoard..")
		err := universalFuncs.ApplyYaml(universalFuncs.GetParentDir() + "/deployment/UniverseDashBoard.yml")
		if err != nil {
			log.Println("启动UniverseDashBoard失败，请检查CubeUniverse项目文件是否完好！\n", err)
		}
	}

	if !controlBackend {
		log.Println("启动CubeUniverse-operator..")
		err := universalFuncs.ApplyYaml(universalFuncs.GetParentDir() + "/deployment/UniverseOperator.yml")
		if err != nil {
			log.Println("启动UniverseOperator失败，请检查CubeUniverse项目文件是否完好！\n", err)
		}
	}
}

func buildCeph() {
	log.Println("启动ceph组件..")
	operator, rbdplugin, mon, mgr, osd := universalFuncs.CheckCephComponent(clientSet)
	if !operator {

	}

}
