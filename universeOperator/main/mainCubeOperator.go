package main

import (
	"CubeUniverse/universalFuncs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	kit "main/cubeOperatorKit"
	"time"
)

var port string = "8890"

func webInit() {
	var r *gin.Engine = gin.Default()
	r = CollectRoute(r) //一次性注册完路由
	//选择监听端口
	panic(r.Run(":" + port))
}

// UniverseVersion CubeUniverse版本号
const UniverseVersion = "dev0.1"

var UUID = uuid.New().String()

func main() {
	//测试用:TODO
	webInit()
	//记得删除

	kit.Init()

	// pod互斥锁，保证同时只有一个operator执行功能
	for {
		locked, _, lockTime := universalFuncs.CheckInUse(kit.ClientSet, "operator-mutex")
		if !locked || time.Now().Sub(lockTime).Seconds() > 5 {
			universalFuncs.SetInUse(kit.ClientSet, "operator-mutex", UUID)
			break
		}
		time.Sleep(3 * time.Second)
	}
	// 启动心跳go程
	go universalFuncs.HeartBeat(kit.ClientSet, "operator-mutex", UUID)

	println("\n ██████╗██╗   ██╗██████╗ ███████╗██╗   ██╗███╗   ██╗██╗██╗   ██╗███████╗██████╗ ███████╗███████╗\n██╔════╝██║   ██║██╔══██╗██╔════╝██║   ██║████╗  ██║██║██║   ██║██╔════╝██╔══██╗██╔════╝██╔════╝\n██║     ██║   ██║██████╔╝█████╗  ██║   ██║██╔██╗ ██║██║██║   ██║█████╗  ██████╔╝███████╗█████╗  \n██║     ██║   ██║██╔══██╗██╔══╝  ██║   ██║██║╚██╗██║██║╚██╗ ██╔╝██╔══╝  ██╔══██╗╚════██║██╔══╝  \n╚██████╗╚██████╔╝██████╔╝███████╗╚██████╔╝██║ ╚████║██║ ╚████╔╝ ███████╗██║  ██║███████║███████╗\n ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝\n")

	log.SetPrefix("[UniverseOperator]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("正在加载UniverseOperator", UniverseVersion)

	//开启协程运行监听web
	go webInit()

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
