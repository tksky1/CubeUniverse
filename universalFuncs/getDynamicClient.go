package universalFuncs

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"log"
)

// GetDynamicClient 加载kubeconfig信息，然后获得DynamicClient
func GetDynamicClient() *dynamic.DynamicClient {
	log.Println("加载kubeconfig鉴权..")
	config, _ := rest.InClusterConfig()
	client, err := dynamic.NewForConfig(config)

	if err != nil {
		log.Panic(err)
	}

	return client
}
