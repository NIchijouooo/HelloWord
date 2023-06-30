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
	router.POST("/api/v2/webHmiPage/getIotWebHmiPageInfo", c.GetIotWebHmiPageInfo)
}

func (c *WebHmiPageController) GetIotWebHmiPageInfo(ctx *gin.Context) {

	webHmiPageDeviceModel := models.WebHmiPageDeviceModel{}

	if err := ctx.ShouldBindBodyWith(&webHmiPageDeviceModel, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if webHmiPageDeviceModel.WebHmiPageCode == "" && webHmiPageDeviceModel.DeviceId == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "参数错误",
			Data:    "",
		})
		return
	}
	webHmiPageId, token, err := c.repo.GetIotWebHmiPageInfo(webHmiPageDeviceModel)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	result := map[string]interface{}{
		"webHmiPageId": webHmiPageId,
		"token":        token,
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: result,
	})
}
