package universalFuncs

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"strings"
)

// CheckClusterHealth 检查集群的健康状况
func CheckClusterHealth(clientSet *kubernetes.Clientset) bool {
	pods, err := clientSet.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{})
	var coredns, etcd, kubeProxy, kubeApiServer, scheduler bool
	for _, pod := range pods.Items {
		if strings.Index(pod.Name, "coredns") != -1 && pod.Status.Phase == "Running" {
			coredns = true
		} else if strings.Index(pod.Name, "etcd") != -1 && pod.Status.Phase == "Running" {
			etcd = true
		} else if strings.Index(pod.Name, "kube-proxy") != -1 && pod.Status.Phase == "Running" {
			kubeProxy = true
		} else if strings.Index(pod.Name, "kube-apiserver") != -1 && pod.Status.Phase == "Running" {
			kubeApiServer = true
		} else if strings.Index(pod.Name, "kube-scheduler") != -1 && pod.Status.Phase == "Running" {
			scheduler = true
		}
	}
	pods2, err := clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	log.Printf("检测到集群共有 %d 个pod\n", len(pods2.Items))
	if coredns && etcd && kubeProxy && kubeApiServer && scheduler {
		return true
	} else {
		if len(pods.Items) == 0 {
			log.Fatal("错误：CubeUniverse没有访问集群容器的权限，请重新部署集群！")
		}
		log.Printf("集群自检：k8s集群必要组件未正常运行！\ncoredns：%t, etcd: %t, kube-proxy: %t, kube-apiserver: %t, scheduler: %t", coredns, etcd, kubeProxy, kubeApiServer, scheduler)
	}
	return false
}
