package main

import (
	"github.com/gin-gonic/gin"
	"main/impleOPR"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("api/opr/push", impleOPR.PushGetObj)
	return r
}
