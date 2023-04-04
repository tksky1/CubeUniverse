package storageact

import (
	"control-backend/cubeControl"
	"strings"

	"github.com/gin-gonic/gin"
)

func DeletePvc(ctx *gin.Context) {
	var name, namespace, actType string
	json := make(map[string]interface{})
	ctx.BindJSON(&json)
	if value, ok := json["name"].(string); ok {
		name = value
	} else {
		Fail(ctx, nil, "type not correct") //返回错误反馈
		return
	}
	if value, ok := json["namespace"].(string); ok {
		namespace = value
	} else {
		Fail(ctx, nil, "type should be string") //返回错误反馈
		return
	}
	if value, ok := json["X-type"].(string); ok {
		actType = value
	} else {
		Fail(ctx, nil, "type should be string") //返回错误反馈
		return
	}

	// //测试用:TODO
	// Success(ctx, nil, "delete done")
	// return
	// //记得删除

	switch strings.ToLower(actType) {
	case "block", "file":
		if err := cubeControl.DeletePVC(name, namespace); err != nil {
			FailCrea(ctx, nil, err.Error())
			return
		}
		Success(ctx, nil, "delete done")
	case "object":
		if err := cubeControl.DeleteObjectBucket(name, namespace); err != nil {
			FailCrea(ctx, nil, err.Error())
			return
		}
		Success(ctx, nil, "delete done")
	}
}
