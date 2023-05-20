package contorl

import (
	"encoding/json"
	"gateway/device"
	"gateway/device/commInterface"
	"gateway/httpServer/model"
	"gateway/setting"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApiAddCommInterface(context *gin.Context) {

	var Param json.RawMessage
	interfaceInfo := struct {
		Name  string           `json:"name"` // 接口名称
		Type  string           `json:"type"` // 接口类型,比如serial,TcpClient,udp,http
		Param *json.RawMessage `json:"param"`
	}{
		Param: &Param,
	}

	err := context.BindJSON(&interfaceInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "JSON解析错误",
			Data:    "",
		})
		return
	}

	switch interfaceInfo.Type {
	case "localSerial":
		serial := commInterface.SerialInterfaceParam{}
		err = json.Unmarshal(Param, &serial)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[串口]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[串口]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("通信接口[串口]json %+v", serial)
		SerialInterface := &commInterface.CommunicationSerialTemplate{
			Param: serial,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		for _, v := range commInterface.CommunicationSerialMap {
			if v.Param.Name == SerialInterface.Param.Name {
				context.JSON(http.StatusOK, model.ResponseData{
					Code:    "1",
					Message: "通信接口[串口]名称已经存在",
					Data:    "",
				})
				return
			}
		}
		SerialInterface.Open()
		commInterface.CommunicationSerialMap = append(commInterface.CommunicationSerialMap, SerialInterface)
		commInterface.WriteCommSerialInterfaceListToJson()
		commInterface.CommunicationInterfaceMap = append(commInterface.CommunicationInterfaceMap, SerialInterface)
	case "tcpClient":
		TcpClient := commInterface.TcpClientInterfaceParam{}
		err = json.Unmarshal(Param, &TcpClient)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[TCP客户端]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[TCP客户端]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("通信接口[TCP客户端] %+v", TcpClient)
		TcpClientInterface := &commInterface.CommunicationTcpClientTemplate{
			Param: TcpClient,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		for _, v := range commInterface.CommunicationTcpClientMap {
			if (v.Param.Port == TcpClientInterface.Param.Port) && (v.Param.IP == TcpClientInterface.Param.IP) {
				context.JSON(http.StatusOK, model.ResponseData{
					Code:    "1",
					Message: "通信接口[TCP客户端]端口已经使用",
					Data:    "",
				})
				return
			}
		}
		//增加通信接口立马生效
		TcpClientInterface.Open()
		commInterface.CommunicationTcpClientMap = append(commInterface.CommunicationTcpClientMap, TcpClientInterface)
		commInterface.WriteCommTcpClientInterfaceListToJson()
		commInterface.CommunicationInterfaceMap = append(commInterface.CommunicationInterfaceMap, TcpClientInterface)
	case "tcpServer":
		TcpServer := commInterface.TcpServerInterfaceParam{}
		err = json.Unmarshal(Param, &TcpServer)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[TCP服务端]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[TCP服务端]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("通信接口[TCP服务端] %+v", TcpServer)
		TcpServerInterface := &commInterface.CommunicationTcpServerTemplate{
			Param: TcpServer,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		for _, v := range commInterface.CommunicationTcpServerMap {
			if v.Param.Port == TcpServerInterface.Param.Port {
				context.JSON(http.StatusOK, model.ResponseData{
					Code:    "1",
					Message: "通信接口[TCP服务端]监听端口已经使用",
					Data:    "",
				})
				return
			}
		}
		TcpServerInterface.Open()
		commInterface.CommunicationTcpServerMap = append(commInterface.CommunicationTcpServerMap, TcpServerInterface)
		commInterface.WriteCommTcpServerInterfaceListToJson()
		commInterface.CommunicationInterfaceMap = append(commInterface.CommunicationInterfaceMap, TcpServerInterface)
	case "ioOut":
		IoOut := commInterface.IoOutInterfaceParam{}
		err = json.Unmarshal(Param, &IoOut)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[开关量输出]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[开关量输出]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("通信接口[开关量输出] %+v", IoOut)
		IoOutInterface := &commInterface.CommunicationIoOutTemplate{
			Param: IoOut,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		for _, v := range commInterface.CommunicationIoOutMap {
			if v.Param.Name == IoOutInterface.Param.Name {
				context.JSON(http.StatusOK, model.ResponseData{
					Code:    "1",
					Message: "通信接口[开关量输入]名称已经存在",
					Data:    "",
				})
				return
			}
		}
		IoOutInterface.Open()
		commInterface.CommunicationIoOutMap = append(commInterface.CommunicationIoOutMap, IoOutInterface)
		commInterface.WriteCommIoOutInterfaceListToJson()
		commInterface.CommunicationInterfaceMap = append(commInterface.CommunicationInterfaceMap, IoOutInterface)
	case "ioIn":
		IoIn := commInterface.IoInInterfaceParam{}
		err = json.Unmarshal(Param, &IoIn)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[开关量输入]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[开关量输入]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("通信接口[开关量输入] %+v", IoIn)
		IoInInterface := &commInterface.CommunicationIoInTemplate{
			Param: IoIn,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		for _, v := range commInterface.CommunicationIoInMap {
			if v.Param.Name == IoInInterface.Param.Name {
				context.JSON(http.StatusOK, model.ResponseData{
					Code:    "1",
					Message: "通信接口[开关量输入]名称已经存在",
					Data:    "",
				})
				return
			}
		}
		IoInInterface.Open()
		commInterface.CommunicationIoInMap = append(commInterface.CommunicationIoInMap, IoInInterface)
		commInterface.WriteCommIoInInterfaceListToJson()
		commInterface.CommunicationInterfaceMap = append(commInterface.CommunicationInterfaceMap, IoInInterface)
	case "httpSmartNode":
		httpSmartNode := commInterface.HTTPSmartNodeInterfaceParam{}
		err = json.Unmarshal(Param, &httpSmartNode)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[httpSmartNode]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[httpSmartNode]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("通信接口[httpSmartNode] %+v", httpSmartNode)
		HTTPSmartNodeInterface := &commInterface.CommunicationHTTPSmartNodeTemplate{
			Param: httpSmartNode,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		for _, v := range commInterface.CommunicationHTTPSmartNodeMap {
			if v.Name == HTTPSmartNodeInterface.Name {
				context.JSON(http.StatusOK, model.ResponseData{
					Code:    "1",
					Message: "通信接口[httpSmartNode]名称已经存在",
					Data:    "",
				})
				return
			}
		}
		HTTPSmartNodeInterface.Open()
		commInterface.CommunicationHTTPSmartNodeMap = append(commInterface.CommunicationHTTPSmartNodeMap, HTTPSmartNodeInterface)
		commInterface.WriteCommHTTPSmartNodeInterfaceListToJson()
		commInterface.CommunicationInterfaceMap = append(commInterface.CommunicationInterfaceMap, HTTPSmartNodeInterface)
	case "httpTianGang":
		httpTianGang := commInterface.HTTPTianGangInterfaceParam{}
		err = json.Unmarshal(Param, &httpTianGang)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[httpTianGang]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[httpTianGang]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("通信接口[httpTianGang] %+v", httpTianGang)
		HTTPTianGangInterface := &commInterface.CommunicationHTTPTianGangTemplate{
			Param: httpTianGang,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		for _, v := range commInterface.CommunicationHTTPTianGangMap {
			if v.Name == HTTPTianGangInterface.Name {
				context.JSON(http.StatusOK, model.ResponseData{
					Code:    "1",
					Message: "通信接口[httpTianGang]名称已经存在",
					Data:    "",
				})
				return
			}
		}
		HTTPTianGangInterface.Open()
		commInterface.CommunicationHTTPTianGangMap = append(commInterface.CommunicationHTTPTianGangMap, HTTPTianGangInterface)
		commInterface.WriteCommHTTPTianGangInterfaceListToJson()
		commInterface.CommunicationInterfaceMap = append(commInterface.CommunicationInterfaceMap, HTTPTianGangInterface)
	case "s7":
		s7 := commInterface.S7InterfaceParam{}
		err = json.Unmarshal(Param, &s7)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[S7]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[S7]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("通信接口[S7] %+v", s7)
		s7Interface := &commInterface.CommunicationS7Template{
			Param: s7,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		for _, v := range commInterface.CommunicationS7Map {
			if v.Name == s7Interface.Name {
				context.JSON(http.StatusOK, model.ResponseData{
					Code:    "1",
					Message: "通信接口[s7]名称已经存在",
					Data:    "",
				})
				return
			}
		}
		s7Interface.Open()
		commInterface.CommunicationS7Map = append(commInterface.CommunicationS7Map, s7Interface)
		commInterface.WriteCommS7InterfaceListToJson()
		commInterface.CommunicationInterfaceMap = append(commInterface.CommunicationInterfaceMap, s7Interface)
	case "sac009":
		sac009 := commInterface.SAC009InterfaceParam{}
		err = json.Unmarshal(Param, &sac009)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[SAC009]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[SAC009]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("通信接口[SAC009] %+v", sac009)
		sac009Interface := &commInterface.CommunicationSAC009Template{
			Param: sac009,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		for _, v := range commInterface.CommunicationSAC009Map {
			if v.Name == sac009Interface.Name {
				context.JSON(http.StatusOK, model.ResponseData{
					Code:    "1",
					Message: "通信接口[SAC009]名称已经存在",
					Data:    "",
				})
				return
			}
		}
		sac009Interface.Open()
		commInterface.CommunicationSAC009Map = append(commInterface.CommunicationSAC009Map, sac009Interface)
		commInterface.WriteCommSAC009InterfaceListToJson()
		commInterface.CommunicationInterfaceMap = append(commInterface.CommunicationInterfaceMap, sac009Interface)
	case "hdkj":
		hdkj := commInterface.HDKJInterfaceParam{}
		err = json.Unmarshal(Param, &hdkj)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[HDKJ]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[HDKJ]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("通信接口[HDKJ] %+v", hdkj)
		hdkjInterface := &commInterface.CommunicationHDKJTemplate{
			Param: hdkj,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		for _, v := range commInterface.CommunicationHDKJMap {
			if v.Name == hdkjInterface.Name {
				context.JSON(http.StatusOK, model.ResponseData{
					Code:    "1",
					Message: "通信接口[HDKJ]名称已经存在",
					Data:    "",
				})
				return
			}
		}
		hdkjInterface.Open()
		commInterface.CommunicationHDKJMap = append(commInterface.CommunicationHDKJMap, hdkjInterface)
		commInterface.WriteCommHDKJInterfaceListToJson()
		commInterface.CommunicationInterfaceMap = append(commInterface.CommunicationInterfaceMap, hdkjInterface)
	case "dr504":
		dr504 := commInterface.DR504InterfaceParam{}
		err = json.Unmarshal(Param, &dr504)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[DR504]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[SAC009]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("通信接口[DR504] %+v", dr504)
		dr504Interface := &commInterface.CommunicationDR504Template{
			Param: dr504,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		for _, v := range commInterface.CommunicationDR504Map {
			if v.Name == dr504Interface.Name {
				context.JSON(http.StatusOK, model.ResponseData{
					Code:    "1",
					Message: "通信接口[DR504]名称已经存在",
					Data:    "",
				})
				return
			}
		}
		dr504Interface.Open()
		commInterface.CommunicationDR504Map = append(commInterface.CommunicationDR504Map, dr504Interface)
		commInterface.WriteCommDR504InterfaceListToJson()
		commInterface.CommunicationInterfaceMap = append(commInterface.CommunicationInterfaceMap, dr504Interface)
	case "stl":
		stl := commInterface.STLInterfaceParam{}
		err = json.Unmarshal(Param, &stl)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[STL]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[STL]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("通信接口[STL] %+v", stl)
		stlInterface := &commInterface.CommunicationSTLTemplate{
			Param: stl,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		for _, v := range commInterface.CommunicationSTLMap {
			if v.Name == stlInterface.Name {
				context.JSON(http.StatusOK, model.ResponseData{
					Code:    "1",
					Message: "通信接口[STL]名称已经存在",
					Data:    "",
				})
				return
			}
		}
		stlInterface.Open()
		commInterface.CommunicationSTLMap = append(commInterface.CommunicationSTLMap, stlInterface)
		commInterface.WriteCommSTLInterfaceListToJson()
		commInterface.CommunicationInterfaceMap = append(commInterface.CommunicationInterfaceMap, stlInterface)
	case "mbTCPClient":
		mbTCP := commInterface.MBTCPInterfaceParam{}
		err = json.Unmarshal(Param, &mbTCP)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[MBTCP]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[MBTCP]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("通信接口[MBTCP] %+v", mbTCP)
		mbTCPInterface := &commInterface.CommunicationMBTCPTemplate{
			Param: mbTCP,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		for _, v := range commInterface.CommunicationMBTCPMap {
			if v.Name == mbTCPInterface.Name {
				context.JSON(http.StatusOK, model.ResponseData{
					Code:    "1",
					Message: "通信接口[MBTCP]名称已经存在",
					Data:    "",
				})
				return
			}
		}

		mbTCPInterface.Open() //LUO add 20230405

		commInterface.CommunicationMBTCPMap = append(commInterface.CommunicationMBTCPMap, mbTCPInterface)
		commInterface.WriteCommMBTCPInterfaceListToJson()
		commInterface.CommunicationInterfaceMap = append(commInterface.CommunicationInterfaceMap, mbTCPInterface)
	case "mbRTUClient":
		mbRTU := commInterface.MBRTUInterfaceParam{}
		err = json.Unmarshal(Param, &mbRTU)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[MBRTU]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[MBRTU]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("通信接口[MBTCP] %+v", mbRTU)
		mbRTUInterface := &commInterface.CommunicationMBRTUTemplate{
			Param: mbRTU,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		for _, v := range commInterface.CommunicationMBRTUMap {
			if v.Name == mbRTUInterface.Name {
				context.JSON(http.StatusOK, model.ResponseData{
					Code:    "1",
					Message: "通信接口[MBRTU]名称已经存在",
					Data:    "",
				})
				return
			}
		}
		commInterface.CommunicationMBRTUMap = append(commInterface.CommunicationMBRTUMap, mbRTUInterface)
		commInterface.WriteCommMBRTUInterfaceListToJson()
		commInterface.CommunicationInterfaceMap = append(commInterface.CommunicationInterfaceMap, mbRTUInterface)
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "添加通信接口成功",
		Data:    "",
	})
}

func ApiModifyCommInterface(context *gin.Context) {

	var Param json.RawMessage
	interfaceInfo := struct {
		Name  string           `json:"name"` // 接口名称
		Type  string           `json:"type"` // 接口类型,比如serial,TcpClient,udp,http
		Param *json.RawMessage `json:"param"`
	}{
		Param: &Param,
	}

	err := context.BindJSON(&interfaceInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "JSON解析错误",
			Data:    "",
		})
		return
	}

	switch interfaceInfo.Type {
	case "localSerial":
		serial := commInterface.SerialInterfaceParam{}
		err = json.Unmarshal(Param, &serial)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[串口]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[串口]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("修改通信接口[串口]%+v", serial)
		SerialInterface := &commInterface.CommunicationSerialTemplate{
			Param: serial,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}
		index := -1
		for k, v := range commInterface.CommunicationSerialMap {
			if v.Name == SerialInterface.Name {
				index = k
				v.Close()
			}
		}
		if index != -1 {
			commInterface.CommunicationSerialMap[index].Close()
			for _, v := range device.CollectInterfaceMap.Coll {
				if v.CommInterfaceName == SerialInterface.Name {
					v.CommInterface = SerialInterface
					v.CommInterfaceUpdate <- true
				}
			}
			commInterface.CommunicationSerialMap[index] = SerialInterface
			commInterface.WriteCommSerialInterfaceListToJson()
			SerialInterface.Open()

			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "修改通信接口[串口]成功",
				Data:    "",
			})
			return
		}
	case "tcpClient":
		TcpClient := commInterface.TcpClientInterfaceParam{}
		err = json.Unmarshal(Param, &TcpClient)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[TCP客户端]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[tcp客户端]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("修改通信接口[TcpClient]%+v", TcpClient)
		TcpClientInterface := &commInterface.CommunicationTcpClientTemplate{
			Param: TcpClient,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		index := -1
		for k, v := range commInterface.CommunicationTcpClientMap {
			if v.Name == TcpClientInterface.Name {
				index = k
				v.Close()
				break
			}
		}
		if index != -1 {
			for _, v := range device.CollectInterfaceMap.Coll {
				if v.CommInterfaceName == TcpClientInterface.Name {
					v.CommInterface = TcpClientInterface
					v.CommInterfaceUpdate <- true
				}
			}
			commInterface.CommunicationTcpClientMap[index] = TcpClientInterface
			commInterface.WriteCommTcpClientInterfaceListToJson()
			TcpClientInterface.Open()
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "修改通信接口[TcpClient]成功",
				Data:    "",
			})
			return
		}
	case "tcpServer":
		TcpServer := commInterface.TcpServerInterfaceParam{}
		err = json.Unmarshal(Param, &TcpServer)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[TCP服务端]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[TCP服务端]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("修改通信接口[TcpServer] %+v", TcpServer)
		TcpServerInterface := &commInterface.CommunicationTcpServerTemplate{
			Param: TcpServer,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		index := -1
		for k, v := range commInterface.CommunicationTcpServerMap {
			if v.Name == TcpServerInterface.Name {
				index = k
				v.Close()
			}
		}
		if index != -1 {
			commInterface.CommunicationTcpServerMap[index] = TcpServerInterface
			for _, v := range device.CollectInterfaceMap.Coll {
				if v.CommInterfaceName == TcpServerInterface.Name {
					v.CommInterface = TcpServerInterface
					v.CommInterfaceUpdate <- true
				}
			}
			commInterface.WriteCommTcpServerInterfaceListToJson()
			TcpServerInterface.Open()
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "通信接口[TCP服务端]成功",
				Data:    "",
			})
			return
		}
	case "ioOut":
		IoOut := commInterface.IoOutInterfaceParam{}
		err = json.Unmarshal(Param, &IoOut)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[开关量输出]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[开关量输出]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("修改通信接口[开关量输出]%+v", IoOut)
		IoOutInterface := &commInterface.CommunicationIoOutTemplate{
			Param: IoOut,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		index := -1
		for k, v := range commInterface.CommunicationIoOutMap {
			if v.Name == IoOutInterface.Name {
				index = k
				v.Close()
			}
		}
		if index != -1 {
			commInterface.CommunicationIoOutMap[index] = IoOutInterface
			for _, v := range device.CollectInterfaceMap.Coll {
				if v.CommInterfaceName == IoOutInterface.Name {
					v.CommInterface = IoOutInterface
					v.CommInterfaceUpdate <- true
				}
			}
			commInterface.WriteCommIoOutInterfaceListToJson()
			IoOutInterface.Open()
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "通信接口[开关量输出]成功",
				Data:    "",
			})
			return
		}
	case "ioIn":
		IoIn := commInterface.IoInInterfaceParam{}
		err = json.Unmarshal(Param, &IoIn)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[开关量输入]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[开关量输入]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("修改通信接口[开关量输入] %+v", IoIn)
		IoInInterface := &commInterface.CommunicationIoInTemplate{
			Param: IoIn,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}
		index := -1
		for k, v := range commInterface.CommunicationIoInMap {
			if v.Name == IoInInterface.Name {
				index = k
				v.Close()
			}
		}
		if index != -1 {
			commInterface.CommunicationIoInMap[index] = IoInInterface
			for _, v := range device.CollectInterfaceMap.Coll {
				if v.CommInterfaceName == IoInInterface.Name {
					v.CommInterface = IoInInterface
					v.CommInterfaceUpdate <- true
				}
			}
			commInterface.WriteCommIoInInterfaceListToJson()
			IoInInterface.Open()
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "修改通信接口[开关量输入]成功",
				Data:    "",
			})
			return
		}
	case "httpTianGang":
		httpTianGang := commInterface.HTTPTianGangInterfaceParam{}
		err = json.Unmarshal(Param, &httpTianGang)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[httpTianGang]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[httpTianGang]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("修改通信接口[httpTianGang]%+v", httpTianGang)
		HTTPTianGangInterface := &commInterface.CommunicationHTTPTianGangTemplate{
			Param: httpTianGang,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		index := -1
		for k, v := range commInterface.CommunicationHTTPTianGangMap {
			if v.Name == HTTPTianGangInterface.Name {
				index = k
				v.Close()
				break
			}
		}
		if index != -1 {
			for _, v := range device.CollectInterfaceMap.Coll {
				if v.CommInterfaceName == HTTPTianGangInterface.Name {
					v.CommInterface = HTTPTianGangInterface
					v.CommInterfaceUpdate <- true
				}
			}
			commInterface.CommunicationHTTPTianGangMap[index] = HTTPTianGangInterface
			commInterface.WriteCommHTTPTianGangInterfaceListToJson()
			HTTPTianGangInterface.Open()
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "修改通信接口[HTTPTianGang]成功",
				Data:    "",
			})
			return
		}
	case "s7":
		s7 := commInterface.S7InterfaceParam{}
		err = json.Unmarshal(Param, &s7)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[s7]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[s7]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("修改通信接口[s7]%+v", s7)
		s7Interface := &commInterface.CommunicationS7Template{
			Param: s7,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		index := -1
		for k, v := range commInterface.CommunicationS7Map {
			if v.Name == s7Interface.Name {
				index = k
				v.Close()
				break
			}
		}
		if index != -1 {
			for _, v := range device.CollectInterfaceMap.Coll {
				if v.CommInterfaceName == s7Interface.Name {
					v.CommInterface = s7Interface
					v.CommInterfaceUpdate <- true
				}
			}
			commInterface.CommunicationS7Map[index] = s7Interface
			commInterface.WriteCommS7InterfaceListToJson()
			s7Interface.Open()
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "修改通信接口[s7]成功",
				Data:    "",
			})
			return
		}
	case "sac009":
		sac009 := commInterface.SAC009InterfaceParam{}
		err = json.Unmarshal(Param, &sac009)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[SAC009]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[SAC009]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("修改通信接口[SAC009]%+v", sac009)
		sac009Interface := &commInterface.CommunicationSAC009Template{
			Param: sac009,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		index := -1
		for k, v := range commInterface.CommunicationSAC009Map {
			if v.Name == sac009Interface.Name {
				index = k
				v.Close()
				break
			}
		}
		if index != -1 {
			for _, v := range device.CollectInterfaceMap.Coll {
				if v.CommInterfaceName == sac009Interface.Name {
					v.CommInterface = sac009Interface
					v.CommInterfaceUpdate <- true
				}
			}
			commInterface.CommunicationSAC009Map[index] = sac009Interface
			commInterface.WriteCommSAC009InterfaceListToJson()
			sac009Interface.Open()
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "修改通信接口[SAC009]成功",
				Data:    "",
			})
			return
		}
	case "hdkj":
		hdkj := commInterface.HDKJInterfaceParam{}
		err = json.Unmarshal(Param, &hdkj)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[HDKJ]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[HDKJ]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("修改通信接口[HDKJ]%+v", hdkj)
		hdkjInterface := &commInterface.CommunicationHDKJTemplate{
			Param: hdkj,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		index := -1
		for k, v := range commInterface.CommunicationHDKJMap {
			if v.Name == hdkjInterface.Name {
				index = k
				v.Close()
				break
			}
		}
		if index != -1 {
			for _, v := range device.CollectInterfaceMap.Coll {
				if v.CommInterfaceName == hdkjInterface.Name {
					v.CommInterface = hdkjInterface
					v.CommInterfaceUpdate <- true
				}
			}
			commInterface.CommunicationHDKJMap[index] = hdkjInterface
			commInterface.WriteCommHDKJInterfaceListToJson()
			hdkjInterface.Open()
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "修改通信接口[HDKJ]成功",
				Data:    "",
			})
			return
		}
	case "dr504":
		dr504 := commInterface.DR504InterfaceParam{}
		err = json.Unmarshal(Param, &dr504)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[DR504]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[DR504]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("修改通信接口[DR504]%+v", dr504)
		dr504Interface := &commInterface.CommunicationDR504Template{
			Param: dr504,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		index := -1
		for k, v := range commInterface.CommunicationDR504Map {
			if v.Name == dr504Interface.Name {
				index = k
				v.Close()
				break
			}
		}
		if index != -1 {
			for _, v := range device.CollectInterfaceMap.Coll {
				if v.CommInterfaceName == dr504Interface.Name {
					v.CommInterface = dr504Interface
					v.CommInterfaceUpdate <- true
				}
			}
			commInterface.CommunicationDR504Map[index] = dr504Interface
			commInterface.WriteCommDR504InterfaceListToJson()
			dr504Interface.Open()
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "修改通信接口[DR504]成功",
				Data:    "",
			})
			return
		}
	case "mbTCPClient":
		mbTCP := commInterface.MBTCPInterfaceParam{}
		err = json.Unmarshal(Param, &mbTCP)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[MBTCP]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[MBTCP]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("修改通信接口[MBTCP]%+v", mbTCP)
		mbTCPInterface := &commInterface.CommunicationMBTCPTemplate{
			Param: mbTCP,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		index := -1
		for k, v := range commInterface.CommunicationMBTCPMap {
			if v.Name == mbTCPInterface.Name {
				index = k
				v.Close()
				break
			}
		}
		if index != -1 {
			for _, v := range device.CollectInterfaceMap.Coll {
				if v.CommInterfaceName == mbTCPInterface.Name {
					v.CommInterface = mbTCPInterface
					v.CommInterfaceUpdate <- true
				}
			}
			commInterface.CommunicationMBTCPMap[index] = mbTCPInterface
			commInterface.WriteCommMBTCPInterfaceListToJson()
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "修改通信接口[MBTCP]成功",
				Data:    "",
			})
			return
		}
	case "mbRTUClient":
		mbRTU := commInterface.MBRTUInterfaceParam{}
		err = json.Unmarshal(Param, &mbRTU)
		if err != nil {
			setting.ZAPS.Errorf("通信接口[MBRTU]JSON格式化错误 %v", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "通信接口[MBRTU]JSON解析错误",
				Data:    "",
			})
			return
		}
		setting.ZAPS.Debugf("修改通信接口[MBRTU]%+v", mbRTU)
		mbRTUInterface := &commInterface.CommunicationMBRTUTemplate{
			Param: mbRTU,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		index := -1
		for k, v := range commInterface.CommunicationMBRTUMap {
			if v.Name == mbRTUInterface.Name {
				index = k
				v.Close()
				break
			}
		}
		if index != -1 {
			for _, v := range device.CollectInterfaceMap.Coll {
				if v.CommInterfaceName == mbRTUInterface.Name {
					v.CommInterface = mbRTUInterface
					v.CommInterfaceUpdate <- true
				}
			}
			commInterface.CommunicationMBRTUMap[index] = mbRTUInterface
			commInterface.WriteCommMBRTUInterfaceListToJson()
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "修改通信接口[MBRTU]成功",
				Data:    "",
			})
			return
		}
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "1",
		Message: "通信接口名称不存在",
		Data:    "",
	})
	return
}

func ApiDeleteCommInterface(context *gin.Context) {

	interfaceInfo := struct {
		Name string `json:"name"`
	}{}

	err := context.BindJSON(&interfaceInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "JSON解析错误",
			Data:    "",
		})
		return
	}

	_, ok := device.CollectInterfaceMap.Coll[interfaceInfo.Name]
	if ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "通信接口被使用，不可以删除",
			Data:    "",
		})
		return
	}

	for k, v := range commInterface.CommunicationSerialMap {
		if v.Name == interfaceInfo.Name {
			v.Close()
			commInterface.CommunicationSerialMap = append(commInterface.CommunicationSerialMap[:k], commInterface.CommunicationSerialMap[k+1:]...)
			commInterface.WriteCommSerialInterfaceListToJson()

			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "",
				Data:    "",
			})
			return
		}
	}

	for k, v := range commInterface.CommunicationTcpClientMap {
		if v.Name == interfaceInfo.Name {
			v.Close()
			commInterface.CommunicationTcpClientMap = append(commInterface.CommunicationTcpClientMap[:k], commInterface.CommunicationTcpClientMap[k+1:]...)
			commInterface.WriteCommTcpClientInterfaceListToJson()

			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "",
				Data:    "",
			})
			return
		}
	}

	for k, v := range commInterface.CommunicationTcpServerMap {
		if v.Name == interfaceInfo.Name {
			v.Close()
			commInterface.CommunicationTcpServerMap = append(commInterface.CommunicationTcpServerMap[:k], commInterface.CommunicationTcpServerMap[k+1:]...)
			commInterface.WriteCommTcpServerInterfaceListToJson()

			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "",
				Data:    "",
			})
			return
		}
	}

	for k, v := range commInterface.CommunicationIoOutMap {
		if v.Name == interfaceInfo.Name {
			v.Close()
			commInterface.CommunicationIoOutMap = append(commInterface.CommunicationIoOutMap[:k], commInterface.CommunicationIoOutMap[k+1:]...)
			commInterface.WriteCommIoOutInterfaceListToJson()

			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "",
				Data:    "",
			})
			return
		}
	}

	for k, v := range commInterface.CommunicationIoInMap {
		if v.Name == interfaceInfo.Name {
			v.Close()
			commInterface.CommunicationIoInMap = append(commInterface.CommunicationIoInMap[:k], commInterface.CommunicationIoInMap[k+1:]...)
			commInterface.WriteCommIoInInterfaceListToJson()

			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "",
				Data:    "",
			})
			return
		}
	}

	for k, v := range commInterface.CommunicationHTTPTianGangMap {
		if v.Name == interfaceInfo.Name {
			v.Close()
			commInterface.CommunicationHTTPTianGangMap = append(commInterface.CommunicationHTTPTianGangMap[:k], commInterface.CommunicationHTTPTianGangMap[k+1:]...)
			commInterface.WriteCommHTTPTianGangInterfaceListToJson()

			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "",
				Data:    "",
			})
			return
		}
	}

	for k, v := range commInterface.CommunicationS7Map {
		if v.Name == interfaceInfo.Name {
			v.Close()
			commInterface.CommunicationS7Map = append(commInterface.CommunicationS7Map[:k], commInterface.CommunicationS7Map[k+1:]...)
			commInterface.WriteCommS7InterfaceListToJson()

			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "",
				Data:    "",
			})
			return
		}
	}

	for k, v := range commInterface.CommunicationSAC009Map {
		if v.Name == interfaceInfo.Name {
			v.Close()
			commInterface.CommunicationSAC009Map = append(commInterface.CommunicationSAC009Map[:k], commInterface.CommunicationSAC009Map[k+1:]...)
			commInterface.WriteCommSAC009InterfaceListToJson()

			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "",
				Data:    "",
			})
			return
		}
	}

	for k, v := range commInterface.CommunicationHDKJMap {
		if v.Name == interfaceInfo.Name {
			v.Close()
			commInterface.CommunicationHDKJMap = append(commInterface.CommunicationHDKJMap[:k], commInterface.CommunicationHDKJMap[k+1:]...)
			commInterface.WriteCommHDKJInterfaceListToJson()

			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "",
				Data:    "",
			})
			return
		}
	}

	for k, v := range commInterface.CommunicationDR504Map {
		if v.Name == interfaceInfo.Name {
			v.Close()
			commInterface.CommunicationDR504Map = append(commInterface.CommunicationDR504Map[:k], commInterface.CommunicationDR504Map[k+1:]...)
			commInterface.WriteCommDR504InterfaceListToJson()

			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "通信接口删除成功",
				Data:    "",
			})
			return
		}
	}

	for k, v := range commInterface.CommunicationMBTCPMap {
		if v.Name == interfaceInfo.Name {
			v.Close()
			commInterface.CommunicationMBTCPMap = append(commInterface.CommunicationMBTCPMap[:k], commInterface.CommunicationMBTCPMap[k+1:]...)
			commInterface.WriteCommMBTCPInterfaceListToJson()

			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "通信接口删除成功",
				Data:    "",
			})
			return
		}
	}

	for k, v := range commInterface.CommunicationMBRTUMap {
		if v.Name == interfaceInfo.Name {
			v.Close()
			commInterface.CommunicationMBRTUMap = append(commInterface.CommunicationMBRTUMap[:k], commInterface.CommunicationMBRTUMap[k+1:]...)
			commInterface.WriteCommMBRTUInterfaceListToJson()

			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "通信接口删除成功",
				Data:    "",
			})
			return
		}
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "1",
		Message: "通信接口不存在",
		Data:    "",
	})
	return
}

func ApiGetCommInterface(context *gin.Context) {

	type CommunicationInterfaceTemplate struct {
		Name   string                      `json:"name"`   // 接口名称
		Type   string                      `json:"type"`   // 接口类型,比如serial,TcpClient,udp,http
		Status commInterface.ConnectStatus `json:"status"` //接口状态
		Param  interface{}                 `json:"param"`  // 接口参数
	}

	interfaceMap := make([]CommunicationInterfaceTemplate, 0)

	for _, v := range commInterface.CommunicationSerialMap {
		CommunicationInterface := CommunicationInterfaceTemplate{
			Name:   v.Name,
			Type:   v.Type,
			Status: v.Status,
			Param:  v.Param,
		}
		interfaceMap = append(interfaceMap, CommunicationInterface)
	}

	for _, v := range commInterface.CommunicationTcpClientMap {
		CommunicationInterface := CommunicationInterfaceTemplate{
			Name:   v.Name,
			Type:   v.Type,
			Status: v.Status,
			Param:  v.Param,
		}
		interfaceMap = append(interfaceMap, CommunicationInterface)
	}

	for _, v := range commInterface.CommunicationTcpServerMap {
		CommunicationInterface := CommunicationInterfaceTemplate{
			Name:   v.Name,
			Type:   v.Type,
			Status: v.Status,
			Param:  v.Param,
		}
		interfaceMap = append(interfaceMap, CommunicationInterface)
	}

	for _, v := range commInterface.CommunicationIoOutMap {
		CommunicationInterface := CommunicationInterfaceTemplate{
			Name:   v.Name,
			Type:   v.Type,
			Status: v.Status,
			Param:  v.Param,
		}
		interfaceMap = append(interfaceMap, CommunicationInterface)
	}

	for _, v := range commInterface.CommunicationIoInMap {
		CommunicationInterface := CommunicationInterfaceTemplate{
			Name:   v.Name,
			Type:   v.Type,
			Status: v.Status,
			Param:  v.Param,
		}
		interfaceMap = append(interfaceMap, CommunicationInterface)
	}

	for _, v := range commInterface.CommunicationHTTPSmartNodeMap {
		CommunicationInterface := CommunicationInterfaceTemplate{
			Name:   v.Name,
			Type:   v.Type,
			Status: v.Status,
			Param:  v.Param,
		}
		interfaceMap = append(interfaceMap, CommunicationInterface)
	}

	for _, v := range commInterface.CommunicationHTTPTianGangMap {
		CommunicationInterface := CommunicationInterfaceTemplate{
			Name:   v.Name,
			Type:   v.Type,
			Status: v.Status,
			Param:  v.Param,
		}
		interfaceMap = append(interfaceMap, CommunicationInterface)
	}

	for _, v := range commInterface.CommunicationS7Map {
		CommunicationInterface := CommunicationInterfaceTemplate{
			Name:   v.Name,
			Type:   v.Type,
			Status: v.Status,
			Param:  v.Param,
		}
		interfaceMap = append(interfaceMap, CommunicationInterface)
	}

	for _, v := range commInterface.CommunicationSAC009Map {
		CommunicationInterface := CommunicationInterfaceTemplate{
			Name:   v.Name,
			Type:   v.Type,
			Status: v.Status,
			Param:  v.Param,
		}
		interfaceMap = append(interfaceMap, CommunicationInterface)
	}

	for _, v := range commInterface.CommunicationHDKJMap {
		CommunicationInterface := CommunicationInterfaceTemplate{
			Name:   v.Name,
			Type:   v.Type,
			Status: v.Status,
			Param:  v.Param,
		}
		interfaceMap = append(interfaceMap, CommunicationInterface)
	}

	for _, v := range commInterface.CommunicationDR504Map {
		CommunicationInterface := CommunicationInterfaceTemplate{
			Name:   v.Name,
			Type:   v.Type,
			Status: v.Status,
			Param:  v.Param,
		}
		interfaceMap = append(interfaceMap, CommunicationInterface)
	}

	for _, v := range commInterface.CommunicationMBTCPMap {
		CommunicationInterface := CommunicationInterfaceTemplate{
			Name:   v.Name,
			Type:   v.Type,
			Status: v.Status,
			Param:  v.Param,
		}
		interfaceMap = append(interfaceMap, CommunicationInterface)
	}

	for _, v := range commInterface.CommunicationMBRTUMap {
		CommunicationInterface := CommunicationInterfaceTemplate{
			Name:   v.Name,
			Type:   v.Type,
			Status: v.Status,
			Param:  v.Param,
		}
		interfaceMap = append(interfaceMap, CommunicationInterface)
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取通信接口成功",
		Data:    interfaceMap,
	})
}

func ApiGetCommInterfaceProtocol(context *gin.Context) {

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取通信接口支持协议",
		Data:    commInterface.CommIntefaceProtocols,
	})
}
