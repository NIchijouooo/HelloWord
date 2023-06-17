package controllers

/**
20230605
*/
import (
	"fmt"
	"gateway/httpServer/model"
	"gateway/models"
	"gateway/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义字典类型管理的控制器
type RealtimeDataController struct {
	repo      *repositories.RealtimeDataRepository
	repoHis   *repositories.HistoryDataRepository
	repoPoint *repositories.DevicePointRepository
}

func NewRealtimeDataController() *RealtimeDataController {
	return &RealtimeDataController{
		repo: repositories.NewRealtimeDataRepository(),
		repoHis:   repositories.NewHistoryDataRepository(),
		repoPoint: repositories.NewDevicePointRepository()}
}

func (ctrl *RealtimeDataController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/realtimeData/GetRealtimeDataYxByID", ctrl.GetRealtimeDataYxByID)
	router.POST("/api/v2/realtimeData/GetRealtimeDataYcByID", ctrl.GetRealtimeDataYcByID)
	router.POST("/api/v2/realtimeData/GetRealtimeDataSettingByID", ctrl.GetRealtimeDataSettingByID)
	router.POST("/api/v2/realtimeData/GetRealtimeDataYxListByID", ctrl.GetRealtimeDataYxListByID)
	router.POST("/api/v2/realtimeData/GetRealtimeDataYcListByID", ctrl.GetRealtimeDataYcListByID)
	router.POST("/api/v2/realtimeData/GetRealtimeDataSettingListByID", ctrl.GetRealtimeDataSettingListByID)
	router.POST("/api/v2/realtimeData/GetPointsByDeviceId", ctrl.GetPointsByDeviceId)
	router.POST("/api/v2/realtimeData/GetDeviceByDevLabel", ctrl.GetDeviceByDevLabel)
	// 注册其他路由...
}

type ParamRealtimeData struct {
	ID         uint64  `form:"id"`
	Code       int     `form:"code"`
	Codes      string  `form:"codes"`
	DeviceName string  `form:"deviceName"`
	DeviceId   int     `form:"DeviceId"`
	DeviceIds  string  `form:"deviceIds"`
	Value      float64 `form:"value"`
	ChannelId  uint64  `form:"channelId"`
	Type       string  `form:"type"`
	Date       string  `form:"date"`
	StartTime  int64   `form:"startTime"`
	EndTime    int64   `form:"endTime"`
	Interval   string  `form:"interval"`
	PointType  string  `form:"pointType"`
	Lable      string  `form:"lable"`
}

/*
*
根据设备id，点位类型获取命令参数属性
*/
func (c *RealtimeDataController) GetPointsByDeviceId(ctx *gin.Context) {
	var realtimeData ParamRealtimeData
	if err := ctx.Bind(&realtimeData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}

	var yx []*models.EmDeviceModelCmdParam
	yx = c.repoPoint.GetPointsByDeviceId(realtimeData.PointType, realtimeData.DeviceId, realtimeData.Code)
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		yx,
	})
}

/*
*
根据设备id，点位类型获取命令参数属性
*/
func (c *RealtimeDataController) GetDeviceByDevLabel(ctx *gin.Context) {
	var realtimeData ParamRealtimeData
	if err := ctx.Bind(&realtimeData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}

	var yx []*models.EmDevice
	yx = c.repoPoint.GetDeviceByDevLabel(realtimeData.Lable)
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		yx,
	})
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
	yx, err := c.repo.GetYxById(realtimeData.DeviceId, realtimeData.Code)
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

	var yxList []*models.PointParam
	yxList, err := c.repoHis.GetYxLogByDeviceIdsCodes(realtimeData.DeviceIds, realtimeData.Codes, realtimeData.Interval, realtimeData.StartTime, realtimeData.EndTime)
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

	var yxList []*models.PointParam
	yxList, err := c.repoHis.GetYcLogByDeviceIdsCodes(realtimeData.DeviceIds, realtimeData.Codes, realtimeData.Interval, realtimeData.StartTime, realtimeData.EndTime)
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
	yxList, err := c.repoHis.GetSettingLogByDeviceIdsCodes(realtimeData.DeviceIds, realtimeData.Codes, realtimeData.Interval, realtimeData.StartTime, realtimeData.EndTime)
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
	yx, err := c.repo.GetYcById(realtimeData.DeviceId, realtimeData.Code)
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
	yx, err := c.repo.GetSettingById(realtimeData.DeviceId, realtimeData.Code)
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
