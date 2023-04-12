package impleOPR

import (
	kit "main/cubeOperatorKit"

	"github.com/gin-gonic/gin"
)

func OssList(ctx *gin.Context) {
	var namespace, bucketClaimName string
	var tag string = "" //给tag一个默认值，更安全
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
	//实现list
	if tag != "" { //说明传入了tag
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
