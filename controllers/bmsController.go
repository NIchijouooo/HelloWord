package controllers

import (
	"encoding/json"
	"gateway/httpServer/model"
	"gateway/models"
	"gateway/models/ReturnModel"
	"gateway/models/query"
	repositories "gateway/repositories"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type BmsController struct {
	hisRepo      *repositories.HistoryDataRepository
	dictDataRepo *repositories.DictDataRepository
	auxRepo      *repositories.AuxiliaryRepository
	emRepo       *repositories.EmRepository
	ruRepo       *repositories.RuleHistoryRepository
	limitRepo    *repositories.LimitConfigRepository
	bmsRepo      *repositories.BmsRepository
	realRepo     *repositories.RealtimeDataRepository
}

func NewBmsController() *BmsController {
	return &BmsController{hisRepo: repositories.NewHistoryDataRepository(),
		dictDataRepo: repositories.NewDictDataRepository(),
		auxRepo:      repositories.NewAuxiliaryRepository(),
		emRepo:       repositories.NewEmRepository(),
		ruRepo:       repositories.NewRuleHistoryRepository(),
		limitRepo:    repositories.NewLimitConfigRepository(),
		bmsRepo:      repositories.NewBmsRepository(),
		realRepo:     repositories.NewRealtimeDataRepository()}
}
func (c *BmsController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/bms/getDeviceTreeByDeviceType", c.GetDeviceTreeByDeviceType)
	router.POST("/api/v2/bms/getYcLogById", c.GetYcLogById)
	router.POST("/api/v2/bms/getHistoryYcByDeviceIdCodes", c.GetHistoryYcByDeviceIdCodes)
	router.POST("/api/v2/bms/getBmsYcMaxAndMinListByDeviceIdCodes", c.GetBmsYcMaxAndMinListByDeviceIdCodes)
	router.POST("/api/v2/bms/getBmsDevices", c.GetBmsDevices)
	router.POST("/api/v2/bms/getDayElectricityChartByDeviceId", c.GetDayElectricityChartByDeviceId) //获取日电量曲线
	router.POST("/api/v2/bms/getDevicesStatus", c.GetDevicesStatus)                                 //获取设备状态

}

// GetDeviceTreeByDeviceType /获取设备类型下的所有设备数据,返回树形结构
func (c *BmsController) GetDeviceTreeByDeviceType(ctx *gin.Context) {
	type tmpQuery struct {
		DeviceType string `json:"deviceType"`
	}
	var queryParam tmpQuery
	if err := ctx.Bind(&queryParam); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	deviceList, err := c.bmsRepo.GetBmsDeviceList(queryParam.DeviceType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	nodeMap := make(map[int]*ReturnModel.TreeDevice)
	var rootNodes []*ReturnModel.TreeDevice
	//构建节点映射表
	for _, node := range deviceList {
		nodeMap[node.Id] = node
	}
	//构建父子关系
	for _, node := range deviceList {
		parentId := node.ParentId
		parentNode, err := nodeMap[parentId]
		if err {
			parentNode.Children = append(parentNode.Children, *node)
		} else {
			rootNodes = append(rootNodes, node)
		}
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取信息成功！",
		Data:    rootNodes,
	})
	return
}

// GetYcLogById 批量获取遥信息GetYcLogById,用于·历史数据
func (c *BmsController) GetYcLogById(ctx *gin.Context) {
	type ycQeury struct {
		DeviceId  int    `form:"DeviceId"`
		Codes     string `form:"codes"`
		StartTime int64  `form:"StartTime"`
		EndTime   int64  `form:"EndTime"`
	}
	var ycQuery ycQeury
	//将传过来的请求体解析到ycQuery中
	if err := ctx.Bind(&ycQuery); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	//查询历史数据
	var ycLog []map[string]interface{} //这里需要解开
	//ycLog, err := c.hisRepo.GetYcLogById(ycQuery.DeviceId, ycQuery.Codes, ycQuery.StartTime, ycQuery.EndTime)
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取信息成功！",
		Data:    ycLog,
	})
}

// GetHistoryYcByDeviceIdCodes 根据选择的codes返回对应的历史数据
func (c *BmsController) GetHistoryYcByDeviceIdCodes(ctx *gin.Context) {
	var ycQuery query.QueryTaoData
	//解析json
	if err := ctx.Bind(&ycQuery); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	returnMap, err := c.hisRepo.GetCharData(ycQuery)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    nil,
		})
	} else {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "获取信息成功！",
			Data:    returnMap,
		})
	}
}

// GetBmsYcMaxAndMinListByDeviceIdCodes BMS批量获取最高最低的遥测信息
func (c *BmsController) GetBmsYcMaxAndMinListByDeviceIdCodes(ctx *gin.Context) {
	var ycData models.YcData
	if err := ctx.Bind(&ycData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	if ycData.DeviceId == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "缺少设备Id",
			Data:    "",
		})
		return
	}
	//把所有的BMS测点查出来
	var dictLabels = []string{
		"energy_storage_bms_max_voltage_code",        //BMS设备最高电压测点
		"energy_storage_bms_min_voltage_code",        //BMS设备最低电压测点
		"energy_storage_bms_max_voltage_serial_code", //BMS设备最高电压序号测点
		"energy_storage_bms_min_voltage_serial_code", //BMS设备最低电压序号测点
		"energy_storage_bms_max_temp_code",           //BMS设备最高温度测点
		"energy_storage_bms_min_temp_code",           //BMS设备最低温度测点
		"energy_storage_bms_max_temp_serial_code",    //BMS设备最高温度序号测点
		"energy_storage_bms_min_temp_serial_code",    //BMS设备最低温度序号测点
		"energy_storage_bms_max_soc_code",            //BMS设备最高SOC测点
		"energy_storage_bms_min_soc_code",            //BMS设备最低SOC测点
		"energy_storage_bms_max_soc_serial_code",     //BMS设备最高SOC序号测点
		"energy_storage_bms_min_soc_serial_code",     //BMS设备最低SOC序号测点
		"energy_storage_bms_max_soh_code",            //BMS设备最高SOH测点
		"energy_storage_bms_min_soh_code",            //BMS设备最低SOH测点
		"energy_storage_bms_max_soh_serial_code",     //BMS设备最高SOH序号测点
		"energy_storage_bms_min_soh_serial_code",     //BMS设备最低SOH序号测点
	}
	//1.查询字典对应的测点code
	DictDataList, err := c.dictDataRepo.GetDictDataByDictLabel(dictLabels)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	var codeList []string
	for _, v := range DictDataList {
		codeList = append(codeList, v.DictValue)
	}
	//2.查询最新一条历史数据
	codes := strings.Join(codeList, ",")
	ycDataList, err := c.hisRepo.GetLastYcListByCode(strconv.Itoa(ycData.DeviceId), codes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//3.查询em_device_model_cmd_param表取出单位
	paramList, err := c.emRepo.GetEmDeviceModelCmdParamListByDeviceIdCodes(ycData.DeviceId, codeList)
	paramMap := make(map[string]*models.EmDeviceModelCmdParam)
	if err == nil && len(paramList) > 0 {
		for _, v := range paramList {
			item := v
			paramMap[v.Name] = &item //name就是测点code
		}
	}
	//4.查询告警信息---暂时先不用，自己遍历=====================================================
	//var param models.RuleHistoryParam
	//deviceIds := []int{ycData.DeviceId}
	//param.DeviceIds = deviceIds
	//param.Codes, err = utils.StringArrayToIntArray(codeList) //测点code
	//param.Tag = "0"                                          //最新告警
	//list, total, err := c.ruRepo.GetRuleHistoryList(param)
	//var ruleHistoryMap = make(map[string]*models.EmRuleHistoryModel)
	//if total != 0 { //存在数据
	//	//for循环list，以键值对形式存储到map中
	//	for _, v := range list {
	//		item := v
	//		ruleHistoryMap[strconv.Itoa(v.PropertyCode)] = &item
	//	}
	//}
	//=====================================================
	//5.获取越线配置信息=============================================
	//根据设备id获取设备信息
	emDevice, err := c.emRepo.GetEmDeviceById(ycData.DeviceId)
	limitConfigList, err := c.limitRepo.GetLimitConfigListByDeviceTypeAndCodes(emDevice.DeviceType, codeList)
	//根据code按照键值对存储到map中
	var limitConfigMap = make(map[string]*ReturnModel.LimitScope)
	if err == nil && len(limitConfigList) > 0 {
		for _, v := range limitConfigList {
			limitScope := &ReturnModel.LimitScope{
				PropertyCode: v.PropertyCode,
				NotifyMin:    v.NotifyMin,
				NotifyMax:    v.NotifyMax,
				SeriousMax:   v.SeriousMax,
				SeriousMin:   v.SeriousMin,
				SecondaryMax: v.SecondaryMax,
				SecondaryMin: v.SecondaryMin,
				UrgentMax:    v.UrgentMax,
				UrgentMin:    v.UrgentMin,
			}
			limitConfigMap[v.PropertyCode] = limitScope
		}
	}
	//=============================================
	//字典数据转成map
	ycMap := make(map[string]*models.YcData)
	for _, v := range ycDataList {
		ycMap[strconv.Itoa(v.Code)] = v
	}
	var resData []ReturnModel.YcData
	//循环字典数据，赋值
	for _, v := range DictDataList {
		var yc ReturnModel.YcData
		tmpYcData := ycMap[v.DictValue] //赋值数据
		//单位获取--------------------------------------
		tmpUnit := paramMap[v.DictValue] //赋值单位
		if tmpUnit != nil {              //存在测点数据，取出单位
			yc.Unit = tmpUnit.Unit //单位还要再去查表
		}
		//告警等级--------------------------------------
		//tmpLevel := ruleHistoryMap[v.DictValue] //告警数据
		//if tmpLevel != nil {
		//	yc.Level = tmpLevel.Level
		//} else {
		//	yc.Level = -1 //-1无告警
		//}

		//遥测数据--------------------------------------
		if tmpYcData != nil { //没有数据赋值为空
			//越线配置--------------------------------------
			tmpLimit := limitConfigMap[v.DictValue] //越线配置
			yc.Level = -1
			if tmpLimit != nil {
				notifyMini, n1 := strconv.ParseFloat(tmpLimit.NotifyMin, 64)
				notifyMax, n2 := strconv.ParseFloat(tmpLimit.NotifyMax, 64)
				secondaryMin, s1 := strconv.ParseFloat(tmpLimit.SecondaryMin, 64)
				secondaryMax, s2 := strconv.ParseFloat(tmpLimit.SecondaryMax, 64)
				seriousMin, s3 := strconv.ParseFloat(tmpLimit.SeriousMin, 64)
				seriousMax, s4 := strconv.ParseFloat(tmpLimit.SeriousMax, 64)
				urgentMin, u1 := strconv.ParseFloat(tmpLimit.UrgentMin, 64)
				urgentMax, u2 := strconv.ParseFloat(tmpLimit.UrgentMax, 64)
				if n1 == nil && n2 == nil && s1 == nil && s2 == nil && s3 == nil && s4 == nil && u1 == nil && u2 == nil {
					if tmpYcData.Value >= notifyMini && tmpYcData.Value <= notifyMax { //越线判断
						yc.Level = 0
					} else if tmpYcData.Value >= secondaryMin && tmpYcData.Value <= secondaryMax {
						yc.Level = 1
					} else if tmpYcData.Value >= seriousMin && tmpYcData.Value <= seriousMax {
						yc.Level = 2
					} else if tmpYcData.Value >= urgentMin && tmpYcData.Value <= urgentMax {
						yc.Level = 3
					} else {
						yc.Level = -1
					}
				}
				yc.LimitScope = *tmpLimit //配置范围
			}
			yc.DeviceId = tmpYcData.DeviceId
			yc.Code = tmpYcData.Code
			yc.Value = tmpYcData.Value
			yc.Name = tmpYcData.Name
			yc.Ts = tmpYcData.Ts
			yc.Sort = v.DictSort   //根据code去拿数据
			yc.Alias = v.DictLabel //根据code去拿数据
		} else { //没有数据就补-
			yc.DeviceId = ycData.DeviceId
			code, err := strconv.Atoi(v.DictValue) //字符串转换int
			if err == nil {
				yc.Code = code
			}
			yc.Name = ""
			yc.Unit = ""
			yc.Value = 0
			yc.Sort = v.DictSort   //根据code去拿数据
			yc.Alias = v.DictLabel //根据code拿数据
			yc.LimitScope = ReturnModel.LimitScope{}
		}
		resData = append(resData, yc)
	}
	//排序 //按sort排序
	sort.Slice(resData, func(i, j int) bool {
		return resData[i].Sort < resData[j].Sort
	})
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取信息成功！",
		Data:    resData,
	})
	return
}

// GetBmsDevices 获取BMS设备
func (c *BmsController) GetBmsDevices(ctx *gin.Context) {
	var dictLabels = []string{"energy_storage_bms_device_label"}
	DictDataList, err := c.dictDataRepo.GetDictDataByDictLabel(dictLabels)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	if len(DictDataList) > 0 {
		var codeList []string
		for _, v := range DictDataList {
			codeList = append(codeList, v.DictValue)
		}
		deviceList, err := c.auxRepo.GetAuxiliaryDevice(codeList[0])
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "获取信息成功！",
			Data:    deviceList,
		})
		return
	} else {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "无数据！",
			Data:    "",
		})
		return
	}
}

// GetDayElectricityChartByDeviceId 获取每天的充放电量数据
func (c *BmsController) GetDayElectricityChartByDeviceId(ctx *gin.Context) {
	type Res struct {
		Name string             `json:"name"`
		Data []repositories.Res `json:"data"`
	}

	var queryData query.QueryTaoData
	if err := ctx.Bind(&queryData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//必要条件校验
	if len(queryData.DeviceIds) == 0 || queryData.StartTime == 0 || queryData.EndTime == 0 || queryData.IntervalType == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "缺少必要参数",
			Data:    "",
		})
		return
	}
	//间隔时间段
	intervalStr := "s"
	if queryData.IntervalType == 2 {
		intervalStr = "m"
	} else if queryData.IntervalType == 3 {
		intervalStr = "h"
	} else if queryData.IntervalType > 3 {
		intervalStr = "d"
	}
	tableName := "charge_discharge"
	if queryData.IntervalType == 3 {
		tableName = "charge_discharge_hour"
	}
	var result []Res
	var resC Res
	resC.Name = "充电"
	//充电
	tdDataLisC, err := c.realRepo.GetElectricityChartByDeviceIds(queryData.DeviceIds, "charge_capacity", queryData.StartTime, queryData.EndTime, intervalStr, tableName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resC.Data = tdDataLisC
	var resD Res
	resD.Name = "放电"
	//放电
	tdDataLisD, err := c.realRepo.GetElectricityChartByDeviceIds(queryData.DeviceIds, "discharge_capacity", queryData.StartTime, queryData.EndTime, intervalStr, tableName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resD.Data = tdDataLisD
	result = append(result, resC)
	result = append(result, resD)
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取信息成功！",
		Data:    result,
	})
}

// GetDevicesStatus 获取设备状态信息
func (c *BmsController) GetDevicesStatus(ctx *gin.Context) {
	var queryData query.QueryTaoData
	if err := ctx.Bind(&queryData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//必要条件校验
	if queryData.DeviceId == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "缺少必要参数",
			Data:    "",
		})
		return
	}
	//1.根据code查询设备测点信息
	dictType := "energy_product_code_setting"     //类型状态
	dictLabel := "energy_storage_bms_status_code" //状态测点
	dict, err := c.dictDataRepo.SelectDictValue(dictType, dictLabel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if dict == (models.DictData{}) {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "无数据",
			Data:    "",
		})
		return
	}
	//2.通过测点信息和设备id去获取最新测点值
	ycData, err := c.hisRepo.GetLastYcListByCode(strconv.Itoa(queryData.DeviceId), dict.DictValue) //查询最新测点值
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if ycData == nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "无数据",
			Data:    "",
		})
		return
	}
	//3.获取设备状态json
	dictStatus := "bms_run_status" //设备运行状态
	dictStatusLabel := "bms"       //状态json
	dictRes, err := c.dictDataRepo.SelectDictValue(dictStatus, dictStatusLabel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//4.解析成map
	var status map[string]string
	errJson := json.Unmarshal([]byte(dictRes.DictValue), &status)
	if errJson != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "json解析失败",
			Data:    "",
		})
		return
	}
	//5.赋值
	for _, yc := range ycData { //读字典赋值
		name := status[strconv.Itoa(int(yc.Value))] //不要小数点
		if name != "" {
			yc.Name = name
		}
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取信息成功！",
		Data:    ycData,
	})
}
