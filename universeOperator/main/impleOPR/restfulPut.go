package impleOPR

import (
	kit "main/cubeOperatorKit"

	"github.com/gin-gonic/gin"
)

func OssPut(ctx *gin.Context) {
	var namespace, bucketClaimName, key string
	var value []byte
	jsons := make(gin.H)
	ctx.BindJSON(&jsons)
	if valueStr, ok := jsons["namespace"].(string); ok { //读取namspace
		namespace = valueStr
	} else {
		Fail(ctx, nil, "namespace should be string") //返回错误反馈
		return
	}
	if valueStr, ok := jsons["name"].(string); ok { //读取bucketname
		bucketClaimName = valueStr
	} else {
		Fail(ctx, nil, "name should be string") //返回错误反馈
		return
	}
	if valueStr, ok := jsons["key"].(string); ok { //读取key
		key = valueStr
	} else {
		Fail(ctx, nil, "key should be string") //返回错误反馈
		return
	}
	//读取传入的value
	//对于value数据，判断其为string还是[]byte
	if valueStr, ok := jsons["value"].(string); ok {
		value = []byte(valueStr)
	} else {
		valueByte, err := jsons["value"].([]byte)
		if err {
			value = valueByte
		} else {
			Fail(ctx, nil, "value should be string or []byte") //返回错误反馈
			return
		}
	}
	//实现put
	err := kit.PutObject(namespace, bucketClaimName, key, &value)
	if err != nil {
		FailUnac(ctx, nil, "Fail Put OBJ: "+err.Error())
	}
	Success(ctx, nil, "Put success")
}
