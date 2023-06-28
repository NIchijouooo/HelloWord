package controllers

import (
	"fmt"
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
	router.POST("/api/v2/em/saveLimitConfig", c.SaveLimitConfig)
	router.POST("/api/v2/em/deleteLimitConfig", c.DeleteLimitConfig)
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
		//20230628同一省份同一月不能重复,查id不同的其他数据，
		configurationList, err := c.repo.GetLimitConfigListCheckById(limitConfig.Id, limitConfig.PropertyCode)
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
				"不能重复添加该点位数据",
				"",
			})
			return
		}

		updateFlag, err := c.repo.UpdateLimitConfig(limitConfig)
		if err != nil {
			ctx.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "error" + err.Error(),
				Data:    "",
			})
			return
		}
		if updateFlag > 0 {
			saveLimitRule(limitConfig, false)
		}
	} else {
		//20230628同一点位不能重复
		configurationList, err := c.repo.GetLimitConfigListList(limitConfig.PropertyCode)
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
				"不能重复添加该点位数据",
				"",
			})
			return
		}

		limitConfig.CreateTime = time.Now().Format(time.DateTime)
		insertFlag, err := c.repo.InsertLimitConfig(limitConfig)
		if err != nil {
			ctx.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "error" + err.Error(),
				Data:    "",
			})
			return
		}
		if insertFlag > 0 {
			saveLimitRule(limitConfig, true)
		}
	}

	_, err := c.repo.UpdateLimitConfig(limitConfig)
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

	deleteFlag, err := c.repo.UpdateLimitConfig(updateLimitConfig)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}

	if deleteFlag > 0 {
		notifyRuleId := updateLimitConfig.NotifyRuleId
		secondaryRuleId := updateLimitConfig.SecondaryRuleId
		seriousRuleId := updateLimitConfig.SeriousRuleId
		urgentRuleId := updateLimitConfig.UrgentRuleId
		ruleIdList := []int{notifyRuleId, secondaryRuleId, seriousRuleId, urgentRuleId}
		_, _ = repositories.NewRuleRepository().DeleteRuleList(ruleIdList)
	}

	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: "1",
	})
}

func saveLimitRule(limitConfig *models.LimitConfigVo, isInsert bool) {

	ruleList := make([]*models.EmRuleModel, 0)

	deviceType := limitConfig.DeviceType
	propertyCode := limitConfig.PropertyCode

	notifyMin := limitConfig.NotifyMin
	notifyMax := limitConfig.NotifyMax
	notifyRuleId := limitConfig.NotifyRuleId
	ruleList = processRuleData(ruleList, notifyRuleId, deviceType, propertyCode, 0, notifyMin, notifyMax)

	secondaryMin := limitConfig.SecondaryMin
	secondaryMax := limitConfig.SecondaryMax
	secondaryRuleId := limitConfig.SecondaryRuleId
	ruleList = processRuleData(ruleList, secondaryRuleId, deviceType, propertyCode, 1, secondaryMin, secondaryMax)

	seriousMin := limitConfig.SeriousMin
	seriousMax := limitConfig.SeriousMax
	seriousRuleId := limitConfig.SeriousRuleId
	ruleList = processRuleData(ruleList, seriousRuleId, deviceType, propertyCode, 2, seriousMin, seriousMax)

	urgentMin := limitConfig.UrgentMin
	urgentMax := limitConfig.UrgentMax
	urgentRuleId := limitConfig.UrgentRuleId
	ruleList = processRuleData(ruleList, urgentRuleId, deviceType, propertyCode, 3, urgentMin, urgentMax)

	if len(ruleList) > 0 {
		updateFlag := 0
		if isInsert {
			updateFlag, _ = repositories.NewRuleRepository().InsertRuleList(ruleList)
		} else {
			updateFlag, _ = repositories.NewRuleRepository().UpdateRuleList(ruleList)
		}
		if updateFlag > 0 {
			for _, rule := range ruleList {
				id := rule.Id
				level := rule.Level
				switch level {
				case 0:
					limitConfig.NotifyRuleId = id
				case 1:
					limitConfig.SecondaryRuleId = id
				case 2:
					limitConfig.SeriousRuleId = id
				case 3:
					limitConfig.UrgentRuleId = id
				default:
					break
				}
			}
		}
	}
}

func processRuleData(ruleList []*models.EmRuleModel, ruleId int, deviceType string, propertyCode string,
	level int, min string, max string) []*models.EmRuleModel {
	notifyMinCondition := fmt.Sprintf("product.${%s:\"XXX\"}.${%s:\"XXX\"}<%s", deviceType, propertyCode, min)
	notifyMaxCondition := fmt.Sprintf("product.${%s:\"XXX\"}.${%s:\"XXX\"}>%s", deviceType, propertyCode, max)
	content := fmt.Sprintf("type=logic\\nif\\n%s&&%s\\nthen\\nrecord.recover.notMeet", notifyMinCondition, notifyMaxCondition)
	notifyRule := models.EmRuleModel{}
	notifyRule.Content = content
	notifyRule.Name = "越限报警"
	notifyRule.TypeClassification = 1
	notifyRule.EnableFlag = 1
	notifyRule.Level = level
	notifyRule.Description = "越限报警"
	notifyRule.CreateTime = time.Now().Format(time.DateTime)
	notifyRule.UpdateTime = time.Now().Format(time.DateTime)
	if ruleId > 0 {
		notifyRule.Id = ruleId
	}

	ruleList = append(ruleList, &notifyRule)

	return ruleList
}
