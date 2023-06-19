package controllers

import (
	"gateway/httpServer/model"
	repositories "gateway/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeviceController struct {
	repo *repositories.DeviceRepository
}

func NewDeviceController() *DeviceController {
	return &DeviceController{
		repo: repositories.NewDeviceRepository(),
	}
}

func (c *DeviceController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/device/test", c.GetAllCommInterfaceProtocols)
}

func (c *DeviceController) GetAllCommInterfaceProtocols(ctx *gin.Context) {
	c.repo.GetEmDevice()
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "1",
		Message: "成功",
		Data:    1,
	})
}
