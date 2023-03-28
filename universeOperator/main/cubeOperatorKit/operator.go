// UniverseOperator 常驻运行 监测集群状态和自动维护
package cubeOperatorKit

import (
	"CubeUniverse/universalFuncs"
	"k8s.io/client-go/kubernetes"
	"log"
	"time"
)

var ClientSet *kubernetes.Clientset

func Init() {

	ClientSet = universalFuncs.GetClientSet()
	sessionCacheMap = make(map[[16]byte]*SessionAndBucketName)
	for {
		time.Sleep(5 * time.Second)
		operator, dashboard, controlBackend, builder := universalFuncs.CheckCubeUniverseComponent(ClientSet)
		if builder {
			continue
		}
		cephOperator, rbdplugin, mon, mgr, osd := universalFuncs.CheckCephComponent(ClientSet)
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
