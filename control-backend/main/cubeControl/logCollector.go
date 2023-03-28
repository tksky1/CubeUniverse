package cubeControl

import (
	"context"
	"io"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"strings"
)

// GetLog Get
func GetLog() (log *CephLog, retErr error) {
	selector := labels.Set{"app": "cubeuniverse"}.AsSelector()

	// 创建一个PodLogOptions对象，用于指定要获取的log的选项
	opts := &coreV1.PodLogOptions{
		Follow:    false,
		Container: "",
	}

	// 使用CoreV1Client的Pods方法获取符合条件的Pod的log
	pods, err := ClientSet.CoreV1().Pods("cubeuniverse").List(context.Background(),
		metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return nil, err
	}
	cephLog := &CephLog{}

	// 遍历每个Pod，获取其日志并打印到控制台
	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, "operator") {
			podReq := ClientSet.CoreV1().Pods("cubeuniverse").GetLogs(pod.Name, opts)
			podLogs, err := podReq.Stream(context.Background())
			outLog, _ := io.ReadAll(podLogs)
			if err != nil {
				retErr = err
			}
			if string(outLog) != "" {
				cephLog.Operator = string(outLog)
			}
		} else if strings.Contains(pod.Name, "backend") {
			podReq := ClientSet.CoreV1().Pods("cubeuniverse").GetLogs(pod.Name, opts)
			podLogs, err := podReq.Stream(context.Background())
			outLog, _ := io.ReadAll(podLogs)
			if err != nil {
				retErr = err
			}
			if string(outLog) != "" {
				cephLog.Backend = string(outLog)
			}
		}
	}
	return cephLog, nil

}
