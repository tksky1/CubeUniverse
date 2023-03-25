package storageact

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 统一返回格式
func Response(ctx *gin.Context, httpStatus int, code int, data gin.H, msg string) {
	ctx.JSON(httpStatus, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

func Success(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusOK, 200, data, msg)
}

func Fail(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusBadRequest, 400, data, msg)
}

// 数据格式有误406
func FailUnac(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusNotAcceptable, 406, data, msg)
}

// 创建/删除命名空间出错405
func FailCrea(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusMethodNotAllowed, 405, data, msg)
}
