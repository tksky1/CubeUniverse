package universalFuncs

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"strings"
)

// CheckCubeUniverseComponent 检测CubeUniverse组件的情况
func CheckCubeUniverseComponent(clientSet *kubernetes.Clientset) (bool, bool, bool, bool) {
	pods, _ := clientSet.CoreV1().Pods("cubeuniverse").List(context.TODO(), metav1.ListOptions{})
	var operator, dashboard, controlBackend, universeBuilder bool
	for _, pod := range pods.Items {
		if strings.Index(pod.Name, "operator") != -1 && pod.Status.Phase == "Running" {
			operator = true
		} else if strings.Index(pod.Name, "dashboard") != -1 && pod.Status.Phase == "Running" {
			dashboard = true
		} else if strings.Index(pod.Name, "control-backend") != -1 && pod.Status.Phase == "Running" {
			controlBackend = true
		} else if strings.Index(pod.Name, "builder") != -1 {
			universeBuilder = true
		}
	}
	if !(operator && dashboard && controlBackend) {
		log.Printf("CubeUniverse自检：\noperator：%t, dashboard: %t, controlBackend: %t, universeBuilder: %t", operator, dashboard, controlBackend, universeBuilder)
	}
	return operator, dashboard, controlBackend, universeBuilder
}

// CheckMLStatus 检查机器学习组件状态
func CheckMLStatus(clientSet *kubernetes.Clientset) (kafka bool, ml bool) {
	pods, _ := clientSet.CoreV1().Pods("cubeuniverse").List(context.TODO(), metav1.ListOptions{})
	for _, pod := range pods.Items {
		if strings.Index(pod.Name, "kafka") != -1 && pod.Status.Phase == "Running" {
			kafka = true
		} else if strings.Index(pod.Name, "ml") != -1 && pod.Status.Phase == "Running" {
			ml = true
		}
	}
	return kafka, ml
}
