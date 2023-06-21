package controllers

import (
	"fmt"
	"gateway/httpServer/model"
	"gateway/models"
	repositories "gateway/repositories"
	"gateway/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type DeviceController struct {
	repo             *repositories.DeviceRepository
	repoRealTimeData *repositories.RealtimeDataRepository
	repoDictData     *repositories.DictDataRepository
}

func NewDeviceController() *DeviceController {
	return &DeviceController{
		repo:             repositories.NewDeviceRepository(),
		repoRealTimeData: repositories.NewRealtimeDataRepository(),
		repoDictData:     repositories.NewDictDataRepository(),
	}
}

func (c *DeviceController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/device/ctrlDevice", c.CtrlDevice)
	router.POST("/api/v2/device/getDeviceListByType", c.getDeviceListByType)
	router.POST("/api/v2/device/getDeviceInfo", c.GetDeviceInfo)
	router.POST("/api/v2/device/getChart", c.GetChart)
	router.POST("/api/v2/device/getGenerateElectricityChart", c.GetGenerateElectricityChart)
	router.POST("/api/v2/device/getProfitChart", c.GetProfitChart)
}

// 未完成
func (ctrl *DeviceController) CtrlDevice(ctx *gin.Context) {
	var param models.CtrlInfo
	if err := ctx.Bind(&param); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	device, err := repositories.NewEmRepository().GetEmDeviceById(param.DeviceId)
	if device == nil || err != nil {
		msg := "未查询到设备"
		if err != nil {
			msg = "error" + err.Error()
		}
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: msg,
			Data:    "",
		})
		return
	}
}

/*
*
根据设备类型获取设备列表
deviceType必传
*/
func (c *DeviceController) getDeviceListByType(ctx *gin.Context) {
	var param models.DeviceParam
	if err := ctx.Bind(&param); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	deviceType := param.DeviceType
	if len(deviceType) == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "参数错误",
			Data:    "",
		})
		return
	}
	list, _ := c.repo.GetDeviceListByType(param)
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: list,
	})
}

func (c *DeviceController) GetDeviceInfo(ctx *gin.Context) {
	type Param struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	}
	type Res struct {
		Id     int     `json:"id"`
		Name   string  `json:"name"`
		Label  string  `json:"label"`
		Params []Param `json:"params"`
	}
	var res []Res
	var deviceParam models.DeviceParam
	if err := ctx.Bind(&deviceParam); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	// 获取设备列表
	deviceList, _ := c.repo.GetDeviceListByType(deviceParam)
	// 查字典
	dictDataList, _ := c.repoDictData.GetDictDataByDictType(deviceParam.DeviceType)
	// 构建数据
	for _, device := range deviceList {
		var v Res
		v.Id = device.Id
		v.Name = device.Name
		v.Label = device.Label
		// 查数据
		for _, dict := range dictDataList {
			var p Param
			p.Name = dict.DictLabel
			// 查时序库
			code, err := strconv.Atoi(dict.DictValue)
			if err != nil {
				fmt.Println(err)
				return
			}
			yc, _ := c.repoRealTimeData.GetYcById(device.Id, code)
			p.Value = yc.Value
			v.Params = append(v.Params, p)
		}
		res = append(res, v)
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "成功",
		Data:    res,
	})

}

func (c *DeviceController) GetChart(ctx *gin.Context) {
	type Param struct {
		DeviceType    string `json:"deviceType"`
		PointDictType string `json:"pointDictType"`
	}
	type Res struct {
		DeviceName string             `json:"deviceName"`
		PointName  string             `json:"pointName"`
		Data       []repositories.Res `json:"data"`
	}
	var param Param
	var result []Res
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	// 找设备
	var deviceParam models.DeviceParam
	deviceParam.DeviceType = param.DeviceType
	deviceList, _ := c.repo.GetDeviceListByType(deviceParam)
	// 找点位的code
	dictDataList, _ := c.repoDictData.GetDictDataByDictType(param.PointDictType)
	// 查时序数据
	for _, device := range deviceList {
		var res Res
		res.DeviceName = device.Name
		for _, dictData := range dictDataList {
			res.PointName = dictData.DictLabel
			tdDataList, _ := c.repoRealTimeData.GetChartByDeviceIdAndCode(device.Id, dictData.DictValue)
			res.Data = tdDataList
		}
		result = append(result, res)
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "成功",
		Data:    result,
	})
}

func (c *DeviceController) GetGenerateElectricityChart(ctx *gin.Context) {
	type Param struct {
		DeviceType string `json:"deviceType"`
	}
	type Res struct {
		Name string             `json:"name"`
		Sum  float32            `json:"sum"`
		Data []repositories.Res `json:"data"`
	}
	// 找设备
	var deviceParam models.DeviceParam
	var param Param
	var result []Res
	var ids []int
	if err := ctx.Bind(&param); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	deviceParam.DeviceType = param.DeviceType
	deviceList, _ := c.repo.GetDeviceListByType(deviceParam)
	for _, device := range deviceList {
		ids = append(ids, device.Id)
	}
	var resC Res
	resC.Name = "充电"
	tdDataListC, _ := c.repoRealTimeData.GetGenerateElectricityChartByDeviceIds(ids)
	resC.Data = tdDataListC
	// 计算充电累计
	sumC := c.repoRealTimeData.GetGenerateElectricitySumByDeviceIds(ids)
	resC.Sum = sumC.Val
	result = append(result, resC)
	var resD Res
	resD.Name = "放电"
	tdDataListD, _ := c.repoRealTimeData.GetGenerateElectricityChartByDeviceIds(ids)
	resD.Data = tdDataListD
	// 计算放电累计
	sumD := c.repoRealTimeData.GetGenerateElectricitySumByDeviceIds(ids)
	resD.Sum = sumD.Val
	result = append(result, resD)

	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "成功",
		Data:    result,
	})
}

func (c *DeviceController) GetProfitChart(ctx *gin.Context) {
	type Param struct {
		DeviceType string `json:"deviceType"`
		DateType   string `json:"dateType"`
	}
	type Res struct {
		Name string             `json:"name"`
		Sum  float32            `json:"sum"`
		Data []repositories.Res `json:"data"`
	}
	// 找设备
	var deviceParam models.DeviceParam
	var param Param
	var result []Res
	var ids []int

	var startTimeC time.Time
	var endTimeC time.Time
	var startTimeD time.Time
	var endTimeD time.Time
	var nameC string
	var nameD string
	if err := ctx.Bind(&param); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	deviceParam.DeviceType = param.DeviceType
	deviceList, _ := c.repo.GetDeviceListByType(deviceParam)
	for _, device := range deviceList {
		ids = append(ids, device.Id)
	}
	// 7天/月/年
	switch param.DateType {
	case "week":
		// 本周一
		startTimeC = utils.GetFirstDateOfWeek(time.Now())
		// 本周日
		endTimeC = utils.GetLastDateOfWeek(time.Now())
		// 上周一
		startTimeD = utils.GetLastWeekFirstDate(time.Now())
		// 上周日
		endTimeD = utils.GetLastDateOfWeek(startTimeD)
		nameC = "本周"
		nameD = "上周"
	case "month":
		// 本月第一天
		startTimeC = utils.GetFirstDateOfMonth(time.Now())
		// 本月最后一天
		endTimeC = utils.GetLastDateOfMonth(time.Now())
		// 上月第一天
		startTimeD = utils.GetLastWeekFirstDate(time.Now())
		// 上周最后一天
		endTimeD = utils.GetLastDateOfWeek(startTimeD)
		nameC = "本月"
		nameD = "上月"
	case "year":
		// 今年第一天
		startTimeC = utils.GetFirstDateOfWeek(time.Now())
		// 今年最后一天
		endTimeC = utils.GetLastDateOfWeek(time.Now())
		// 去年第一天
		startTimeD = utils.GetLastWeekFirstDate(time.Now())
		// 去年最后一天
		endTimeD = utils.GetLastDateOfWeek(startTimeD)
		nameC = "今年"
		nameD = "去年"
	default:

	}
	var resC Res
	resC.Name = nameC
	tdDataListC, _ := c.repoRealTimeData.GetProfitChartByDeviceIds(ids, startTimeC.UnixMilli(), endTimeC.UnixMilli())
	resC.Data = tdDataListC
	// 计算本周收益累计
	sumC := c.repoRealTimeData.GetProfitSumByDeviceIds(ids, startTimeC.UnixMilli(), endTimeC.UnixMilli())
	resC.Sum = sumC.Val
	result = append(result, resC)
	var resD Res
	resD.Name = nameD
	tdDataListD, _ := c.repoRealTimeData.GetProfitChartByDeviceIds(ids, startTimeD.UnixMilli(), endTimeD.UnixMilli())
	resD.Data = tdDataListD
	// 计算下周累计
	sumD := c.repoRealTimeData.GetProfitSumByDeviceIds(ids, startTimeD.UnixMilli(), endTimeD.UnixMilli())
	resD.Sum = sumD.Val
	result = append(result, resD)
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "成功",
		Data:    result,
	})
}
