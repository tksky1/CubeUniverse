package storageact

import (
	"control-backend/cubeControl"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
)

type info struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Volume    string `json:"volume"`
	CreaTime  string `json:"time"`
}

func ListPvc(ctx *gin.Context) {
	var actType string
	var resMap gin.H = make(gin.H)
	actType = ctx.GetHeader("X-type")

	//仅用于测试，最后删去再替换:TODO
	for index := 0; index < 10; index++ {
		var responseStru info
		responseStru.Name = actType + "testname-" + string(index+int('1'))
		responseStru.Namespace = actType + "testnamespace-" + string(index+int('1'))
		responseStru.Volume = string(index + int('1'))
		responseStru.CreaTime = fmt.Sprint("2006-01-02 15:04:0", index)
		key := fmt.Sprint(actType, index)
		resMap[key] = responseStru
	}
	Success(ctx, resMap, "all info") //测试返回
	return
	//记得删除

	switch strings.ToLower(actType) {
	case "block":
		pvcList, err := cubeControl.ListBlockSystemPVC()
		if err != nil { //取出信息报错
			FailCrea(ctx, nil, err.Error())
			return
		}
		//遍历切片
		for index, pvcTar := range pvcList {
			var responseStru info
			responseStru.Name = pvcTar.Name
			responseStru.Namespace = pvcTar.Namespace
			responseStru.Volume = string(pvcTar.Spec.Resources.Requests[corev1.ResourceStorage].Format)
			responseStru.CreaTime = pvcTar.CreationTimestamp.Format("2006-01-02 15:04:05") //只展示年月日
			key := fmt.Sprint(actType, index)                                              //产生Map的独特key
			resMap[key] = responseStru
		}
		//将数据输入到响应，具体来说是响应的data字段里面，会依次产生info对象，
		Success(ctx, resMap, "all info")
	case "file":
		pvcList, err := cubeControl.ListFileSystemPVC()
		if err != nil { //取出信息报错
			FailCrea(ctx, nil, err.Error())
			return
		}
		//遍历切片
		for index, pvcTar := range pvcList {
			var responseStru info
			responseStru.Name = pvcTar.Name
			responseStru.Namespace = pvcTar.Namespace
			responseStru.Volume = string(pvcTar.Spec.Resources.Requests[corev1.ResourceStorage].Format)
			responseStru.CreaTime = pvcTar.CreationTimestamp.Format("2006-01-02 15:04:05") //只展示年月日
			key := fmt.Sprint(actType, index)                                              //产生Map的独特key
			resMap[key] = responseStru
		}
		//将数据输入到响应，具体来说是响应的data字段里面，会依次产生info对象，
		Success(ctx, resMap, "all info")
	case "object":
		pvcList, err := cubeControl.ListObjectBucketClaim()
		if err != nil { //取出信息报错
			FailCrea(ctx, nil, err.Error())
			return
		}
		//遍历切片
		for index, pvcTar := range pvcList {
			key := fmt.Sprint(actType, index) //产生Map的独特key
			resMap[key] = pvcTar
		}
		Success(ctx, resMap, "all info")
	}
}
