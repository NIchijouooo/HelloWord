package controllers

import (
	"gateway/httpServer/model"
	"gateway/models"
	"gateway/repositories"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"time"
)

type LimitConfigController struct {
	repo *repositories.LimitConfigRepository
}

func NewLimitConfigController() *LimitConfigController {
	return &LimitConfigController{
		repo: repositories.NewLimitConfigRepository(),
	}
}

func (c *LimitConfigController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/em/getLimitConfigListByDeviceType", c.GetLimitConfigListByDeviceType)
	router.POST("/api/v2/em/saveLimitConfigList", c.SaveLimitConfig)
	router.POST("/api/v2/em/deleteLimitConfigList", c.DeleteLimitConfig)
}

func (c *LimitConfigController) GetLimitConfigListByDeviceType(ctx *gin.Context) {
	limitConfig := new(models.LimitConfigVo)
	if err := ctx.ShouldBindBodyWith(&limitConfig, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ruleModelJsonList, err := c.repo.GetLimitConfigListByDeviceType(limitConfig.DeviceType)
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
		Data: ruleModelJsonList,
	})
}

func (c *LimitConfigController) SaveLimitConfig(ctx *gin.Context) {
	limitConfig := new(models.LimitConfigVo)
	if err := ctx.ShouldBindBodyWith(&limitConfig, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := limitConfig.Id
	limitConfig.UpdateTime = time.Now().Format(time.DateTime)
	if id > 0 {
		_, err := c.repo.UpdateLimitConfig(limitConfig)
		if err != nil {
			ctx.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "error" + err.Error(),
				Data:    "",
			})
			return
		}
	} else {
		limitConfig.CreateTime = time.Now().Format(time.DateTime)
		_, err := c.repo.InsertLimitConfig(limitConfig)
		if err != nil {
			ctx.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "error" + err.Error(),
				Data:    "",
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: "1",
	})
}

func (c *LimitConfigController) DeleteLimitConfig(ctx *gin.Context) {
	limitConfig := new(models.LimitConfigVo)
	if err := ctx.ShouldBindBodyWith(&limitConfig, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := limitConfig.Id

	updateLimitConfig := new(models.LimitConfigVo)
	updateLimitConfig.Id = id
	updateLimitConfig.DelFlag = 1
	updateLimitConfig.UpdateTime = time.Now().Format(time.DateTime)

	_, err := c.repo.UpdateLimitConfig(updateLimitConfig)
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
		Data: "1",
	})
}
