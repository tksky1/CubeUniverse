package universalFuncs

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

// PatchYaml 将Yaml解析并创建对应资源 使用PATCH方法
func PatchYaml(filename string, namespace string) error {
	log.Println("部署：", filename)

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = PatchBytes(bytes, namespace)
	return err
}

// PatchCrdFromYaml 用于创建带CRD的yaml资源 使用PATCH方法
func PatchCrdFromYaml(yamlFilePath string, nameSpace string, clientSet *kubernetes.Clientset, dd *dynamic.DynamicClient) error {

	log.Println("准备ceph组件", yamlFilePath, "..")

	fileBytes, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return err
	}
	err = PatchCrdFromBytes(fileBytes, nameSpace, clientSet, dd)
	return err
}
