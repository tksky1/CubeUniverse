// @Description  UniverseOperator 常驻运行 监测集群状态和自动维护
package main

import (
	"CubeUniverse/universalFuncs"
	"k8s.io/client-go/kubernetes"
	"log"
	"time"
)

// UniverseVersion CubeUniverse版本号
const UniverseVersion = "dev0.1"

var clientSet *kubernetes.Clientset

func main() {
	println("\n ██████╗██╗   ██╗██████╗ ███████╗██╗   ██╗███╗   ██╗██╗██╗   ██╗███████╗██████╗ ███████╗███████╗\n██╔════╝██║   ██║██╔══██╗██╔════╝██║   ██║████╗  ██║██║██║   ██║██╔════╝██╔══██╗██╔════╝██╔════╝\n██║     ██║   ██║██████╔╝█████╗  ██║   ██║██╔██╗ ██║██║██║   ██║█████╗  ██████╔╝███████╗█████╗  \n██║     ██║   ██║██╔══██╗██╔══╝  ██║   ██║██║╚██╗██║██║╚██╗ ██╔╝██╔══╝  ██╔══██╗╚════██║██╔══╝  \n╚██████╗╚██████╔╝██████╔╝███████╗╚██████╔╝██║ ╚████║██║ ╚████╔╝ ███████╗██║  ██║███████║███████╗\n ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝\n")

	log.SetPrefix("[UniverseOperator]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("正在加载UniverseOperator", UniverseVersion)

	clientSet = universalFuncs.GetClientSet()

	for {
		time.Sleep(5 * time.Second)
		operator, dashboard, controlBackend, builder := universalFuncs.CheckCubeUniverseComponent(clientSet)
		if builder {
			continue
		}
		cephOperator, rbdplugin, mon, mgr, osd := universalFuncs.CheckCephComponent(clientSet)
		if !(operator && dashboard && controlBackend && cephOperator && rbdplugin && mon && mgr && osd) {
			log.Println("监测到集群未完全运行，启动UniverseBuilder..")
			err := universalFuncs.ApplyYaml(universalFuncs.GetParentDir()+"/deployment/UniverseBuilder.yml", "cubeuniverse")
			if err != nil {
				log.Panic("启动UniverseBuilder失败，请检查CubeUniverse项目文件是否完好！\n", err)
			}
			time.Sleep(15 * time.Second)
		}
	}

}
