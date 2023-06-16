package device

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway/device/commInterface"
	"gateway/device/eventBus"
	"gateway/protocol/dlt645"
	"gateway/setting"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type CommunicationCmdTemplate struct {
	CollInterfaceName string //采集接口名称
	DeviceName        string //采集接口下设备名称
	FunName           string
	FunPara           string
}

type CommunicationDirectDataReqTemplate struct {
	CollInterfaceName string //采集接口名称
	Data              []byte
}

type CommunicationDirectDataAckTemplate struct {
	CollInterfaceName string //采集接口名称
	Status            bool
	Data              []byte
}

type CommunicationMessageTemplate struct {
	TimeStamp string `json:"timeStamp"` //时间戳
	Direction int    `json:"direction"` //数据方向
	Label     string `json:"label"`     //数据标识
	Content   string `json:"content"`   //数据内容
}

//type CommunicationRxTemplate struct {
//	Status bool
//	RxBuf  []byte
//}

type CommunicationManageTemplate struct {
	EmergencyRequestChan  chan CommunicationCmdTemplate
	EmergencyAckChan      chan ReceiveDataTemplate
	CommonRequestChan     chan CommunicationCmdTemplate
	DirectDataRequestChan chan CommunicationDirectDataReqTemplate
	DirectDataAckChan     chan CommunicationDirectDataAckTemplate
	PacketChan            chan []byte
	CommMessage           []CommunicationMessageTemplate
	QuitChan              chan bool
}

const (
	MessageDirection_TX = 1
	MessageDirection_RX = 0
)

const (
	CommunicationState_Start int = iota
	CommunicationState_Generate
	CommunicationState_Send       //命令发送
	CommunicationState_Wait       //命令等待接收
	CommunicationState_WaitSucess //命令接收成功
	CommunicationState_WaitFail   //命令接收失败
	CommunicationState_Stop
	CommunicationState_DirectDataSend //透传数据发送
	CommunicationState_DirectDataWait //透传数据等待接收
	CommunicationState_DirectDataStop //透传数据任务停止
)

func NewCommunicationManageTemplate() *CommunicationManageTemplate {

	template := &CommunicationManageTemplate{
		EmergencyRequestChan:  make(chan CommunicationCmdTemplate, 1),
		CommonRequestChan:     make(chan CommunicationCmdTemplate, 100),
		EmergencyAckChan:      make(chan ReceiveDataTemplate, 1),
		DirectDataRequestChan: make(chan CommunicationDirectDataReqTemplate, 1),
		DirectDataAckChan:     make(chan CommunicationDirectDataAckTemplate, 1),
		PacketChan:            make(chan []byte, 100), //最多连续接收100帧数据
		CommMessage:           make([]CommunicationMessageTemplate, 0),
		QuitChan:              make(chan bool, 2),
	}

	return template
}

func (c *CommunicationManageTemplate) CommunicationManageMessageAdd(collName string, dir int, buf []byte) {
	CommunicationMessage := CommunicationMessageTemplate{
		TimeStamp: time.Now().Format("2006-01-02 15:04:05.1234"),
		Direction: dir,
		Content:   fmt.Sprintf("%X", buf),
	}

	CollectInterfaceMap.Lock.Lock()
	coll, ok := CollectInterfaceMap.Coll[collName]
	if !ok {
		CollectInterfaceMap.Lock.Unlock()
		return
	}
	CollectInterfaceMap.Lock.Unlock()

	err := coll.MessageEventBus.Publish(collName, CommunicationMessage)
	if err != nil {
		//setting.ZAPS.Debugf("采集接口[%s]发布报文消息失败[%v]", collName, err)
	}

}

func (c *CommunicationManageTemplate) CommunicationManageAddCommon(cmd CommunicationCmdTemplate) {

	c.CommonRequestChan <- cmd
}

func (c *CommunicationManageTemplate) CommunicationManageAddEmergency(cmd CommunicationCmdTemplate) ReceiveDataTemplate {

	c.EmergencyRequestChan <- cmd

	return <-c.EmergencyAckChan
}

func (c *CommunicationManageTemplate) CommunicationManageAddDirectData(req CommunicationDirectDataReqTemplate) CommunicationDirectDataAckTemplate {

	c.DirectDataRequestChan <- req

	return <-c.DirectDataAckChan
}

func (c *CommunicationManageTemplate) CommunicationManageProcessReceiveData(ctx context.Context, comm commInterface.CommunicationInterface) {

	//阻塞读
	txBuf := make([]byte, 1024)
	setting.ZAPS.Debugf("通信接口[%s]接收数据协程1/4进入", comm.GetName())
	for {
		select {
		case <-ctx.Done():
			setting.ZAPS.Debugf("通信接口[%s]接收数据协程1/4退出", comm.GetName())
			return
		case <-c.QuitChan:
			setting.ZAPS.Debugf("通信接口[%s]接收数据协程1/4退出", comm.GetName())
			return
		default:
			//阻塞读
			rxBuf, _ := comm.ReadData(txBuf)
			rxBufCnt := len(rxBuf)
			if rxBufCnt > 0 {
				//setting.ZAPS.Debugf("通信接口[%s]curRxBufCnt %v", comm.GetName(), rxBufCnt)
				//setting.ZAPS.Debugf("通信接口[%s]CurRxBuf %X", comm.GetName(), rxBuf[:rxBufCnt])

				//追加接收的数据到接收缓冲区
				c.PacketChan <- rxBuf[:rxBufCnt]
				//清除本次接收数据
				rxBufCnt = 0
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func CommunicationUpdateNodeState(node *DeviceNodeTemplate, rxStatus bool, collName string, eventBus *eventBus.Bus, offLineCnt int) {
	node.CommTotalCnt++
	if rxStatus == true {
		//设备从离线变成上线
		if node.CommStatus == "offLine" {
			content := CollectInterfaceEventTemplate{
				Topic:    "onLine",
				CollName: collName,
				NodeName: node.Name,
				Content:  node.Name,
			}
			err := eventBus.Publish("onLine", content)
			if err != nil {
				setting.ZAPS.Debugf("采集接口[%s]发布节点[%s]上线消息", collName, node.Name)
			}
		} else {
			content := CollectInterfaceEventTemplate{
				Topic:    "update",
				CollName: collName,
				NodeName: node.Name,
				Content:  node.Name,
			}
			err := eventBus.Publish("update", content)
			if err != nil {
				//setting.ZAPS.Debugf("采集接口[%s]发布节点[%s]属性更新消息", collName, node.Name)
			}
		}

		node.CommSuccessCnt++
		node.CurCommFailCnt = 0
		node.CommStatus = "onLine"
		node.LastCommRTC = time.Now().Format("2006-01-02 15:04:05")
	} else {
		node.CurCommFailCnt++
		if node.CurCommFailCnt >= offLineCnt {
			node.CurCommFailCnt = 0
			//设备从上线变成离线
			if node.CommStatus == "onLine" {
				content := CollectInterfaceEventTemplate{
					Topic:    "offLine",
					CollName: collName,
					NodeName: node.Name,
					Content:  node.Name,
				}
				err := eventBus.Publish("offLine", content)
				if err != nil {
					setting.ZAPS.Debugf("采集接口[%s]发布节点[%s]离线消息", collName, node.Name)
				}
			}
			node.CommStatus = "offLine"
		}
	}
}

func (c *CommunicationManageTemplate) CommunicationStateMachine(cmd CommunicationCmdTemplate,
	collName string,
	commInterface commInterface.CommunicationInterface,
	node *DeviceNodeTemplate,
	eventBus *eventBus.Bus,
	lStateMap map[string]*lua.LState,
	offLineCnt int) ReceiveDataTemplate {
	rxResult := ReceiveDataTemplate{
		Status: false,
	}

	commState := CommunicationState_Start

	commStep := 0
	txBuf := make([]byte, 0)
	continues := false
	startT := time.Now() //计算当前时间
	func() {
		for {
			switch commState {
			case CommunicationState_Start:
				{
					commState = CommunicationState_Generate
				}
			case CommunicationState_Generate:
				{
					//--------------组包---------------------------
					result := false
					nodeVariables := make(map[string]interface{})
					for _, v := range node.Properties {
						if len(v.Value) > 0 {
							nodeVariables[v.Name] = v.Value[len(v.Value)-1].Value
						}
					}
					variablesJson, _ := json.Marshal(nodeVariables)
					if cmd.FunName == "GetDeviceRealVariables" {
						txBuf, result, continues = node.GenerateGetRealVariables(lStateMap[node.TSL], node.Addr, commStep, string(variablesJson))
						if result == false {
							commState = CommunicationState_Stop
						} else {
							commState = CommunicationState_Send
							commStep++
						}
					} else {
						txBuf, result, continues = node.DeviceCustomCmd(lStateMap[node.TSL], node.Addr,
							cmd.FunName,
							cmd.FunPara,
							commStep,
							string(variablesJson))
						if result == false {
							commState = CommunicationState_Stop
						} else {
							commState = CommunicationState_Send
							commStep++
						}
					}
				}
			case CommunicationState_Send:
				{
					//---------------发送前清空接收缓存-------------------------
					for i := 0; i < len(c.PacketChan); i++ {
						<-c.PacketChan
					}
					//---------------发送-------------------------
					_, _ = commInterface.WriteData(txBuf)
					node.CommTotalCnt++
					setting.ZAPS.Infof("采集接口[%s]设备[%s]发送数据[%d:%X]", collName, node.Name, len(txBuf), txBuf)
					c.CommunicationManageMessageAdd(collName, MessageDirection_TX, txBuf)
					commState = CommunicationState_Wait
				}
			case CommunicationState_Wait:
				{
					//阻塞读
					rxBuf := make([]byte, 256)
					rxTotalBuf := make([]byte, 0)
					rxBufCnt := 0
					rxTotalBufCnt := 0
					var timeout int
					timeout, _ = strconv.Atoi(commInterface.GetTimeOut())
					timerOut := time.NewTimer(time.Duration(timeout) * time.Millisecond)
					func() {
						for {
							select {
							//继续接收数据
							case rxBuf = <-c.PacketChan:
								{
									rxBufCnt = len(rxBuf)
									if rxBufCnt > 0 {
										rxTotalBufCnt += rxBufCnt
										//追加接收的数据到接收缓冲区
										rxTotalBuf = append(rxTotalBuf, rxBuf[:rxBufCnt]...)
									}
									//setting.ZAPS.Debugf("采集接口[%s]rxTotalCnt[%v]:rxTotalBuf[%X]", collName, rxTotalBufCnt, rxTotalBuf)
									//清除本次接收数据
									rxBufCnt = 0
									rxBuf = rxBuf[0:0]
								}
							//是否接收超时
							case <-timerOut.C:
								{
									timerOut.Stop()
									if len(rxTotalBuf) > 0 {
										c.CommunicationManageMessageAdd(collName, MessageDirection_RX, rxTotalBuf)
									}
									setting.ZAPS.Infof("采集接口[%s]设备[%s]接收超时 接收数据[%d:%X]", collName, node.Name, len(rxTotalBuf), rxTotalBuf)

									node.CurCommFailCnt++
									if node.CurCommFailCnt >= offLineCnt {
										node.CurCommFailCnt = 0
										//设备从上线变成离线
										if node.CommStatus == "onLine" {
											content := CollectInterfaceEventTemplate{
												Topic:    "offLine",
												CollName: collName,
												NodeName: node.Name,
												Content:  node.Name,
											}
											err := eventBus.Publish("offLine", content)
											if err != nil {
												setting.ZAPS.Debugf("采集接口[%s]发布节点[%s]离线消息", collName, node.Name)
											}
										}
										node.CommStatus = "offLine"
									}
									rxTotalBufCnt = 0
									rxTotalBuf = rxTotalBuf[0:0]

									commState = CommunicationState_WaitFail
									return
								}
							//是否正确收到数据包
							case rxStatus := <-node.AnalysisRx(lStateMap[node.TSL], node.Addr,
								node.Properties, rxTotalBuf, rxTotalBufCnt):
								{
									timerOut.Stop()
									setting.ZAPS.Infof("采集服务[%s]设备[%s]接收成功 接收数据[%d:%X]", collName, node.Name, len(rxTotalBuf), rxTotalBuf)

									rxResult.Status = rxStatus.Status
									rxResult.Properties = rxStatus.Properties

									c.CommunicationManageMessageAdd(collName, MessageDirection_RX, rxTotalBuf)

									//设备从离线变成上线
									if node.CommStatus == "offLine" {
										content := CollectInterfaceEventTemplate{
											Topic:    "onLine",
											CollName: collName,
											NodeName: node.Name,
											Content:  node.Name,
										}
										err := eventBus.Publish("onLine", content)
										if err != nil {
											setting.ZAPS.Debugf("采集接口[%s]发布节点[%s]上线消息", collName, node.Name)
										}
									}

									if continues == false {
										content := CollectInterfaceEventTemplate{
											Topic:    "update",
											CollName: collName,
											NodeName: node.Name,
											Content:  node.Name,
										}
										err := eventBus.Publish("update", content)
										if err != nil {
											//setting.ZAPS.Debugf("采集接口[%s]发布节点[%s]属性更新消息", collName, node.Name)
										}
									}

									node.CommSuccessCnt++
									node.CurCommFailCnt = 0
									node.CommStatus = "onLine"
									node.LastCommRTC = time.Now().Format("2006-01-02 15:04:05")

									rxTotalBufCnt = 0
									rxTotalBuf = rxTotalBuf[0:0]

									commState = CommunicationState_WaitSucess
									return
								}
							}
						}
					}()
				}
			case CommunicationState_WaitSucess:
				{
					//通信帧延时
					var interval int
					interval, _ = strconv.Atoi(commInterface.GetInterval())
					time.Sleep(time.Duration(interval) * time.Millisecond)
					commState = CommunicationState_Stop
				}
			case CommunicationState_WaitFail:
				{
					commState = CommunicationState_Stop
				}
			case CommunicationState_Stop:
				{
					tc := time.Since(startT) //计算耗时
					setting.ZAPS.Debugf("采集服务[%s]本次采集用时%s", collName, tc)
					if continues == true {
						commState = CommunicationState_Start
					} else {
						return
					}
				}
			}
		}
	}()

	return rxResult
}

func (c *CommunicationManageTemplate) CommunicationStateMachineIoIn(cmd CommunicationCmdTemplate, collName string,
	commInterface commInterface.CommunicationInterface, node *DeviceNodeTemplate,
	eventBus *eventBus.Bus, lStateMap map[string]*lua.LState, offLineCnt int) ReceiveDataTemplate {
	rxResult := ReceiveDataTemplate{
		Status: false,
	}

	txBuf := make([]byte, 0)
	result := false
	if cmd.FunName == "GetDeviceRealVariables" {
		txBuf, result, _ = node.GenerateGetRealVariables(lStateMap[node.TSL], node.Addr, 0, "")
		if result == false {
			setting.ZAPS.Errorf("%v:GetRealVariables fail", collName)
			return rxResult
		}
	} else {
		txBuf, result, _ = node.DeviceCustomCmd(lStateMap[node.TSL], node.Addr, cmd.FunName, cmd.FunPara, 0, "")
		if result == false {
			setting.ZAPS.Errorf("%v:DeviceCustomCmd fail", collName)
			return rxResult
		}
	}
	//commInterface.WriteData(txBuf)

	//阻塞读
	txBuf = append(txBuf, []byte(node.Addr)...)
	rxBuf, err := commInterface.ReadData(txBuf)
	if err != nil {
		return rxResult
	}
	rxBufCnt := len(rxBuf)
	if rxBufCnt > 0 {
		rxStatus := <-node.AnalysisRx(lStateMap[node.TSL], node.Addr, node.Properties, rxBuf[:rxBufCnt], rxBufCnt)
		rxResult = rxStatus
	}

	CommunicationUpdateNodeState(node, rxResult.Status, collName, eventBus, offLineCnt)

	return rxResult
}

func (c *CommunicationManageTemplate) CommunicationStateMachineIoOut(cmd CommunicationCmdTemplate, collName string,
	commInterface commInterface.CommunicationInterface, node *DeviceNodeTemplate,
	eventBus *eventBus.Bus, lStateMap map[string]*lua.LState, offLineCnt int) ReceiveDataTemplate {

	type IoOutTemplate struct {
		Name  string `json:"name"`
		Value []byte `json:"value"`
	}

	rxResult := ReceiveDataTemplate{
		Status: false,
	}

	rxStatus := ReceiveDataTemplate{
		Status: true,
	}
	if cmd.FunName == "GetDeviceRealVariables" || cmd.FunName == "GetRealVariables" {
		txBuf, result, _ := node.GenerateGetRealVariables(lStateMap[node.TSL], node.Addr, 0, "")
		if result == false {
			setting.ZAPS.Errorf("%v:GetRealVariables fail", collName)
			return rxResult
		}

		ioOutRequest := IoOutTemplate{
			Name:  node.Addr,
			Value: txBuf,
		}

		ioOutJson, _ := json.Marshal(&ioOutRequest)
		//阻塞读
		ackData, err := commInterface.ReadData(ioOutJson)
		if err != nil {
			return rxResult
		}

		ioOutAck := IoOutTemplate{}
		err = json.Unmarshal(ackData, &ioOutAck)
		if err != nil {
			return rxResult
		}

		rxStatus = <-node.AnalysisRx(lStateMap[node.TSL], node.Addr, node.Properties, ioOutAck.Value, len(ioOutAck.Value))
	} else {
		txBuf, result, _ := node.DeviceCustomCmd(lStateMap[node.TSL], node.Addr, cmd.FunName, cmd.FunPara, 0, "")
		if result == false {
			setting.ZAPS.Errorf("%v:DeviceCustomCmd fail", collName)
			return rxResult
		}

		ioOutRequest := IoOutTemplate{
			Name:  node.Addr,
			Value: txBuf,
		}

		ioOutJson, _ := json.Marshal(&ioOutRequest)
		//阻塞写
		_, err := commInterface.WriteData(ioOutJson)
		if err != nil {
			return rxResult
		}

	}

	CommunicationUpdateNodeState(node, rxStatus.Status, collName, eventBus, offLineCnt)

	rxResult = rxStatus

	return rxResult
}

func (c *CommunicationManageTemplate) CommunicationStateMachineHTTP(cmd CommunicationCmdTemplate, collName string,
	comm commInterface.CommunicationInterface, node *DeviceNodeTemplate,
	eventBus *eventBus.Bus, lStateMap map[string]*lua.LState, offLineCnt int) ReceiveDataTemplate {

	txBuf := make([]byte, 0)

	rxResult := ReceiveDataTemplate{
		Status: false,
	}
	rxStatus := false
	if comm.GetType() == commInterface.CommTypeHTTPSmartNode {

		if cmd.FunName == "GetDeviceRealVariables" || cmd.FunName == "GetRealVariables" {

			queryVariantData := struct {
				DeviceCode string   `json:"deviceCode"`
				Variants   []string `json:"variants"`
			}{}

			type QueryAckVariantTemplate struct {
				Code  string      `json:"code"`
				Time  int64       `json:"time"`
				Value interface{} `json:"value"`
			}

			queryAck := struct {
				ErrCode      int                       `json:"errCode"`
				ErrMsg       string                    `json:"errMsg"`
				VariantDatas []QueryAckVariantTemplate `json:"variantDatas"`
			}{}

			queryVariantData.DeviceCode = node.Addr
			for _, v := range node.Properties {
				queryVariantData.Variants = append(queryVariantData.Variants, v.Name)
			}

			buf, err := json.Marshal(queryVariantData)
			if err != nil {
				setting.ZAPS.Errorf("采集服务[%s]GetRealVariables函数调用失败", collName)
				return rxResult
			}
			txBuf = append(txBuf, buf...)
			rxBuf, err := comm.ReadData(txBuf)
			if err != nil {
				return rxResult
			}
			rxBufCnt := len(rxBuf)
			if rxBufCnt == 0 {
				return rxResult
			}
			setting.ZAPS.Debugf("采集服务[%s]接收数据 %v", collName, string(rxBuf))

			err = json.Unmarshal(rxBuf, &queryAck)
			if err != nil {
				setting.ZAPS.Errorf("采集服务[%s]接收数据JSON格式化错误 %v", collName)
				rxStatus = false
			} else {
				if queryAck.ErrCode == 0 {
					rxStatus = true
					value := TSLPropertyValueTemplate{}
					for _, v := range queryAck.VariantDatas {
						for k, p := range node.Properties {
							if v.Code == p.Name {
								switch v.Value.(type) {
								case string:
									switch p.Type {
									case PropertyTypeInt32:
										value.Value, _ = strconv.ParseUint(v.Value.(string), 10, 32)
									case PropertyTypeUInt32:
										value.Value, _ = strconv.ParseInt(v.Value.(string), 10, 32)
									case PropertyTypeDouble:
										value.Value, _ = strconv.ParseFloat(v.Value.(string), 64)
									case PropertyTypeString:
										value.Value = v.Value.(string)
									}
								case int32:
									switch p.Type {
									case PropertyTypeInt32:
										value.Value, _ = v.Value.(int32)
									case PropertyTypeUInt32:
										value.Value, _ = v.Value.(uint32)
									case PropertyTypeDouble:
										value.Value, _ = v.Value.(float64)
									case PropertyTypeString:
										value.Value = strconv.Itoa(v.Value.(int))
									}
								}

								value.Explain = ""
								value.TimeStamp = time.Unix(v.Time, 0)

								if len(node.Properties[k].Value) < VariableMaxCnt {
									node.Properties[k].Value = append(node.Properties[k].Value, value)
								} else {
									node.Properties[k].Value = node.Properties[k].Value[1:]
									node.Properties[k].Value = append(node.Properties[k].Value, value)
								}
							}
						}
					}
				}
			}
			if rxStatus == true {
				//setting.ZAPS.Debugf("采集服务[%s]节点[%s]接收成功 接收数据[%d:%v]", collName, node.Name, rxBufCnt, string(rxBuf[:rxBufCnt]))
				c.CommunicationManageMessageAdd(collName, MessageDirection_RX, rxBuf[:rxBufCnt])
			}
		} else if cmd.FunName == "SetVariables" {
			wVariables := make(map[string]interface{})
			setting.ZAPS.Debugf("采集服务[%s]SetVariables函数调用参数%v", collName, cmd.FunPara)
			err := json.Unmarshal([]byte(cmd.FunPara), &wVariables)
			if err != nil {
				setting.ZAPS.Errorf("采集服务[%s]SetVariables函数调用失败", collName)
				return rxResult
			}
			model, ok := TSLLuaMap[node.TSL]
			if !ok {
				setting.ZAPS.Errorf("采集服务[%s]采集模型不存在", collName)
				return rxResult
			}

			queryVariantData := struct {
				DeviceCode string                 `json:"deviceCode"`
				WriteData  map[string]interface{} `json:"writeData"`
			}{
				DeviceCode: node.Addr,
				WriteData:  make(map[string]interface{}),
			}

			for _, v := range model.Properties {
				variable, ok := wVariables[v.Name]
				if !ok {
					continue
				}
				queryVariantData.WriteData[v.Name] = variable
			}
			buf, err := json.Marshal(queryVariantData)
			if err != nil {
				setting.ZAPS.Errorf("采集服务[%s]SetVariables函数调用失败", collName)
				return rxResult
			}
			txBuf := make([]byte, 0)
			txBuf = append(txBuf, buf...)
			rxBuf, err := comm.WriteData(txBuf)
			if err != nil {
				return rxResult
			}
			rxBufCnt := len(rxBuf)
			if rxBufCnt == 0 {
				return rxResult
			}

			type QueryAckVariantTemplate struct {
				Code  string      `json:"code"`
				Time  int64       `json:"time"`
				Value interface{} `json:"value"`
			}

			queryAck := struct {
				ErrCode      int                       `json:"errCode"`
				ErrMsg       string                    `json:"errMsg"`
				WriteResult  map[string]interface{}    `json:"writeResult"`
				VariantDatas []QueryAckVariantTemplate `json:"variantDatas"`
			}{
				WriteResult: make(map[string]interface{}),
			}

			err = json.Unmarshal(rxBuf, &queryAck)
			if err != nil {
				setting.ZAPS.Errorf("采集服务[%s]接收数据JSON格式化错误 %v", collName)
				rxStatus = false
			} else {
				if queryAck.ErrCode == 0 {
					rxStatus = true
					//value := TSLPropertyValueTemplate{}
					//for _, v := range queryAck.VariantDatas {
					//	for k, p := range node.Properties {
					//		if v.Code == p.Name {
					//			switch v.Value.(type) {
					//			case string:
					//				switch p.Type {
					//				case PropertyTypeInt32:
					//					value.Value, _ = strconv.ParseUint(v.Value.(string), 10, 32)
					//				case PropertyTypeUInt32:
					//					value.Value, _ = strconv.ParseInt(v.Value.(string), 10, 32)
					//				case PropertyTypeDouble:
					//					value.Value, _ = strconv.ParseFloat(v.Value.(string), 64)
					//				case PropertyTypeString:
					//					value.Value = v.Value.(string)
					//				}
					//			case int32:
					//				switch p.Type {
					//				case PropertyTypeInt32:
					//					value.Value, _ = v.Value.(int32)
					//				case PropertyTypeUInt32:
					//					value.Value, _ = v.Value.(uint32)
					//				case PropertyTypeDouble:
					//					value.Value, _ = v.Value.(float64)
					//				case PropertyTypeString:
					//					value.Value = strconv.Itoa(v.Value.(int))
					//				}
					//			}
					//
					//			value.Explain = ""
					//			value.TimeStamp = time.Unix(v.Time, 0)
					//
					//			if len(node.Properties[k].Value) < VariableMaxCnt {
					//				node.Properties[k].Value = append(node.Properties[k].Value, value)
					//			} else {
					//				node.Properties[k].Value = node.Properties[k].Value[1:]
					//				node.Properties[k].Value = append(node.Properties[k].Value, value)
					//			}
					//		}
					//	}
					//}
				}
			}
			if rxStatus == true {
				//setting.ZAPS.Debugf("采集服务[%s]节点[%s]接收成功 接收数据[%d:%v]", collName, node.Name, rxBufCnt, string(rxBuf[:rxBufCnt]))
				c.CommunicationManageMessageAdd(collName, MessageDirection_RX, rxBuf[:rxBufCnt])
			}
		}

		rxResult.Status = rxStatus
		//rxResult.Properties = rxBuf
	} else if comm.GetType() == commInterface.CommTypeHTTPTianGang {
		queryVariantData := struct {
			DeviceCode string `json:"deviceCode"`
		}{}

		type QueryAckVariantTemplate struct {
			CustomerCode     string `json:"CUSTOMERCODE"`       //表号
			SysReadTime      string `json:"SYS_READ_TIME"`      //读取时间
			InFlow           string `json:"IN_FLOW"`            //累计流量
			OutFlow          string `json:"OUT_FLOW"`           //反向累计流量
			FlowSpeed        string `json:"FLOW_SPEED"`         //瞬时流量
			RunningStateName string `json:"RUNNING_STATE_NAME"` //状态
			Pressure1        string `json:"PRESSURE1"`          //压力
			MeterSubType     string `json:"METER_SUB_TYPE"`     //仪表类型
		}

		queryAck := struct {
			Result    string                    `json:"result"`
			ErrorCode int                       `json:"errorcode"`
			Msg       string                    `json:"msg"`
			Data      []QueryAckVariantTemplate `json:"data"`
		}{}

		if (cmd.FunName == "GetDeviceRealVariables") || (cmd.FunName == "GetRealVariables") {
			queryVariantData.DeviceCode = node.Addr
			buf, err := json.Marshal(queryVariantData)
			if err != nil {
				setting.ZAPS.Errorf("采集服务[%s]GetRealVariables函数调用失败", collName)
				return rxResult
			}
			txBuf = append(txBuf, buf...)
		}
		c.CommunicationManageMessageAdd(collName, MessageDirection_TX, txBuf)
		rxBuf, err := comm.ReadData(txBuf)
		if err != nil {
			return rxResult
		}
		rxBufCnt := len(rxBuf)
		if rxBufCnt == 0 {
			return rxResult
		}
		setting.ZAPS.Debugf("采集服务[%s]接收数据 %v", collName, string(rxBuf))

		err = json.Unmarshal(rxBuf, &queryAck)
		if err != nil {
			setting.ZAPS.Errorf("采集服务[%s]接收数据JSON格式化错误 %v", collName, err)
			rxStatus = false
		} else {
			if queryAck.ErrorCode == 0 {
				rxStatus = true
				value := TSLPropertyValueTemplate{}
				for _, v := range queryAck.Data {
					for k, p := range node.Properties {
						if p.Name == "PostiveFlux" {
							switch p.Type {
							case PropertyTypeInt32:
								value.Value, _ = strconv.ParseUint(v.InFlow, 10, 32)
							case PropertyTypeUInt32:
								value.Value, _ = strconv.ParseInt(v.InFlow, 10, 32)
							case PropertyTypeDouble:
								value.Value, _ = strconv.ParseFloat(v.InFlow, 64)
							case PropertyTypeString:
								value.Value = v.InFlow
							}
						}
						value.Index = len(node.Properties[k].Value)
						value.TimeStamp, err = time.ParseInLocation("2006-01-02 15:04:05", v.SysReadTime, time.Local)
						if err != nil {
							setting.ZAPS.Debugf("time解析错误 %v", err)
						}

						if len(node.Properties[k].Value) < VariableMaxCnt {
							node.Properties[k].Value = append(node.Properties[k].Value, value)
						} else {
							node.Properties[k].Value = node.Properties[k].Value[1:]
							node.Properties[k].Value = append(node.Properties[k].Value, value)
						}

						setting.ZAPS.Debugf("property %+v", node.Properties[k].Value)
					}
				}
			}
		}
		if rxStatus == true {
			setting.ZAPS.Debugf("采集服务[%s]节点[%s]接收成功 接收数据[%d:%v]", collName, node.Name, rxBufCnt, string(rxBuf[:rxBufCnt]))
			c.CommunicationManageMessageAdd(collName, MessageDirection_RX, rxBuf[:rxBufCnt])
		}
		rxResult.Status = rxStatus
		//rxResult.RxBuf = rxBuf
	}

	CommunicationUpdateNodeState(node, rxStatus, collName, eventBus, offLineCnt)

	return rxResult
}

func (c *CommunicationManageTemplate) CommunicationStateMachineS7(cmd CommunicationCmdTemplate, collName string,
	comm commInterface.CommunicationInterface, node *DeviceNodeTemplate,
	eventBus *eventBus.Bus, lStateMap map[string]*lua.LState, offLineCnt int) ReceiveDataTemplate {

	rxResult := ReceiveDataTemplate{
		Status: false,
	}
	rxStatus := false

	type QueryAckVariantTemplate struct {
		Code  string      `json:"code"`
		Time  time.Time   `json:"time"`
		Value interface{} `json:"value"`
	}

	rxBuf := make([]byte, 0)
	if cmd.FunName == "GetDeviceRealVariables" || cmd.FunName == "GetRealVariables" {
		model, ok := TSLModelS7Map[node.TSL]
		if !ok {
			setting.ZAPS.Errorf("采集服务[%s]采集模型不存在", collName)
			return rxResult
		}
		for _, v := range model.Properties {
			s7Param := commInterface.S7ParamTemplate{
				Name:    v.Name,
				Address: v.Params.DBNumber,
				Start:   v.Params.StartAddr,
				Type:    v.Params.DataType,
			}
			setting.ZAPS.Debugf("采集服务[%s]调用参数%+v", collName, s7Param)
			buf, err := json.Marshal(s7Param)
			if err != nil {
				setting.ZAPS.Errorf("采集服务[%s]调用参数JSON格式化失败", collName)
				return rxResult
			}
			txBuf := make([]byte, 0)
			txBuf = append(txBuf, buf...)
			rxBuf, err = comm.ReadData(txBuf)
			if err != nil {
				setting.ZAPS.Errorf("采集服务[%s]读数据失败", collName)
				return rxResult
			}

			rxBufCnt := len(rxBuf)
			if rxBufCnt == 0 {
				return rxResult
			}

			aData := commInterface.S7ParamTemplate{}
			err = json.Unmarshal(rxBuf, &aData)
			if err != nil {
				setting.ZAPS.Errorf("采集服务[%s]接收数据JSON格式化错误 %v", collName)
				rxStatus = false
			} else {
				rxStatus = true
				value := TSLPropertyValueTemplate{}

				for k, p := range node.Properties {
					if aData.Name == p.Name {
						//setting.ZAPS.Debugf("pName %v", p.Name)
						//setting.ZAPS.Debugf("valueType %T", aData.Value)
						switch aData.Value.(type) {
						case string:
							switch p.Type {
							case PropertyTypeInt32:
								value.Value, _ = strconv.ParseUint(aData.Value.(string), 10, 32)
							case PropertyTypeUInt32:
								value.Value, _ = strconv.ParseInt(aData.Value.(string), 10, 32)
							case PropertyTypeDouble:
								value.Value, _ = strconv.ParseFloat(aData.Value.(string), 64)
							case PropertyTypeString:
								value.Value = aData.Value.(string)
							}
						case int32:
							switch p.Type {
							case PropertyTypeInt32:
								value.Value, _ = aData.Value.(int32)
							case PropertyTypeUInt32:
								value.Value, _ = aData.Value.(uint32)
							case PropertyTypeDouble:
								value.Value, _ = aData.Value.(float64)
							case PropertyTypeString:
								value.Value = strconv.Itoa(aData.Value.(int))
							}
						case float64:
							switch p.Type {
							case PropertyTypeInt32:
								value.Value, _ = aData.Value.(float64)
							case PropertyTypeUInt32:
								value.Value, _ = aData.Value.(float64)
							case PropertyTypeDouble:
								switch v.Decimals {
								case 0:
									value.Value, _ = strconv.ParseFloat(fmt.Sprintf("%d", int64(aData.Value.(float64))), 64)
								case 1:
									value.Value, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", aData.Value.(float64)), 64)
								case 2:
									value.Value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", aData.Value.(float64)), 64)
								case 3:
									value.Value, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", aData.Value.(float64)), 64)
								case 4:
									value.Value, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", aData.Value.(float64)), 64)
								}
							case PropertyTypeString:
								value.Value = strconv.Itoa(aData.Value.(int))
							}
						}

						value.Explain = ""
						value.TimeStamp = time.Now()

						if len(node.Properties[k].Value) < VariableMaxCnt {
							node.Properties[k].Value = append(node.Properties[k].Value, value)
						} else {
							node.Properties[k].Value = node.Properties[k].Value[1:]
							node.Properties[k].Value = append(node.Properties[k].Value, value)
						}
					}
				}
			}
			if rxStatus == true {
				//setting.ZAPS.Debugf("采集服务[%s]节点[%s]接收成功 接收数据[%d:%X]", collName, node.Name, rxBufCnt, rxBuf[:rxBufCnt])
				c.CommunicationManageMessageAdd(collName, MessageDirection_RX, rxBuf[:rxBufCnt])
			}
		}
	} else if cmd.FunName == "SetVariables" {
		wVariables := make(map[string]interface{})
		setting.ZAPS.Debugf("采集服务[%s]SetVariables函数调用参数%v", collName, cmd.FunPara)
		err := json.Unmarshal([]byte(cmd.FunPara), &wVariables)
		if err != nil {
			setting.ZAPS.Errorf("采集服务[%s]SetVariables函数调用失败", collName)
			return rxResult
		}

		model, ok := TSLModelS7Map[node.TSL]
		if !ok {
			setting.ZAPS.Errorf("采集服务[%s]采集模型不存在", collName)
			return rxResult
		}
		for _, v := range model.Properties {
			variable, ok := wVariables[v.Name]
			if !ok {
				continue
			}
			s7Param := commInterface.S7ParamTemplate{
				Name:    v.Name,
				Address: v.Params.DBNumber,
				Start:   v.Params.StartAddr,
				Type:    v.Params.DataType,
				Value:   variable,
			}
			buf, err := json.Marshal(s7Param)
			if err != nil {
				setting.ZAPS.Errorf("采集服务[%s]SetVariables函数调用失败", collName)
				return rxResult
			}
			txBuf := make([]byte, 0)
			txBuf = append(txBuf, buf...)
			rxBuf, err = comm.WriteData(txBuf)
			if err != nil {
				return rxResult
			}
			rxBufCnt := len(rxBuf)
			if rxBufCnt == 0 {
				return rxResult
			}

			aData := commInterface.S7ParamTemplate{}
			err = json.Unmarshal(rxBuf, &aData)
			if err != nil {
				setting.ZAPS.Errorf("采集服务[%s]接收数据JSON格式化错误 %v", collName, err)
				rxStatus = false
			} else {

				rxStatus = true
				value := TSLPropertyValueTemplate{}

				for k, p := range node.Properties {
					if aData.Name == p.Name {
						setting.ZAPS.Debugf("pName %v", p.Name)
						setting.ZAPS.Debugf("valueType %T", aData.Value)
						switch aData.Value.(type) {
						case string:
							switch p.Type {
							case PropertyTypeInt32:
								value.Value, _ = strconv.ParseUint(aData.Value.(string), 10, 32)
							case PropertyTypeUInt32:
								value.Value, _ = strconv.ParseInt(aData.Value.(string), 10, 32)
							case PropertyTypeDouble:
								value.Value, _ = strconv.ParseFloat(aData.Value.(string), 64)
							case PropertyTypeString:
								value.Value = aData.Value.(string)
							}
						case int32:
							switch p.Type {
							case PropertyTypeInt32:
								value.Value, _ = aData.Value.(int32)
							case PropertyTypeUInt32:
								value.Value, _ = aData.Value.(uint32)
							case PropertyTypeDouble:
								value.Value, _ = aData.Value.(float64)
							case PropertyTypeString:
								value.Value = strconv.Itoa(aData.Value.(int))
							}
						case float64:
							switch p.Type {
							case PropertyTypeInt32:
								value.Value, _ = aData.Value.(float64)
							case PropertyTypeUInt32:
								value.Value, _ = aData.Value.(float64)
							case PropertyTypeDouble:
								value.Value, _ = aData.Value.(float64)
							case PropertyTypeString:
								value.Value = strconv.Itoa(aData.Value.(int))
							}
						}

						value.Explain = ""
						value.TimeStamp = time.Now()

						if len(node.Properties[k].Value) < VariableMaxCnt {
							node.Properties[k].Value = append(node.Properties[k].Value, value)
						} else {
							node.Properties[k].Value = node.Properties[k].Value[1:]
							node.Properties[k].Value = append(node.Properties[k].Value, value)
						}
					}
				}
			}
			if rxStatus == true {
				setting.ZAPS.Debugf("采集服务[%s]节点[%s]接收成功 接收数据[%d:%X]", collName, node.Name, rxBufCnt, rxBuf[:rxBufCnt])
				c.CommunicationManageMessageAdd(collName, MessageDirection_RX, rxBuf[:rxBufCnt])
			}
		}
	}

	rxResult.Status = rxStatus
	//rxResult.Properties = rxBuf

	CommunicationUpdateNodeState(node, rxStatus, collName, eventBus, offLineCnt)

	return rxResult
}

func (c *CommunicationManageTemplate) CommunicationStateMachineSAC009(cmd CommunicationCmdTemplate, collName string,
	comm commInterface.CommunicationInterface, node *DeviceNodeTemplate,
	eventBus *eventBus.Bus, lStateMap map[string]*lua.LState, offLineCnt int) ReceiveDataTemplate {

	rxResult := ReceiveDataTemplate{
		Status: false,
	}
	rxStatus := false

	type QueryAckVariantTemplate struct {
		Code  string      `json:"code"`
		Time  time.Time   `json:"time"`
		Value interface{} `json:"value"`
	}

	rxBuf := make([]byte, 0)
	if cmd.FunName == "GetDeviceRealVariables" || cmd.FunName == "GetRealVariables" {
		model, ok := TSLLuaMap[node.TSL]
		if !ok {
			setting.ZAPS.Errorf("采集服务[%s]采集模型不存在", collName)
			return rxResult
		}
		SAC009Node := commInterface.SAC009NodeTemplate{
			IMEI:       node.Addr,
			Properties: make(map[string]commInterface.SAC009PropertyTemplate),
		}
		for _, v := range model.Properties {
			property := commInterface.SAC009PropertyTemplate{
				Name: v.Name,
			}
			SAC009Node.Properties[v.Name] = property
		}
		//		setting.ZAPS.Debugf("sca0009Properties %v", SAC009Node.Properties)
		buf, err := json.Marshal(SAC009Node)
		if err != nil {
			setting.ZAPS.Errorf("采集服务[%s]GetRealVariables函数JSON格式化错误 %v", collName, err)
			return rxResult
		}
		txBuf := make([]byte, 0)
		txBuf = append(txBuf, buf...)
		rxBuf, err = comm.ReadData(txBuf)
		if err != nil {
			return rxResult
		}

		rxBufCnt := len(rxBuf)
		if rxBufCnt == 0 {
			return rxResult
		}

		ackProperties := make(map[string]commInterface.SAC009PropertyTemplate)
		err = json.Unmarshal(rxBuf, &ackProperties)
		if err != nil {
			setting.ZAPS.Errorf("采集服务[%s]接收数据JSON格式化错误 %v", collName)
			rxStatus = false
		} else {
			rxStatus = true
			value := TSLPropertyValueTemplate{}

			for k, p := range node.Properties {
				property, ok := ackProperties[p.Name]
				if !ok {
					continue
				}
				//setting.ZAPS.Debugf("pName %v", p.Name)
				//setting.ZAPS.Debugf("valueType %T", property.Value)
				switch property.Value.(type) {
				case string:
					switch p.Type {
					case PropertyTypeInt32:
						value.Value, _ = strconv.ParseUint(property.Value.(string), 10, 32)
					case PropertyTypeUInt32:
						value.Value, _ = strconv.ParseInt(property.Value.(string), 10, 32)
					case PropertyTypeDouble:
						value.Value, _ = strconv.ParseFloat(property.Value.(string), 64)
					case PropertyTypeString:
						value.Value = property.Value.(string)
					}
				case int32:
					switch p.Type {
					case PropertyTypeInt32:
						value.Value, _ = property.Value.(int32)
					case PropertyTypeUInt32:
						value.Value, _ = property.Value.(uint32)
					case PropertyTypeDouble:
						value.Value, _ = property.Value.(float64)
					case PropertyTypeString:
						value.Value = strconv.Itoa(property.Value.(int))
					}
				case float64:
					switch p.Type {
					case PropertyTypeInt32:
						value.Value, _ = property.Value.(float64)
					case PropertyTypeUInt32:
						value.Value, _ = property.Value.(float64)
					case PropertyTypeDouble:
						decimals := 0
						for _, v := range model.Properties {
							if v.Name == p.Name {
								decimals = v.Decimals
							}
						}
						if decimals > 0 {
							value.Value = property.Value.(float64) / math.Pow(10, float64(decimals))
						} else {
							value.Value = property.Value.(float64)
						}
					case PropertyTypeString:
						value.Value = strconv.Itoa(property.Value.(int))
					}
				}
				value.Explain = ""
				value.TimeStamp = property.TimeStamp

				if len(node.Properties[k].Value) < VariableMaxCnt {
					node.Properties[k].Value = append(node.Properties[k].Value, value)
				} else {
					node.Properties[k].Value = node.Properties[k].Value[1:]
					node.Properties[k].Value = append(node.Properties[k].Value, value)
				}
			}
		}
		if rxStatus == true {
			setting.ZAPS.Debugf("采集服务[%s]节点[%s]接收成功 接收数据[%d]", collName, node.Name, rxBufCnt)
			c.CommunicationManageMessageAdd(collName, MessageDirection_RX, rxBuf[:rxBufCnt])
		}
	} else if cmd.FunName == "SetVariables" {
		wVariables := make(map[string]interface{})
		setting.ZAPS.Debugf("采集服务[%s]SetVariables函数调用参数%v", collName, cmd.FunPara)
		err := json.Unmarshal([]byte(cmd.FunPara), &wVariables)
		if err != nil {
			setting.ZAPS.Errorf("采集服务[%s]SetVariables函数调用失败", collName)
			return rxResult
		}
		model, ok := TSLLuaMap[node.TSL]
		if !ok {
			setting.ZAPS.Errorf("采集服务[%s]采集模型不存在", collName)
			return rxResult
		}
		SAC009Node := commInterface.SAC009NodeTemplate{
			IMEI:       node.Addr,
			Properties: make(map[string]commInterface.SAC009PropertyTemplate),
		}
		for _, v := range model.Properties {
			variable, ok := wVariables[v.Name]
			if !ok {
				continue
			}
			property := commInterface.SAC009PropertyTemplate{
				Name:  v.Name,
				Value: variable,
			}
			SAC009Node.Properties[v.Name] = property
		}
		buf, err := json.Marshal(SAC009Node)
		if err != nil {
			setting.ZAPS.Errorf("采集服务[%s]SetVariables函数调用失败", collName)
			return rxResult
		}
		txBuf := make([]byte, 0)
		txBuf = append(txBuf, buf...)
		rxBuf, err = comm.WriteData(txBuf)
		if err != nil {
			return rxResult
		}
		rxStatus = true
	}

	rxResult.Status = rxStatus
	CommunicationUpdateNodeState(node, rxStatus, collName, eventBus, offLineCnt)

	return rxResult
}

func (c *CommunicationManageTemplate) CommunicationStateMachineDR504(cmd CommunicationCmdTemplate,
	collName string,
	commInterface commInterface.CommunicationInterface,
	node *DeviceNodeTemplate,
	eventBus *eventBus.Bus,
	lStateMap map[string]*lua.LState,
	offLineCnt int) ReceiveDataTemplate {
	rxResult := ReceiveDataTemplate{
		Status: false,
	}

	commState := CommunicationState_Start

	commStep := 0
	txBuf := make([]byte, 0)
	continues := false
	startT := time.Now() //计算当前时间
	addr := ""
	if len(node.Addr) > 20 {
		addr = node.Addr[20:]
	}
	func() {
		for {
			switch commState {
			case CommunicationState_Start:
				{
					commState = CommunicationState_Generate
				}
			case CommunicationState_Generate:
				{
					//--------------组包---------------------------
					result := false

					nodeVariables := make(map[string]interface{})
					for _, v := range node.Properties {
						if len(v.Value) > 0 {
							nodeVariables[v.Name] = v.Value[len(v.Value)-1]
						}
					}
					variablesJson, _ := json.Marshal(nodeVariables)
					if cmd.FunName == "GetDeviceRealVariables" {
						txBuf, result, continues = node.GenerateGetRealVariables(lStateMap[node.TSL], addr, commStep, string(variablesJson))
						if result == false {
							commState = CommunicationState_Stop
						} else {
							commState = CommunicationState_Send
							commStep++
						}
					} else {
						txBuf, result, continues = node.DeviceCustomCmd(lStateMap[node.TSL], addr,
							cmd.FunName,
							cmd.FunPara,
							commStep,
							string(variablesJson))
						if result == false {
							commState = CommunicationState_Stop
						} else {
							commState = CommunicationState_Send
							commStep++
						}
					}
				}
			case CommunicationState_Send:
				{
					//---------------发送前清空接收缓存-------------------------
					for i := 0; i < len(c.PacketChan); i++ {
						<-c.PacketChan
					}
					//---------------发送-------------------------
					cmd := struct {
						SN   string `json:"SN"`
						Data []byte `json:"Data"`
					}{
						Data: txBuf,
					}
					if len(node.Addr) > 20 {
						cmd.SN = node.Addr[:20]
					}
					cmdData, _ := json.Marshal(&cmd)
					commInterface.WriteData(cmdData)
					node.CommTotalCnt++
					setting.ZAPS.Infof("采集接口[%s]设备[%s]发送数据[%d:%X]", collName, node.Name, len(txBuf), txBuf)
					c.CommunicationManageMessageAdd(collName, MessageDirection_TX, txBuf)
					commState = CommunicationState_Wait
				}
			case CommunicationState_Wait:
				{
					//阻塞读
					rxBuf := make([]byte, 256)
					rxTotalBuf := make([]byte, 0)
					rxBufCnt := 0
					rxTotalBufCnt := 0
					var timeout int
					timeout, _ = strconv.Atoi(commInterface.GetTimeOut())
					timerOut := time.NewTimer(time.Duration(timeout) * time.Millisecond)
					func() {
						for {
							select {
							//继续接收数据
							case rxBuf = <-c.PacketChan:
								{
									rxBufCnt = len(rxBuf)
									if rxBufCnt > 0 {
										rxTotalBufCnt += rxBufCnt
										//追加接收的数据到接收缓冲区
										rxTotalBuf = append(rxTotalBuf, rxBuf[:rxBufCnt]...)
										//清除本次接收数据
										rxBufCnt = 0
										rxBuf = rxBuf[0:0]
									}
								}
							//是否接收超时
							case <-timerOut.C:
								{
									timerOut.Stop()
									if len(rxTotalBuf) > 0 {
										c.CommunicationManageMessageAdd(collName, MessageDirection_RX, rxTotalBuf)
									}
									setting.ZAPS.Infof("采集接口[%s]设备[%s]接收超时 接收数据[%d:%X]", collName, node.Name, len(rxTotalBuf), rxTotalBuf)

									node.CurCommFailCnt++
									if node.CurCommFailCnt >= offLineCnt {
										node.CurCommFailCnt = 0
										//设备从上线变成离线
										if node.CommStatus == "onLine" {
											content := CollectInterfaceEventTemplate{
												Topic:    "offLine",
												CollName: collName,
												NodeName: node.Name,
												Content:  node.Name,
											}
											err := eventBus.Publish("offLine", content)
											if err != nil {
												setting.ZAPS.Debugf("采集接口[%s]发布节点[%s]离线消息", collName, node.Name)
											}
										}
										node.CommStatus = "offLine"
									}
									rxTotalBufCnt = 0
									rxTotalBuf = rxTotalBuf[0:0]

									commState = CommunicationState_WaitFail
									return
								}
							//是否正确收到数据包
							case rxStatus := <-node.AnalysisRx(lStateMap[node.TSL], addr,
								node.Properties, rxTotalBuf, rxTotalBufCnt):
								{
									timerOut.Stop()
									setting.ZAPS.Infof("采集服务[%s]设备[%s]接收成功 接收数据[%d:%X]", collName, node.Name, len(rxTotalBuf), rxTotalBuf)

									rxResult.Status = rxStatus.Status
									rxResult.Properties = rxStatus.Properties

									c.CommunicationManageMessageAdd(collName, MessageDirection_RX, rxTotalBuf)

									//设备从离线变成上线
									if node.CommStatus == "offLine" {
										content := CollectInterfaceEventTemplate{
											Topic:    "onLine",
											CollName: collName,
											NodeName: node.Name,
											Content:  node.Name,
										}
										err := eventBus.Publish("onLine", content)
										if err != nil {
											setting.ZAPS.Debugf("采集接口[%s]发布节点[%s]上线消息", collName, node.Name)
										}
									}

									if continues == false {
										content := CollectInterfaceEventTemplate{
											Topic:    "update",
											CollName: collName,
											NodeName: node.Name,
											Content:  node.Name,
										}
										err := eventBus.Publish("update", content)
										if err != nil {
											setting.ZAPS.Debugf("采集接口[%s]发布节点[%s]属性更新消息", collName, node.Name)
										}
									}

									node.CommSuccessCnt++
									node.CurCommFailCnt = 0
									node.CommStatus = "onLine"
									node.LastCommRTC = time.Now().Format("2006-01-02 15:04:05")

									rxTotalBufCnt = 0
									rxTotalBuf = rxTotalBuf[0:0]
									commState = CommunicationState_WaitSucess
									return
								}
							}
						}
					}()
				}
			case CommunicationState_WaitSucess:
				{
					//通信帧延时
					var interval int
					interval, _ = strconv.Atoi(commInterface.GetInterval())
					time.Sleep(time.Duration(interval) * time.Millisecond)
					commState = CommunicationState_Stop
				}
			case CommunicationState_WaitFail:
				{
					commState = CommunicationState_Stop
				}
			case CommunicationState_Stop:
				{
					tc := time.Since(startT) //计算耗时
					setting.ZAPS.Debugf("采集服务[%s]本次采集用时%s", collName, tc)
					if continues == true {
						commState = CommunicationState_Start
					} else {
						return
					}
				}
			}
		}
	}()

	return rxResult
}

func (c *CommunicationManageTemplate) CommunicationStateMachineModbusTCP(cmd CommunicationCmdTemplate, collName string,
	comm commInterface.CommunicationInterface, node *DeviceNodeTemplate,
	eventBus *eventBus.Bus, lStateMap map[string]*lua.LState, offLineCnt int) ReceiveDataTemplate {

	rxResult := ReceiveDataTemplate{
		Status: false,
	}
	rxStatus := false

	type QueryAckVariantTemplate struct {
		Code  string      `json:"code"`
		Time  time.Time   `json:"time"`
		Value interface{} `json:"value"`
	}

	//now := time.Now()
	//var next time.Time
	//var sub time.Duration
	//count := 0

	rxBuf := make([]byte, 0)
	if cmd.FunName == "GetDeviceRealVariables" || cmd.FunName == "GetRealVariables" {
		model, ok := TSLModbusMap[node.TSL]
		if !ok {
			setting.ZAPS.Errorf("采集服务[%s]采集模型不存在", collName)
			return rxResult
		}

		for _, v := range model.Cmd {
			//now = time.Now()
			//setting.ZAPS.Error("1 ----->", now.Format("2006-01-02 15:04:05"))
			addr, _ := strconv.Atoi(node.Addr)
			mbTCPParam := commInterface.MBTCPInterfaceParamTemplate{
				SlaveID:      byte(addr),
				FuncCode:     byte(v.FunCode),
				StartRegAddr: v.StartRegAddr,
				TotalRegCnt:  v.RegCnt,
			}
			param := commInterface.MBTCPParamTemplate{}
			for _, r := range v.Registers {
				param.Rule = r.RuleType
				param.RegAddr = r.RegAddr
				param.Type = r.Type
				param.RegCnt = r.RegCnt
				mbTCPParam.Param = append(mbTCPParam.Param, param)
			}
			buf, err := json.Marshal(mbTCPParam)
			if err != nil {
				setting.ZAPS.Errorf("采集服务[%s]调用参数JSON格式化失败", collName)
				continue
			}
			txBuf := make([]byte, 0)
			txBuf = append(txBuf, buf...)

			rxBuf, err = comm.ReadData(txBuf)
			if err != nil {
				setting.ZAPS.Errorf("采集服务[%s]读数据失败 %v", collName, err)
				continue
			}

			rxBufCnt := len(rxBuf)
			if rxBufCnt == 0 {
				continue
			}

			//next = time.Now()
			//setting.ZAPS.Error("2 ----->", next.Format("2006-01-02 15:04:05"), "  ", next.Sub(now))
			aData := make([]commInterface.MBTCPParamTemplate, 0)
			err = json.Unmarshal(rxBuf, &aData)
			if err != nil {
				setting.ZAPS.Errorf("采集服务[%s]接收数据JSON格式化错误 %v", collName, err)
				rxStatus = false
				continue
			} else {
				//next = time.Now()
				//setting.ZAPS.Error("3 ----->", next.Format("2006-01-02 15:04:05"), "  ", next.Sub(now))
				rxStatus = true
				//value := TSLPropertyValueTemplate{}
				aDtaLen := len(aData)
				aDataCount := aDtaLen / 20
				A1 := 0
				A2 := 0
				//fmt.Println("------> go Count ", aDataCount)

				var wg sync.WaitGroup
				var mutex sync.Mutex

				for i := 0; i <= aDataCount; i++ {
					if 0 == aDataCount {
						A1 = 0
						A2 = aDtaLen
					} else if i == aDataCount {
						A1 = i * 20
						A2 = aDtaLen
					} else {
						A1 = i * 20
						A2 = i*20 + 20
					}

					wg.Add(1) // 添加一个计数器

					go func(a1, a2 int) {
						defer wg.Done() // goroutine 结束时计数器减 1

						value := TSLPropertyValueTemplate{}
						for _, d := range aData[a1:a2] {
							for _, r := range v.Registers {
								if r.RegAddr == d.RegAddr {
									for k, p := range node.Properties {
										if r.Name == p.Name {
											//setting.ZAPS.Debugf("pName %v", p.Name)
											//setting.ZAPS.Debugf("type %T,value %v", d.Value, d)

											if true == r.BitOffsetSw {
												if _, ok := d.Value.(float64); ok {
													value.Value = uint32(d.Value.(float64)) & (0x01 << r.BitOffset)
												} else {
													value.Value = -1
												}

											} else {
												var fValue float64
												if r.Formula != "" {
													fStr := r.Formula
													if strings.Contains(fStr, "t") {
														switch d.Value.(type) {
														case uint32:
															fStr = strings.ReplaceAll(fStr, "t", fmt.Sprintf("%d", d.Value.(uint32)))
														case int32:
															fStr = strings.ReplaceAll(fStr, "t", fmt.Sprintf("%d", d.Value.(int32)))
														case float64:
															fStr = strings.ReplaceAll(fStr, "t", fmt.Sprintf("%f", d.Value.(float64)))
														case string:
														}
													}
													err, value := setting.FormulaRun(fStr)
													if err != nil {
														fValue = d.Value.(float64)
													} else {
														fValue = value
													}
												} else {
													switch d.Value.(type) {
													case uint32:
														fValue = float64(d.Value.(uint32))
													case int32:
														fValue = float64(d.Value.(int32))
													case float64:
														fValue = d.Value.(float64)
													case string:
													}
												}

												switch p.Type {
												case PropertyTypeInt32:
													value.Value = (int32)(fValue)
												case PropertyTypeUInt32:
													value.Value = (uint32)(fValue)
												case PropertyTypeDouble:
													//判断modbus通信接口返回的数据是否已经是float，如果是则保留精度,不是则换算
													if strings.Contains(r.RuleType, "Float") {
														value.Value, _ = strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(r.Decimals)+"f", fValue), 64)
													} else {
														if r.Decimals > 0 {
															value.Value = fValue / math.Pow10(r.Decimals)
														} else {
															value.Value = fValue
														}
													}
												case PropertyTypeString:
													value.Value = d.Value.(string)
												}
											}
											value.Explain = ""
											value.TimeStamp = time.Now()

											mutex.Lock()
											if len(node.Properties[k].Value) < VariableMaxCnt {
												node.Properties[k].Value = append(node.Properties[k].Value, value)
											} else {
												node.Properties[k].Value = node.Properties[k].Value[1:]
												node.Properties[k].Value = append(node.Properties[k].Value, value)
											}
											mutex.Unlock()
											//break
										}
									}
									break
								}
							}
						}
					}(A1, A2)
				}

				// 等待所有的 goroutine 都结束
				wg.Wait()

				//next = time.Now()
				//count++
				//sub = next.Sub(now)
				//setting.ZAPS.Error("4 ----->", count, "  ", next.Format("2006-01-02 15:04:05"), "  ", sub)
			}
			if rxStatus == true {
				//setting.ZAPS.Debugf("采集服务[%s]节点[%s]接收成功 接收数据[%d:%X]", collName, node.Name, rxBufCnt, rxBuf[:rxBufCnt])
				c.CommunicationManageMessageAdd(collName, MessageDirection_RX, rxBuf[:rxBufCnt])
			}
		}
	} else if cmd.FunName == "SetVariables" {

		wVariables := make(map[string]interface{})
		setting.ZAPS.Debugf("采集服务[%s]SetVariables函数调用参数%v", collName, cmd.FunPara)
		err := json.Unmarshal([]byte(cmd.FunPara), &wVariables)
		if err != nil {
			setting.ZAPS.Errorf("采集服务[%s]SetVariables 函数调用失败", collName)
			return rxResult
		}

		model, ok := TSLModbusMap[node.TSL]
		if !ok {
			setting.ZAPS.Errorf("采集服务[%s]采集模型不存在", collName)
			return rxResult
		}

		SlaveAddr, _ := strconv.Atoi(node.Addr) //gwai add 2023-04-21
		//properties := model.GetTSLModbusModelProperties()  //gwai del 2023-04-21
		for _, cmd := range model.Cmd { //gwai add 2023-04-21
			for _, v := range cmd.Registers {
				variable, ok := wVariables[v.Name]
				if !ok {
					continue
				}

				var FuncCode byte = 0x10 //gwai add 2023-04-21
				if cmd.FunCode == 1 {
					FuncCode = 0x0F
				}

				mbTCPParam := commInterface.MBTCPInterfaceParamTemplate{
					SlaveID:      byte(SlaveAddr),
					FuncCode:     FuncCode,
					StartRegAddr: v.RegAddr,
					TotalRegCnt:  v.RegCnt,
				}
				param := commInterface.MBTCPParamTemplate{
					Rule:    v.RuleType,
					RegAddr: v.RegAddr,
					Type:    v.Type,
					RegCnt:  v.RegCnt,
				}
				switch variable.(type) {
				case uint32:
					param.Value = (float64)(variable.(uint32))
				case int32:
					param.Value = (float64)(variable.(int32))
				case float64:
					param.Value = variable
				case string:
					param.Value, _ = strconv.ParseFloat(variable.(string), 64)
				}

				//setting.ZAPS.Debugf("type %T", variable)
				mbTCPParam.Param = append(mbTCPParam.Param, param)

				buf, err := json.Marshal(mbTCPParam)
				if err != nil {
					setting.ZAPS.Errorf("采集服务[%s]SetVariables函数调用失败", collName)
					return rxResult
				}
				txBuf := buf
				rxBuf, err = comm.WriteData(txBuf)
				if err != nil {
					return rxResult
				}
				rxBufCnt := len(rxBuf)
				if rxBufCnt == 0 {
					return rxResult
				}

				aData := make([]commInterface.MBTCPParamTemplate, 0)
				err = json.Unmarshal(rxBuf, &aData)
				if err != nil {
					setting.ZAPS.Errorf("采集服务[%s]接收数据JSON格式化错误 %v", collName, err)
					rxStatus = false
				} else {
					rxStatus = true
				}
				if rxStatus == true {
					//setting.ZAPS.Debugf("采集服务[%s]节点[%s]接收成功 接收数据[%d:%X]", collName, node.Name, rxBufCnt, rxBuf[:rxBufCnt])
					c.CommunicationManageMessageAdd(collName, MessageDirection_RX, rxBuf[:rxBufCnt])
				}
			}
		}
	}

	rxResult.Status = rxStatus
	//rxResult.Properties = rxBuf

	CommunicationUpdateNodeState(node, rxStatus, collName, eventBus, offLineCnt)

	return rxResult
}

func (c *CommunicationManageTemplate) CommunicationStateMachineHDKJ(cmd CommunicationCmdTemplate, collName string,
	comm commInterface.CommunicationInterface, node *DeviceNodeTemplate,
	eventBus *eventBus.Bus, lStateMap map[string]*lua.LState, offLineCnt int) ReceiveDataTemplate {

	rxResult := ReceiveDataTemplate{
		Status: false,
	}
	rxStatus := false

	type QueryAckVariantTemplate struct {
		Code  string      `json:"code"`
		Time  time.Time   `json:"time"`
		Value interface{} `json:"value"`
	}

	rxBuf := make([]byte, 0)
	if cmd.FunName == "GetDeviceRealVariables" || cmd.FunName == "GetRealVariables" {
		model, ok := TSLLuaMap[node.TSL]
		if !ok {
			setting.ZAPS.Errorf("采集服务[%s]采集模型不存在", collName)
			return rxResult
		}
		HDKJNode := commInterface.HDKJNodeTemplate{
			IMEI:       node.Addr,
			Properties: make(map[string]commInterface.HDKJPropertyTemplate),
		}
		for _, v := range model.Properties {
			property := commInterface.HDKJPropertyTemplate{
				Name: v.Name,
			}
			HDKJNode.Properties[v.Name] = property
		}
		//		setting.ZAPS.Debugf("sca0009Properties %v", SAC009Node.Properties)
		buf, err := json.Marshal(HDKJNode)
		if err != nil {
			setting.ZAPS.Errorf("采集服务[%s]GetRealVariables函数JSON格式化错误 %v", collName, err)
			return rxResult
		}
		txBuf := make([]byte, 0)
		txBuf = append(txBuf, buf...)
		rxBuf, err = comm.ReadData(txBuf)
		if err != nil {
			rxStatus = false
		}

		rxBufCnt := len(rxBuf)
		if rxBufCnt == 0 {
			rxStatus = false
		}

		ackProperties := make(map[string]commInterface.HDKJPropertyTemplate)
		err = json.Unmarshal(rxBuf, &ackProperties)
		if err != nil {
			setting.ZAPS.Errorf("采集服务[%s]接收数据JSON格式化错误 %v", collName, err)
			rxStatus = false
		} else {
			rxStatus = true
			value := TSLPropertyValueTemplate{}

			for k, p := range node.Properties {
				property, ok := ackProperties[p.Name]
				if !ok {
					continue
				}
				switch property.Value.(type) {
				case string:
					switch p.Type {
					case PropertyTypeInt32:
						value.Value, _ = strconv.ParseUint(property.Value.(string), 10, 32)
					case PropertyTypeUInt32:
						value.Value, _ = strconv.ParseInt(property.Value.(string), 10, 32)
					case PropertyTypeDouble:
						value.Value, _ = strconv.ParseFloat(property.Value.(string), 64)
					case PropertyTypeString:
						value.Value = property.Value.(string)
					}
				case int32:
					switch p.Type {
					case PropertyTypeInt32:
						value.Value, _ = property.Value.(int32)
					case PropertyTypeUInt32:
						value.Value, _ = property.Value.(uint32)
					case PropertyTypeDouble:
						value.Value, _ = property.Value.(float64)
					case PropertyTypeString:
						value.Value = strconv.Itoa(property.Value.(int))
					}
				case float64:
					switch p.Type {
					case PropertyTypeInt32:
						value.Value, _ = property.Value.(float64)
					case PropertyTypeUInt32:
						value.Value, _ = property.Value.(float64)
					case PropertyTypeDouble:
						decimals := 0
						for _, v := range model.Properties {
							if v.Name == p.Name {
								decimals = v.Decimals
							}
						}
						if decimals > 0 {
							value.Value = property.Value.(float64) / math.Pow(10, float64(decimals))
						} else {
							value.Value = property.Value.(float64)
						}
					case PropertyTypeString:
						value.Value = strconv.Itoa(property.Value.(int))
					}
				}
				value.Explain = property.Explain
				value.TimeStamp = property.TimeStamp

				if len(node.Properties[k].Value) < VariableMaxCnt {
					node.Properties[k].Value = append(node.Properties[k].Value, value)
				} else {
					node.Properties[k].Value = node.Properties[k].Value[1:]
					node.Properties[k].Value = append(node.Properties[k].Value, value)
				}
			}
		}
		if rxStatus == true {
			setting.ZAPS.Debugf("采集服务[%s]节点[%s]接收成功 接收数据[%d]", collName, node.Name, rxBufCnt)
			c.CommunicationManageMessageAdd(collName, MessageDirection_RX, rxBuf[:rxBufCnt])
		}
	} else if cmd.FunName == "SetVariables" {
		rxStatus = true
	}

	rxResult.Status = rxStatus
	CommunicationUpdateNodeState(node, rxStatus, collName, eventBus, offLineCnt)

	return rxResult
}

func (c *CommunicationManageTemplate) CommunicationStateMachineDLT64507(cmd CommunicationCmdTemplate, collName string,
	comm commInterface.CommunicationInterface, node *DeviceNodeTemplate,
	eventBus *eventBus.Bus, lStateMap map[string]*lua.LState, offLineCnt int) ReceiveDataTemplate {

	rxResult := ReceiveDataTemplate{
		Status: false,
	}

	rxStatus := false

	type packAndProperties struct {
		BlockRead     byte //是否按块读写 0是单个标识读取  1是按块读取
		frame         dlt645.D07PackFrame
		propertiesMap []*TSLDLT6452007PropertyTemplate
	}

	rxBuf := make([]byte, 256)
	packAndPropertiesMap := make([]packAndProperties, 0)

	if cmd.FunName == "GetDeviceRealVariables" || cmd.FunName == "GetRealVariables" {
		model, ok := TSLDLT6452007Map[node.TSL]
		if !ok {
			setting.ZAPS.Errorf("采集服务[%s]采集模型不存在", collName)
			return rxResult
		}

		for _, v := range model.Cmd {
			if v.BlockRead == 0 { //如果不是块读。那么按ID一个一个读本组内的规则数据
				for _, f := range v.Properties {
					d07Frame := dlt645.D07PackFrame{
						RulerId:  f.RulerId,
						CtrlCode: 0x11,
						DataLen:  0,
						Address:  node.Addr,
						Data:     nil,
					}
					packAndPropertiesMap = append(packAndPropertiesMap, packAndProperties{v.BlockRead, d07Frame, []*TSLDLT6452007PropertyTemplate{f}})
				}
			} else { //如果是按块读取，那么直接发块读取命令
				d07Frame := dlt645.D07PackFrame{
					RulerId:  v.BlockRulerId,
					CtrlCode: 0x11,
					DataLen:  0,
					Address:  node.Addr,
					Data:     nil,
				}
				pMap := make([]*TSLDLT6452007PropertyTemplate, 0)
				for _, f := range v.Properties {
					pMap = append(pMap, f)
				}

				packAndPropertiesMap = append(packAndPropertiesMap, packAndProperties{v.BlockRead, d07Frame, pMap})
			}
		}
	}

	for _, curPackAndProperties := range packAndPropertiesMap { //读取数据
		if txBuf, ok := curPackAndProperties.frame.PackD07FrameByData(); ok == nil {
			comm.WriteData(txBuf)

			//阻塞读
			rxTotalBuf := make([]byte, 0)
			rxBufCnt := 0
			rxTotalBufCnt := 0
			var timeout int
			timeout, _ = strconv.Atoi(comm.GetTimeOut())
			timerOut := time.NewTimer(time.Duration(timeout) * time.Millisecond)

			func() {
				for {
					select {
					//继续接收数据
					case rxBuf = <-c.PacketChan:
						{
							rxBufCnt = len(rxBuf)
							if rxBufCnt > 0 {
								rxTotalBufCnt += rxBufCnt
								//追加接收的数据到接收缓冲区
								rxTotalBuf = append(rxTotalBuf, rxBuf[:rxBufCnt]...)
								//清除本次接收数据
								rxBufCnt = 0
								rxBuf = rxBuf[0:0]
							}

							if rxTotalBuf[len(rxTotalBuf)-1] == 0x16 {
								if unpackFrame, ok := dlt645.UnpackD07Frame(rxTotalBuf); ok == nil {
									timerOut.Stop()
									setting.ZAPS.Infof("采集服务[%s]设备[%s]接收成功 接收数据[%d:%X]", collName, node.Name, len(rxTotalBuf), rxTotalBuf)
									c.CommunicationManageMessageAdd(collName, MessageDirection_RX, rxTotalBuf)
									rxStatus = true

									var d07Data dlt645.TransD07DataTemplate
									for _, properties := range curPackAndProperties.propertiesMap {
										Offset := 0
										if curPackAndProperties.BlockRead != 0 {
											Offset = properties.BlockAddOffset
										}
										if properties.RulerAddOffset > 0 {
											Offset += properties.RulerAddOffset - 1
										}
										if Offset > len(unpackFrame.Data) {
											setting.ZAPS.Error("偏移地址配置错误")
											return
										}

										d07Data.TransDir = dlt645.ED07TransF2u
										d07Data.Frame = unpackFrame.Data[Offset:]
										if properties.Type == PropertyTypeString {
											d07Data.User = ""
										} else if properties.Type == PropertyTypeDouble {
											d07Data.User = 0.0
										} else {
											d07Data.User = 0.0
										}

										value := TSLPropertyValueTemplate{}
										if dd, ok := dlt645.UnpackD07ByFormat(properties.Format, &d07Data); ok == nil {
											for k, p := range node.Properties {
												if properties.Name == p.Name {

													value.Value = dd.User
													value.Explain = ""
													value.TimeStamp = time.Now()

													if len(node.Properties[k].Value) < VariableMaxCnt {
														node.Properties[k].Value = append(node.Properties[k].Value, value)
													} else {
														node.Properties[k].Value = node.Properties[k].Value[1:]
														node.Properties[k].Value = append(node.Properties[k].Value, value)
													}
												}
											}
										}
									}
									return
								}
							}
						}
					//是否接收超时
					case <-timerOut.C:
						{
							timerOut.Stop()
							if len(rxTotalBuf) > 0 {
								c.CommunicationManageMessageAdd(collName, MessageDirection_RX, rxTotalBuf)
							}
							setting.ZAPS.Infof("采集接口[%s]设备[%s]接收超时 接收数据[%d:%X]", collName, node.Name, len(rxTotalBuf), rxTotalBuf)
							rxStatus = false

							rxTotalBufCnt = 0
							rxTotalBuf = rxTotalBuf[0:0]

							return
						}
					}
				}
			}()

		}

	}

	rxResult.Status = rxStatus
	CommunicationUpdateNodeState(node, rxStatus, collName, eventBus, offLineCnt)

	return rxResult
}

func (c *CommunicationManageTemplate) CommunicationDirectDataStateMachine(req CommunicationDirectDataReqTemplate, commInterface commInterface.CommunicationInterface) CommunicationDirectDataAckTemplate {

	ack := CommunicationDirectDataAckTemplate{
		Status: false,
		Data:   make([]byte, 0),
	}

	commState := CommunicationState_DirectDataSend

	func() {
		for {
			switch commState {
			case CommunicationState_DirectDataSend:
				{
					//---------------发送-------------------------
					commInterface.WriteData(req.Data)
					commState = CommunicationState_DirectDataWait
				}
			case CommunicationState_DirectDataWait:
				{
					//阻塞读
					rxBuf := make([]byte, 256)
					rxTotalBuf := make([]byte, 0)
					rxBufCnt := 0
					rxTotalBufCnt := 0
					var timeout int
					timeout, _ = strconv.Atoi(commInterface.GetTimeOut())
					timerOut := time.NewTimer(time.Duration(timeout) * time.Millisecond)
					func() {
						for {
							select {
							//继续接收数据
							case rxBuf = <-c.PacketChan:
								{
									rxBufCnt = len(rxBuf)
									if rxBufCnt > 0 {
										rxTotalBufCnt += rxBufCnt
										//追加接收的数据到接收缓冲区
										rxTotalBuf = append(rxTotalBuf, rxBuf[:rxBufCnt]...)
										//清除本次接收数据
										rxBufCnt = 0
										rxBuf = rxBuf[0:0]
									}
								}
							//是否接收超时
							case <-timerOut.C:
								{
									timerOut.Stop()
									ack.Data = append(ack.Data, rxTotalBuf[:rxTotalBufCnt]...)
									if len(ack.Data) == 0 {
										ack.Status = false
									} else {
										ack.Status = true
									}
									rxTotalBufCnt = 0
									rxTotalBuf = rxTotalBuf[0:0]
									commState = CommunicationState_DirectDataStop
									return
								}
							}
						}
					}()
				}
			case CommunicationState_DirectDataStop:
				{
					return
				}
			}
		}
	}()

	return ack
}
