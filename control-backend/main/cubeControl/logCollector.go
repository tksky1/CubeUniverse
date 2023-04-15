package cubeControl

import (
	"context"
	"io"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"log"
	"strings"
	"sync"
	"time"
)

// CephLogNow 从这里读取log，记得先检查锁
var CephLogNow = CephLog{}
var LogMutex sync.Mutex

// GetLog 循环刷新Operator和Backend的Log，展示给用户
func GetLog() {

	for {
		time.Sleep(3 * time.Second)

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
			log.Println(err)
			continue
		}

		LogMutex.Lock()

		// 遍历每个Pod，获取其日志并打印到控制台
		for _, pod := range pods.Items {
			if strings.Contains(pod.Name, "operator") {
				podReq := ClientSet.CoreV1().Pods("cubeuniverse").GetLogs(pod.Name, opts)
				podLogs, err := podReq.Stream(context.Background())
				if err != nil {
					log.Println(err)
					continue
				}
				outLog, err := io.ReadAll(podLogs)
				if err != nil {
					log.Println(err)
				}
				if string(outLog) != "" {
					CephLogNow.Operator = getLast100Lines(string(outLog))
				}
			} else if strings.Contains(pod.Name, "backend") {
				podReq := ClientSet.CoreV1().Pods("cubeuniverse").GetLogs(pod.Name, opts)
				podLogs, err := podReq.Stream(context.Background())
				if err != nil {
					log.Println(err)
					continue
				}
				outLog, _ := io.ReadAll(podLogs)
				if string(outLog) != "" {
					CephLogNow.Backend = getLast100Lines(string(outLog))
				}
			}
		}
		LogMutex.Unlock()

	}

}

func getLast100Lines(str string) string {
	lines := strings.Split(str, "\n")
	count := len(lines)
	if count <= 100 {
		return str
	}
	last100 := lines[count-100:]
	return "\n" + strings.Join(last100, "\n")
}
