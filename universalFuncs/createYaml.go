package universalFuncs

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
)

// CreateYaml 将Yaml解析并创建对应资源 使用CREATE方法
func CreateYaml(filename string, namespace string) error {
	log.Println("部署：", filename)

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = CreateBytes(bytes, namespace)
	return err
}

// CreateCrdFromYaml 用于创建带CRD的yaml资源 使用CREATE方法
func CreateCrdFromYaml(yamlFilePath string, nameSpace string, clientSet *kubernetes.Clientset, dd *dynamic.DynamicClient) error {

	log.Println("准备ceph组件", yamlFilePath, "..")

	fileBytes, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return err
	}
	err = CreateCrdFromBytes(fileBytes, nameSpace, clientSet, dd)
	return err
}
