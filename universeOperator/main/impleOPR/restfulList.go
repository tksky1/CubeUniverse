package impleOPR

import (
	kit "main/cubeOperatorKit"

	"github.com/gin-gonic/gin"
)

func OssList(ctx *gin.Context) {
	var namespace, bucketClaimName string
	var tag string = "" //给tag一个默认值，更安全
	namespace = ctx.Query("namespace")
	bucketClaimName = ctx.Query("name")
	tag = ctx.Query("tag")
	//实现list
	if tag != "" { //说明传入了tag
		if valueArr, err := kit.ListObjectByTag(namespace, bucketClaimName, tag); err != nil {
			FailCrea(ctx, nil, "list err: "+err.Error())
			return
		} else {
			Success(ctx, gin.H{"keys": valueArr}, "list key success: namespace="+namespace+" name="+bucketClaimName+"tag= "+tag)
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
