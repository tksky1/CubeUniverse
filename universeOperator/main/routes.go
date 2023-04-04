package main

import (
	"github.com/gin-gonic/gin"
	"main/impleOPR"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("api/opr/putget", impleOPR.PushGetDeleteListObj)
	r.GET("api/opr/wsputget", impleOPR.ConstPushGetDeleteList)

	return r
}
