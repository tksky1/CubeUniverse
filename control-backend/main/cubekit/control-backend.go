package cubekit

import (
	"CubeUniverse/universalFuncs"
	"log"
	"time"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// UniverseVersion CubeUniverse版本号
const UniverseVersion = "dev0.1"

var ClientSet *kubernetes.Clientset
var DynamicClient *dynamic.DynamicClient

<<<<<<< HEAD:control-backend/main/cubekit/control-backend.go
func Init() {

=======
func main() {
>>>>>>> parent of 594495b (chore: add test function):control-backend/main/control-backend.go
	log.SetPrefix("[ControlBackend]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("正在加载ControlBackend", UniverseVersion)

	ClientSet = universalFuncs.GetClientSet()
	DynamicClient = universalFuncs.GetDynamicClient()

	for { //等待osd加载完成
		_, _, _, _, osdStat := universalFuncs.CheckCephComponent(ClientSet)
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

}
