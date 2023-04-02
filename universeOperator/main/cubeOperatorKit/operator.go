package cubeOperatorKit

import (
	"CubeUniverse/universalFuncs"
	"k8s.io/client-go/kubernetes"
)

// UniverseOperator 常驻运行 监测集群状态和自动维护

var ClientSet *kubernetes.Clientset

func Init() {

	ClientSet = universalFuncs.GetClientSet()
	sessionCacheMap = make(map[[16]byte]*SessionAndBucketName)

}
