package impleOPR

import (
	"io"
	kit "main/cubeOperatorKit"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func PutGetDeleteListObj(ctx *gin.Context) {
	//根据post请求的body体来解析参数，仅支持json两种格式
	var namespace, bucketClaimName, key, actType, blockStr, indexBlock, tag string
	tag = ""                      //给tag一个默认值，更安全
	var blockNum, indexNum = 1, 0 //记录分块数量，默认为一

	//var value []byte

	actType = ctx.Query("X-action")
	namespace = ctx.Query("namespace")
	bucketClaimName = ctx.Query("name")
	key = ctx.Query("key")
	blockStr = ctx.Query("block")
	indexBlock = ctx.Query("index")
	tag = ctx.Query("tag")

	//if valueStr, ok := jsons["X-action"].(string); ok { //检测用户使用什么方法
	//	actType = valueStr
	//} else {
	//	Fail(ctx, nil, "X-action should be string and not nil") //返回错误反馈
	//	return
	//}
	//if valueStr, ok := jsons["namespace"].(string); ok {
	//	namespace = valueStr
	//} else {
	//	Fail(ctx, nil, "namespace should be string") //返回错误反馈
	//	return
	//}
	//if valueStr, ok := jsons["name"].(string); ok {
	//	bucketClaimName = valueStr
	//} else {
	//	Fail(ctx, nil, "name should be string") //返回错误反馈
	//	return
	//}
	//if valueStr, ok := jsons["key"].(string); ok || strings.ToLower(actType) == "list" { //用户调用list方法的时候key可以为空
	//	key = valueStr
	//} else {
	//	Fail(ctx, nil, "key should be string") //返回错误反馈
	//	return
	//}
	//
	//if valueStr, ok := jsons["block"].(string); ok { //加入分块的机制的块数，运行用户选择数据的分块运输块数
	//	blockStr = valueStr
	//}
	//if valueStr, ok := jsons["index"].(string); ok { //加入分块的机制的索引，运行用户选择数据的分块运输索引值
	//	indexBlock = valueStr
	//}
	//if valueStr, ok := jsons["tag"].(string); ok { //得到tag标签的值，此为可选参数
	//	tag = valueStr
	//}
	//对于分块数，如果没写的话默认为1
	if blockStr == "" {
		blockNum = 1
	} else {
		var err error = nil
		blockNum, err = strconv.Atoi(blockStr)
		if err != nil {
			Fail(ctx, nil, "block should be string represent a number") //返回错误反馈，block应该代表整数
			return
		}
	}
	//对于索引值，如果没写的话默认为0
	if indexBlock == "" {
		indexNum = 0
	} else {
		var err error = nil
		indexNum, err = strconv.Atoi(indexBlock)
		if err != nil {
			Fail(ctx, nil, "index should be string represent a number") //返回错误反馈，block应该代表整数
			return
		}
	}
	//保证索引比分块小
	if indexNum >= blockNum {
		Fail(ctx, nil, "index out of range") //返回错误反馈，block应该代表整数
		return
	}
	////对于value数据，判断其为string还是[]byte
	//if valueStr, ok := jsons["value"].(string); ok {
	//	value = []byte(valueStr)
	//} else {
	//	valueByte, err := jsons["value"].([]byte)
	//	if err {
	//		value = valueByte
	//	}
	//	if !err && actType == "put" {
	//		Fail(ctx, nil, "value should be string or []byte") //返回错误反馈
	//		return
	//	}
	//}

	switch strings.ToLower(actType) {
	case "put":
		reader := ctx.Request.Body.(io.Reader)
		err := kit.PutObject(namespace, bucketClaimName, key, ctx.Request.ContentLength, &reader)
		if err != nil {
			FailUnac(ctx, nil, "Fail Put OBJ: "+err.Error())
		}
		Success(ctx, nil, "Put success")

	case "get":

		value, err := kit.GetObject(namespace, bucketClaimName, key)
		if err != nil {
			FailUnac(ctx, nil, err.Error())
			return
		}

		max := len(*value)
		quantity := max / blockNum
		var value2Str string
		if indexNum == blockNum-1 {
			value2Str = string((*value)[indexNum*quantity:])
		} else {
			value2Str = string((*value)[indexNum*quantity : quantity*(indexNum+1)])
		}

		Success(ctx, gin.H{"value" + strconv.Itoa(indexNum): value2Str, "key": key, "namespace": namespace, "name": bucketClaimName}, "obj value")
	case "delete":
		if err := kit.DeleteObject(namespace, bucketClaimName, key); err != nil {
			FailCrea(ctx, nil, "delete err: "+err.Error())
			return
		} else {
			Success(ctx, nil, "delete success:namespace="+namespace+" name="+bucketClaimName+" key="+key)
			return
		}
	case "list":
		if tag != "" {
			if valueArr, err := kit.ListObjectByTag(namespace, bucketClaimName, tag); err != nil {
				FailCrea(ctx, nil, "list err: "+err.Error())
				return
			} else {
				Success(ctx, gin.H{"keys": valueArr}, "list key success: namespace="+namespace+" name="+bucketClaimName)
				return
			}
		}
		if valueArr, err := kit.ListObjectFromBucket(namespace, bucketClaimName); err != nil {
			FailCrea(ctx, nil, "list err: "+err.Error())
			return
		} else {
			Success(ctx, gin.H{"keys": valueArr}, "list key success: namespace="+namespace+" name="+bucketClaimName)
			return
		}
	}
}
