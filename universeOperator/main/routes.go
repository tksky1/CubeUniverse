package main

import (
	"github.com/gin-gonic/gin"
	"main/impleOPR"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("api/opr/pushget", impleOPR.PushGetDeleteListObj)
	r.GET("api/opr/wspushget", impleOPR.ConstPushGetDeleteList)

	return r
}
