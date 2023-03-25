package storageact

import (
	"control-backend/apikit"
	"github.com/gin-gonic/gin"
	"strings"
)

// 处理的get请求的Header必须由X-type字段和X-action字段
func CheckOrCreate(ctx *gin.Context) {
	//根据get请求的type与action数据来选择调用函数
	get_type := ctx.GetHeader("X-type")
	get_action := ctx.GetHeader("X-action")
	get_type = strings.ToLower(get_type)
	get_action = strings.ToLower(get_action)
	switch get_type {
	//如果是块存储的请求
	case "block":
		if get_action == "check" {
			checkFunc(get_type, apikit.CheckBlockStorage, ctx) //检查
		} else if get_action == "create" {
			createFunc(get_type, apikit.CreateBlockStorage, ctx) //创建
		}
	case "file":
		if get_action == "check" {
			checkFunc(get_type, apikit.CheckFileSystemStorage, ctx)
		} else if get_action == "create" {
			createFunc(get_type, apikit.CreateFileSystemStorage, ctx)
		}
	case "object":
		if get_action == "check" {
			checkFunc(get_type, apikit.CheckObjectStorage, ctx)
		} else if get_action == "create" {
			createFunc(get_type, apikit.CreateObjectStorage, ctx)
		}
	}
	//完成三种存储模式的检查与创建工作
}

// 处理的请求必须带有name、namespace、volume字段,均为string
func Pvcreq(ctx *gin.Context) {
	//根据post请求body体参数来解析数据,支持postform和json两种格式
	name := ctx.PostForm("name")
	namespace := ctx.PostForm("namespace")
	volume := ctx.PostForm("volume")

	if name == "" && namespace == "" && volume == "" {
		json := make(map[string]interface{})
		ctx.BindJSON(&json)
		name = json["name"].(string)
		namespace = json["namespace"].(string)
		volume = json["volume"].(string)
	}

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
