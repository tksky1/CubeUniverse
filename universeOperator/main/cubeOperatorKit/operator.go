package cubeOperatorKit

import (
	"CubeUniverse/universalFuncs"
	"github.com/Shopify/sarama"
	"k8s.io/client-go/kubernetes"
)

// UniverseOperator 常驻运行 监测集群状态和自动维护

var ClientSet *kubernetes.Clientset
var Producer *sarama.SyncProducer

func Init() {

	ClientSet = universalFuncs.GetClientSet()
	sessionCacheMap = make(map[[16]byte]*SessionAndBucketName)

}
