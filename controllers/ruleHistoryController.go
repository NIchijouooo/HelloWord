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
	router.POST("/api/v2/ruleHistory/getRuleHistoryStatistic", c.getRuleHistoryStatistic)
	router.POST("/api/v2/ruleHistory/updateRuleHistory", c.updateRuleHistory)
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
	list, total, err := c.repo.GetRuleHistoryList(param)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	// 将查询结果返回给客户端
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: gin.H{
			"data":  list,
			"total": total,
		},
	})
}

/*
*
获取历史告警统计数据
*/
func (c *RuleHistoryController) getRuleHistoryStatistic(ctx *gin.Context) {
	var param models.RuleHistoryParam
	if err := ctx.Bind(&param); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	data, err := c.repo.GetRuleHistoryStatistic(param)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	// 将查询结果返回给客户端
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: data,
	})
}

/*
*
更新历史告警
*/
func (c *RuleHistoryController) updateRuleHistory(ctx *gin.Context) {
	var param models.EmRuleHistoryModel
	if err := ctx.Bind(&param); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	id := param.Id
	if id == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "参数错误",
		})
		return
	}
	_, err := c.repo.UpdateRuleHistoryTag(&param)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	// 将查询结果返回给客户端
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: param,
	})
}
