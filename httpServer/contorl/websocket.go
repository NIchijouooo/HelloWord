package contorl

import (
	"encoding/json"
	"gateway/device"
	"gateway/device/eventBus"
	"gateway/report/mqttRT"
	"gateway/setting"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Connection struct {
	Conn     *websocket.Conn
	inChan   chan []byte
	outChan  chan []byte
	exitChan chan bool
	isClose  bool
}

const (
	LogMsgType_Collect int = iota
	LogMsgType_Report
	LogMsgType_System
)

var (
	Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func NewConnection(conn *websocket.Conn) (*Connection, error) {
	wsConn := &Connection{
		Conn:     conn,
		inChan:   make(chan []byte, 1000),
		outChan:  make(chan []byte, 1000),
		exitChan: make(chan bool, 1),
		isClose:  false,
	}
	return wsConn, nil
}

func InitWebsocket(context *gin.Context) {
	wbConn, err := Upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		setting.ZAPS.Errorf("websocket协议升级失败 %v", err)
		context.JSON(http.StatusOK, gin.H{
			"code":    "1",
			"message": "websocket协议升级失败",
			"data":    "",
		})
		return
	}

	name := context.Query("name")
	msgType := context.Query("type")
	if msgType == "" {
		setting.ZAPS.Errorf("websocket请求参数type为空")
		context.JSON(http.StatusOK, gin.H{
			"code":    "1",
			"message": "websocket请求参数type为空",
			"data":    "",
		})
		return
	}
	tyInt, _ := strconv.Atoi(msgType)
	param := struct {
		Type int    `json:"type"`
		Name string `json:"name"`
	}{
		Type: tyInt,
		Name: name,
	}

	ticker := time.NewTicker(5 * time.Second)
	sub := eventBus.NewSub()
	switch param.Type {
	case LogMsgType_Collect:
		{
			device.CollectInterfaceMap.Lock.RLock()
			coll, ok := device.CollectInterfaceMap.Coll[param.Name]
			device.CollectInterfaceMap.Lock.RUnlock()

			if ok {
				coll.MessageEventBus.Subscribe(param.Name, sub)
			}
		}
	case LogMsgType_Report:
		{
			for _, v := range mqttRT.ReportServiceParamListRT.ServiceList {
				if v.GWParam.ServiceName == param.Name {
					v.MessageEventBus.Subscribe(param.Name, sub)
				}
			}
		}
	case LogMsgType_System:
		{
			setting.ZapEventBus.Subscribe("zap", sub)
		}
	}

	conn, _ := NewConnection(wbConn)
	wbConn.SetCloseHandler(conn.CloseMessage)
	go conn.ReadMessage(name, sub)

	for {
		select {
		case <-ticker.C:
			{
				//setting.ZAPS.Infof("WebSocket[%s]5秒打印一次消息", name)
			}
		case subMsg := <-sub.Out():
			{
				res, err := json.Marshal(subMsg)
				if err != nil {
					setting.ZAPS.Errorf("WebSocket[%s]发送数据格式化错误%v", name, err)
					continue
				}
				err = conn.Conn.WriteMessage(websocket.TextMessage, res)
				if err != nil {
					setting.ZAPS.Errorf("WebSocket[%s]发送数据错误%v", name, err)
					if err.Error() == "connection is closed" {
						setting.ZAPS.Errorf("WebSocket[%s]关闭连接", name)
						conn.isClose = true
						_ = conn.Conn.Close()
						return
					}
				}
			}
		case rdMsg := <-conn.inChan:
			{
				if len(rdMsg) == 0 {
					continue
				}
				setting.ZAPS.Infof("WebSocket[%s]接收数据[%s]", name, string(rdMsg))
				if strings.Contains(string(rdMsg), "close") {
					conn.isClose = true
					_ = conn.Conn.Close()
					return
				}
			}
		case <-conn.exitChan:
			{
				_ = conn.Conn.Close()
				return
			}
		}
	}
}

func (conn *Connection) ReadMessage(name string, sub eventBus.Sub) {
	for {
		if conn.isClose {
			setting.ZAPS.Errorf("WebSocket[%s]连接关闭", name)

			setting.ZapEventBus.UnSubscribe("zap", sub)
			/* 取消订阅采集接口 */
			device.CollectInterfaceMap.Lock.Lock()
			for _, v := range device.CollectInterfaceMap.Coll {
				v.MessageEventBus.UnSubscribe(v.CollInterfaceName, sub)
			}
			device.CollectInterfaceMap.Lock.Unlock()
			/* 取消订阅mqttRT上报服务 */
			for _, v := range mqttRT.ReportServiceParamListRT.ServiceList {
				v.MessageEventBus.UnSubscribe(v.GWParam.ServiceName, sub)
			}

			return
		}

		_, rdMsg, err := conn.Conn.ReadMessage()
		if err != nil && err.Error() == "websocket: close 1001 (going away)" {
			setting.ZAPS.Errorf("WebSocket[%s]连接关闭", name)
			setting.ZapEventBus.UnSubscribe("zap", sub)
			/* 取消订阅采集接口 */
			device.CollectInterfaceMap.Lock.Lock()
			for _, v := range device.CollectInterfaceMap.Coll {
				v.MessageEventBus.UnSubscribe(v.CollInterfaceName, sub)
			}
			device.CollectInterfaceMap.Lock.Unlock()
			/* 取消订阅mqttRT上报服务 */
			for _, v := range mqttRT.ReportServiceParamListRT.ServiceList {
				v.MessageEventBus.UnSubscribe(v.GWParam.ServiceName, sub)
			}
			return
		} else if err != nil {
			setting.ZAPS.Errorf("WebSocket[%s]读取消息失败 %v", name, err)
			continue
		}
		conn.inChan <- rdMsg
	}
}

func (conn *Connection) CloseMessage(code int, text string) error {

	conn.isClose = true

	conn.exitChan <- true

	setting.ZAPS.Infof("websocket关闭代码[%v]信息[%v]", code, text)
	_ = conn.Conn.Close()
	return nil
}
