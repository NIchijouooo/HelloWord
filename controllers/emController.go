package controllers

import (
	"encoding/json"
	"fmt"
	"gateway/device"
	"gateway/httpServer/model"
	"gateway/models"
	repositories "gateway/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type EmController struct {
	repo *repositories.EmRepository
}

func NewEMController() *EmController {
	return &EmController{
		repo: repositories.NewEmRepository(),
	}
}

func (c *EmController) RegisterRoutes(router *gin.RouterGroup) {
	// 通道
	router.POST("/api/v2/em/getAllCommInterfaceProtocols", c.GetAllCommInterfaceProtocols)
	router.POST("/api/v2/em/addCommInterface", c.AddCommInterface)
	router.POST("/api/v2/em/communication", c.GetCommInterfaces)
	router.DELETE("/api/v2/em/delComInterface", c.DelComInterface)
	router.PUT("/api/v2/em/updateCommInterface", c.UpdateCommInterface)
	// 采集接口
	router.POST("/api/v2/em/addCollInterface", c.AddCollInterface)
	router.POST("/api/v2/em/updateCollInterface", c.UpdateCollInterface)
	// 设备
	router.POST("/api/v2/em/addEmDevice", c.AddEmDevice)
	router.POST("/api/v2/em/updateEmDevice", c.AddEmDevice)
	router.POST("/api/v2/em/delEmDevice", c.AddEmDevice)
	// 设备模型
	router.POST("/api/v2/em/addEmDeviceModel", c.AddEmDeviceModel)
	router.POST("/api/v2/em/addEmDeviceModelCmd", c.AddEmDeviceModelCmd)
	router.POST("/api/v2/em/addEmDeviceModelCmdParam", c.AddEmDeviceModelCmdParam)
	router.POST("/api/v2/em/getEmDeviceModelCmdParamListByName", c.GetEmDeviceModelCmdParamListByName)
	router.POST("/api/v2/em/getEmDeviceModelCmdParamListByDeviceId", c.GetEmDeviceModelCmdParamListByDeviceId)
}

func (c *EmController) GetAllCommInterfaceProtocols(ctx *gin.Context) {
	rows, err := c.repo.GetAllCommInterfaceProtocols()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "1",
		Message: "查询成功",
		Data:    rows,
	})
}

// AddCommInterface 新增通道
func (c *EmController) AddCommInterface(ctx *gin.Context) {
	var localSerial models.LocalSerial
	var TCP models.TCPClient
	var commInterface models.CommInterface
	var data []byte
	if localSerialErr := ctx.ShouldBindBodyWith(&localSerial, binding.JSON); localSerialErr == nil {
		commInterface.CommInterfaceProtocol = localSerial.Type
		commInterface.Name = localSerial.Name
		data, _ = json.Marshal(localSerial.Param)
	} else if TCPErr := ctx.ShouldBindBodyWith(&TCP, binding.JSON); TCPErr == nil {
		commInterface.CommInterfaceProtocol = TCP.Type
		commInterface.Name = TCP.Name
		data, _ = json.Marshal(TCP.Param)
	}
	// 判断是否有重名的通道，有就直接返回
	commInterfaceByName, _ := c.repo.GetCommInterfaceByName(commInterface.Name)
	if commInterfaceByName != nil {
		//ctx.JSON(http.StatusOK, model.ResponseData{
		//	Code:    "1",
		//	Message: "通道名已存在，添加失败",
		//})
		return
	}

	protocolByName, _ := c.repo.GetCommInterfaceProtocolByName(commInterface.CommInterfaceProtocol)
	commInterface.CommInterfaceProtocolId = protocolByName.Id
	commInterface.Data = string(data)
	err := c.repo.AddCommInterface(&commInterface)
	if err != nil {
		fmt.Println("error:", err)
	}
	//ctx.JSON(http.StatusOK, model.ResponseData{
	//	Code:    "0",
	//	Message: "添加通道成功",
	//})
}

// GetCommInterfaces 获取所有通道
func (c *EmController) GetCommInterfaces(ctx *gin.Context) {
	var res []interface{}
	rows, err := c.repo.GetAllCommInterfaces()
	for _, row := range rows {
		var t models.TCPClient
		var l models.LocalSerial

		if row.CommInterfaceProtocol == "localSerial" || row.CommInterfaceProtocol == "mbRTUClient" {
			l.Id = row.Id
			l.Type = row.CommInterfaceProtocol
			l.Name = row.Name
			err = json.Unmarshal([]byte(row.Data), &l.Param)
			if err != nil {
				fmt.Println("error:", err)
			}
			res = append(res, l)
		} else {
			t.Id = row.Id
			t.Type = row.CommInterfaceProtocol
			t.Name = row.Name
			err = json.Unmarshal([]byte(row.Data), &t.Param)
			if err != nil {
				fmt.Println("error:", err)
			}
			res = append(res, t)
		}
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "1",
		Message: "成功",
		Data:    res,
	})
}

// UpdateCommInterface 修改通道
func (c *EmController) UpdateCommInterface(ctx *gin.Context) {
	var localSerial models.LocalSerial
	var TCP models.TCPClient
	var commInterface models.CommInterface
	var data []byte
	if localSerialErr := ctx.ShouldBindBodyWith(&localSerial, binding.JSON); localSerialErr == nil {
		commInterface.Id = localSerial.Id
		commInterface.CommInterfaceProtocol = localSerial.Type
		commInterface.Name = localSerial.Name
		data, _ = json.Marshal(localSerial.Param)
	} else if TCPErr := ctx.ShouldBindBodyWith(&TCP, binding.JSON); TCPErr == nil {
		commInterface.Id = TCP.Id
		commInterface.CommInterfaceProtocol = TCP.Type
		commInterface.Name = TCP.Name
		data, _ = json.Marshal(TCP.Param)
	}
	protocolByName, _ := c.repo.GetCommInterfaceProtocolByName(commInterface.CommInterfaceProtocol)
	commInterface.CommInterfaceProtocolId = protocolByName.Id
	commInterface.Data = string(data)
	err := c.repo.UpdateCommInterface(&commInterface)
	if err != nil {
		fmt.Println("error:", err)
	}
	//ctx.JSON(http.StatusOK, model.ResponseData{
	//	Code:    "0",
	//	Message: "成功",
	//	Data:    commInterface,
	//})
	return
}

// DelComInterface 删除通道
func (c *EmController) DelComInterface(ctx *gin.Context) {
	var tmp struct {
		Name string `json:"name"`
	}
	if err := ctx.ShouldBindBodyWith(&tmp, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	commInterfaceByName, _ := c.repo.GetCommInterfaceByName(tmp.Name)

	if err := c.repo.DelCommInterface(commInterfaceByName.Id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//ctx.JSON(http.StatusOK, model.ResponseData{
	//	Code:    "1",
	//	Message: "删除通信接口成功",
	//})
	return
}

// AddCollInterface 新增采集接口
func (c *EmController) AddCollInterface(ctx *gin.Context) {
	var addEmCollInterface models.AddEmCollInterface
	var emCollInterface models.EmCollInterface
	var data []byte
	if err := ctx.ShouldBindBodyWith(&addEmCollInterface, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 判断采集接口是否有重名，有就直接返回
	emCollInterface.Name = addEmCollInterface.CollInterfaceName
	emCollInterfaceByName, _ := c.repo.GetCollInterfaceByName(emCollInterface.Name)
	if emCollInterfaceByName != nil {
		//ctx.JSON(http.StatusOK, model.ResponseData{
		//	Code:    "0",
		//	Message: "通道名已存在，添加失败",
		//})
		return
	}

	emCollInterface.OfflinePeriod = addEmCollInterface.OfflinePeriod
	emCollInterface.PollPeriod = addEmCollInterface.PollPeriod
	commInterfaceByName, _ := c.repo.GetCommInterfaceByName(addEmCollInterface.CommInterfaceName)
	if commInterfaceByName == nil {
		//ctx.JSON(http.StatusOK, model.ResponseData{
		//	Code:    "0",
		//	Message: "通信接口名未找到",
		//})
		return
	}
	emCollInterface.CommInterfaceId = commInterfaceByName.Id
	data, _ = json.Marshal(addEmCollInterface)
	emCollInterface.Data = string(data)
	// 写入
	err := c.repo.AddCollInterface(&emCollInterface)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//ctx.JSON(http.StatusOK, model.ResponseData{
	//	Code:    "1",
	//	Message: "添加采集接口成功",
	//})
	return
}

func (c *EmController) UpdateCollInterface(ctx *gin.Context) {
	var addEmCollInterface models.AddEmCollInterface
	var data []byte
	err := ctx.ShouldBindBodyWith(&addEmCollInterface, binding.JSON)
	if err != nil {
		return
	}
	// 查名字获取id
	var emCollInterface models.EmCollInterface
	emCollInterface.OfflinePeriod = addEmCollInterface.OfflinePeriod
	emCollInterface.PollPeriod = addEmCollInterface.PollPeriod
	collInterfaceByName, _ := c.repo.GetCollInterfaceByName(addEmCollInterface.CollInterfaceName)
	if collInterfaceByName == nil {
		return
	}
	emCollInterface.Id = collInterfaceByName.Id
	emCollInterface.Name = addEmCollInterface.CollInterfaceName
	data, _ = json.Marshal(addEmCollInterface)
	emCollInterface.Data = string(data)
	err = c.repo.UpdateCollInterface(&emCollInterface)
	if err != nil {
		return
	}
	return
}

func (c *EmController) DeleteCollInterface(ctx *gin.Context) {
	var tmp struct {
		Name string `json:"name"`
	}
	if err := ctx.ShouldBindBodyWith(&tmp, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	collInterfaceByName, _ := c.repo.GetCollInterfaceByName(tmp.Name)

	if err := c.repo.DeleteCollInterface(collInterfaceByName.Id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	return
}

func (c *EmController) AddEmDevice(ctx *gin.Context) {
	var emDevice models.EmDevice
	var addEmDevice models.AddEmDevice
	var data []byte

	if err := ctx.ShouldBindBodyWith(&addEmDevice, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	emDevice.Name = addEmDevice.Name
	emDevice.Label = addEmDevice.Label
	emDevice.Addr = addEmDevice.Addr
	// 判断是否有重名的设备，重名直接返回
	emDevice.Name = addEmDevice.Name
	emDeviceByName, _ := c.repo.GetEmDeviceByName(emDevice.Name)
	if emDeviceByName != nil {
		//ctx.JSON(http.StatusOK, model.ResponseData{
		//	Code:    "0",
		//	Message: "设备名已存在，添加失败",
		//})
		return
	}
	// 判断是否有对应的采集接口
	emCollInterfaceByName, _ := c.repo.GetCollInterfaceByName(addEmDevice.InterfaceName)
	if emCollInterfaceByName == nil {
		//ctx.JSON(http.StatusOK, model.ResponseData{
		//	Code:    "0",
		//	Message: "采集接口名不存在，添加失败",
		//})
		return
	}
	emDevice.CollInterfaceId = emCollInterfaceByName.Id
	// 判断是否有对应的设备模型
	emDeviceModelByName, _ := c.repo.GetEmDeviceModelByName(addEmDevice.Tsl)
	if emDeviceModelByName == nil {
		//ctx.JSON(http.StatusOK, model.ResponseData{
		//	Code:    "0",
		//	Message: "设备模型名不存在，添加失败",
		//})
		return
	}
	emDevice.ModelId = emDeviceModelByName.Id

	data, _ = json.Marshal(addEmDevice)
	emDevice.Data = string(data)
	err := c.repo.AddEmDevice(&emDevice)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	//ctx.JSON(http.StatusOK, model.ResponseData{
	//	Code:    "1",
	//	Message: "添加设备成功",
	//})
	return
}

func (c *EmController) AddEmDeviceFromXlsx(name string, tsl string, addr string, label string, collInterface string) {
	var data []byte
	var emDevice models.EmDevice
	emDevice.Name = name
	emDevice.Label = label
	emDevice.Addr = addr

	// 通过模型名找模型id
	emDeviceModelByName, err := c.repo.GetEmDeviceModelByName(tsl)
	emDevice.ModelId = emDeviceModelByName.Id
	// 通过采集接口名找接口id
	collInterfaceByName, err := c.repo.GetCollInterfaceByName(collInterface)
	emDevice.CollInterfaceId = collInterfaceByName.Id
	data, err = json.Marshal(emDevice)
	emDevice.Data = string(data)
	err = c.repo.AddEmDevice(&emDevice)
	if err != nil {
		return
	}
	return
}

func (c *EmController) UpdateEmDevice(ctx *gin.Context) {
	var addEmDevice models.AddEmDevice
	var data []byte
	err := ctx.ShouldBindBodyWith(&addEmDevice, binding.JSON)
	if err != nil {
		return
	}
	// 查名字获取id
	var emDevice models.EmDevice
	emDevice.Name = addEmDevice.Name
	emDevice.Label = addEmDevice.Label

	emCollInterfaceByName, _ := c.repo.GetCollInterfaceByName(addEmDevice.InterfaceName)
	emDevice.CollInterfaceId = emCollInterfaceByName.Id
	emDeviceModelByName, _ := c.repo.GetEmDeviceModelByName(addEmDevice.Tsl)
	emDevice.ModelId = emDeviceModelByName.Id
	emDevice.Addr = addEmDevice.Addr

	EmDeviceByName, _ := c.repo.GetEmDeviceByName(addEmDevice.Name)
	if EmDeviceByName == nil {
		return
	}
	emDevice.Id = EmDeviceByName.Id
	data, _ = json.Marshal(addEmDevice)
	emDevice.Data = string(data)
	err = c.repo.UpdateEmDevice(&emDevice)
	if err != nil {
		return
	}
	return
}

func (c *EmController) DeleteEmDevice(ctx *gin.Context) {
	var tmp struct {
		DeviceNames   []string `json:"deviceNames"`
		InterfaceName string   `json:"interface_name"`
	}
	if err := ctx.ShouldBindBodyWith(&tmp, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, name := range tmp.DeviceNames {
		emDeviceByName, _ := c.repo.GetEmDeviceByName(name)
		c.repo.DeleteEmDevice(emDeviceByName.Id)
	}
	return
}

func (c *EmController) AddEmDeviceModel(ctx *gin.Context) {
	var emDeviceModel models.EmDeviceModel
	var addEmDeviceModel models.AddEmDeviceModel
	var data []byte

	if err := ctx.ShouldBindBodyWith(&addEmDeviceModel, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	emDeviceModel.Name = addEmDeviceModel.Name
	emDeviceModel.Label = addEmDeviceModel.Label
	emDeviceModel.Type = addEmDeviceModel.Type
	// 判断设备模型是否有重名，有就直接返回
	emDeviceModel.Name = addEmDeviceModel.Name
	emDeviceModelByName, _ := c.repo.GetEmDeviceModelByName(emDeviceModel.Name)
	if emDeviceModelByName != nil {
		//ctx.JSON(http.StatusOK, model.ResponseData{
		//	Code:    "0",
		//	Message: "设备模型已存在，添加失败",
		//})
		return
	}

	data, _ = json.Marshal(addEmDeviceModel)
	emDeviceModel.Data = string(data)
	err := c.repo.AddEmDeviceModel(&emDeviceModel)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	//ctx.JSON(http.StatusOK, model.ResponseData{
	//	Code:    "1",
	//	Message: "添加设备模型成功",
	//})
	return
}

func (c *EmController) UpdateEmDeviceModel(ctx *gin.Context) {
	var addEmDeviceModel models.AddEmDeviceModel
	var data []byte
	err := ctx.ShouldBindBodyWith(&addEmDeviceModel, binding.JSON)
	if err != nil {
		return
	}
	// 查名字获取id
	var emDeviceModel models.EmDeviceModel
	emDeviceModel.Type = addEmDeviceModel.Type
	emDeviceModel.Name = addEmDeviceModel.Name
	emDeviceModel.Label = addEmDeviceModel.Label
	emDeviceModelByName, _ := c.repo.GetEmDeviceModelByName(addEmDeviceModel.Name)
	if emDeviceModelByName == nil {
		return
	}
	emDeviceModel.Id = emDeviceModelByName.Id
	emDeviceModel.Name = addEmDeviceModel.Name
	data, _ = json.Marshal(addEmDeviceModel)
	emDeviceModel.Data = string(data)
	err = c.repo.UpdateEmDeviceModel(&emDeviceModel)
	if err != nil {
		return
	}
	return
}

func (c *EmController) DeleteEmDeviceModel(ctx *gin.Context) {
	var tmp struct {
		Name string `json:"name"`
	}
	if err := ctx.ShouldBindBodyWith(&tmp, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	emDeviceModelByName, _ := c.repo.GetEmDeviceModelByName(tmp.Name)

	if emDeviceModelByName == nil {
		return
	}

	// 先查是否被device绑定，如果被绑定，不删，直接返回
	deviceList, _ := c.repo.GetEmDeviceByModelId(emDeviceModelByName.Id)
	if len(deviceList) == 0 {
		// 删除模型
		if emDeviceModelByName != nil {
			c.repo.DeleteEmDeviceModel(emDeviceModelByName.Id)
			// 删除cmd
			cmdList, _ := c.repo.GetEmDeviceModelCmdByModelId(emDeviceModelByName.Id)
			for _, cmd := range cmdList {
				c.repo.DeleteEmDeviceModelCmd(cmd.Id)
				// 删除param
				paramList, _ := c.repo.GetEmDeviceModelCmdParamByCmdId(cmd.Id)
				for _, param := range paramList {
					c.repo.DeleteEmDeviceModelCmdParam(param.Id)
				}
			}
		}
	}
	return
}

func (c *EmController) AddEmDeviceModelCmd(ctx *gin.Context) {
	var emDeviceModelCmd models.EmDeviceModelCmd
	var addEmDeviceModelCmd models.AddEmDeviceModelCmd
	var data []byte

	if err := ctx.ShouldBindBodyWith(&addEmDeviceModelCmd, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	emDeviceModelCmd.Name = addEmDeviceModelCmd.Name
	emDeviceModelCmd.Label = addEmDeviceModelCmd.Label
	// 判断设备模型命令是否有重名，有就直接返回
	emDeviceModelCmd.Name = addEmDeviceModelCmd.Name
	emDeviceModelCmdByName, _ := c.repo.GetEmDeviceModelCmdByName(emDeviceModelCmd.Name)
	if emDeviceModelCmdByName != nil {
		//ctx.JSON(http.StatusOK, model.ResponseData{
		//	Code:    "0",
		//	Message: "设备模型命令已存在，添加失败",
		//})
		return
	}
	// 查询对应的模型
	emDeviceModelByName, _ := c.repo.GetEmDeviceModelByName(addEmDeviceModelCmd.TslName)
	if emDeviceModelByName == nil {
		//ctx.JSON(http.StatusOK, model.ResponseData{
		//	Code:    "0",
		//	Message: "设备模型不存在，添加失败",
		//})
		return
	}
	data, _ = json.Marshal(addEmDeviceModelCmd)
	emDeviceModelCmd.Data = string(data)
	emDeviceModelCmd.DeviceModelId = emDeviceModelByName.Id

	err := c.repo.AddEmDeviceModelCmd(&emDeviceModelCmd)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	//ctx.JSON(http.StatusOK, model.ResponseData{
	//	Code:    "1",
	//	Message: "添加设备模型命令成功",
	//})
	return
}

func (c *EmController) AddEmDeviceModelCmdFromXlsx(cmd interface{}, protocol string, tslName string) {
	var data []byte
	switch protocol {
	case "modbus":
		var emDeviceModelCmd models.EmDeviceModelCmd
		tslModbusCmdTemplate := cmd.(device.TSLModbusCmdTemplate)
		emDeviceModelCmd.Name = tslModbusCmdTemplate.Name
		emDeviceModelCmd.Label = tslModbusCmdTemplate.Label
		// 通过模型名找模型id
		emDeviceModelByName, err := c.repo.GetEmDeviceModelByName(tslName)
		emDeviceModelCmd.DeviceModelId = emDeviceModelByName.Id
		data, err = json.Marshal(cmd)
		emDeviceModelCmd.Data = string(data)
		err = c.repo.AddEmDeviceModelCmd(&emDeviceModelCmd)
		if err != nil {
			return
		}
	case "dlt645":
		emDeviceModelCmd := cmd.(device.TSLDLT6452007CmdTemplate)
		fmt.Println(emDeviceModelCmd.Name)
	default:
		return
	}
	return
}

func (c *EmController) UpdateEmDeviceModelCmd(ctx *gin.Context) {
	var addEmDeviceModelCmd models.AddEmDeviceModelCmd
	var data []byte
	err := ctx.ShouldBindBodyWith(&addEmDeviceModelCmd, binding.JSON)
	if err != nil {
		return
	}
	// 查名字获取id
	var emDeviceModelCmd models.EmDeviceModelCmd
	emDeviceModelCmd.Name = addEmDeviceModelCmd.Name
	emDeviceModelCmd.Label = addEmDeviceModelCmd.Label
	emDeviceModelCmdByName, _ := c.repo.GetEmDeviceModelCmdByName(addEmDeviceModelCmd.Name)

	emDeviceModelByName, _ := c.repo.GetEmDeviceModelByName(addEmDeviceModelCmd.TslName)
	emDeviceModelCmd.DeviceModelId = emDeviceModelByName.Id
	emDeviceModelCmd.Id = emDeviceModelCmdByName.Id
	data, _ = json.Marshal(addEmDeviceModelCmd)
	emDeviceModelCmd.Data = string(data)
	err = c.repo.UpdateEmDeviceModelCmd(&emDeviceModelCmd)
	if err != nil {
		return
	}
	return
}

func (c *EmController) DeleteEmDeviceModelCmd(ctx *gin.Context) {
	var tmp = struct {
		TSLName string   `json:"tslName"`
		Names   []string `json:"names"`
	}{}
	if err := ctx.ShouldBindBodyWith(&tmp, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, name := range tmp.Names {
		emDeviceModelCmdByName, _ := c.repo.GetEmDeviceModelCmdByName(name)
		c.repo.DeleteEmDeviceModelCmd(emDeviceModelCmdByName.Id)
		// 删除cmd下的param
		paramList, _ := c.repo.GetEmDeviceModelCmdParamByCmdId(emDeviceModelCmdByName.Id)
		for _, param := range paramList {
			c.repo.DeleteEmDeviceModelCmdParam(param.Id)
		}
	}
	return
}

func (c *EmController) AddEmDeviceModelCmdParam(ctx *gin.Context) {
	var emDeviceModelCmdParam models.EmDeviceModelCmdParam
	var addEmDeviceModelCmdParam models.AddEmDeviceModelCmdParam
	var data []byte

	if err := ctx.ShouldBindBodyWith(&addEmDeviceModelCmdParam, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	emDeviceModelCmdParam.Name = addEmDeviceModelCmdParam.Name
	emDeviceModelCmdParam.Label = addEmDeviceModelCmdParam.Label
	emDeviceModelCmdParam.IotDataType = addEmDeviceModelCmdParam.IotDataType
	// 判断设备模型是否有重名，有就直接返回
	emDeviceModelCmdParam.Name = addEmDeviceModelCmdParam.Name
	emDeviceModelCmdParamByName, _ := c.repo.GetEmDeviceModelCmdParamByName(emDeviceModelCmdParam.Name)
	if emDeviceModelCmdParamByName != nil {
		//ctx.JSON(http.StatusOK, model.ResponseData{
		//	Code:    "0",
		//	Message: "设备模型命令参数，添加失败",
		//})
		return
	}
	// 找对应的cmd
	emDeviceModelCmdByName, _ := c.repo.GetEmDeviceModelCmdByName(addEmDeviceModelCmdParam.CmdName)
	if emDeviceModelCmdByName == nil {
		//ctx.JSON(http.StatusOK, model.ResponseData{
		//	Code:    "0",
		//	Message: "设备模型命令名不存在，添加失败",
		//})
		return
	}
	emDeviceModelCmdParam.DeviceModelCmdId = emDeviceModelCmdByName.Id

	data, _ = json.Marshal(addEmDeviceModelCmdParam)
	emDeviceModelCmdParam.Data = string(data)

	err := c.repo.AddEmDeviceModelCmdParam(&emDeviceModelCmdParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	//ctx.JSON(http.StatusOK, model.ResponseData{
	//	Code:    "1",
	//	Message: "添加设备模型命令参数成功",
	//})
	return
}

func (c *EmController) AddEmDeviceModelCmdParamFromXlsx(property interface{}, protocol string, cmdName string) {
	var data []byte
	switch protocol {
	case "modbus":
		var emDeviceModelCmdParam models.EmDeviceModelCmdParam
		tslModbusPropertyTemplate := property.(device.TSLModbusPropertyTemplate)
		emDeviceModelCmdParam.Name = tslModbusPropertyTemplate.Name
		emDeviceModelCmdParam.Label = tslModbusPropertyTemplate.Label
		// 通过模型名找模型id
		emDeviceModelCmdByName, err := c.repo.GetEmDeviceModelCmdByName(cmdName)
		emDeviceModelCmdParam.DeviceModelCmdId = emDeviceModelCmdByName.Id
		data, err = json.Marshal(property)
		emDeviceModelCmdParam.Data = string(data)
		err = c.repo.AddEmDeviceModelCmdParam(&emDeviceModelCmdParam)
		if err != nil {
			return
		}
	case "dlt645":
		emDeviceModelCmdParam := property.(device.TSLDLT6452007PropertyTemplate)
		fmt.Println(emDeviceModelCmdParam.Name)
	default:
		return
	}
	return
}

func (c *EmController) UpdateEmDeviceModelCmdParam(ctx *gin.Context) {
	var addEmDeviceModelCmdParam models.AddEmDeviceModelCmdParam
	var data []byte
	err := ctx.ShouldBindBodyWith(&addEmDeviceModelCmdParam, binding.JSON)
	if err != nil {
		return
	}
	// 查名字获取id
	var emDeviceModelCmdParam models.EmDeviceModelCmdParam
	emDeviceModelCmdParam.Name = addEmDeviceModelCmdParam.Name
	emDeviceModelCmdParam.Label = addEmDeviceModelCmdParam.Label
	emDeviceModelCmdParam.IotDataType = addEmDeviceModelCmdParam.IotDataType
	emDeviceModelCmdParamByName, _ := c.repo.GetEmDeviceModelCmdParamByName(addEmDeviceModelCmdParam.Name)
	emDeviceModelCmdParam.Id = emDeviceModelCmdParamByName.Id

	emDeviceModelCmdByName, _ := c.repo.GetEmDeviceModelCmdByName(addEmDeviceModelCmdParam.CmdName)
	emDeviceModelCmdParam.DeviceModelCmdId = emDeviceModelCmdByName.Id
	data, _ = json.Marshal(addEmDeviceModelCmdParam)
	emDeviceModelCmdParam.Data = string(data)
	err = c.repo.UpdateEmDeviceModelCmdParam(&emDeviceModelCmdParam)
	if err != nil {
		return
	}
	return
}

func (c *EmController) DeleteEmDeviceModelCmdParam(ctx *gin.Context) {
	var tmp = struct {
		TSLName       string   `json:"tslName"`
		CmdName       string   `json:"cmdName"`
		PropertyNames []string `json:"propertyNames"`
	}{}
	if err := ctx.ShouldBindBodyWith(&tmp, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, name := range tmp.PropertyNames {
		emDeviceModelCmdByNameParam, _ := c.repo.GetEmDeviceModelCmdParamByName(name)
		if emDeviceModelCmdByNameParam == nil {
			return
		}
		c.repo.DeleteEmDeviceModelCmdParam(emDeviceModelCmdByNameParam.Id)
	}
	return
}

// 根据设备名称获取所有模型
func (c *EmController) GetEmDeviceModelCmdParamListByName(ctx *gin.Context) {
	var tmp struct {
		Name string `json:"name"`
	}
	if err := ctx.ShouldBindBodyWith(&tmp, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var deviceCmdParamList []models.EmDeviceModelCmdParam
	deviceCmdParamList, _ = c.repo.GetEmDeviceModelCmdParamListByName(tmp.Name)
	if deviceCmdParamList == nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "无数据",
		})
	}
	var retList []models.AddEmDeviceModelCmdParam
	for i := 0; i < len(deviceCmdParamList); i++ {
		//取出data
		var addEmDeviceModelCmdParam models.AddEmDeviceModelCmdParam
		if err := json.Unmarshal([]byte(deviceCmdParamList[i].Data), &addEmDeviceModelCmdParam); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		retList = append(retList, addEmDeviceModelCmdParam)
	}

	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "成功",
		Data:    retList,
	})
}

// 根据设备ID获取所有模型
func (c *EmController) GetEmDeviceModelCmdParamListByDeviceId(ctx *gin.Context) {
	var tmp struct {
		DeviceId int `json:"deviceId"`
	}
	if err := ctx.ShouldBindBodyWith(&tmp, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var deviceCmdParamList []models.EmDeviceModelCmdParam
	deviceCmdParamList, _ = c.repo.GetEmDeviceModelCmdParamListByDeviceId(tmp.DeviceId)
	if deviceCmdParamList == nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "无数据",
		})
	}
	var retList []models.AddEmDeviceModelCmdParam
	for i := 0; i < len(deviceCmdParamList); i++ {
		//取出data
		var addEmDeviceModelCmdParam models.AddEmDeviceModelCmdParam
		if err := json.Unmarshal([]byte(deviceCmdParamList[i].Data), &addEmDeviceModelCmdParam); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		retList = append(retList, addEmDeviceModelCmdParam)
	}

	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "成功",
		Data:    retList,
	})
}
