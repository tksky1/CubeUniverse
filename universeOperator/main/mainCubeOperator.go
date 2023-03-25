package main

import (
	"CubeUniverse/universalFuncs"
	"log"
	kit "main/cubeOperatorKit"
	"time"
)

func main() {
	kit.Init()
	for {
		time.Sleep(5 * time.Second)
		operator, dashboard, controlBackend, builder := universalFuncs.CheckCubeUniverseComponent(kit.ClientSet)
		if builder {
			continue
		}
		cephOperator, rbdplugin, mon, mgr, osd := universalFuncs.CheckCephComponent(kit.ClientSet)
		if !(operator && dashboard && controlBackend && cephOperator && rbdplugin && mon && mgr && osd) {
			log.Println("监测到集群未完全运行，启动UniverseBuilder..")
			err := universalFuncs.PatchYaml(universalFuncs.GetParentDir()+"/deployment/UniverseBuilder.yml", "cubeuniverse")
			if err != nil {
				log.Panic("启动UniverseBuilder失败，请检查CubeUniverse项目文件是否完好！\n", err)
			}
			time.Sleep(15 * time.Second)
		}
	}

}
