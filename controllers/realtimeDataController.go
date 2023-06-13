package controllers

/**
20230605
*/
import (
	"fmt"
	"gateway/httpServer/model"
	"gateway/models"
	repositories "gateway/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 定义字典类型管理的控制器
type RealtimeDataController struct {
	repo    *repositories.RealtimeDataRepository
	repoHis *repositories.HistoryDataRepository
}

func NewRealtimeDataController() *RealtimeDataController {
	return &RealtimeDataController{repo: repositories.NewRealtimeDataRepository(), repoHis: repositories.NewHistoryDataRepository()}
}

func (ctrl *RealtimeDataController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/realtimeData/GetRealtimeDataYxByID", ctrl.GetRealtimeDataYxByID)
	router.POST("/api/v2/realtimeData/GetRealtimeDataYcByID", ctrl.GetRealtimeDataYcByID)
	router.POST("/api/v2/realtimeData/GetRealtimeDataSettingByID", ctrl.GetRealtimeDataSettingByID)
	router.POST("/api/v2/realtimeData/GetRealtimeDataYxListByID", ctrl.GetRealtimeDataYxListByID)
	router.POST("/api/v2/realtimeData/GetRealtimeDataYcListByID", ctrl.GetRealtimeDataYcListByID)
	router.POST("/api/v2/realtimeData/GetRealtimeDataSettingListByID", ctrl.GetRealtimeDataSettingListByID)
	// 注册其他路由...
}

type ParamRealtimeData struct {
	ID         uint64  `form:"id"`
	Code       int     `form:"code"`
	Codes      string  `form:"codes"`
	DeviceName string  `form:"deviceName"`
	DeviceID   int     `form:"deviceID"`
	Value      float64 `form:"value"`
	ChannelID  uint64  `form:"channelID"`
	Type       string  `form:"type"`
	Date       string  `form:"date"`
	StartTime  int64   `form:"startTime"`
	EndTime    int64   `form:"endTime"`
}

func (c *RealtimeDataController) GetRealtimeDataYxByID(ctx *gin.Context) {
	var realtimeData ParamRealtimeData
	if err := ctx.Bind(&realtimeData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}

	var yx models.YxData
	yx, err := c.repo.GetYxById(realtimeData.DeviceID, realtimeData.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		yx,
	})
}

func (c *RealtimeDataController) GetRealtimeDataYxListByID(ctx *gin.Context) {
	var realtimeData ParamRealtimeData
	if err := ctx.Bind(&realtimeData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	fmt.Println(realtimeData.StartTime)
	fmt.Println(realtimeData.EndTime)

	var yxList []*models.YxData
	yxList, err := c.repoHis.GetYxLogById(realtimeData.DeviceID, realtimeData.Codes, realtimeData.StartTime, realtimeData.EndTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		yxList,
	})
}

func (c *RealtimeDataController) GetRealtimeDataYcListByID(ctx *gin.Context) {
	var realtimeData ParamRealtimeData
	if err := ctx.Bind(&realtimeData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	fmt.Println(realtimeData.StartTime)
	fmt.Println(realtimeData.EndTime)

	var yxList []*models.YcData
	yxList, err := c.repoHis.GetYcLogById(realtimeData.DeviceID, realtimeData.Codes, realtimeData.StartTime, realtimeData.EndTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		yxList,
	})
}
func (c *RealtimeDataController) GetRealtimeDataSettingListByID(ctx *gin.Context) {
	var realtimeData ParamRealtimeData
	if err := ctx.Bind(&realtimeData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	fmt.Println(realtimeData.StartTime)
	fmt.Println(realtimeData.EndTime)

	var yxList []*models.SettingData
	yxList, err := c.repoHis.GetSettingLogById(realtimeData.DeviceID, realtimeData.Codes, realtimeData.StartTime, realtimeData.EndTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		yxList,
	})
}
func (c *RealtimeDataController) GetRealtimeDataYcByID(ctx *gin.Context) {
	var realtimeData ParamRealtimeData
	if err := ctx.Bind(&realtimeData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}

	var yx models.YcData
	yx, err := c.repo.GetYcById(realtimeData.DeviceID, realtimeData.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		yx,
	})
}
func (c *RealtimeDataController) GetRealtimeDataSettingByID(ctx *gin.Context) {
	var realtimeData ParamRealtimeData
	if err := ctx.Bind(&realtimeData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}

	var yx models.SettingData
	yx, err := c.repo.GetSettingById(realtimeData.DeviceID, realtimeData.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		yx,
	})
}
