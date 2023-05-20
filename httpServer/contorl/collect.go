package contorl

import (
	"gateway/device"
	"gateway/httpServer/model"
	"gateway/setting"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

func ApiAddCollectInterface(context *gin.Context) {
	interfaceInfo := struct {
		CollectInterfaceName string `json:"collInterfaceName"` // 采集接口名字
		CommInterfaceName    string `json:"commInterfaceName"` // 通信接口名字
		ProtocolTypeName     string `json:"protocolTypeName"`  // 协议名称  gwai add 2023-05-10
		PollPeriod           int    `json:"pollPeriod"`
		OfflinePeriod        int    `json:"offlinePeriod"`
	}{}

	err := context.ShouldBindJSON(&interfaceInfo)
	if err != nil {
		setting.ZAPS.Error("增加采集接口JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "增加采集接口JSON格式化错误",
			Data:    "",
		})
		return
	}

	err = device.AddCollectInterface(interfaceInfo.CollectInterfaceName,
		interfaceInfo.CommInterfaceName,
		interfaceInfo.ProtocolTypeName,
		interfaceInfo.PollPeriod,
		interfaceInfo.OfflinePeriod)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Data:    "",
			Message: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "增加采集接口成功",
		Data:    "",
	})
}

func ApiModifyCollectInterface(context *gin.Context) {

	interfaceInfo := struct {
		CollectInterfaceName string `json:"collInterfaceName"` // 采集接口名字
		CommInterfaceName    string `json:"commInterfaceName"` // 通信接口名字
		ProtocolTypeName     string `json:"protocolTypeName"`  // 协议名称  gwai add 2023-05-10
		PollPeriod           int    `json:"pollPeriod"`
		OfflinePeriod        int    `json:"offlinePeriod"`
	}{}

	err := context.ShouldBindJSON(&interfaceInfo)
	if err != nil {
		setting.ZAPS.Error("修改采集接口JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改采集接口JSON格式化错误",
			Data:    "",
		})
		return
	}

	err = device.ModifyCollectInterface(interfaceInfo.CollectInterfaceName,
		interfaceInfo.CommInterfaceName,
		interfaceInfo.ProtocolTypeName,
		interfaceInfo.PollPeriod,
		interfaceInfo.OfflinePeriod)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Data:    "",
			Message: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "修改采集接口成功",
		Data:    "",
	})
}

func ApiDeleteCollectInterface(context *gin.Context) {

	interfaceInfo := struct {
		Name string `json:"name"` // 采集接口名字
	}{}

	err := context.ShouldBindJSON(&interfaceInfo)
	if err != nil {
		setting.ZAPS.Error("删除采集接口JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除采集接口JSON格式化错误",
			Data:    "",
		})
		return
	}

	err = device.DeleteCollectInterface(interfaceInfo.Name)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Data:    "",
			Message: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "删除采集接口成功",
		Data:    "",
	})
}

func ApiGetCollectInterface(context *gin.Context) {

	type DeviceNodeTemplate struct {
		Index          int    `json:"index"`          //设备偏移量
		Name           string `json:"name"`           //设备名称
		Label          string `json:"label"`          //设备标签
		Addr           string `json:"addr"`           //设备地址
		TSL            string `json:"tsl"`            //设备物模型
		LastCommRTC    string `json:"lastCommRTC"`    //最后一次通信时间戳
		CommTotalCnt   int    `json:"commTotalCnt"`   //通信总次数
		CommSuccessCnt int    `json:"commSuccessCnt"` //通信成功次数
		CurCommFailCnt int    `json:"-"`              //当前通信失败次数
		CommStatus     string `json:"commStatus"`     //通信状态
		CommMaxTime    int    `json:"commMaxTime"`    //通信最长用时
		CommMinTime    int    `json:"commMinTime"`    //通信最短用时
	}

	//采集接口模板
	collectInterface := struct {
		DeviceNodeCnt       int                  `json:"deviceNodeCnt"`       //设备数量
		DeviceNodeOnlineCnt int                  `json:"deviceNodeOnlineCnt"` //设备在线数量
		DeviceNodeMap       []DeviceNodeTemplate `json:"deviceNodeMap"`       //节点表
	}{
		DeviceNodeMap: make([]DeviceNodeTemplate, 0),
	}

	sName := context.Query("name")

	coll, ok := device.CollectInterfaceMap.Coll[sName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Data:    "",
			Message: "采集接口名称不存在",
		})
		return
	}

	collectInterface.DeviceNodeCnt = coll.DeviceNodeCnt
	collectInterface.DeviceNodeOnlineCnt = coll.DeviceNodeOnlineCnt

	node := DeviceNodeTemplate{}
	for _, v := range coll.DeviceNodeMap {
		node.Index = v.Index
		node.Name = v.Name
		node.Label = v.Label
		node.Addr = v.Addr
		node.TSL = v.TSL
		node.LastCommRTC = v.LastCommRTC
		node.CommTotalCnt = v.CommTotalCnt
		node.CommSuccessCnt = v.CommSuccessCnt
		node.CurCommFailCnt = v.CurCommFailCnt
		node.CommStatus = v.CommStatus
		collectInterface.DeviceNodeMap = append(collectInterface.DeviceNodeMap, node)
	}

	//排序，方便前端页面显示
	sort.Slice(collectInterface.DeviceNodeMap, func(i, j int) bool {
		iName := collectInterface.DeviceNodeMap[i].Index
		jName := collectInterface.DeviceNodeMap[j].Index
		return iName > jName
	})

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取单个采集接口信息成功",
		Data:    collectInterface,
	})
}

func ApiGetCollectInterfaces(context *gin.Context) {

	type InterfaceParamTemplate struct {
		Index               int    `json:"index"`
		CollInterfaceName   string `json:"collInterfaceName"`   // 采集接口
		CommInterfaceName   string `json:"commInterfaceName"`   // 通信接口
		ProtocolTypeName    string `json:"protocolTypeName"`    // 协议名称  gwai add 2023-05-10
		PollPeriod          int    `json:"pollPeriod"`          // 采集周期
		OfflinePeriod       int    `json:"offlinePeriod"`       // 离线超时周期
		DeviceNodeCnt       int    `json:"deviceNodeCnt"`       // 设备数量
		DeviceNodeOnlineCnt int    `json:"deviceNodeOnlineCnt"` // 设备在线数量
	}

	interfaces := make([]InterfaceParamTemplate, 0)
	for _, v := range device.CollectInterfaceMap.Coll {
		if v == nil {
			continue
		}
		Param := InterfaceParamTemplate{
			Index:               v.Index,
			CollInterfaceName:   v.CollInterfaceName,
			CommInterfaceName:   v.CommInterfaceName,
			ProtocolTypeName:    v.ProtocolTypeName, // 协议名称  gwai add 2023-05-10
			PollPeriod:          v.PollPeriod,
			OfflinePeriod:       v.OfflinePeriod,
			DeviceNodeCnt:       v.DeviceNodeCnt,
			DeviceNodeOnlineCnt: v.DeviceNodeOnlineCnt,
		}
		interfaces = append(interfaces, Param)
	}

	//排序，方便前端页面显示
	sort.Slice(interfaces, func(i, j int) bool {
		iName := interfaces[i].Index
		jName := interfaces[j].Index
		return iName > jName
	})

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取所有采集接口成功",
		Data:    interfaces,
	})

}
