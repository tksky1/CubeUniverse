package cubeOperatorKit

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"log"
	"strings"
)

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
				jsonStored, _ := simplejson.NewJson([]byte(storedObject))
				storedArray := jsonStored.MustStringArray()
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
