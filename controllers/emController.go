package controllers

import (
	"encoding/json"
	"fmt"
	"gateway/models"
	repositories "gateway/repositories"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
)

type EmController struct {
	repo *repositories.EmRepository
}

func NewCommInterfaceProtocolController() *EmController {
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
	router.POST("/api/v2/em/updateCommInterface", c.UpdateCommInterface)
	// 采集接口
	router.POST("/api/v2/em/addCollInterface", c.AddCollInterface)
	// 设备
	router.POST("/api/v2/em/addEmDevice", c.AddEmDevice)
	router.POST("/api/v2/em/updateEmDevice", c.AddEmDevice)
	router.POST("/api/v2/em/delEmDevice", c.AddEmDevice)
	// 设备模型
	router.POST("/api/v2/em/addEmDeviceModel", c.AddEmDeviceModel)
	router.POST("/api/v2/em/addEmDeviceModelCmd", c.AddEmDeviceModelCmd)
	router.POST("/api/v2/em/addEmDeviceModelCmdParam", c.AddEmDeviceModelCmdParam)
}

func (c *EmController) GetAllCommInterfaceProtocols(ctx *gin.Context) {
	rows, err := c.repo.GetAllCommInterfaceProtocols()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, rows)
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
		ctx.JSON(http.StatusOK, "通道名已存在，添加失败")
		return
	}

	protocolByName, _ := c.repo.GetCommInterfaceProtocolByName(commInterface.CommInterfaceProtocol)
	commInterface.CommInterfaceProtocolId = protocolByName.Id
	commInterface.Data = string(data)
	err := c.repo.AddCommInterface(&commInterface)
	if err != nil {
		fmt.Println("error:", err)
	}
	ctx.JSON(http.StatusOK, "添加通道成功")
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
	ctx.JSON(http.StatusOK, res)
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
	ctx.JSON(http.StatusOK, commInterface)
}

// DelComInterface 删除通道
func (c *EmController) DelComInterface(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := c.repo.DelCommInterface(id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "ComInterface deleted successfully"})
}

// AddCollInterface 新增采集接口
func (c *EmController) AddCollInterface(ctx *gin.Context) {
	var addEmCollInterface models.AddEmCollInterface
	var emCollInterface models.EmCollInterface
	var data []byte
	if err := ctx.ShouldBindJSON(&addEmCollInterface); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 判断采集接口是否有重名，有就直接返回
	emCollInterface.Name = addEmCollInterface.CollInterfaceName
	emCollInterfaceByName, _ := c.repo.GetCollInterfaceByName(emCollInterface.Name)
	if emCollInterfaceByName != nil {
		ctx.JSON(http.StatusOK, "通道名已存在，添加失败")
		return
	}

	emCollInterface.OfflinePeriod = addEmCollInterface.OfflinePeriod
	emCollInterface.PollPeriod = addEmCollInterface.PollPeriod
	commInterfaceByName, _ := c.repo.GetCommInterfaceByName(addEmCollInterface.CommInterfaceName)
	if commInterfaceByName == nil {
		ctx.JSON(http.StatusOK, "通信接口名未找到")
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
	ctx.JSON(http.StatusOK, "添加采集接口成功")
	return
}

func (c *EmController) AddEmDevice(ctx *gin.Context) {
	var emDevice models.EmDevice
	var addEmDevice models.AddEmDevice
	var data []byte

	if err := ctx.ShouldBindJSON(&addEmDevice); err != nil {
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
		ctx.JSON(http.StatusOK, "设备名已存在，添加失败")
		return
	}
	// 判断是否有对应的采集接口
	emCollInterfaceByName, _ := c.repo.GetCollInterfaceByName(addEmDevice.InterfaceName)
	if emCollInterfaceByName == nil {
		ctx.JSON(http.StatusOK, "采集接口名不存在，添加失败")
		return
	}
	emDevice.CollInterfaceId = emCollInterfaceByName.Id
	// 判断是否有对应的设备模型
	emDeviceModelByName, _ := c.repo.GetEmDeviceModelByName(addEmDevice.Tsl)
	if emDeviceModelByName == nil {
		ctx.JSON(http.StatusOK, "设备模型名不存在，添加失败")
		return
	}
	emDevice.ModelId = emDeviceModelByName.Id

	data, _ = json.Marshal(addEmDevice)
	emDevice.Data = string(data)
	err := c.repo.AddEmDevice(&emDevice)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, "添加设备成功")
	return
}

func (c *EmController) AddEmDeviceModel(ctx *gin.Context) {
	var emDeviceModel models.EmDeviceModel
	var addEmDeviceModel models.AddEmDeviceModel
	var data []byte

	if err := ctx.ShouldBindJSON(&addEmDeviceModel); err != nil {
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
		ctx.JSON(http.StatusOK, "设备模型已存在，添加失败")
		return
	}

	data, _ = json.Marshal(addEmDeviceModel)
	emDeviceModel.Data = string(data)
	err := c.repo.AddEmDeviceModel(&emDeviceModel)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, "添加设备模型成功")
	return
}

func (c *EmController) AddEmDeviceModelCmd(ctx *gin.Context) {
	var emDeviceModelCmd models.EmDeviceModelCmd
	var addEmDeviceModelCmd models.AddEmDeviceModelCmd
	var data []byte

	if err := ctx.ShouldBindJSON(&addEmDeviceModelCmd); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	emDeviceModelCmd.Name = addEmDeviceModelCmd.Name
	emDeviceModelCmd.Label = addEmDeviceModelCmd.Label
	// 判断设备模型命令是否有重名，有就直接返回
	emDeviceModelCmd.Name = addEmDeviceModelCmd.Name
	emDeviceModelCmdByName, _ := c.repo.GetEmDeviceModelCmdByName(emDeviceModelCmd.Name)
	if emDeviceModelCmdByName != nil {
		ctx.JSON(http.StatusOK, "设备模型命令已存在，添加失败")
		return
	}
	// 查询对应的模型
	emDeviceModelByName, _ := c.repo.GetEmDeviceModelByName(addEmDeviceModelCmd.TslName)
	if emDeviceModelByName == nil {
		ctx.JSON(http.StatusOK, "设备模型不存在，添加失败")
		return
	}
	data, _ = json.Marshal(addEmDeviceModelCmd)
	emDeviceModelCmd.Data = string(data)
	emDeviceModelCmd.DeviceModelId = emDeviceModelByName.Id

	err := c.repo.AddEmDeviceModelCmd(&emDeviceModelCmd)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, "添加设备模型命令成功")
	return
}

func (c *EmController) AddEmDeviceModelCmdParam(ctx *gin.Context) {
	var emDeviceModelCmdParam models.EmDeviceModelCmdParam
	var addEmDeviceModelCmdParam models.AddEmDeviceModelCmdParam
	var data []byte

	if err := ctx.ShouldBindJSON(&addEmDeviceModelCmdParam); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	emDeviceModelCmdParam.Name = addEmDeviceModelCmdParam.Name
	emDeviceModelCmdParam.Label = addEmDeviceModelCmdParam.Label
	// 判断设备模型是否有重名，有就直接返回
	emDeviceModelCmdParam.Name = addEmDeviceModelCmdParam.Name
	emDeviceModelCmdParamByName, _ := c.repo.GetEmDeviceModelCmdParamByName(emDeviceModelCmdParam.Name)
	if emDeviceModelCmdParamByName != nil {
		ctx.JSON(http.StatusOK, "设备模型命令参数，添加失败")
		return
	}
	// 找对应的cmd
	emDeviceModelCmdByName, _ := c.repo.GetEmDeviceModelCmdByName(addEmDeviceModelCmdParam.CmdName)
	if emDeviceModelCmdByName == nil {
		ctx.JSON(http.StatusOK, "设备模型命令名不存在，添加失败")
		return
	}
	emDeviceModelCmdParam.DeviceModelCmdId = emDeviceModelCmdByName.Id

	data, _ = json.Marshal(addEmDeviceModelCmdParam)
	emDeviceModelCmdParam.Data = string(data)

	err := c.repo.AddEmDeviceModelCmdParam(&emDeviceModelCmdParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, "添加设备模型命令参数成功")
	return
}
