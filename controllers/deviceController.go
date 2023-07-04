package controllers

import (
	"fmt"
	"gateway/device"
	"gateway/httpServer/model"
	"gateway/models"
	"gateway/models/query"
	repositories "gateway/repositories"
	"gateway/utils"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
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
	router.POST("/api/v2/device/getDeviceListByTypeList", c.getDeviceListByTypeList)
	router.POST("/api/v2/device/getDeviceInfo", c.GetDeviceInfo)
	router.POST("/api/v2/device/getChart", c.GetChart)
	router.POST("/api/v2/device/getGenerateElectricityChart", c.GetGenerateElectricityChart)
	router.POST("/api/v2/device/getProfitChart", c.GetProfitChart)
	router.POST("/api/v2/device/getCtrlHistoryList", c.GetCtrlHistoryList)
}

// CtrlDevice 遥控摇调
func (c *DeviceController) CtrlDevice(ctx *gin.Context) {
	var param models.CtrlInfo
	if err := ctx.Bind(&param); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	serviceInfo := struct {
		CollInterfaceName string                 `json:"collInterfaceName"`
		DeviceName        string                 `json:"deviceName"`
		ServiceName       string                 `json:"serviceName"`
		ServiceParam      map[string]interface{} `json:"serviceParam"`
	}{}

	// 查点
	deviceParam := c.repo.GetDeviceModelCmdParam(param.ParamId)
	serviceInfo.ServiceName = "SetVariables"
	serviceInfo.DeviceName = deviceParam.DeviceName
	serviceInfo.CollInterfaceName = deviceParam.CollName
	serviceParam := make(map[string]interface{})
	serviceParam[deviceParam.ParamName] = param.Value
	serviceInfo.ServiceParam = serviceParam

	device.CollectInterfaceMap.Lock.Lock()
	coll, ok := device.CollectInterfaceMap.Coll[serviceInfo.CollInterfaceName]
	device.CollectInterfaceMap.Lock.Unlock()
	if !ok {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "device is not exist",
			Data:    "",
		})
		return
	}

	cmd := device.CommunicationCmdTemplate{}
	cmd.CollInterfaceName = serviceInfo.CollInterfaceName
	cmd.DeviceName = serviceInfo.DeviceName
	cmd.FunName = serviceInfo.ServiceName
	paramStr, _ := json.Marshal(serviceInfo.ServiceParam)
	cmd.FunPara = string(paramStr)
	cmdRX := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
	var emCtrlHistory models.EmCtrlHistory
	emCtrlHistory.DeviceId = param.DeviceId
	emCtrlHistory.ParamId = param.ParamId
	emCtrlHistory.CtrlUserName = param.CtrlUserName
	emCtrlHistory.Value = fmt.Sprint(param.Value)
	emCtrlHistory.CreateTime = time.Now().String()
	emCtrlHistory.UpdateTime = time.Now().String()
	if cmdRX.Status == true {
		emCtrlHistory.CtrlStatus = 1
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "控制指令发送成功",
			Data:    "",
		})
	} else {
		emCtrlHistory.CtrlStatus = 0
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "device is not return",
			Data:    "",
		})
	}
	c.repo.AddCtrlHistory(&emCtrlHistory)
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
	list, _ := c.repo.GetDeviceListByType(param)
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "成功",
		Data:    list,
	})
}

func (c *DeviceController) getDeviceListByTypeList(ctx *gin.Context) {
	var param []string
	var res []models.EmDevice
	if err := ctx.Bind(&param); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	for _, item := range param {
		var deviceType models.DeviceParam
		deviceType.DeviceType = item
		list, _ := c.repo.GetDeviceListByType(deviceType)
		res = append(res, list...)
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "成功",
		Data:    res,
	})
}

func (c *DeviceController) GetDeviceInfo(ctx *gin.Context) {
	type Param struct {
		Name  string      `json:"name"`
		Label string      `json:"label"`
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
		v.Name = device.Label
		v.Label = device.Name
		// 查数据
		for _, dict := range dictDataList {
			var p Param
			p.Name = dict.DictValue
			p.Label = dict.DictLabel
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
	now := time.Now().In(utils.Location)
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UnixMilli()
	endTime := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()).Add(-time.Second).UnixMilli()
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
			tdDataList, _ := c.repoRealTimeData.GetChartByDeviceIdAndCode(device.Id, dictData.DictValue, startTime, endTime)
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
	tdDataListC, _ := c.repoRealTimeData.GetGenerateElectricityChartByDeviceIds(ids, "charge_capacity", "d")
	resC.Data = tdDataListC
	// 计算充电累计
	sumC := c.repoRealTimeData.GetGenerateElectricitySumByDeviceIds(ids)
	resC.Sum = sumC.Val
	result = append(result, resC)
	var resD Res
	resD.Name = "放电"
	tdDataListD, _ := c.repoRealTimeData.GetGenerateElectricityChartByDeviceIds(ids, "discharge_capacity", "d")
	resD.Data = tdDataListD
	// 计算放电累计
	sumD := c.repoRealTimeData.GetReleaseElectricitySumByDeviceIds(ids)
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
	interval := "1d"

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
		startTimeD = utils.GetLastMonthFirstDate(time.Now())
		// 上月最后一天
		endTimeD = utils.GetLastDateOfLastMonth(time.Now())
		nameC = "本月"
		nameD = "上月"
	case "year":
		// 今年第一天
		startTimeC = utils.GetFirstDateOfYear(time.Now())
		// 今年最后一天
		endTimeC = utils.GetLastDateOfYear(time.Now())
		// 去年第一天
		startTimeD = utils.GetFirstDateOfFirstYear(time.Now())
		// 去年最后一天
		endTimeD = utils.GetLastDateOfLastYear(time.Now())
		nameC = "今年"
		nameD = "去年"
		interval = "1n"
	default:

	}
	var resC Res
	resC.Name = nameC
	tdDataListC, _ := c.repoRealTimeData.GetProfitChartByDeviceIds(ids, startTimeC.UnixMilli(), endTimeC.UnixMilli(), interval)
	resC.Data = tdDataListC
	// 计算本周收益累计
	sumC := c.repoRealTimeData.GetProfitSumByDeviceIds(ids, startTimeC.UnixMilli(), endTimeC.UnixMilli())
	resC.Sum = sumC.Val
	result = append(result, resC)
	var resD Res
	resD.Name = nameD
	tdDataListD, _ := c.repoRealTimeData.GetProfitChartByDeviceIds(ids, startTimeD.UnixMilli(), endTimeD.UnixMilli(), interval)
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

func (c *DeviceController) GetCtrlHistoryList(ctx *gin.Context) {
	var (
		res   []repositories.CtrlHistory
		total int64
	)
	var PageBase query.PageBase
	if err := ctx.Bind(&PageBase); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			err.Error(),
			"",
		})
		return
	}

	res, total, err := c.repo.GetCtrlHistoryList(PageBase.PageNum, PageBase.PageSize)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			err.Error(),
			"",
		})
		return
	}

	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"成功",
		gin.H{
			"data":  res,
			"total": total,
		},
	})
}
