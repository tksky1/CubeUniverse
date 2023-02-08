package universalFuncs

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

// GetClientSet 加载kubeconfig信息，然后获得ClientSet
func GetClientSet() *kubernetes.Clientset {
	log.Println("加载kubeconfig鉴权..")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Println("Warning：未检测到pod运行时，程序未运行在pod中！")
		log.Println("Warning：尝试读取环境变量KUBECONFIG..")
		config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
		if err != nil {
			config, err = clientcmd.BuildConfigFromFlags("", "/etc/kubernetes/admin.conf")
			if err != nil {
				config, err = clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
				if err != nil {
					log.Fatalln("错误：未能从环境变量读取kubeconfig！请确保使用云原生部署方式，而非直接运行程序！")
				}
			}
		}
	}

	log.Println("读取kubeconfig成功")
	var client *kubernetes.Clientset
	client, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("错误：未能加载ClientSet！请确保使用云原生部署方式，而非直接运行程序！")
	}
	return client
}
