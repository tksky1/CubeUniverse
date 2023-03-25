package storageact

import (
	"regexp"

	"github.com/gin-gonic/gin"
)

// 处理的请求必须带有name、namespace、volume字段,均为string
func Pvcreq(ctx *gin.Context) {
	//根据post请求body体参数来解析数据,支持postform和json两种格式
	name := ctx.PostForm("name")
	namespace := ctx.PostForm("namespace")
	volume := ctx.PostForm("volume")

	if name == "" && namespace == "" && volume == "" {
		json := make(map[string]interface{})
		ctx.BindJSON(&json)
		name = json["name"].(string)
		namespace = json["namespace"].(string)
		volume = json["volume"].(string)
	}
	match, _ := regexp.MatchString("[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*", name)
	if !match {

	}
}
