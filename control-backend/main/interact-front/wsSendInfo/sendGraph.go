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

// 定义channel,字节数组，并初始化,设置缓冲区大小为10
// 不能放在循环体里面，不然每次都会被初始化丢失数据,我这里为了方便就直接放在全局变量了
var readInfo chan []byte = make(chan []byte, 10)

type drawData struct {
	ReadBytesPerSec       int `json:"ReadBytesPerSec"`       // 每秒读的bytes
	ReadOperationsPerSec  int `json:"ReadOperationsPerSec"`  // 每秒读操作数
	WriteBytesPerSec      int `json:"WriteBytesPerSec"`      // 每秒写的bytes
	WriteOperationPerSec  int `json:"WriteOperationPerSec"`  // 每秒写操作数
	RecoveringBytesPerSec int `json:"RecoveringBytesPerSec"` // 每秒恢复数据流量
}

// 做一个map的缓存队列,最大存储1000个map信息
var queue = list.New()

// 用于在协程中完成读取websocket的操作
func readMsg(ws *websocket.Conn, readInfo chan []byte) {
	for {
		mt, msg, errRead := ws.ReadMessage()
		if errRead != nil {
			fmt.Println("err : " + errRead.Error())
			ws.WriteMessage(websocket.TextMessage, []byte("session over"))
			//关闭websocket
			ws.Close()
			return
		}
		//如果读取成功,将读取结果发送到channel中
		if errRead == nil && mt == websocket.TextMessage {
			fmt.Println("done : get msg :---" + string(msg))
			readInfo <- msg
		}
	}

}

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
	ws.SetReadLimit(1024 * 10) // 设置读取缓冲区大小为1024字节
	defer ws.Close()

	//向缓存先填入10个信息
	for i := 1; i <= 10; i++ {
		resMap := make(gin.H)
		resMap["CephHosts"], err = cubeControl.GetCephHosts()
		if err != nil {
			log.Println("调用ceph数据出错：" + err.Error())
		}
		resMap["inQuorumMonitors"], resMap["outQuorumMonitor"], err = cubeControl.GetCephMonitor()
		if err != nil {
			log.Println("调用ceph数据出错：" + err.Error())
		}
		resMap["CephOSDs"], err = cubeControl.GetCephOSD()
		if err != nil {
			log.Println("调用ceph数据出错：" + err.Error())
		}
		resMap["CephPools"], err = cubeControl.GetCephPool()
		if err != nil {
			log.Println("调用ceph数据出错：" + err.Error())
		}
		resMap["CephPerformance"], err = cubeControl.GetCephPerformance()
		if err != nil {
			log.Println("调用ceph数据出错：" + err.Error())
		}
		var dataArray [11]drawData
		for times := 0; times < 10; {
			performanceArray, err := cubeControl.GetCephPerformance()
			if err != nil {
				log.Println("调用getCephPerformace err: " + err.Error())
				continue
			}
			var data drawData
			data.ReadBytesPerSec = performanceArray.ReadBytesPerSec
			data.ReadOperationsPerSec = performanceArray.ReadOperationsPerSec
			data.RecoveringBytesPerSec = performanceArray.RecoveringBytesPerSec
			data.WriteBytesPerSec = performanceArray.WriteBytesPerSec
			data.WriteOperationPerSec = performanceArray.WriteOperationPerSec
			dataArray[times] = data
			times++
		}
		//加入一个画数据表所用的数据
		resMap["CephDrawData"] = dataArray

		//日志信息的结构体
		cubeControl.LogMutex.Lock()
		logStruct := cubeControl.CephLogNow
		resMap["Operatorlog"] = logStruct.Operator
		resMap["Backendlog"] = logStruct.Backend
		cubeControl.LogMutex.Unlock()
		//将map入队
		queue.PushBack(resMap)
	}
	//开启一个协程填入缓存数据
	go func() {
		for {
			//如果链表过长休眠等待资源消耗
			if queue.Len() >= 20 {
				time.Sleep(2 * time.Second)
			}
			//500毫秒的间歇
			time.Sleep(500 * time.Millisecond)
			resMap := make(gin.H)
			resMap["CephHosts"], err = cubeControl.GetCephHosts()
			if err != nil {
				log.Println("调用ceph数据出错：" + err.Error())
			}
			resMap["inQuorumMonitors"], resMap["outQuorumMonitors"], err = cubeControl.GetCephMonitor()
			if err != nil {
				log.Println("调用ceph数据出错：" + err.Error())
			}
			resMap["CephOSDs"], err = cubeControl.GetCephOSD()
			if err != nil {
				log.Println("调用ceph数据出错：" + err.Error())
			}
			resMap["CephPools"], err = cubeControl.GetCephPool()
			if err != nil {
				log.Println("调用ceph数据出错：" + err.Error())
			}
			resMap["CephPerformance"], err = cubeControl.GetCephPerformance()
			if err != nil {
				log.Println("调用ceph数据出错：" + err.Error())
			}

			var dataArray [11]drawData
			for times := 0; times < 10; {
				performanceArray, err := cubeControl.GetCephPerformance()
				if err != nil {
					log.Println("调用getCephPerformace err: " + err.Error())
					continue
				}
				var data drawData
				data.ReadBytesPerSec = performanceArray.ReadBytesPerSec
				data.ReadOperationsPerSec = performanceArray.ReadOperationsPerSec
				data.RecoveringBytesPerSec = performanceArray.RecoveringBytesPerSec
				data.WriteBytesPerSec = performanceArray.WriteBytesPerSec
				data.WriteOperationPerSec = performanceArray.WriteOperationPerSec
				dataArray[times] = data
				times++
			}
			//加入一个画数据表所用的数据
			resMap["CephDrawData"] = dataArray

			//日志信息的结构体
			cubeControl.LogMutex.Lock()
			logStruct := cubeControl.CephLogNow
			resMap["Operatorlog"] = logStruct.Operator
			resMap["Backendlog"] = logStruct.Backend
			cubeControl.LogMutex.Unlock()
			//将map入队
			queue.PushBack(resMap)
		}
	}()
	//进入死循环
	for {
		//从缓存中拿数据
		//缓存吃空了
		if queue.Len() <= 0 {
			time.Sleep(1 * time.Second) //歇一下
			continue
		}
		time.Sleep(2 * time.Second)
		resMap := queue.Front().Value //取出队头元素
		queue.Remove(queue.Front())   //删除队头，即出队

		//开启读取协程
		go readMsg(ws, readInfo)

		//通过select语句进行channel读取
		select {
		case msg1 := <-readInfo:
			fmt.Println("in readInfo---" + string(msg1))
			jsons := make(map[string]interface{})
			//将json字符串解析
			if errCtx := json.Unmarshal(msg1, &jsons); errCtx != nil {
				log.Println(errCtx.Error())
			}
			//说明要结束
			if value, ok := jsons["over"].(string); ok && value == "yes" {
				if err := ws.WriteMessage(websocket.TextMessage, []byte("bye")); err != nil {
					log.Println(err.Error())
					ws.Close()
				}
				time.Sleep(500 * time.Millisecond)
				//关闭连接
				if err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
					log.Println(err.Error())
					ws.Close()
				}
				return //结束死循环
			} else if ok { //如果value不为yes，但没有报错
				//格式化json
				msg, _ := json.Marshal(&resMap)
				//发送数据
				if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
					//发送错误，记录到日志
					log.Println(err.Error())
					ws.Close()
				}
			}
		default: //默认没有从管道接收到数据的情况下，发送数据
			//格式化json
			msg, _ := json.Marshal(&resMap)
			//发送数据
			if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				//发送错误，记录到日志
				log.Println(err.Error())
				ws.Close()
			}
		}

	}
}
