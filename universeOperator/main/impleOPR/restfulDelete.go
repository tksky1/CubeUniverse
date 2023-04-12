package impleOPR

import (
	kit "main/cubeOperatorKit"

	"github.com/gin-gonic/gin"
)

func OssDelete(ctx *gin.Context) {
	var namespace, bucketClaimName, key string
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
	//实现delete
	if err := kit.DeleteObject(namespace, bucketClaimName, key); err != nil {
		FailCrea(ctx, nil, "delete err: "+err.Error())
		return
	} else {
		Success(ctx, nil, "delete success:namespace="+namespace+" name="+bucketClaimName+" key="+key)
		return
	}
}
