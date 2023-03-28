package universalFuncs

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"strings"
)

// CheckCephComponent 检测Ceph组件的情况
func CheckCephComponent(clientSet *kubernetes.Clientset) (operator bool, rbdplugin bool, mon bool, mgr bool, osd bool) {
	pods, _ := clientSet.CoreV1().Pods("rook-ceph").List(context.TODO(), metav1.ListOptions{})

	for _, pod := range pods.Items {
		if strings.Index(pod.Name, "operator") != -1 && pod.Status.Phase == "Running" {
			operator = true
		} else if strings.Index(pod.Name, "rbdplugin") != -1 && strings.Index(pod.Name, "pre") == -1 && pod.Status.Phase == "Running" {
			rbdplugin = true
		} else if strings.Index(pod.Name, "mon") != -1 && strings.Index(pod.Name, "pre") == -1 && pod.Status.Phase == "Running" {
			mon = true
		} else if strings.Index(pod.Name, "mgr") != -1 && strings.Index(pod.Name, "pre") == -1 && pod.Status.Phase == "Running" {
			mgr = true
		} else if strings.Index(pod.Name, "osd") != -1 && strings.Index(pod.Name, "pre") == -1 && pod.Status.Phase == "Running" {
			osd = true
		}
	}
	if !(operator && rbdplugin && mon && mgr && osd) {
		log.Printf("Ceph自检：\noperator：%t, rbdplugin: %t, mon: %t, mgr: %t, osd: %t", operator, rbdplugin, mon, mgr, osd)
	}
	return operator, rbdplugin, mon, mgr, osd
}

func CheckMysqlStat(clientSet *kubernetes.Clientset) bool {
	pods, _ := clientSet.CoreV1().Pods("rook-ceph").List(context.TODO(), metav1.ListOptions{})
	for _, pod := range pods.Items {
		if strings.Index(pod.Name, "mysql") != -1 && pod.Status.Phase == "Running" {
			return true
		}
	}
	return false
}
