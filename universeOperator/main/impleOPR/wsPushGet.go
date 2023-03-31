package impleOPR

import (
	"encoding/json"
	"fmt"
	"log"
	kit "main/cubeOperatorKit"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 通过websocket实现pushget
func ConstPushGet(ctx *gin.Context) {
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
			//检测msg，好debug
			fmt.Println("done : get msg :---" + string(msg))
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

// // 测试方法:TODO
// func pushgetImpleTest(jsons gin.H, ws *websocket.Conn) {
// 	var namespace, bucketClaimName, key, actType, blockStr string
// 	var blockNum int = 1
// 	var value []byte
// 	if valueStr, ok := jsons["namespace"].(string); ok {
// 		namespace = valueStr
// 	} else {
// 		ws.WriteMessage(websocket.TextMessage, []byte("namespace should be string")) //返回错误反馈
// 		return
// 	}
// 	if valueStr, ok := jsons["name"].(string); ok {
// 		bucketClaimName = valueStr
// 	} else {
// 		ws.WriteMessage(websocket.TextMessage, []byte("name should be string")) //返回错误反馈
// 		return
// 	}
// 	if valueStr, ok := jsons["key"].(string); ok {
// 		key = valueStr
// 	} else {
// 		ws.WriteMessage(websocket.TextMessage, []byte("key should be string")) //返回错误反馈
// 		return
// 	}
// 	if valueStr, ok := jsons["X-action"].(string); ok {
// 		actType = valueStr
// 	} else {
// 		ws.WriteMessage(websocket.TextMessage, []byte("X-action should be string"))
// 		return
// 	}
// 	if valueStr, ok := jsons["block"].(string); ok { //加入分块的机制，运行用户选择数据的分块运输
// 		blockStr = valueStr
// 	}
// 	//对于分块数，如果没写的话默认为1
// 	if blockStr == "" {
// 		blockNum = 1
// 	} else {
// 		var err error = nil
// 		blockNum, err = strconv.Atoi(blockStr)
// 		if err != nil {
// 			ws.WriteMessage(websocket.TextMessage, []byte("block should be string represent a number"))
// 			return
// 		}
// 	}
// 	//对于value数据，判断其为string还是[]byte
// 	if valueStr, ok := jsons["value"].(string); ok {
// 		value = []byte(valueStr)
// 	} else {
// 		valueByte, err := jsons["value"].([]byte)
// 		if err { //如果是byte
// 			value = valueByte
// 		}
// 		if !err && actType == "push" {
// 			ws.WriteMessage(websocket.TextMessage, []byte("value should be string or []byte"))
// 			return
// 		}
// 	}

// 	switch strings.ToLower(actType) {
// 	case "push":
// 		fmt.Printf("%s-%s-%s-%s", namespace, bucketClaimName, key, string(value))
// 		ws.WriteMessage(websocket.TextMessage, []byte("put success"))
// 		return

// 	case "get":
// 		value := "this is test value 23333"
// 		fmt.Println([]byte(value))
// 		//根据block进行分组
// 		for index, valueBytes := range splitArray([]byte(value), blockNum) {
// 			valueMap := map[string][]byte{
// 				"value" + strconv.Itoa(index): valueBytes,
// 				"key":                         []byte(key),
// 				"namespace":                   []byte(namespace),
// 			}
// 			valueJson, _ := json.Marshal(&valueMap)
// 			fmt.Println(valueMap["value"])
// 			ws.WriteMessage(websocket.TextMessage, valueJson)
// 		}
// 		return
// 	}

// }

// // 记得删除
func pushgetImple(jsons gin.H, ws *websocket.Conn) {
	var namespace, bucketClaimName, key, actType, blockStr string
	var blockNum int = 1
	var value []byte
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
	if valueStr, ok := jsons["key"].(string); ok {
		key = valueStr
	} else {
		ws.WriteMessage(websocket.TextMessage, []byte("key should be string")) //返回错误反馈
		return
	}
	if valueStr, ok := jsons["X-action"].(string); ok {
		actType = valueStr
	} else {
		ws.WriteMessage(websocket.TextMessage, []byte("X-action should be string"))
		return
	}
	if valueStr, ok := jsons["block"].(string); ok { //加入分块的机制，运行用户选择数据的分块运输
		blockStr = valueStr
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
	//对于value数据，判断其为string还是[]byte
	if valueStr, ok := jsons["value"].(string); ok {
		value = []byte(valueStr)
	} else {
		valueByte, err := jsons["value"].([]byte)
		if err { //如果是byte
			value = valueByte
		}
		if !err && actType == "push" {
			ws.WriteMessage(websocket.TextMessage, []byte("value should be string or []byte"))
			return
		}
	}

	switch strings.ToLower(actType) {
	case "push":
		if err := kit.PutObject(namespace, bucketClaimName, key, value); err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("Fail Put OBJ: "+err.Error()))
		}
		ws.WriteMessage(websocket.TextMessage, []byte("put success"))
		return

	case "get":
		value, err := kit.GetObject(namespace, bucketClaimName, key)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("Fail Put OBJ: "+err.Error()))
			return
		}
		//根据block进行分组
		for index, valueBytes := range splitArray([]byte(value), blockNum) {
			valueMap := map[string][]byte{
				"value" + strconv.Itoa(index): valueBytes,
				"key":                         []byte(key),
				"namespace":                   []byte(namespace),
			}
			valueJson, _ := json.Marshal(&valueMap)
			fmt.Println(valueMap["value"])
			ws.WriteMessage(websocket.TextMessage, valueJson)
		}
		return
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
