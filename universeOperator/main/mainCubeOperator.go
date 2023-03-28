package main

import (
	"github.com/gin-gonic/gin"
	kit "main/cubeOperatorKit"
)

var port string = "8888"

func webInit() {
	var r *gin.Engine = gin.Default()
	r = CollectRoute(r) //一次性注册完路由
	//选择监听端口
	panic(r.Run(":" + port))
}

func main() {
	//开启协程运行监听web
	go webInit()
	//剩下完成初始化
	kit.Init()
	
}
