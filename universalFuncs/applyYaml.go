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

// ApplyYaml 将Yaml解析并创建对应资源
func ApplyYaml(filename string, namespace string) error {
	log.Println("部署：", filename)

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = ApplyBytes(bytes, namespace)
	return err
}

// ApplyCrdFromYaml 用于创建带CRD的yaml资源
func ApplyCrdFromYaml(yamlFilePath string, nameSpace string, clientSet *kubernetes.Clientset, dd *dynamic.DynamicClient) error {

	log.Println("准备ceph组件", yamlFilePath, "..")

	fileBytes, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return err
	}
	err = ApplyCrdFromBytes(fileBytes, nameSpace, clientSet, dd)
	return err
}
