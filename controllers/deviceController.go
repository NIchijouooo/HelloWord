package controllers

import (
	"gateway/httpServer/model"
	"gateway/models"
	repositories "gateway/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
*
设备相关控制器
*/
	repositories "gateway/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeviceController struct {
	repo *repositories.DeviceRepository
}

func NewDeviceController() *DeviceController {
	return &DeviceController{repo: repositories.NewDeviceRepository()}
}

func (c *DeviceController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/device/test", c.GetAllCommInterfaceProtocols)
	router.POST("/api/v2/device/ctrlDevice", c.CtrlDevice)
	router.POST("/api/v2/device/getDeviceListByType", c.getDeviceListByType)
}

// 未完成
func (ctrl *DeviceController) CtrlDevice(ctx *gin.Context) {
	var param models.CtrlInfo
	if err := ctx.Bind(&param); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	device, err := repositories.NewEmRepository().GetEmDeviceById(param.DeviceId)
	if device == nil || err != nil {
		msg := "未查询到设备"
		if err != nil {
			msg = "error" + err.Error()
		}
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: msg,
			Data:    "",
		})
		return
	}

}

/*
*
根据设备类型获取设备列表
deviceType必传
*/
func (c *DeviceController) getDeviceListByType(ctx *gin.Context) {
	var param models.DeviceParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	deviceType := param.DeviceType
	if len(deviceType) == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "参数错误",
			Data:    "",
		})
		return
	}
	list, _ := c.repo.GetDeviceListByType(param)
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: list,
	return &DeviceController{
		repo: repositories.NewDeviceRepository(),
	}
}


func (c *DeviceController) GetAllCommInterfaceProtocols(ctx *gin.Context) {
	c.repo.GetEmDevice()
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "1",
		Message: "成功",
		Data:    1,
	})
}
