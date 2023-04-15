package cubeOperatorKit

import (
	"CubeUniverse/universalFuncs"
	"encoding/json"
	"github.com/bitly/go-simplejson"
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
	producer, err = InitKafkaProducer()
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
		count := len(jsonNow.MustArray())
		for j := 0; j < count; j++ {
			word := jsonNow.GetIndex(j).MustString()
			storedObject, err := GetObject(namespace, bucketClaim, "cubeuniverse/"+word)
			if storedObject == nil || err != nil {
				storeJsonByte, _ := json.Marshal([]string{objectKey})
				err := PutObject(namespace, bucketClaim, "cubeuniverse/"+word, storeJsonByte)
				if err != nil {
					log.Println("向对象存储桶声明"+bucketClaim+"写入索引失败：", err)
					return
				}
			} else {
				jsonStored, _ := simplejson.NewJson(storedObject)
				storedArray := jsonStored.MustStringArray()
				for _, theKey := range storedArray {
					if theKey == objectKey {
						return
					}
				}
				storedArray = append(storedArray, objectKey)
				storeJsonByte, _ := json.Marshal(storedArray)
				err := PutObject(namespace, bucketClaim, "cubeuniverse/"+word, storeJsonByte)
				if err != nil {
					log.Println("向对象存储桶声明"+bucketClaim+"写入索引失败：", err)
					return
				}
			}
		}
	}
}
