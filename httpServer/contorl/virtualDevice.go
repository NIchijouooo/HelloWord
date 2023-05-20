package contorl

import (
	"fmt"
	"gateway/httpServer/model"
	"gateway/setting"
	"gateway/virtual"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ApiAddVirtualDevice(context *gin.Context) {
	params := struct {
		DeviceName  string `json:"name"`
		DeviceLabel string `json:"label"`
	}{}

	err := context.BindJSON(&params)
	if err != nil {
		setting.ZAPS.Errorf("添加虚拟设备格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加虚拟设备格式化错误",
			Data:    "",
		})
		return
	}

	_, ok := virtual.VirtualDevice.Nodes[params.DeviceName]
	if ok {
		setting.ZAPS.Error("添加虚拟设备格式化错误,设备名称已经存在")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加虚拟设备格式化错误,设备名称已经存在",
			Data:    "",
		})
		return
	}

	err = virtual.VirtualDevice.VirtualDeviceAddNode(params.DeviceName, params.DeviceLabel)
	if err != nil {
		setting.ZAPS.Errorf("添加虚拟设备错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: fmt.Sprintf("添加虚拟设备错误,%v", err),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: fmt.Sprintf("添加虚拟设备成功"),
		Data:    "",
	})
}

func ApiModifyVirtualDevice(context *gin.Context) {
	params := struct {
		DeviceName  string `json:"name"`
		DeviceLabel string `json:"label"`
	}{}

	err := context.BindJSON(&params)
	if err != nil {
		setting.ZAPS.Errorf("修改虚拟设备格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改虚拟设备格式化错误",
			Data:    "",
		})
		return
	}

	_, ok := virtual.VirtualDevice.Nodes[params.DeviceName]
	if !ok {
		setting.ZAPS.Error("修改虚拟设备格式化错误,设备名称不存在")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改虚拟设备格式化错误,设备名称不存在",
			Data:    "",
		})
		return
	}

	err = virtual.VirtualDevice.VirtualDeviceModifyNode(params.DeviceName, params.DeviceLabel)
	if err != nil {
		setting.ZAPS.Errorf("修改虚拟设备错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: fmt.Sprintf("修改虚拟设备错误,%v", err),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: fmt.Sprintf("修改虚拟设备成功"),
		Data:    "",
	})
}

func ApiDeleteVirtualDevice(context *gin.Context) {
	params := struct {
		DeviceNames []string `json:"names"`
	}{}

	err := context.BindJSON(&params)
	if err != nil {
		setting.ZAPS.Errorf("删除虚拟设备格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除虚拟设备格式化错误",
			Data:    "",
		})
		return
	}

	err = virtual.VirtualDevice.VirtualDeviceDeleteNodes(params.DeviceNames)
	if err != nil {
		setting.ZAPS.Errorf("删除虚拟设备错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: fmt.Sprintf("删除虚拟设备错误,%v", err),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: fmt.Sprintf("删除虚拟设备成功"),
		Data:    "",
	})
}

func ApiGetVirtualDevice(context *gin.Context) {

	type NodeTemplate struct {
		Index int    `json:"index"`
		Name  string `json:"name"`
		Label string `json:"label"`
	}

	nodes := make([]NodeTemplate, 0)
	for _, v := range virtual.VirtualDevice.Nodes {
		node := NodeTemplate{
			Index: v.Index,
			Name:  v.Name,
			Label: v.Label,
		}
		nodes = append(nodes, node)
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: fmt.Sprintf("获取虚拟设备成功"),
		Data:    nodes,
	})
}

func ApiAddVirtualDeviceProperty(context *gin.Context) {

	params := struct {
		DeviceName string                          `json:"deviceName"`
		Property   virtual.VirtualPropertyTemplate `json:"property"`
	}{}

	err := context.BindJSON(&params)
	if err != nil {
		setting.ZAPS.Errorf("添加属性格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加属性格式化错误",
			Data:    "",
		})
		return
	}

	node, ok := virtual.VirtualDevice.Nodes[params.DeviceName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加属性错误,设备名称不存在",
			Data:    "",
		})
		return
	}

	err = node.VirtualDeviceAddProperty(params.Property)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: fmt.Sprintf("添加属性错误,%v", err),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "添加属性成功",
		Data:    "",
	})
}

func ApiModifyVirtualDeviceProperty(context *gin.Context) {

	params := struct {
		DeviceName string                          `json:"deviceName"`
		Property   virtual.VirtualPropertyTemplate `json:"property"`
	}{}

	err := context.BindJSON(&params)
	if err != nil {
		setting.ZAPS.Errorf("修改属性格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改属性格式化错误",
			Data:    "",
		})
		return
	}

	node, ok := virtual.VirtualDevice.Nodes[params.DeviceName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改属性错误,设备名称不存在",
			Data:    "",
		})
		return
	}

	err = node.VirtualDeviceModifyProperty(params.Property)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: fmt.Sprintf("修改属性错误,%v", err.Error()),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "修改属性成功",
		Data:    "",
	})
}

func ApiDeleteVirtualDeviceProperties(context *gin.Context) {

	params := struct {
		DeviceName    string   `json:"deviceName"`
		PropertyNames []string `json:"propertyNames"`
	}{}

	err := context.BindJSON(&params)
	if err != nil {
		setting.ZAPS.Errorf("删除属性格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除属性格式化错误",
			Data:    "",
		})
		return
	}

	node, ok := virtual.VirtualDevice.Nodes[params.DeviceName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除属性错误,设备名称不存在",
			Data:    "",
		})
		return
	}

	err = node.VirtualDeviceDeleteProperties(params.PropertyNames)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: fmt.Sprintf("删除属性错误,%v", err.Error()),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "删除属性成功",
		Data:    "",
	})
}

func ApiGetVirtualDeviceProperties(context *gin.Context) {

	deviceName := context.Query("deviceName")

	node, ok := virtual.VirtualDevice.Nodes[deviceName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取属性错误,设备名称不存在",
			Data:    "",
		})
		return
	}

	properties := make([]virtual.VirtualPropertyTemplate, 0)

	for _, v := range node.Properties {
		properties = append(properties, *v)
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取属性成功",
		Data:    properties,
	})
}
