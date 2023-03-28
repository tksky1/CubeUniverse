package storageact

import (
	"control-backend/cubeControl"
	"github.com/gin-gonic/gin"
	"strings"
)

// 处理的get请求的Header必须由X-type字段和X-action字段
func CheckStorage(ctx *gin.Context) {
	//根据get请求的type与action数据来选择调用函数
	get_type := ctx.GetHeader("X-type")
	get_action := ctx.GetHeader("X-action")
	get_type = strings.ToLower(get_type)
	get_action = strings.ToLower(get_action)

	//测试用:TODO
	Success(ctx, nil, strings.ToUpper(string(get_type[0]))+"CS") //返回创建成功 block/file created success缩写BCS,FCS
	return
	//记得删除

	switch get_type {
	//如果是块存储的请求
	case "block":
		if get_action == "check" {
			checkFunc(get_type, cubeControl.CheckBlockStorage, ctx) //检查
		} else if get_action == "create" {
			createFunc(get_type, cubeControl.CreateBlockStorage, ctx) //创建
		}
	case "file":
		if get_action == "check" {
			checkFunc(get_type, cubeControl.CheckFileSystemStorage, ctx)
		} else if get_action == "create" {
			createFunc(get_type, cubeControl.CreateFileSystemStorage, ctx)
		}
	case "object":
		if get_action == "check" {
			checkFunc(get_type, cubeControl.CheckObjectStorage, ctx)
		} else if get_action == "create" {
			createFunc(get_type, cubeControl.CreateObjectStorage, ctx)
		}
	}
	//完成三种存储模式的检查与创建工作
}

// 内部工具方法
func createFunc(creatype string, imple func() error, ctx *gin.Context) {
	err := imple()  //调用创建函数，返回错误信息
	if err != nil { //如果有错误不为空则发送400并在msg中加入错误描述信息
		if strings.HasPrefix(strings.ToLower(err.Error()), creatype+" alrea") {
			Fail(ctx, nil, strings.ToUpper(string(creatype[0]))+"AE") //BAE FAE 等缩写
		} else {
			Fail(ctx, nil, err.Error()) //其他错误返回错误信息即可
		}
		return
	}
	Success(ctx, nil, strings.ToUpper(string(creatype[0]))+"CS") //返回创建成功 block/file created success缩写BCS,FCS
}

func checkFunc(checktype string, imple func() bool, ctx *gin.Context) {
	res := imple()
	if res {
		Success(ctx, nil, "provide "+checktype) //如果提供在返回json的msg字段中加入provide block
		return
	}
	Success(ctx, nil, "not provided") //如果不提供则返回not provided
}
