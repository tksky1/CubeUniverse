package universalFuncs

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"time"
)

// 这个go文件用于实现pod间的互斥机制，类似于mutex.Lock和Unlock，但也依赖于心跳和时间

// SetInUse 标注该key已被占用，占用者为uuid。类似于mutex.lock
func SetInUse(clientSet *kubernetes.Clientset, key string, uuid string) {
	data := make(map[string]string)
	data["uuid"] = uuid
	data["time"] = time.Now().Format(time.RFC3339)
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      key,
			Namespace: "cubeuniverse",
		},
		Data: data,
	}
	dataBytes, _ := json.Marshal(data)
	_, err1 := clientSet.CoreV1().ConfigMaps("cubeuniverse").Patch(context.Background(), key,
		types.MergePatchType, []byte(fmt.Sprintf(`{"data":%s}`, string(dataBytes))), metav1.PatchOptions{})
	if err1 != nil {
		_, err2 := clientSet.CoreV1().ConfigMaps("cubeuniverse").Create(context.Background(), configMap, metav1.CreateOptions{})
		if err2 != nil {
			log.Println("设置pod互斥锁失败：", err1.Error(), err2.Error())
		}
	}
	pod, err := GetPodNow(clientSet)
	if err != nil {
		log.Println("设置pod互斥锁失败：", err.Error())
	}
	pod.Labels["status"] = "ready"
	_, err = clientSet.CoreV1().Pods(pod.Namespace).Update(context.Background(), pod, metav1.UpdateOptions{})
	if err != nil {
		log.Println("设置pod互斥锁失败：", err.Error())
	}

}

// CheckInUse 检查key对应的互斥锁是否已被占用，如果已占用，返回占用者uuid、占用时间
func CheckInUse(clientSet *kubernetes.Clientset, key string) (locked bool, uuid string, lockTime time.Time) {
	locked = true
	cm, err := clientSet.CoreV1().ConfigMaps("cubeuniverse").Get(context.Background(), key, metav1.GetOptions{})
	if err != nil {
		return false, "", time.Now()
	}
	if cm == nil {
		return false, "", time.Now()
	}
	uuid = cm.Data["uuid"]
	timeString := cm.Data["time"]
	lockTime, _ = time.Parse(time.RFC3339, timeString)
	return locked, uuid, lockTime
}

// GetPodNow 返回运行本代码的pod
func GetPodNow(clientset *kubernetes.Clientset) (*v1.Pod, error) {
	pod, err := clientset.CoreV1().Pods("cubeuniverse").Get(context.Background(), os.Getenv("HOSTNAME"),
		metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pod, nil
}
