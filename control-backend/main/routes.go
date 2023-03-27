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
	r.GET("api/storage/check", middleware.AuthMiddleware(), storageact.CheckStorage) //检测是否已经某种块存储
	r.POST("api/storage/pvc", middleware.AuthMiddleware(), storageact.PvcCrea)       //用户pvc调用创建新存储
	r.POST("api/storage/pvcpatch", middleware.AuthMiddleware(), storageact.PvcPatch) //patch请求进行更新数据
	r.DELETE("api/storage/pvc", middleware.AuthMiddleware(), storageact.DeletePvc)   //对申请的存储删除
	r.GET("api/storage/pvc", middleware.AuthMiddleware(), storageact.ListPvc)        //查看所有申请的容器
	return r
}
