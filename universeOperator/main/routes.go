package main

import (
	"github.com/gin-gonic/gin"
	"main/impleOPR"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("osspgdl", impleOPR.PutGetDeleteListObj)
	r.GET("osswspgdl", impleOPR.ConstPushGetDeleteList)
	//传统restful请求
	r.GET("oss", impleOPR.OssGet)
	r.POST("oss", impleOPR.OssPut)
	r.DELETE("oss", impleOPR.OssDelete)
	r.GET("osslist", impleOPR.OssList)
	return r
}
