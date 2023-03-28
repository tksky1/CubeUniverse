package universalFuncs

import (
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
			cmd := exec.Command(os.Args[0])
			log.Println("被抢占，重启..")
			_ = cmd.Start()
			os.Exit(0)
		}
		SetInUse(clientSet, key, uuid)
	}

}
