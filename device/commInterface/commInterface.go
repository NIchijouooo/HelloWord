package commInterface

import (
	"encoding/json"
	"gateway/setting"
	"gateway/utils"
)

type CommunicationInterface interface {
	Open() bool
	Close() bool
	WriteData(data []byte) ([]byte, error)
	ReadData(data []byte) ([]byte, error)
	GetName() string
	GetTimeOut() string
	GetInterval() string
	GetType() int
}

type ProtocolTemplate struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

type ConnectStatus int

const (
	CommIsUnConnect ConnectStatus = iota
	CommIsConnect
)

const (
	CommTypeIoIn int = iota
	CommTypeIoOut
	CommTypeSerial
	CommTypeTcpClient
	CommTypeTcpServer
	CommTypeHTTPSmartNode
	CommTypeHTTPTianGang
	CommTypeS7
	CommTypeSAC009
	CommTypeSTL
	CommTypeDR504
	CommTypeModbusTCP
	CommTypeHDKJ
	CommTypeModbusRTU
)

var CommIntefaceProtocols = make([]ProtocolTemplate, 0)

//通信接口Map
var CommunicationInterfaceMap = make([]CommunicationInterface, 0)

func CommInterfaceInit() {

	ReadCommInterfaceProtocol()

	//获取串口通信接口参数
	if ReadCommSerialInterfaceListFromJson() == true {
		for _, v := range CommunicationSerialMap {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	//获取TcpClient通信接口参数
	if ReadCommTcpClientInterfaceListFromJson() == true {
		for _, v := range CommunicationTcpClientMap {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	//获取TcpServer通信接口参数
	if ReadCommTcpServerInterfaceListFromJson() == true {
		for _, v := range CommunicationTcpServerMap {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	//获取开关量输出通信接口参数
	if ReadCommIoOutInterfaceListFromJson() == true {
		for _, v := range CommunicationIoOutMap {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	//获取开关量输入通信接口参数
	if ReadCommIoInInterfaceListFromJson() == true {
		for _, v := range CommunicationIoInMap {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	//获取HTTPSmartNode通信接口参数
	if ReadCommHTTPSmartNodeInterfaceListFromJson() == true {
		for _, v := range CommunicationHTTPSmartNodeMap {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	//获取HTTPTianGang通信接口参数
	if ReadCommHTTPTianGangInterfaceListFromJson() == true {
		for _, v := range CommunicationHTTPTianGangMap {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	//获取S7通信接口参数
	if ReadCommS7InterfaceListFromJson() == true {
		for _, v := range CommunicationS7Map {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	//获取天睿信诚分体空调通信接口参数
	if ReadCommSAC009InterfaceListFromJson() == true {
		for _, v := range CommunicationSAC009Map {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	//获取济南有人DR504通信接口参数
	if ReadCommDR504InterfaceListFromJson() == true {
		for _, v := range CommunicationDR504Map {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	//获取STL通信接口参数
	if ReadCommSTLInterfaceListFromJson() == true {
		for _, v := range CommunicationSTLMap {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	//获取ModbusTCP通信接口参数
	if ReadCommMBTCPInterfaceListFromJson() == true {
		for _, v := range CommunicationMBTCPMap {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	//获取ModbusRTU通信接口参数
	if ReadCommMBRTUInterfaceListFromJson() == true {
		for _, v := range CommunicationMBRTUMap {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	//获取和达科技水表通信接口参数
	if ReadCommHDKJInterfaceListFromJson() == true {
		for _, v := range CommunicationHDKJMap {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	for _, v := range CommunicationInterfaceMap {
		v.Open()
	}
}

func ReadCommInterfaceProtocol() {
	data, err := utils.FileRead("./selfpara/commInterfaceProtocol.json")
	if err != nil {
		setting.ZAPS.Debugf("打开通信接口协议配置json文件失败 %v", err)
		return
	}
	err = json.Unmarshal(data, &CommIntefaceProtocols)
	if err != nil {
		setting.ZAPS.Errorf("通信接口协议配置json文件格式化失败 %v", err)
		return
	}
	setting.ZAPS.Debugf("打开通信接口协议配置json文件成功")
}
