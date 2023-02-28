package universalFuncs

import (
	"io/ioutil"
	syaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	sigyaml "sigs.k8s.io/yaml"
)

import (
	"bytes"
	"context"
	"fmt"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
)

type ExecuteYaml struct {
	applyYaml string
	namespace string
}

var clientset *kubernetes.Clientset

func NewYaml(applyYaml, namespace string) *ExecuteYaml {
	return &ExecuteYaml{
		applyYaml: applyYaml,
		namespace: namespace,
	}
}

func (y *ExecuteYaml) GtGVR(gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {

	gr, err := restmapper.GetAPIGroupResources(clientset.Discovery())
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	mapper := restmapper.NewDiscoveryRESTMapper(gr)

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	return mapping.Resource, nil
}

func (y *ExecuteYaml) UpdateFromYaml() error {

	dynameicclient := GetDynamicClient()
	d := yaml.NewYAMLOrJSONDecoder(bytes.NewBufferString(y.applyYaml), 4096)

	for {
		var rawObj runtime.RawExtension
		err := d.Decode(&rawObj)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("decode is err %v", err)
		}

		obj, _, err := syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return fmt.Errorf("rawobj is err%v", err)
		}

		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return fmt.Errorf("tounstructured is err %v", err)
		}

		unstructureObj := &unstructured.Unstructured{Object: unstructuredMap}
		gvr, err := y.GtGVR(unstructureObj.GroupVersionKind())
		if err != nil {
			return err
		}
		unstructuredYaml, err := sigyaml.Marshal(unstructureObj)
		if err != nil {
			return fmt.Errorf("unable to marshal resource as yaml: %w", err)
		}
		_, getErr := dynameicclient.Resource(gvr).Namespace(y.namespace).Get(context.Background(), unstructureObj.GetName(), metav1.GetOptions{})
		if getErr != nil {
			_, createErr := dynameicclient.Resource(gvr).Namespace(y.namespace).Create(context.Background(), unstructureObj, metav1.CreateOptions{})
			if createErr != nil {
				println("create err!")
				return createErr
			}
		}

		force := true
		if y.namespace == unstructureObj.GetNamespace() {

			_, err = dynameicclient.Resource(gvr).
				Namespace(y.namespace).
				Patch(context.Background(),
					unstructureObj.GetName(),
					types.ApplyPatchType,
					unstructuredYaml, metav1.PatchOptions{
						FieldManager: unstructureObj.GetName(),
						Force:        &force,
					})

			if err != nil {
				return fmt.Errorf("unable to patch resource: %w", err)
			}

		} else {

			_, err = dynameicclient.Resource(gvr).
				Patch(context.Background(),
					unstructureObj.GetName(),
					types.ApplyPatchType,
					unstructuredYaml, metav1.PatchOptions{
						Force:        &force,
						FieldManager: unstructureObj.GetName(),
					})
			if err != nil {
				return fmt.Errorf("ns is nil unable to patch resource: %w", err)
			}

		}

	}
	return nil

}

// ReadFile 从文件读取字符串
func ReadFile(fileName string) string {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		return ""
	}
	return string(f)
}

func ApplyYamlv2(filepath string, clientSet *kubernetes.Clientset, namespace string) error {
	//fmt.Println(s)
	clientset = clientSet
	ey := NewYaml(ReadFile(filepath), namespace)
	return ey.UpdateFromYaml()
}
