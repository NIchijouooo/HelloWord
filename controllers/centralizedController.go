package controllers

import (
	"gateway/httpServer/model"
	"gateway/models"
	"gateway/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义集控管理的控制器
type CentralizedController struct {
	repo *repositories.CentralizedRepository
}

func NewCentralizedController() *CentralizedController {
	return &CentralizedController{repo: repositories.NewCentralizedRepository()}
}

func (ctrl *CentralizedController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/auxiliary/getPolicyList", ctrl.getPolicyList)
	router.POST("/api/v2/auxiliary/createPolicy", ctrl.createPolicy)
	router.POST("/api/v2/auxiliary/updatePolicy", ctrl.updatePolicy)
	router.POST("/api/v2/auxiliary/deletePolicy", ctrl.deletePolicy)
}

// 请求参数
type Param struct {
	Id       int `form:"id"`
	PageNum  int `form:"pageNum"`
	PageSize int `form:"pageSize"`
}

// 获取策略列表
func (c *CentralizedController) getPolicyList(ctx *gin.Context) {
	policyList, err := c.repo.GetPolicyList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"获取信息成功！",
		policyList,
	})
	return
}

// 新增策略配置
func (c *CentralizedController) createPolicy(ctx *gin.Context) {
	var policyData models.EmStrategy
	if err := ctx.ShouldBindJSON(&policyData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	if err := c.repo.CreatePolicy(&policyData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"ok",
		policyData,
	})
}

// 修改策略配置
func (c *CentralizedController) updatePolicy(ctx *gin.Context) {
	var policyData models.EmStrategy
	if err := ctx.ShouldBindJSON(&policyData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	if err := c.repo.UpdatePolicy(&policyData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"ok",
		policyData,
	})
}

// 删除策略配置
func (c *CentralizedController) deletePolicy(ctx *gin.Context) {
	var paramData Param
	if err := ctx.Bind(&paramData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	if err := c.repo.DeletePolicy(paramData.Id); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"ok",
		"",
	})
}
