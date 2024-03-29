package cubeOperatorKit

import (
	"CubeUniverse/universalFuncs"
	"bytes"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"io"
	"log"
	"strings"
	"time"
)

// StartML 用于轮询ML组件状况然后启动队列监听
func StartML() {
	for {
		time.Sleep(3 * time.Second)
		kafka, ml := universalFuncs.CheckMLStatus(ClientSet)
		if kafka && ml {
			break
		}
	}
	log.Println("ML加载完成，开始初始化与ML的交互..")
	var err error
	Producer, err = InitKafkaProducer()
	if err != nil {
		log.Println("准备ML出错：" + err.Error())
		StartML()
		return
	}
	consumer, err := InitKafkaConsumer()
	if err != nil {
		log.Println("准备ML出错：" + err.Error())
		StartML()
		return
	}
	err = ConsumerStartListening(*consumer, ConsumeML)
	log.Println("准备ML出错：" + err.Error())
	StartML()
}

// ConsumeML 作为参数提供给消息处理
func ConsumeML(key string, value string) {

	log.Println("开始处理消息" + key)

	index1 := strings.IndexByte(key, '%')
	namespace := key[:index1]
	index2 := strings.IndexByte(key, ':')
	bucketClaim := key[index1+1 : index2]
	objectKey := key[index2+1:]
	theJson, err := simplejson.NewJson([]byte(value))
	if err != nil {
		log.Println("处理kafka返回值失败：", err)
		return
	}
	for i := 1; i <= 3; i++ {
		jsonNow := theJson.GetIndex(i)
		stringNow := jsonNow.MustString()
		outputs := strings.Split(stringNow, ", ")
		for _, word := range outputs {
			storedObject, err := GetObject(namespace, bucketClaim, "cubeuniverse/"+word)
			if storedObject == nil || err != nil {
				storeJsonByte, _ := json.Marshal([]string{objectKey})
				reader := bytes.NewReader(storeJsonByte)
				reader2 := io.Reader(reader)
				err := PutObject(namespace, bucketClaim, "cubeuniverse/"+word, 0, &reader2)
				if err != nil {
					log.Println("向对象存储桶声明"+bucketClaim+"写入索引失败：", err)
					return
				}
			} else {
				jsonStored, _ := simplejson.NewJson(*storedObject)
				storedArray := jsonStored.MustStringArray()
				for _, theKey := range storedArray {
					if theKey == objectKey {
						return
					}
				}
				storedArray = append(storedArray, objectKey)
				storeJsonByte, _ := json.Marshal(storedArray)
				reader := bytes.NewReader(storeJsonByte)
				reader2 := io.Reader(reader)
				err := PutObject(namespace, bucketClaim, "cubeuniverse/"+word, 0, &reader2)
				if err != nil {
					log.Println("向对象存储桶声明"+bucketClaim+"写入索引失败：", err)
					return
				}
			}
		}
	}

	log.Println("处理" + key + "完成")

}
