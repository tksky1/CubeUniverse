package main

import (
	"control-backend/interact-front/storageact"
	"control-backend/login-kit/controller"
	"control-backend/login-kit/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	// r.POST("api/auth/register", controller.Register)//不再提供注册服务
	r.POST("api/auth/login", controller.Login)
	r.POST("api/auth/altpas", controller.AlterPasswd)
	r.GET("api/auth/info", middleware.AuthMiddleware(), controller.Info)
	r.GET("api/storiage/check", middleware.AuthMiddleware(), storageact.CheckOrCreate) //检测是否已经某种块存储
	r.POST("api/storiage/pvc", middleware.AuthMiddleware())                            //用户pvc调用
	return r
}
