package controllers

import (
	"gateway/httpServer/model"
	"gateway/models"
	"gateway/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 定义配置中心的控制器
type ConfigurationCenterController struct {
	repo *repositories.ConfigurationCenterRepository
}

func NewConfigurationCenterController() *ConfigurationCenterController {
	return &ConfigurationCenterController{repo: repositories.NewConfigurationCenterRepository()}
}

func (c *ConfigurationCenterController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/configurationCenter/addConfiguration", c.AddConfiguration)
	router.POST("/api/v2/configurationCenter/updateConfiguration", c.UpdateConfiguration)
	router.POST("/api/v2/configurationCenter/getConfigurationList", c.GetConfigurationList)
	router.POST("/api/v2/configurationCenter/deleteConfiguration", c.DeleteConfiguration)
}

type ConfigurationParamType struct {
	Month    string `form:"month"`
	Province string `form:"province"`
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize"`
}

// 新增电费配置
func (c *ConfigurationCenterController) AddConfiguration(ctx *gin.Context) {
	var configuration models.EmConfiguration
	if err := ctx.ShouldBindJSON(&configuration); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"0",
			"error" + err.Error(),
			"",
		})
		return
	}
	//20230628同一省份同一月不能重复
	configurationList, _, err := c.repo.GetConfigurationList(configuration.Province, configuration.Month, 1, 10)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}

	if len(configurationList) > 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"不能重复添加该省份当月数据",
			"",
		})
		return
	}

	if err := c.repo.AddConfiguration(&configuration); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"0",
			"error" + err.Error(),
			"",
		})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		configuration,
	})
}

// 更新电费配置
func (c *ConfigurationCenterController) UpdateConfiguration(ctx *gin.Context) {
	var configuration models.EmConfiguration
	if err := ctx.ShouldBindJSON(&configuration); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"0",
			"error" + err.Error(),
			"",
		})
		return
	}
	//20230628同一省份同一月不能重复,查id不同的其他数据，
	configurationList, err := c.repo.GetConfigurationCheckById(configuration.Id, configuration.Province, configuration.Month)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}

	if len(configurationList) > 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"不能重复添加该省份当月数据",
			"",
		})
		return
	}

	if err := c.repo.UpdateConfiguration(&configuration); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"0",
			"error" + err.Error(),
			"",
		})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		configuration,
	})
}

// 获取电费配置列表
func (c *ConfigurationCenterController) GetConfigurationList(ctx *gin.Context) {
	var configurationList []models.EmConfiguration
	var paramType ConfigurationParamType
	if err := ctx.Bind(&paramType); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}

	configurationList, total, err := c.repo.GetConfigurationList(paramType.Province, paramType.Month, paramType.PageNum, paramType.PageSize)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}

	// 将查询结果返回给客户端
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		gin.H{
			"data":  configurationList,
			"total": total,
		},
	})
}

// 删除电费配置
func (c *ConfigurationCenterController) DeleteConfiguration(ctx *gin.Context) {
	var configuration models.EmConfiguration
	if err := ctx.Bind(&configuration); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	if err := c.repo.DeleteConfiguration(configuration.Id); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"0",
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
