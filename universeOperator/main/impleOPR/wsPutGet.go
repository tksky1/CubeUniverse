package impleOPR

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	kit "main/cubeOperatorKit"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 通过websocket实现pushget
func ConstPushGetDeleteList(ctx *gin.Context) {
	var upgrade = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}, //对于长连接的每次发送都会校验用户
	}
	ws, err := upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{
			"code": 405,
			"data": nil,
			"msg":  "can not establish connection",
		})
		return
	}
	ws.SetReadLimit(1024 * 30) // 设置读取缓冲区大小为1024字节
	defer ws.Close()
	for {
		//循环接收message
		mt, msg, errRead := ws.ReadMessage()
		if errRead != nil { //如果读取出错，报出错误
			fmt.Println("err : " + errRead.Error())
			ws.WriteMessage(websocket.TextMessage, []byte("session over"))
			return
		}
		//如果读取成功,将读取解析为json
		if errRead == nil && mt == websocket.TextMessage {
			// //检测msg，好debug
			// fmt.Println("done : get msg :---" + string(msg))
			//转化为json格式
			jsons := make(map[string]interface{})
			//将json字符串解析
			if errCtx := json.Unmarshal(msg, &jsons); errCtx != nil {
				log.Print(errCtx.Error())
			}
			//说明要结束
			if value, ok := jsons["over"].(string); ok && value == "yes" {
				ws.WriteMessage(websocket.TextMessage, []byte("bye"))
				time.Sleep(500 * time.Millisecond)
				//关闭连接
				ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return //结束死循环
			} else { //如果value不为yes，则检查是需要调用什么方法
				//调用事件处理函数
				pushgetImple(jsons, ws)

			}
		}

	}
}

func pushgetImple(jsons gin.H, ws *websocket.Conn) {
	var namespace, bucketClaimName, key, actType, blockStr, indexBlock, tag string
	tag = "" //给tag一个默认值
	var blockNum, indexNum int = 1, 0
	var value []byte
	if valueStr, ok := jsons["X-action"].(string); ok { //确定方法
		actType = valueStr
	} else {
		ws.WriteMessage(websocket.TextMessage, []byte("X-action should be string"))
		return
	}
	if valueStr, ok := jsons["namespace"].(string); ok {
		namespace = valueStr
	} else {
		ws.WriteMessage(websocket.TextMessage, []byte("namespace should be string")) //返回错误反馈
		return
	}
	if valueStr, ok := jsons["name"].(string); ok {
		bucketClaimName = valueStr
	} else {
		ws.WriteMessage(websocket.TextMessage, []byte("name should be string")) //返回错误反馈
		return
	}
	if valueStr, ok := jsons["key"].(string); ok || strings.ToLower(actType) == "list" { //调用list方法的时候可以不需要key的存在
		key = valueStr
	} else {
		ws.WriteMessage(websocket.TextMessage, []byte("key should be string")) //返回错误反馈
		return
	}

	if valueStr, ok := jsons["block"].(string); ok { //加入分块的机制，运行用户选择数据的分块运输
		blockStr = valueStr
	}
	if valueStr, ok := jsons["index"].(string); ok { //加入分块的机制的索引，运行用户选择数据的分块运输
		indexBlock = valueStr
	}
	if valueStr, ok := jsons["tag"].(string); ok { //得到tag标签的值，此为可选参数
		tag = valueStr
	}
	//对于分块数，如果没写的话默认为1
	if blockStr == "" {
		blockNum = 1
	} else {
		var err error = nil
		blockNum, err = strconv.Atoi(blockStr)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("block should be string represent a number"))
			return
		}
	}
	//对于索引值，如果没写的话默认为0
	if indexBlock == "" {
		indexNum = 0
	} else {
		var err error = nil
		indexNum, err = strconv.Atoi(indexBlock)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("index should be string represent a number"))
			return
		}
	}
	//保证索引比分块小
	if indexNum >= blockNum {
		ws.WriteMessage(websocket.TextMessage, []byte("index out of range"))
		return
	}

	// //对于value数据，判断其为string还是[]byte
	// if valueStr, ok := jsons["value"].(string); ok {
	// 	value = []byte(valueStr)
	// } else {
	// 	valueByte, err := jsons["value"].([]byte)
	// 	if err { //如果是byte
	// 		value = valueByte
	// 	}
	// 	if !err && actType == "put" {
	// 		ws.WriteMessage(websocket.TextMessage, []byte("value should be string or []byte"))
	// 		return
	// 	}
	// }

	switch strings.ToLower(actType) {
	case "put":
		//更换为流式传输
		var reader io.Reader = ws.UnderlyingConn().(*net.TCPConn)
		if err := kit.PutObject(namespace, bucketClaimName, key, reader); err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("Fail Put OBJ: "+err.Error()))
		}
		ws.WriteMessage(websocket.TextMessage, []byte("put success"))
		return

	case "get":
		value, err := kit.GetObject(namespace, bucketClaimName, key)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("Fail Get OBJ: "+err.Error()))
			return
		}
		//根据block进行分组
		valueBytes := splitArray([]byte(value), blockNum)
		//将数据转为string避免bytes数据被base64编码
		value2Str := valueBytes[indexNum]
		valueMap := map[string]string{
			"value" + strconv.Itoa(indexNum): string(value2Str),
			"key":                            key,
			"namespace":                      namespace,
			"name":                           bucketClaimName,
		}
		//转化为json
		valueJson, _ := json.Marshal(&valueMap)
		ws.WriteMessage(websocket.TextMessage, valueJson) //发送数据
		return
	case "delete":
		if err := kit.DeleteObject(namespace, bucketClaimName, key); err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("Fail delete OBJ: "+err.Error()))
		} else {
			ws.WriteMessage(websocket.TextMessage, []byte("delete success"))
			return
		}
	case "list":
		if tag != "" { //在list调用时，如果传入了tag说明是调用的对象存储
			if valueArr, err := kit.ListObjectByTag(namespace, bucketClaimName, tag); err != nil {
				ws.WriteMessage(websocket.TextMessage, []byte("Fail list OBJ: "+err.Error()))
				return
			} else {
				//将返回数据结构化为json格式
				valueMap := gin.H{
					"value":     valueArr,
					"namespace": namespace,
					"name":      bucketClaimName,
				}
				if valueJson, err := json.Marshal(&valueMap); err != nil {
					log.Print("err json convert: " + err.Error())
					return
				} else {
					ws.WriteMessage(websocket.TextMessage, valueJson)
					return
				}
			}
		}
		if valueArr, err := kit.ListObjectFromBucket(namespace, bucketClaimName); err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("Fail list OBJ: "+err.Error()))
			return
		} else {
			//将返回数据结构化为json格式
			valueMap := gin.H{
				"value":     valueArr,
				"namespace": namespace,
				"name":      bucketClaimName,
			}
			if valueJson, err := json.Marshal(&valueMap); err != nil {
				log.Print("err json convert: " + err.Error())
				return
			} else {
				ws.WriteMessage(websocket.TextMessage, valueJson)
				return
			}

		}
	}

}

// 数组平分
func splitArray(arr []byte, num int) [][]byte {
	max := int(len(arr))
	if max < num {
		return nil
	}
	var segmens = make([][]byte, 0)
	quantity := max / num
	end := int(0)
	for i := int(1); i <= num; i++ {
		qu := i * quantity
		if i != num {
			segmens = append(segmens, arr[i-1+end:qu])
		} else {
			segmens = append(segmens, arr[i-1+end:])
		}
		end = qu - i
	}
	return segmens
}

// 常用的结构体
type InfoNBK struct {
	Namespace       string `json:"namespace"`
	BucketclaimName string `json:"name"`
	Key             string `json:"key"`
}
