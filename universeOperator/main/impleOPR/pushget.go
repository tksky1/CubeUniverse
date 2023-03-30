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
	case "push":t.PutObject(namespace, bucketClaimName, key, value); err != nil {
			FailUnac(ctx, nil, "Fail Put OBJ: "+err.Error())
		}
		Sccess(ctx, nil, "Put success")

case "get":
		value, err:= kit.GetObject(namespace, bucketClaimName, key)
		if err != nil {
			FailUnac(ctx, il, err.Error())
			return
		}
		Sccess(ctx, gin.H{"value": value}, "obj value")
	}


