package impleOPR

import (
	kit "main/cubeOperatorKit"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func PushGetObj(ctx *gin.Context) {
	//根据post请求的body体来解析参数，仅支持json两种格式
	var namespace, bucketClaimName, key, actType, blockStr string
	var blockNum int = 1 //记录分块数量，默认为一

	var value []byte
	jsons := make(gin.H)
	ctx.BindJSON(&jsons)
	if valueStr, ok := jsons["namespace"].(string); ok {
		namespace = valueStr
	} else {
		Fail(ctx, nil, "namespace should be string") //返回错误反馈
		return
	}
	if valueStr, ok := jsons["name"].(string); ok {
		bucketClaimName = valueStr
	} else {
		Fail(ctx, nil, "name should be string") //返回错误反馈
		return
	}
	if valueStr, ok := jsons["key"].(string); ok {
		key = valueStr
	} else {
		Fail(ctx, nil, "key should be string") //返回错误反馈
		return
	}
	if valueStr, ok := jsons["X-action"].(string); ok {
		actType = valueStr
	} else {
		Fail(ctx, nil, "X-action should be string") //返回错误反馈
		return
	}
	if valueStr, ok := jsons["block"].(string); ok { //加入分块的机制，运行用户选择数据的分块运输
		blockStr = valueStr
	}
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
	//对于value数据，判断其为string还是[]byte
	if valueStr, ok := jsons["value"].(string); ok {
		value = []byte(valueStr)
	} else {
		valueByte, err := jsons["value"].([]byte)
		if err {
			value = valueByte
		}
		if !err && actType == "push" {
			Fail(ctx, nil, "value should be string or []byte") //返回错误反馈
			return
		}
	}

	switch strings.ToLower(actType) {
	case "push":
		err := kit.PutObject(namespace, bucketClaimName, key, value)
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

		//返回get得到到对象信息，这里附带其key namespace等，并进行分块
		for index, valueBytes := range splitArray([]byte(value), blockNum) {
			Success(ctx, gin.H{"value" + strconv.Itoa(index): valueBytes, "key": key, "namespace": namespace, "name": bucketClaimName}, "obj value")
		}
	}
}
