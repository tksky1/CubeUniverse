package impleOPR

import (
	"github.com/gin-gonic/gin"
	kit "main/cubeOperatorKit"
	"strings"
)

func PushGetObj(ctx *gin.Context) {
	//根据post请求的body体来解析参数，仅支持json两种格式
	var namespace, bucketClaimName, key, actType string
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
	//对于value数据，判断其为string还是[]byte
	if valueStr, ok := jsons["value"].(string); ok {
		value = []byte(valueStr)
	} else {
		if valueByte, err := jsons["value"].([]byte); err {
			value = valueByte
		} else {
			Fail(ctx, nil, "value should be string or []byte") //返回错误反馈
			return
		}
	}

	switch strings.ToLower(actType) {
	case "push":
		if err := kit.PutObject(namespace, bucketClaimName, key, value); err != nil {
			FailUnac(ctx, nil, "Fail Put OBJ: "+err.Error())
		}
		Success(ctx, nil, "Put success")

	case "get":
		value, err := kit.GetObject(namespace, bucketClaimName, key)
		if err != nil {
			FailUnac(ctx, nil, err.Error())
			return
		}
		Success(ctx, gin.H{"value": value}, "obj value")
	}

}
