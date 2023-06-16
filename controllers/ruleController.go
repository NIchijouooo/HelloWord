package controllers

import (
	"gateway/models"
	"gateway/repositories"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type RuleController struct {
	repo *repositories.RuleRepository
}

func NewRuleController() *RuleController {
	return &RuleController{
		repo: repositories.NewRuleRepository(),
	}
}

func (c *RuleController) RegisterRoutes(router *gin.RouterGroup) {
	//router.POST("/api/v2/em/getAllCommInterfaceProtocols", c.GetAllCommInterfaceProtocols)
	//router.POST("/api/v2/em/addCommInterface", c.AddCommInterface)
	//router.POST("/api/v2/em/communication", c.GetCommInterfaces)
	//router.DELETE("/api/v2/em/delComInterface", c.DelComInterface)
	//router.PUT("/api/v2/em/updateCommInterface", c.UpdateCommInterface)
}

func (c *RuleController) getLimitConfigListByDeviceLabel(ctx *gin.Context) {
	limitConfig := new(models.LimitConfigVo)
	if err := ctx.ShouldBindBodyWith(&limitConfig, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ruleList, err := c.repo.GetRuleByDeviceLabel(limitConfig.DeviceLabel)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(ruleList) == 0 {

	}

}
