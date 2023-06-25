package controllers

import (
	"gateway/httpServer/model"
	"gateway/models"
	"gateway/repositories"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type WebHmiPageController struct {
	repo *repositories.WebHmiPageRepository
}

func NewWebHmiPageController() *WebHmiPageController {
	return &WebHmiPageController{
		repo: repositories.NewWebHmiPageRepository(),
	}
}

func (c *WebHmiPageController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/webHmiPage/getWebHmiPageDeviceInfo", c.GetWebHmiPageDeviceInfo)
	router.POST("/api/v2/webHmiPage/saveWebHmiPageDeviceInfo", c.SaveWebHmiPageDeviceInfo)
}

func (c *WebHmiPageController) GetWebHmiPageDeviceInfo(ctx *gin.Context) {
	webHmiPageDeviceModel := models.WebHmiPageDeviceModel{}
	if err := ctx.ShouldBindBodyWith(&webHmiPageDeviceModel, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deviceId := webHmiPageDeviceModel.DeviceId

	webHmiPageDeviceModelList, err := c.repo.GetWebHmiPageDeviceInfo(deviceId)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: webHmiPageDeviceModelList,
	})
}

func (c *WebHmiPageController) SaveWebHmiPageDeviceInfo(ctx *gin.Context) {

	webHmiPageDeviceModelList := make([]*models.WebHmiPageDeviceModel, 0)

	if err := ctx.ShouldBindBodyWith(&webHmiPageDeviceModelList, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	saveFlag, err := c.repo.SaveWebHmiPageDeviceInfo(webHmiPageDeviceModelList)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: saveFlag,
	})
}
