package universalFuncs

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"os/exec"
	"time"
)

// HeartBeat 在需要互斥锁的程序内运行，维持心跳。如果发现心跳被其他程序破坏，则强制重启整个程序让位。用于go程。
func HeartBeat(clientSet *kubernetes.Clientset, key string, ownUUID string) {

	for {
		time.Sleep(3 * time.Second)
		locked, uuid, _ := CheckInUse(clientSet, key)
		if locked && ownUUID != uuid {
			// 被抢占，重启程序
			pod, err := GetPodNow(clientSet)
			if err != nil {
				log.Println("设置pod互斥锁失败：", err.Error())
			}
			pod.Labels["status"] = "not-ready"
			_, err = clientSet.CoreV1().Pods(pod.Namespace).Update(context.Background(), pod, metav1.UpdateOptions{})
			if err != nil {
				log.Println("设置pod互斥锁失败：", err.Error())
			}

			cmd := exec.Command(os.Args[0])
			log.Println("被抢占，重启..")
			_ = cmd.Start()
			os.Exit(0)
		}
		SetInUse(clientSet, key, uuid)
	}

}
