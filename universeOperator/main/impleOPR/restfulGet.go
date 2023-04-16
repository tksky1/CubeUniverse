package impleOPR

import (
	kit "main/cubeOperatorKit"
	"strconv"

	"github.com/gin-gonic/gin"
)

func OssGet(ctx *gin.Context) {
	var namespace, bucketClaimName, key, blockStr, indexBlock string
	var blockNum, indexNum int = 1, 0 //记录分块数量与用户所需的分块号，默认1、0

	namespace = ctx.Query("namespace")
	bucketClaimName = ctx.Query("name")
	key = ctx.Query("key")
	blockStr = ctx.Query("block")
	indexBlock = ctx.Query("index")

	// jsons := make(gin.H)
	// ctx.BindJSON(&jsons)
	// if valueStr, ok := jsons["namespace"].(string); ok {
	// 	namespace = valueStr
	// } else {
	// 	Fail(ctx, nil, "namespace should be string") //返回错误反馈
	// 	return
	// }
	// if valueStr, ok := jsons["name"].(string); ok {
	// 	bucketClaimName = valueStr
	// } else {
	// 	Fail(ctx, nil, "name should be string") //返回错误反馈
	// 	return
	// }
	// if valueStr, ok := jsons["key"].(string); ok {
	// 	key = valueStr
	// } else {
	// 	Fail(ctx, nil, "key should be string") //返回错误反馈
	// 	return
	// }
	// if valueStr, ok := jsons["block"].(string); ok { //加入分块的机制的块数，运行用户选择数据的分块运输块数
	// 	blockStr = valueStr
	// }
	// if valueStr, ok := jsons["index"].(string); ok { //加入分块的机制的索引，运行用户选择数据的分块运输索引值
	// 	indexBlock = valueStr
	// }

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

	//响应返回get值
	value, err := kit.GetObject(namespace, bucketClaimName, key)
	if err != nil {
		FailUnac(ctx, nil, err.Error())
		return
	}
	//返回get得到到对象信息，这里附带其key namespace等，并进行分块
	//valueBytes := splitArray(value, blockNum, indexNu)
	//将bytes类型数据转为string避免base64转换)
	max := int(len(*value))
	quantity := max / blockNum
	var value2Str string 
	if (indexNum == blockNum-1){ 
		value2Sr = string((*value)[indexNum*quantity:])
	} else {
		alue2Str = string((*value)[indexNum*quantity : quantity*(indexNum+1)])
	}
	uccess(ctx, gin.H{"value" + strconv.Itoa(indexNum): value2Str, "key": key, "namespace": namespace, "name": bucketClaimName}, "obj value")
}
