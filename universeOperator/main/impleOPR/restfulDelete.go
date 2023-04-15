package impleOPR

import (
	kit "main/cubeOperatorKit"

	"github.com/gin-gonic/gin"
)

func OssDelete(ctx *gin.Context) {
	var namespace, bucketClaimName, key string
	namespace = ctx.Query("namespace")
	bucketClaimName = ctx.Query("name")
	key = ctx.Query("key")
	//实现delete
	if err := kit.DeleteObject(namespace, bucketClaimName, key); err != nil {
		FailCrea(ctx, nil, "delete err: "+err.Error())
		return
	} else {
		Success(ctx, nil, "delete success:namespace="+namespace+" name="+bucketClaimName+" key="+key)
		return
	}
}
