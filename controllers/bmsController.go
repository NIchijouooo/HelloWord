package controllers

import (
	"gateway/httpServer/model"
	"gateway/models"
	"gateway/models/ReturnModel"
	"gateway/models/query"
	repositories "gateway/repositories"
	"gateway/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type BmsController struct {
	hisRepo      *repositories.HistoryDataRepository
	dictDataRepo *repositories.DictDataRepository
	auxRepo      *repositories.AuxiliaryRepository
}

func NewBmsController() *BmsController {
	return &BmsController{hisRepo: repositories.NewHistoryDataRepository(), dictDataRepo: repositories.NewDictDataRepository(), auxRepo: repositories.NewAuxiliaryRepository()}
}
func (ctrl *BmsController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/bms/getYcLastByDeviceIdAndDict", ctrl.GetYcLastByDeviceIdAndDict)
	router.POST("/api/v2/bms/getYcLogById", ctrl.GetYcLogById)
	router.POST("/api/v2/bms/getHistoryYcByDeviceIdCodes", ctrl.GetHistoryYcByDeviceIdCodes)
	router.POST("/api/v2/bms/getBmsYcMaxAndMinListByDeviceIdCodes", ctrl.GetBmsYcMaxAndMinListByDeviceIdCodes)
	router.POST("/api/v2/bms/getBmsDevices", ctrl.GetBmsDevices)

}

// 获取设备点位最新的一条非空数据
func (c *BmsController) GetYcLastByDeviceIdAndDict(ctx *gin.Context) {
	//1.根据设备id去查询字典
	//2.将查询出来的所有code拼接成字符串
	deviceId := ctx.Value("deviceId")
	//codes := ctx.Value("codes")
	print(deviceId)
	//codes := "2005,2002"
	//3.然后查询涛思数据库
	//ycList, err := c.hisRepo.GetLastYcListByCode(codes)
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	//ctx.JSON(http.StatusOK, model.ResponseData{
	//	"0",
	//	"获取信息成功！",
	//	ycList,
	//})
}

//批量获取遥信息GetYcLogById,用于·历史数据
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
			"1",
			"error" + err.Error(),
			"",
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
		"0",
		"获取信息成功！",
		ycLog,
	})
}

//根据选择的codes返回对应的历史数据
func (c *BmsController) GetHistoryYcByDeviceIdCodes(ctx *gin.Context) {
	var ycQuery *query.QueryTaoData
	//解析json
	if err := ctx.ShouldBindBodyWith(&ycQuery, binding.JSON); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	//必要条件校验
	if ycQuery.StartTime == 0 || ycQuery.EndTime == 0 || ycQuery.Interval == 0 || ycQuery.IntervalType == 0 || len(ycQuery.CodeList) == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error",
			"缺少必要参数",
		})
		return
	}
	//间隔时间段
	intervalStr := "s"
	if ycQuery.IntervalType == 2 {
		intervalStr = "m"
	} else if ycQuery.IntervalType == 3 {
		intervalStr = "h"
	} else if ycQuery.IntervalType > 3 {
		intervalStr = "d"
	}
	// 将 int 数组转换为字符串数组
	strArr := make([]string, len(ycQuery.CodeList))
	for i, v := range ycQuery.CodeList {
		strArr[i] = strconv.Itoa(v)
	}

	//拼接codelist sql语句
	ycQuery.Codes = strings.Join(strArr, ",")
	//查询历史数据
	ycList, err := c.hisRepo.GetLastYcHistoryByDeviceIdAndCodeList(ycQuery.DeviceId, ycQuery.Codes, ycQuery.StartTime, ycQuery.EndTime, strconv.Itoa(ycQuery.Interval)+intervalStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(ycList) > 0 {
		var xAxisList []string
		// 初始化x轴数据,返回x轴时间对应的历史数据分组,key=x轴,value=x轴对应的历史数据集合
		returnMap := service.GetCharData(xAxisList, ycQuery.StartTime, ycQuery.EndTime, ycQuery.Interval, ycQuery.IntervalType, ycList, ycQuery.CodeList, ycQuery.CodeNameList)

		ctx.JSON(http.StatusOK, model.ResponseData{
			"0",
			"获取信息成功！",
			returnMap,
		})
	} else {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"无数据！",
			"",
		})
	}
}

//BMS批量获取最高最低的遥测信息
func (c *BmsController) GetBmsYcMaxAndMinListByDeviceIdCodes(ctx *gin.Context) {
	var ycData models.YcData
	if err := ctx.Bind(&ycData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
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
	//查询字典code
	DictDataList, err := c.dictDataRepo.GetDictDataByDictLabel(dictLabels)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	var codeList []string
	for _, v := range DictDataList {
		codeList = append(codeList, v.DictValue)
	}
	//查询历史数据
	codes := strings.Join(codeList, ",")
	ycDataList, err := c.hisRepo.GetLastYcListByCode(strconv.Itoa(ycData.DeviceId), codes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//字典数据转成map
	dictList := make(map[string]models.DictData)
	for _, v := range DictDataList {
		dictList[v.DictValue] = v
	}
	var resData []ReturnModel.YcData
	for _, v := range ycDataList {
		tmpDictData := dictList[strconv.Itoa(v.Code)]
		var ycData ReturnModel.YcData
		ycData.DeviceId = v.DeviceId
		ycData.Code = v.Code
		ycData.Value = v.Value
		ycData.Name = v.Name
		ycData.Ts = v.Ts
		ycData.Sort = tmpDictData.DictSort   //根据code去拿数据
		ycData.Alias = tmpDictData.DictLabel //根据code去拿数据
		resData = append(resData, ycData)
	}
	//排序
	sort.Slice(resData, func(i, j int) bool {
		return resData[i].Sort < resData[j].Sort
	})
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"获取信息成功！",
		resData,
	})

	//ycMap := make(map[string]*models.YcData)
	////将遥测数据转成map
	//for _, v := range ycDataList {
	//	ycMap[strconv.Itoa(v.Code)] = v
	//}
	//
	////遍历返回数据
	//resMap := make(map[string]models.YcData)
	//for _, v := range DictDataList {
	//	if ycMap[v.DictValue] != nil {
	//		resMap[v.DictLabel] = *ycMap[v.DictValue]
	//	}
	//}
	//ctx.JSON(http.StatusOK, model.ResponseData{
	//	"0",
	//	"获取信息成功！",
	//	resMap,
	//})
	return
}

//获取BMS设备
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
			"0",
			"获取信息成功！",
			deviceList,
		})
		return
	} else {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"无数据！",
			"",
		})
		return
	}
}
