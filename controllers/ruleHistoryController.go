package controllers

import (
	"gateway/httpServer/model"
	"gateway/models"
	"gateway/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RuleHistoryController struct {
	repo *repositories.RuleHistoryRepository
}

func NewRuleHistoryController() *RuleHistoryController {
	return &RuleHistoryController{
		repo: repositories.NewRuleHistoryRepository(),
	}
}

func (c *RuleHistoryController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/ruleHistory/getRuleHistoryList", c.getRuleHistoryList)
	//router.POST("/api/v2/em/addCommInterface", c.AddCommInterface)
	//router.POST("/api/v2/em/communication", c.GetCommInterfaces)
	//router.DELETE("/api/v2/em/delComInterface", c.DelComInterface)
	//router.PUT("/api/v2/em/updateCommInterface", c.UpdateCommInterface)
}

func (c *RuleHistoryController) getRuleHistoryList(ctx *gin.Context) {
	var param models.RuleHistoryParam
	if err := ctx.Bind(&param); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	deviceId := param.DeviceId
	if deviceId == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "参数错误",
			Data:    "",
		})
		return
	}
	c.repo.GetRuleHistoryList(param)
}
