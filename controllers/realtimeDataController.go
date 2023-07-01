package controllers

/**
20230605
*/
import (
	"fmt"
	"gateway/httpServer/model"
	"gateway/models"
	"gateway/repositories"
	"gateway/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义字典类型管理的控制器
type RealtimeDataController struct {
	repo      *repositories.RealtimeDataRepository
	repoHis   *repositories.HistoryDataRepository
	repoPoint *repositories.DevicePointRepository
	repoData *repositories.DictDataRepository
	repoEm *repositories.EmRepository
}

func NewRealtimeDataController() *RealtimeDataController {
	return &RealtimeDataController{
		repo:      repositories.NewRealtimeDataRepository(),
		repoHis:   repositories.NewHistoryDataRepository(),
		repoPoint: repositories.NewDevicePointRepository(),
		repoData: repositories.NewDictDataRepository(),
		repoEm: repositories.NewEmRepository()}
}

func (ctrl *RealtimeDataController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/realtimeData/GetRealtimeDataYxByID", ctrl.GetRealtimeDataYxByID)
	router.POST("/api/v2/realtimeData/GetRealtimeDataYcByID", ctrl.GetRealtimeDataYcByID)
	router.POST("/api/v2/realtimeData/GetRealtimeDataSettingByID", ctrl.GetRealtimeDataSettingByID)
	router.POST("/api/v2/realtimeData/GetRealtimeDataYxListByID", ctrl.GetRealtimeDataYxListByID)
	router.POST("/api/v2/realtimeData/getRealtimeDataYxListByDevId", ctrl.GetRealtimeDataYxListByDevId)
	router.POST("/api/v2/realtimeData/getRealtimeDataYxListByIdOrCodeList", ctrl.GetRealtimeDataYxListByIdOrCodeList)
	router.POST("/api/v2/realtimeData/GetRealtimeDataYcListByID", ctrl.GetRealtimeDataYcListByID)
	router.POST("/api/v2/realtimeData/getRealtimeDataYcListByDevId", ctrl.GetRealtimeDataYcListByDevId)
	router.POST("/api/v2/realtimeData/getRealtimeDataYcListByIdOrCodeList", ctrl.GetRealtimeDataYcListByIdOrCodeList)
	router.POST("/api/v2/realtimeData/GetRealtimeDataSettingListByID", ctrl.GetRealtimeDataSettingListByID)
	router.POST("/api/v2/realtimeData/GetPointsByDeviceId", ctrl.GetPointsByDeviceId)
	router.POST("/api/v2/realtimeData/GetDeviceByDevLabel", ctrl.GetDeviceByDevLabel)
	router.POST("/api/v2/realtimeData/GetDeviceByDeviceType", ctrl.GetDeviceByDeviceType)
	router.POST("/api/v2/realtimeData/GetRealtimeDataListByDevId", ctrl.GetRealtimeDataListByDevId)
	router.POST("/api/v2/realtimeData/GetProfitPowerRate", ctrl.GetProfitPowerRate)
	// 注册其他路由...
}

type ParamRealtimeData struct {
	ID         uint64  `form:"id"`
	Code       int     `form:"code"`
	Codes      string  `form:"codes"`
	DeviceName string  `form:"deviceName"`
	DeviceType string  `form:"deviceType"`
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
	ProjectId      string  `form:"projectId"`
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
根据设备类型获取命令参数属性
*/
func (c *RealtimeDataController) GetDeviceByDeviceType(ctx *gin.Context) {
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
	yx = c.repoPoint.GetDeviceByDeviceType(realtimeData.DeviceType)
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

/*
*
根据设备id集合和遥信编码集合获取遥信列表
*/
func (c *RealtimeDataController) GetRealtimeDataYxListByIdOrCodeList(ctx *gin.Context) {
	var realtimeData ParamRealtimeData
	if err := ctx.Bind(&realtimeData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}

	yxList, err := c.repo.GetYxListByDevIdsAndCodes(realtimeData.DeviceIds, realtimeData.Codes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: yxList,
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
	var ycList []*models.PointParam
	ycList, err := c.repoHis.GetYcLogByDeviceIdsCodes(realtimeData.DeviceIds, realtimeData.Codes, realtimeData.Interval, realtimeData.StartTime, realtimeData.EndTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		ycList,
	})
}

/*
*
根据设备id集合和遥信编码集合获取遥信列表
*/
func (c *RealtimeDataController) GetRealtimeDataYcListByIdOrCodeList(ctx *gin.Context) {
	var realtimeData ParamRealtimeData
	if err := ctx.Bind(&realtimeData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}

	ycList, err := c.repo.GetYcListByDevIdsAndCodes(realtimeData.DeviceIds, realtimeData.Codes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: ycList,
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

/*
*
获取设备全部遥信实时数据
*/
func (c *RealtimeDataController) GetRealtimeDataYxListByDevId(ctx *gin.Context) {
	var realtimeData ParamRealtimeData
	if err := ctx.Bind(&realtimeData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	deviceId := realtimeData.DeviceId
	if deviceId == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "参数错误",
			Data:    "",
		})
		return
	}
	yxList, _ := c.repo.GetYxListById(deviceId)
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: yxList,
	})
}

/*
*
获取设备全部遥测实时数据
*/
func (c *RealtimeDataController) GetRealtimeDataYcListByDevId(ctx *gin.Context) {
	var realtimeData ParamRealtimeData
	if err := ctx.Bind(&realtimeData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	deviceId := realtimeData.DeviceId
	if deviceId == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "参数错误",
			Data:    "",
		})
		return
	}
	ycList, _ := c.repo.GetYcListById(deviceId)
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: ycList,
	})
}

/*
*
获取设备全部实时数据
*/
func (c *RealtimeDataController) GetRealtimeDataListByDevId(ctx *gin.Context) {
	var realtimeData ParamRealtimeData
	if err := ctx.Bind(&realtimeData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	deviceId := realtimeData.DeviceId
	if deviceId == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "参数错误",
			Data:    "",
		})
		return
	}
	var slice1 []models.PointParam

	yxList, _ := c.repo.GetYxPointParamListById(deviceId)
	ycList, _ := c.repo.GetYcPointParamListById(deviceId)
	settingList, _ := c.repo.GetSettingPointParamListById(deviceId)

	slice1 = append(slice1, yxList[:]...)
	slice1 = append(slice1, ycList[:]...)
	slice1 = append(slice1, settingList[:]...)
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: slice1,
	})
}

/*
*
获取收益，7天，月，年，拿字典中电能表类型的设备，sum.年的按月采样
*/
func (c *RealtimeDataController) GetProfitPowerRate(ctx *gin.Context) {
	var returnMap struct {
		XAxisList []string          `json:"xAxisList"`
		DataMap   map[string][]string `json:"dataMap"`
	}

	var realtimeData ParamRealtimeData
	if err := ctx.Bind(&realtimeData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}

	dictData, err := c.repoData.SelectDictValue("device_type", "电能表")
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	var ids []int
	var intervalType string
	var startTime  int64
	var endTime    int64
	var xAxisList []string
	//查这个类型的设备
	deviceList, err := c.repoEm.GetDeviceList(dictData.DictValue)
	if len(deviceList) <= 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "暂无设备",
			Data:    "",
		})
		return
	}
	for _, device := range deviceList {
		ids = append(ids, device.Id)
	}
	//去taos查这些设备的日数据,0-7天,1-周,2-年
	if realtimeData.Interval == "" {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "参数错误",
			Data:    "",
		})
		return
	}else if realtimeData.Interval == "0" {
		intervalType = "1d"
		//7天
		xAxisList = utils.GetLast7Days();
		startTime, endTime = utils.GetLast7DaysTimestamps()
	}else if realtimeData.Interval == "1"{
		intervalType = "1d"
		//月
		xAxisList = utils.GetCurrentMonthDays();
		// 获取当月的开始时间戳和结束时间戳
		startTime, endTime = utils.GetCurrentMonthTimestamps()
	}else if realtimeData.Interval == "2"{
		intervalType = "1n"
		//年
		xAxisList = utils.GetAllMonths();
		// 获取当年的开始时间戳和结束时间戳
		startTime, endTime = utils.GetCurrentYearTimestamps()
	}
	if intervalType == "" {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "参数错误",
			Data:    "",
		})
		return
	}
	listT, listP, listF, listV, err := c.repo.GetDayProfitByDeviceIds(ids, startTime, endTime, intervalType);
	//拼接x轴
	returnMap.XAxisList = xAxisList

	pMap := make(map[string][]string)
	pMap["listTop"] = listT
	pMap["listPeak"] = listP
	pMap["listFlat"] = listF
	pMap["listValley"] = listV

	returnMap.DataMap = pMap

	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: returnMap,
	})
}
