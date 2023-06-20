package controllers

import (
	"gateway/httpServer/model"
	"gateway/models"
	"gateway/models/ReturnModel"
	"gateway/models/query"
	repositories "gateway/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type BmsController struct {
	hisRepo      *repositories.HistoryDataRepository
	dictDataRepo *repositories.DictDataRepository
	auxRepo      *repositories.AuxiliaryRepository
	emRepo       *repositories.EmRepository
}

func NewBmsController() *BmsController {
	return &BmsController{hisRepo: repositories.NewHistoryDataRepository(),
		dictDataRepo: repositories.NewDictDataRepository(),
		auxRepo:      repositories.NewAuxiliaryRepository(),
		emRepo:       repositories.NewEmRepository()}
}
func (ctrl *BmsController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/bms/getYcLastByDeviceIdAndDict", ctrl.GetYcLastByDeviceIdAndDict)
	router.POST("/api/v2/bms/getYcLogById", ctrl.GetYcLogById)
	router.POST("/api/v2/bms/getHistoryYcByDeviceIdCodes", ctrl.GetHistoryYcByDeviceIdCodes)
	router.POST("/api/v2/bms/getBmsYcMaxAndMinListByDeviceIdCodes", ctrl.GetBmsYcMaxAndMinListByDeviceIdCodes)
	router.POST("/api/v2/bms/getBmsDevices", ctrl.GetBmsDevices)

}

// GetYcLastByDeviceIdAndDict 获取设备点位最新的一条非空数据
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

// 根据选择的codes返回对应的历史数据
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

// BMS批量获取最高最低的遥测信息
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
	//取出单位
	paramList, err := c.emRepo.GetEmDeviceModelCmdParamListByDeviceIdCodes(ycData.DeviceId, codeList)
	paramMap := make(map[string]models.EmDeviceModelCmdParam)
	if err == nil && len(paramList) > 0 {
		for _, v := range paramList {
			paramMap[v.Name] = v //name就是测点code
		}
	}
	//字典数据转成map
	ycMap := make(map[string]*models.YcData)
	for _, v := range ycDataList {
		ycMap[strconv.Itoa(v.Code)] = v
	}
	var resData []ReturnModel.YcData
	//循环字典数据，赋值
	for _, v := range DictDataList {
		tmpYcData := ycMap[v.DictValue]  //赋值数据
		tmpUnit := paramMap[v.DictValue] //赋值单位
		var yc ReturnModel.YcData
		if tmpYcData != nil { //存在测点数据，取出单位
			yc.Uint = tmpUnit.Unit //单位还要再去查表
		}
		if tmpYcData != nil { //没有数据赋值为空
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
			yc.Value = 0
			yc.Name = "-"
			yc.Uint = "-"
			yc.Sort = v.DictSort   //根据code去拿数据
			yc.Alias = v.DictLabel //根据code拿数据
		}
		resData = append(resData, yc)
	}
	//排序 //按sort排序
	sort.Slice(resData, func(i, j int) bool {
		return resData[i].Sort < resData[j].Sort
	})
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"获取信息成功！",
		resData,
	})
	return
}

// 获取BMS设备
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
