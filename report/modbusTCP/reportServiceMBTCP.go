package modbusTCP

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"gateway/device"
	"gateway/device/eventBus"
	"gateway/setting"
	"gateway/utils"
	"github.com/jasonlvhit/gocron"
	modbus "github.com/thinkgos/gomodbus/v2"
	"math"
	"time"
)

type ReportServiceMBTCPTemplate struct {
	Index           int                             `json:"index"`
	GWParam         ReportServiceMBTCPParamTemplate `json:"GWParam"`
	CancelFunc      context.CancelFunc              `json:"-"`
	MessageEventBus eventBus.Bus                    `json:"-"` //通信报文总线
}

type ReportServiceMBTCPParamTemplate struct {
	ServiceName  string `json:"serviceName"`
	IP           string `json:"ip"`
	Port         string `json:"port"`
	ReportStatus string `json:"reportStatus"`
	ReportTime   int    `json:"reportTime"`
	Protocol     string `json:"protocol"`
	Param        struct {
		SlaveID                  int                              `json:"slaveID"`
		TcpServer                *modbus.TCPServer                `json:"-"`
		CoilStatusRegisterStart  int                              `json:"coilStatusRegStart"`
		CoilStatusRegisterCnt    int                              `json:"coilStatusRegCnt"`
		CoilStatusRegisters      map[string]MBTCPRegisterTemplate `json:"coilStatusRegisters,omitempty"`
		InputStatusRegisterStart int                              `json:"inputStatusRegStart"`
		InputStatusRegisterCnt   int                              `json:"inputStatusRegCnt"`
		InputStatusRegisters     map[string]MBTCPRegisterTemplate `json:"inputStatusRegisters,omitempty"`
		HoldingRegisterStart     int                              `json:"holdingRegStart"`
		HoldingRegisterCnt       int                              `json:"holdingRegCnt"`
		HoldingRegisters         map[string]MBTCPRegisterTemplate `json:"holdingRegisters,omitempty"`
		InputRegisterStart       int                              `json:"inputRegStart"`
		InputRegisterCnt         int                              `json:"inputRegCnt"`
		InputRegisters           map[string]MBTCPRegisterTemplate `json:"inputRegisters,omitempty"`
	}
}

type MBTCPRegisterTemplate struct {
	Index        int    `json:"index"`
	RegName      string `json:"regName"`      //寄存器名称
	Label        string `json:"label"`        //寄存器标签
	PropertyType int    `json:"propertyType"` //属性类型（0：外部变量 1：内部变量）
	CollName     string `json:"collName"`     //采集接口名称
	NodeName     string `json:"nodeName"`     //设备名称
	PropertyName string `json:"propertyName"` //属性名称
	RegAddr      int    `json:"regAddr"`      //寄存器地址
	RegCnt       int    `json:"regCnt"`       //寄存器数量
	RegType      int    `json:"regType"`      //寄存器类型
	Rule         string `json:"rule"`         //寄存器解析规则
}

type ReportServiceMBTCPListTemplate struct {
	ServiceList []*ReportServiceMBTCPTemplate
}

const (
	MBRegTypeCoilStatus = iota
	MBRegTypeInputStatus
	MBRegTypeHoldingRegister
	MBRegTypeInputRegister
)

var ReportServiceMBTCPList ReportServiceMBTCPListTemplate
var writeTimer *time.Timer

func ReportServiceMBTCPInit() {

	writeTimer = time.AfterFunc(time.Second, func() {
		ReportServiceMBTCPWriteToJson()
	})
	writeTimer.Stop()

	err := ReportServiceMBTCPReadFromJson()
	if err != nil {
		setting.ZAPS.Infof("上报服务[MBTCP]读取配置文件失败 %v", err)
		return
	}

	for _, v := range ReportServiceMBTCPList.ServiceList {
		//创建服务
		v.NewReportServiceMBTCP()

		// 定义一个cron运行器
		scheduler := gocron.NewScheduler()

		setting.ZAPS.Infof("上报服务[%s]定时更新寄存器任务", v.GWParam.ServiceName)
		_ = scheduler.Every(1).Second().Do(v.GWParam.UpdateReadRegister)

		scheduler.Start()
	}
}

func (r *ReportServiceMBTCPTemplate) NewReportServiceMBTCP() {

	ctx, cancel := context.WithCancel(context.Background())
	r.CancelFunc = cancel

	r.GWParam.Param.TcpServer = modbus.NewTCPServer()
	mTCPNode := modbus.NewNodeRegister(byte(r.GWParam.Param.SlaveID),
		0, 0,
		0, 0,
		0, 0,
		uint16(r.GWParam.Param.HoldingRegisterStart), uint16(r.GWParam.Param.HoldingRegisterCnt))
	r.GWParam.Param.TcpServer.AddNodes(mTCPNode)

	r.GWParam.Param.TcpServer.RegisterFunctionHandler(0x06, r.GWParam.WriteSingleRegisterFunctionHandler)

	go r.ReportServiceMBTCPPoll(ctx, r.GWParam.ServiceName, r.GWParam.Port, r.GWParam.Param.TcpServer)
}

func (r *ReportServiceMBTCPTemplate) ReportServiceMBTCPPoll(ctx context.Context, name string, port string, tcpServer *modbus.TCPServer) {

	go func() {
		err := tcpServer.ListenAndServe(":" + port)
		if err != nil {
			setting.ZAPS.Errorf("上报服务[%v]modbusTCP监听失败 %v", name, err)
			r.GWParam.ReportStatus = "offLine"
			return
		}
	}()
	r.GWParam.ReportStatus = "onLine"
	for {
		select {
		case <-ctx.Done():
			{
				setting.ZAPS.Infof("上报服务[%v]modbusTCP退出监听", name)
				r.GWParam.ReportStatus = "offLine"
				tcpServer.DeleteAllNode()
				err := tcpServer.Close()
				if err != nil {
					setting.ZAPS.Info("上报服务[%v]modbusTCP退出监听失败 %v", name, err)
				}
				return
			}
		}
	}
}

func ReportServiceMBTCPReadFromJson() error {
	data, err := utils.FileRead("./selfpara/reportServiceMBTCP.json")
	if err != nil {
		setting.ZAPS.Errorf("上报服务[MBTCP]配置json文件读取失败 %v", err)
		return err
	}

	//mbTCP := ReportServiceMBTCPListTemplate{}
	err = json.Unmarshal(data, &ReportServiceMBTCPList)
	if err != nil {
		setting.ZAPS.Errorf("上报服务[MBTCP]配置json文件格式化失败 %v", err)
		return err
	}

	setting.ZAPS.Info("上报服务[MBTCP]配置json文件读取成功")

	return nil
}

func ReportServiceMBTCPWriteParam() {
	writeTimer.Reset(time.Second)
}

func ReportServiceMBTCPWriteToJson() {
	utils.DirIsExist("./selfpara")
	sJson, _ := json.Marshal(ReportServiceMBTCPList)
	err := utils.FileWrite("./selfpara/reportServiceMBTCP.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("上报服务[MBTCP]配置json文件写入失败")
		return
	}
	setting.ZAPS.Debugf("上报服务[MBTCP]配置json文件写入成功")
}

func (r *ReportServiceMBTCPParamTemplate) UpdateReadRegister() {

	slaveNode, err := r.Param.TcpServer.GetNode(byte(r.Param.SlaveID))
	if err != nil {
		setting.ZAPS.Errorf("上报服务[%v]modbusTCP查找SlaveID错误 %v", r.ServiceName, err)
		return
	}

	for _, v := range r.Param.HoldingRegisters {
		coll, ok := device.CollectInterfaceMap.Coll[v.CollName]
		if !ok {
			continue
		}
		node, ok := coll.DeviceNodeMap[v.NodeName]
		if !ok {
			continue
		}
		var value float64
		for _, p := range node.Properties {
			if p.Name == v.PropertyName {
				if len(p.Value) == 0 {
					continue
				}
				switch p.Value[len(p.Value)-1].Value.(type) {
				case uint32:
					value = float64(p.Value[len(p.Value)-1].Value.(uint32))
				case int32:
					value = float64(p.Value[len(p.Value)-1].Value.(int32))
				case float64:
					value = p.Value[len(p.Value)-1].Value.(float64)
				}

				buf := make([]byte, 0)
				if v.RegCnt == 1 {
					switch v.Rule {
					case "Int_AB":
						fallthrough
					case "Int_BA":
					}
					regBytes := make([]byte, 2)
					binary.LittleEndian.PutUint16(regBytes, uint16(value))
					buf = append(buf, regBytes...)
				} else if v.RegCnt == 2 {
					switch v.Rule {
					case "Long_ABCD":
						fallthrough
					case "Long_BADC":
						fallthrough
					case "Long_DCBA":
						fallthrough
					case "Long_CDAB":
						regBytes := make([]byte, 4)
						binary.BigEndian.PutUint32(regBytes, uint32(value))
						buf = append(buf, regBytes...)
					case "Float_ABCD":
						fallthrough
					case "Float_BADC":
						fallthrough
					case "Float_DCBA":
						fallthrough
					case "Float_CDAB":
						regBytes := make([]byte, 4)
						bits := math.Float32bits(float32(value))
						binary.BigEndian.PutUint32(regBytes, bits)
						buf = append(buf, regBytes...)
					}
				} else if v.RegCnt == 4 {
					switch v.Rule {
					case "Double_ABCDEFGH":
						regBytes := make([]byte, 8)
						bits := math.Float64bits(value)
						binary.BigEndian.PutUint64(regBytes, bits)
						buf = append(buf, regBytes...)
					}
				}
				err = slaveNode.WriteHoldingsBytes(uint16(v.RegAddr), uint16(len(buf)/2), buf)
			}
		}
	}
}

func (r *ReportServiceMBTCPParamTemplate) WriteSingleRegisterFunctionHandler(reg *modbus.NodeRegister, data []byte) ([]byte, error) {

	if len(data) != modbus.FuncWriteMinSize {
		return nil, &modbus.ExceptionError{modbus.ExceptionCodeIllegalDataValue}
	}

	address := binary.BigEndian.Uint16(data)
	quality := 1
	valBuf := make([]uint16, 0)
	val := binary.BigEndian.Uint16(data[2:])

	valBuf = append(valBuf, val)
	for _, v := range r.Param.HoldingRegisters {
		if uint16(v.RegAddr) != address {
			continue
		}

		coll, ok := device.CollectInterfaceMap.Coll[v.CollName]
		if !ok {
			continue
		}
		node, ok := coll.DeviceNodeMap[v.NodeName]
		if !ok {
			continue
		}
		for _, p := range node.Properties {
			if p.Name == v.PropertyName {
				err := reg.WriteHoldings(address, valBuf)
				if err != nil {
					setting.ZAPS.Errorf("WriteHoldingsBytes错误 %v", err)
				}
				setting.ZAPS.Debugf("address:%d quality:%d valBuf:%x", address, quality, valBuf)

				//从采集服务中找到相应节点
				cmd := device.CommunicationCmdTemplate{}
				cmd.CollInterfaceName = coll.CollInterfaceName
				cmd.DeviceName = node.Name
				cmd.FunName = "SetVariables"
				valueMap := make(map[string]interface{})
				valueMap[v.PropertyName] = val
				paramStr, _ := json.Marshal(valueMap)
				cmd.FunPara = string(paramStr)

				ackData := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
				if ackData.Status {
					return data, nil
				} else {
					return nil, nil
				}
			}
		}
	}

	return nil, nil
}
