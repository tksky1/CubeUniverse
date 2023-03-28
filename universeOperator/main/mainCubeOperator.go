package main

import (
	"CubeUniverse/universalFuncs"
	"github.com/google/uuid"
	"log"
	kit "main/cubeOperatorKit"
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
	//开启协程运行监听web
	go webInit()
	//剩下完成初始化
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

}
