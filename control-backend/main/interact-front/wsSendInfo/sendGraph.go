package wsSendInfo

import (
	"container/list"
	"control-backend/cubeControl"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 做一个map的缓存队列,最大存储1000个map信息
var queue = list.New()

func ConstSend(ctx *gin.Context) {
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
	ws.SetReadLimit(1024) // 设置读取缓冲区大小为1024字节
	defer ws.Close()

	//测试部分：TODO
	for i := 1; i <= 10; i++ {
		//500毫秒的间歇
		time.Sleep(500 * time.Millisecond)
		resMap := make(gin.H)
		resMap["CephHosts"] = "CephHosts"
		resMap["inQuorumMonitor"], resMap["outQuorumMonitor"] = "inQuorumMonitor", "outQuorumMonitor"
		resMap["CephOSD"] = "CephOSD"
		resMap["CephPool"] = "CephPool"
		resMap["CephPerformance"] = "CephPerformance"
		//将map入队
		queue.PushBack(resMap)
	}
	for {
		//从缓存中拿数据
		resMap := queue.Front().Value //取出队头元素
		//测试就不出队了
		// queue.Remove(queue.Front())   //删除队头，即出队
		time.Sleep(1 * time.Second)
		println("working on sending")
		//格式化json
		msg, _ := json.Marshal(&resMap)
		//发送数据
		if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			//发送错误，记录到日志
			log.Print(err.Error())

		}
		//设置读取用户请求时间，也能实现心跳
		ws.SetReadDeadline(time.Now().Add(4 * time.Second))
		//读取用户返回数据，用于用户主动断开连接，以及长时间无用户响应而终止
		mt, msg, errRead := ws.ReadMessage()
		if errRead != nil {
			fmt.Println("err : " + errRead.Error())
			ws.WriteMessage(websocket.TextMessage, []byte("session over"))
			return
		}
		if errRead == nil && mt == websocket.TextMessage {
			jsons := make(map[string]interface{})
			//将json字符串解析
			if errCtx := json.Unmarshal(msg, &jsons); errCtx != nil {
				log.Print(errCtx.Error())
			}
			//说明要结束
			if value, ok := jsons["over"].(string); ok && value == "yes" {
				//关闭连接
				ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return //结束死循环
			}

		}

	}
	return
	//记得删除

	//向缓存先填入10个信息
	for i := 1; i <= 10; i++ {
		resMap := make(gin.H)
		resMap["CephHosts"], _ = cubeControl.GetCephHosts()
		resMap["inQuorumMonitor"], resMap["outQuorumMonitor"], _ = cubeControl.GetCephMonitor()
		resMap["CephOSD"], _ = cubeControl.GetCephOSD()
		resMap["CephPool"], _ = cubeControl.GetCephPool()
		resMap["CephPerformance"], _ = cubeControl.GetCephPerformance()
		//将map入队
		queue.PushBack(resMap)
	}
	//开启一个协程填入缓存数据
	go func() {
		for {
			//如果链表过长休眠等待资源消耗
			if queue.Len() >= 100 {
				time.Sleep(5 * time.Second)
			}
			//500毫秒的间歇
			time.Sleep(500 * time.Millisecond)
			resMap := make(gin.H)
			resMap["CephHosts"], _ = cubeControl.GetCephHosts()
			resMap["inQuorumMonitor"], resMap["outQuorumMonitor"], _ = cubeControl.GetCephMonitor()
			resMap["CephOSD"], _ = cubeControl.GetCephOSD()
			resMap["CephPool"], _ = cubeControl.GetCephPool()
			resMap["CephPerformance"], _ = cubeControl.GetCephPerformance()
			//将map入队
			queue.PushBack(resMap)
		}
	}()
	//进入死循环
	for {
		time.Sleep(3 * time.Second)
		//从缓存中拿数据
		//缓存吃空了
		if queue.Len() <= 0 {
			time.Sleep(1 * time.Second) //歇一下
		}
		time.Sleep(1 * time.Second)
		resMap := queue.Front().Value //取出队头元素
		queue.Remove(queue.Front())   //删除队头，即出队
		//设置读取用户请求时间，也能实现心跳
		ws.SetReadDeadline(time.Now().Add(4 * time.Second))
		//格式化json
		msg, _ := json.Marshal(&resMap)
		//发送数据
		if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			//发送错误，记录到日志
			log.Print(err.Error())

		}
		//读取用户返回数据，用于用户主动断开连接,以及长时间无用户响应而终止
		mt, msg, err := ws.ReadMessage()
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("session over"))
			return
		}
		if err != nil && mt == websocket.TextMessage {
			jsons := make(map[string]interface{})
			//将json字符串解析
			if errCtx := json.Unmarshal(msg, &jsons); errCtx != nil {
				log.Print(errCtx.Error())
			}
			//说明要结束
			if value, ok := jsons["over"].(string); ok && value == "yes" {
				//关闭连接
				ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return //结束死循环
			}

		}

	}
}
