package controllers

import (
	"gateway/httpServer/model"
	"gateway/models"
	"gateway/repositories"
	"gateway/service/job"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

// 定义集控管理的控制器
type CentralizedController struct {
	repo       *repositories.CentralizedRepository
	deviceRepo *repositories.DeviceRepository
}

func NewCentralizedController() *CentralizedController {
	return &CentralizedController{
		repo:       repositories.NewCentralizedRepository(),
		deviceRepo: repositories.NewDeviceRepository()}
}

func (ctrl *CentralizedController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/centralized/getPolicyList", ctrl.getPolicyList)
	router.POST("/api/v2/centralized/createPolicy", ctrl.createPolicy)
	router.POST("/api/v2/centralized/updatePolicy", ctrl.updatePolicy)
	router.POST("/api/v2/centralized/deletePolicy", ctrl.deletePolicy)
	router.POST("/api/v2/centralized/getDeviceYkYtList", ctrl.getDeviceYkYtList)
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

	// 按照日期时间字符串进行倒序排序
	sort.Slice(policyList, func(i, j int) bool {
		t1, err := time.Parse("2006-01-02 15:04:05", policyList[i].UpdateTime)
		if err != nil {
			// 处理解析错误
			return false
		}

		t2, err := time.Parse("2006-01-02 15:04:05", policyList[j].UpdateTime)
		if err != nil {
			// 处理解析错误
			return false
		}

		return t1.After(t2)
	})

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
	policyData.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	policyData.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
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
	policyData.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	oldData, err := c.repo.UpdatePolicy(&policyData)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	// 原来开启,现在不开,更新成功后拿原来配置的时间重置功率
	if oldData.Status == 1 && policyData.Status == 0 {
		job.ResetPower(oldData)
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
	oldData, err := c.repo.DeletePolicy(paramData.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	// 原来开启,现在删除,删除成功后拿原来配置的时间重置功率
	if oldData.Status == 1 {
		job.ResetPower(oldData)
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"ok",
		"",
	})
}

// 获取实时控制遥控遥调列表
func (c *CentralizedController) getDeviceYkYtList(ctx *gin.Context) {
	YkYtList, err := c.repo.GetDeviceYkYtList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"获取信息成功！",
		YkYtList,
	})
	return
}
