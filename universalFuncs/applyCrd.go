package universalFuncs

import (
	"bytes"
	"context"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/util/flowcontrol"
	"log"
)

func ApplyCrdFromYaml(yamlFilePath string, nameSpace string, clientSet *kubernetes.Clientset, dd *dynamic.DynamicClient) {

	log.Println("准备ceph组件", yamlFilePath, "..")

	filebytes, err := ioutil.ReadFile(yamlFilePath)
	config, _ := rest.InClusterConfig()
	config.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(1000, 1000)
	config.QPS = 1000
	config.Burst = 1000
	if err != nil {
		log.Printf("%v\n", err)
	}

	gr, err := restmapper.GetAPIGroupResources(clientSet.Discovery())
	if err != nil {
		log.Fatal(err)
	}
	mapper := restmapper.NewDiscoveryRESTMapper(gr)
	var dri dynamic.ResourceInterface
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(filebytes), 1000)
	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			break
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			log.Fatal(err)
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			log.Fatal(err)
		}

		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace(nameSpace)
			}
			dri = dd.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = dd.Resource(mapping.Resource)
		}

		obj2, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("CRD: %s/%s created\n", obj2.GetKind(), obj2.GetName())
		}
	}
}
