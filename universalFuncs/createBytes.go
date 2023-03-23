package universalFuncs

import (
	"bytes"
	"context"
	"errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/util/flowcontrol"
	"log"
)

// CreateBytes 将bytes格式的Yaml解析并创建对应资源
// 也可以更新已存在资源
func CreateBytes(yamlBytes []byte, namespace string) error {

	ctx := context.TODO()
	cfg, _ := rest.InClusterConfig()
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return err
	}

	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode(yamlBytes, nil, obj)
	if err != nil {
		return err
	}

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		dr = dyn.Resource(mapping.Resource)
	}

	_, err = dr.Create(ctx, obj, metav1.CreateOptions{
		FieldManager: namespace,
	})

	return err
}

// CreateCrdFromBytes 从Bytes格式的Yaml创建CRD资源
// 也可以用于更新已存在资源
func CreateCrdFromBytes(yamlBytes []byte, nameSpace string, clientSet *kubernetes.Clientset, dd *dynamic.DynamicClient) error {

	config, _ := rest.InClusterConfig()
	config.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(1000, 1000)
	config.QPS = 1000
	config.Burst = 1000

	gr, err := restmapper.GetAPIGroupResources(clientSet.Discovery())
	if err != nil {
		return err
	}
	mapper := restmapper.NewDiscoveryRESTMapper(gr)
	var dri dynamic.ResourceInterface
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(yamlBytes), 1000)
	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			break
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return err
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}

		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace(nameSpace)
			}
			dri = dd.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = dd.Resource(mapping.Resource)
		}

		if !checkNamespaceExist(unstructuredObj.GetNamespace(), clientSet) && unstructuredObj.GetNamespace() != "rook-ceph" && unstructuredObj.GetNamespace() != "" {
			return errors.New("命名空间" + unstructuredObj.GetNamespace() + "不存在！")
		}

		obj2, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
		if err != nil {
			return err
		} else {
			log.Printf("CRD: %s/%s 已创建\n", obj2.GetKind(), obj2.GetName())
		}
	}
	return nil
}
