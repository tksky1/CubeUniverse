// @Description  CubeUniverse入口
package main

import (
	"CubeUniverse/universalFuncs"
	"k8s.io/client-go/kubernetes"
	"log"
)

// UniverseVersion CubeUniverse版本号
const UniverseVersion = "dev0.1"

var clientSet *kubernetes.Clientset

// 主函数。从这里开始运行
func main() {

	// 字符画by patorjk.com
	println("\n ██████╗██╗   ██╗██████╗ ███████╗██╗   ██╗███╗   ██╗██╗██╗   ██╗███████╗██████╗ ███████╗███████╗\n██╔════╝██║   ██║██╔══██╗██╔════╝██║   ██║████╗  ██║██║██║   ██║██╔════╝██╔══██╗██╔════╝██╔════╝\n██║     ██║   ██║██████╔╝█████╗  ██║   ██║██╔██╗ ██║██║██║   ██║█████╗  ██████╔╝███████╗█████╗  \n██║     ██║   ██║██╔══██╗██╔══╝  ██║   ██║██║╚██╗██║██║╚██╗ ██╔╝██╔══╝  ██╔══██╗╚════██║██╔══╝  \n╚██████╗╚██████╔╝██████╔╝███████╗╚██████╔╝██║ ╚████║██║ ╚████╔╝ ███████╗██║  ██║███████║███████╗\n ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝\n")
	log.SetPrefix("[CubeUniverse]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("loading CubeUniverse", UniverseVersion)

	clientSet = universalFuncs.GetClientSet()

	k8sStatus := universalFuncs.CheckClusterHealth(clientSet)
	if !k8sStatus {
		log.Fatal("错误：k8s集群有必要组件未处于运行状态，请恢复集群功能后再试！")
	}
	log.Println("k8s集群自检成功，进行CubeUniverse组件检测..")

	operator, dashboard, controlBackend, _ := universalFuncs.CheckCubeUniverseComponent(clientSet)
	if !(operator && dashboard && controlBackend) {
		log.Println("CubeUniverse组件未完全加载，启动UniverseBuilder..")
		err := universalFuncs.PatchYaml(universalFuncs.GetParentDir()+"/deployment/UniverseBuilder.yml", "cubeuniverse")
		if err != nil {
			log.Panic("启动UniverseBuilder失败，请检查CubeUniverse项目文件是否完好！\n", err)
		}
	} else {
		log.Println("CubeUniverse自检完成，已在正常运行")
	}
}
