package storageact

import (
	"control-backend/cubeControl"
	"github.com/gin-gonic/gin"
	"regexp"
	"strconv"
	"strings"
)

// PvcCrea 处理的请求必须带有name、namespace、volume字段,均为string
func PvcCrea(ctx *gin.Context) {
	//根据post请求body体参数来解析数据,支持postform和json两种格式
	var name, namespace, sVolume, actType, autoScales string
	var volume int     //用于后面volume参数格式转化
	var autoScale bool //自动扩容选项
	//根据post请求body体参数来解析数据,支持postform和json两种格式
	name = ctx.PostForm("name")
	namespace = ctx.PostForm("namespace")
	sVolume = ctx.PostForm("volume")
	actType = ctx.PostForm("X-type")
	autoScales = ctx.PostForm("autoscale")
	if name == "" && namespace == "" && sVolume == "" && actType == "" && autoScales == "" {

		json := make(map[string]interface{})
		ctx.BindJSON(&json)
		if value, ok := json["X-type"].(string); ok {
			actType = value
		} else {
			Fail(ctx, nil, "type should be string") //返回错误反馈
			return
		}
		if value, ok := json["name"].(string); ok {
			name = value
		} else {
			Fail(ctx, nil, "type not correct") //返回错误反馈
			return
		}
		if value, ok := json["namespace"].(string); ok {
			namespace = value
		} else {
			Fail(ctx, nil, "type should be string") //返回错误反馈
			return
		}
		if value, ok := json["volume"].(string); ok || strings.ToLower(actType) == "object" {
			sVolume = value
		} else {
			Fail(ctx, nil, "type should be string") //返回错误反馈
			return
		}

		if value, ok := json["autoscale"].(string); ok || strings.ToLower(actType) == "object" {
			autoScales = value
		} else {
			Fail(ctx, nil, "type should be string") //返回错误反馈
			return
		}
	}
	//
	//如果有autoscale 转化为bool
	if autoScales != "" && strings.ToLower(autoScales) == "true" {
		autoScale = true //如果为true则改为真
	} else {
		autoScale = false //其余情况---为空和为false均改为false
	}
	//如果有volume参数要进行格式转化
	if sVolume != "" {
		volumeCp, err := strconv.Atoi(sVolume)
		if err != nil {
			FailUnac(ctx, nil, "volume not number")
			return
		}
		volume = volumeCp
	}
	//名字要符合标准
	match, _ := regexp.MatchString("[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*", name)
	if !match { //如果名称不符合标准则去除
		FailUnac(ctx, nil, "format err")
		return
	}

	// //测试用记得删除
	// Success(ctx, nil, "create done")
	// return
	// //:TODO

	//判断调用方法
	switch strings.ToLower(actType) {
	case "block":
		if err := cubeControl.CreateBlockPVC(name, namespace, volume, autoScale); err != nil {
			FailCrea(ctx, nil, err.Error())
			return
		}
		Success(ctx, nil, "create done")
	case "file":
		if err := cubeControl.CreateFileSystemPVC(name, namespace, volume, autoScale); err != nil {
			FailCrea(ctx, nil, err.Error())
			return
		}
		Success(ctx, nil, "create done")
	case "object":
		sMaxObject := ctx.PostForm("maxobject")
		sMaxGBsize := ctx.PostForm("maxgbsize")
		json := make(map[string]interface{})
		ctx.BindJSON(&json)
		if sMaxObject == "" && sMaxGBsize == "" {
			if value, ok := json["maxobject"].(string); ok {
				sMaxObject = value
			} else {
				Fail(ctx, nil, "type should be string") //返回错误反馈
				return
			}
			if value, ok := json["maxgbsize"].(string); ok {
				sMaxGBsize = value
			} else {
				Fail(ctx, nil, "type should be string") //返回错误反馈
				return
			}
		}

		maxGBsize, err := strconv.Atoi(sMaxGBsize)
		if err != nil {
			FailUnac(ctx, nil, "maxgbsize not number") //传入的不是数字格式问题
			return
		}
		maxobject, err := strconv.Atoi(sMaxObject)
		if err != nil {
			FailUnac(ctx, nil, "maxobject not number")
			return
		}
		if err := cubeControl.CreateObjectBucket(name, namespace, maxobject, maxGBsize); err != nil {
			FailCrea(ctx, nil, err.Error())
			return
		}
		Success(ctx, nil, "create done")
	}

}

func PvcPatch(ctx *gin.Context) {
	var autoScales string = ""
	//根据post请求body体参数来解析数据,支持postform和json两种格式
	name := ctx.PostForm("name")
	namespace := ctx.PostForm("namespace")
	sVolume := ctx.PostForm("volume")
	autoScales = ctx.PostForm("autoscale")
	actType := ctx.PostForm("X-type")
	var autoScale bool //自动扩容选项
	var volume int     //用于后面volume参数格式转化
	if name == "" && namespace == "" && sVolume == "" && actType == "" && autoScales == "" {
		json := make(map[string]interface{})
		ctx.BindJSON(&json)
		if value, ok := json["X-type"].(string); ok {
			actType = value
		} else {
			Fail(ctx, nil, "type should be string") //返回错误反馈
			return
		}
		if value, ok := json["name"].(string); ok {
			name = value
		} else {
			Fail(ctx, nil, "type not correct") //返回错误反馈
			return

		}
		if value, ok := json["namespace"].(string); ok {
			namespace = value
		} else {
			Fail(ctx, nil, "type should be string") //返回错误反馈
			return
		}
		if value, ok := json["volume"].(string); ok || strings.ToLower(actType) == "object" {
			sVolume = value
		} else {
			Fail(ctx, nil, "type should be string") //返回错误反馈
			return
		}

		if value, ok := json["autoscale"].(string); ok || strings.ToLower(actType) == "object" {
			autoScales = value
		} else {
			Fail(ctx, nil, "type should be string") //返回错误反馈
			return
		}
	}
	//如果有autoscale 转化为bool
	if autoScales != "" && strings.ToLower(autoScales) == "true" {
		autoScale = true //如果为true则改为真
	} else {
		autoScale = false //其余情况---为空和为false均改为false
	}
	//如果有volume参数要进行格式转化
	if sVolume != "" {
		volumeCp, err := strconv.Atoi(sVolume)
		if err != nil {
			FailUnac(ctx, nil, "volume not number")
			return
		}
		volume = volumeCp
	}
	//名字要符合标准
	match, _ := regexp.MatchString("[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*", name)
	if !match { //如果名称不符合标准则去除
		FailUnac(ctx, nil, "format err")
		return
	}
	//// 开发完后删去,仅用于测试:TODO
	//Success(ctx, nil, "test done")
	//return
	////记得删除

	//判断调用方法
	switch strings.ToLower(actType) {
	case "block":
		if err := cubeControl.PatchBlockPVC(name, namespace, volume, autoScale); err != nil {
			FailCrea(ctx, nil, err.Error())
			return
		}
		Success(ctx, nil, "patch done")
	case "file":
		if err := cubeControl.PatchFileSystemPVC(name, namespace, volume, autoScale); err != nil {
			FailCrea(ctx, nil, err.Error())
			return
		}
		Success(ctx, nil, "patch done")
	case "object":
		sMaxObject := ctx.PostForm("maxobject")
		sMaxGBsize := ctx.PostForm("maxgbsize")
		json := make(map[string]interface{})
		ctx.BindJSON(&json)
		if sMaxObject == "" && sMaxGBsize == "" {
			if value, ok := json["maxobject"].(string); ok {
				sMaxObject = value
			} else {
				Fail(ctx, nil, "type should be string") //返回错误反馈
				return
			}
			if value, ok := json["maxgbsize"].(string); ok {
				sMaxGBsize = value
			} else {
				Fail(ctx, nil, "type should be string") //返回错误反馈
				return
			}
		}

		maxGBsize, err := strconv.Atoi(sMaxGBsize)
		if err != nil {
			FailUnac(ctx, nil, "maxgbsize not number") //传入的不是数字格式问题
			return
		}
		maxobject, err := strconv.Atoi(sMaxObject)
		if err != nil {
			FailUnac(ctx, nil, "maxobject not number")
			return
		}
		if err := cubeControl.PatchObjectBucket(name, namespace, maxobject, maxGBsize); err != nil {
			FailCrea(ctx, nil, err.Error())
			return
		}
		Success(ctx, nil, "patch done")
	}
}
