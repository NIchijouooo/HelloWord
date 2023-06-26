package controllers

import (
	"gateway/httpServer/model"
	"gateway/models"
	"gateway/models/ReturnModel"
	"gateway/models/query"
	repositories "gateway/repositories"
	"gateway/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
	"strings"
)

// 定义辅控管理的控制器
type AuxiliaryController struct {
	repo    *repositories.AuxiliaryRepository
	hisRepo *repositories.HistoryDataRepository
	emRepo  *repositories.EmRepository
}

func NewAuxiliaryController() *AuxiliaryController {
	return &AuxiliaryController{repo: repositories.NewAuxiliaryRepository(),
		hisRepo: repositories.NewHistoryDataRepository(),
		emRepo:  repositories.NewEmRepository()}
}

func (ctrl *AuxiliaryController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/auxiliary/getDeviceListByDeviceType", ctrl.GetDeviceListByDeviceType)
	router.GET("/api/v2/auxiliary/getDeviceType", ctrl.GetAuxiliaryDeviceType)
	router.POST("/api/v2/auxiliary/getLastYcByDeviceIdAndCodes", ctrl.GetLastYcByDeviceIdAndCodes)
	router.POST("/api/v2/auxiliary/getHistoryYcByDeviceIdCodes", ctrl.GetHistoryYcByDeviceIdCodes)
	router.POST("/api/v2/auxiliary/getEmDeviceModelCmdParamListByDeviceId", ctrl.GetEmDeviceModelCmdParamListByDeviceId)
	router.POST("/api/v2/auxiliary/getDeviceControlPointList", ctrl.GetDeviceControlPointList)

}

// /获取设备类型下的所有设备数据
func (c *AuxiliaryController) GetDeviceListByDeviceType(ctx *gin.Context) {
	type tmpQuery struct {
		DeviceType string `json:"deviceType"`
	}
	var query tmpQuery
	if err := ctx.Bind(&query); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	deviceList, err := c.repo.GetAuxiliaryDevice(query.DeviceType)
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
}

// 获取所有设备类型
func (c *AuxiliaryController) GetAuxiliaryDeviceType(ctx *gin.Context) {
	deviceTypeList, err := c.repo.GetAuxiliaryDeviceType()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"获取信息成功！",
		deviceTypeList,
	})
	return
}

// 获取最新遥测信息GetLastYcListByCode
// 遥测控制器已经写有
func (c *AuxiliaryController) GetLastYcByDeviceIdAndCodes(ctx *gin.Context) {
	type ycQueryData struct {
		DeviceIds []int `form:"deviceIds"`
		Codes     []int `form:"codes"`
	}
	var ycQuery ycQueryData
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
	//将ycQuery.DeviceIds转换成字符串
	deviceIds := utils.IntArrayToString(ycQuery.DeviceIds, ",")
	codes := utils.IntArrayToString(ycQuery.Codes, ",")
	ycLog, err := c.hisRepo.GetLastYcListByCode(deviceIds, codes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"获取信息成功！",
		ycLog,
	})
}

type ReturnMap struct {
	XAxisList []string          `json:"xAxisList"`
	DataMap   map[int][]float64 `json:"dataMap"`
}

//根据选择的codes返回对应时间的历史数据

func (c *AuxiliaryController) GetHistoryYcByDeviceIdCodes(ctx *gin.Context) {
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

type DeviceModel struct {
	Property struct {
		AccessMode int `json:"accessMode"`
	} `json:"property"`
}

// 根据设备ID获取所有模型
func (c *AuxiliaryController) GetEmDeviceModelCmdParamListByDeviceId(ctx *gin.Context) {
	var tmp struct {
		DeviceId int `json:"deviceId"`
	}
	if err := ctx.ShouldBindBodyWith(&tmp, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	array := []string{"yc"}
	var deviceCmdParamList []models.EmDeviceModelCmdParam
	deviceCmdParamList, _ = c.emRepo.GetYcListByDeviceId(tmp.DeviceId, array)
	if deviceCmdParamList == nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "无数据",
		})
		return
	}

	//收集所有codes查询最新测点信息
	var codes []string
	mapCodes := make(map[string]string) //按name进行label分组
	for i := 0; i < len(deviceCmdParamList); i++ {
		//var deviceModel DeviceModel
		//err := json.Unmarshal([]byte(deviceCmdParamList[i].Data), &deviceModel) //取出不可控的遥测
		//if err != nil {
		//	continue
		//}
		//if deviceModel.Property.AccessMode == 0 { //等于0表示不可控，不等于0表示可控
		codes = append(codes, deviceCmdParamList[i].Name)
		mapCodes[deviceCmdParamList[i].Name] = deviceCmdParamList[i].Label
		//}
	}
	//拿到所有不可控测点的最新数据
	ycList, err := c.hisRepo.GetLastYcListByCode(strconv.Itoa(tmp.DeviceId), strings.Join(codes, ","))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, model := range ycList {
		model.Name = mapCodes[strconv.Itoa(model.Code)]
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "成功",
		Data:    ycList,
	})
}

// 获取实时控制遥控遥调列表
func (c *AuxiliaryController) GetDeviceControlPointList(ctx *gin.Context) {
	type Res struct {
		YcData []ReturnModel.AuxYcData `json:"ycData"`
		YxData []ReturnModel.AuxYxData `json:"yxData"`
	}
	var tmp struct {
		DeviceId int `json:"deviceId"`
	}
	if err := ctx.ShouldBindBodyWith(&tmp, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	array := []string{"yc", "yx"}
	YkYtList, err := c.emRepo.GetYcListByDeviceId(tmp.DeviceId, array)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(YkYtList) == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "无数据",
		})
		return
	}
	//收集所有codes查询最新测点信息
	var YxCodes []string
	var YcCodes []string
	nameMap := make(map[string]models.EmDeviceModelCmdParam)
	for i := 0; i < len(YkYtList); i++ {

		if YkYtList[i].IotDataType == "yc" {
			YcCodes = append(YcCodes, YkYtList[i].Name) //收集遥测code
		} else if YkYtList[i].IotDataType == "yx" {
			YxCodes = append(YxCodes, YkYtList[i].Name) //收集遥信code
		}
		nameMap[YkYtList[i].Name] = YkYtList[i]
	}
	var result Res
	if len(YcCodes) > 0 {
		ycList, err := c.hisRepo.GetLastYcListByCode(strconv.Itoa(tmp.DeviceId), strings.Join(YcCodes, ",")) //拿到所有可控测点的最新数据
		var ycData []ReturnModel.AuxYcData
		if err == nil {
			if len(ycList) > 0 {
				for i := range ycList { //遍历结果，重新赋值实体，然后返回
					var tmpYcData ReturnModel.AuxYcData
					tmpMap := nameMap[strconv.Itoa(ycList[i].Code)] //
					if tmpMap != (models.EmDeviceModelCmdParam{}) { //如果存在对应的键值对 测点重新命名
						tmpYcData.Name = tmpMap.Label
						tmpYcData.Unit = tmpMap.Unit
					} else {
						tmpYcData.Name = ycList[i].Name
					}
					tmpYcData.Code = ycList[i].Code
					tmpYcData.Value = ycList[i].Value
					tmpYcData.Ts = ycList[i].Ts
					ycData = append(ycData, tmpYcData)
				}
			}
		}
		result.YcData = ycData
	}
	if len(YxCodes) > 0 {
		yxList, err := c.hisRepo.GetLastYxListByCode(tmp.DeviceId, strings.Join(YxCodes, ","))
		var yxData []ReturnModel.AuxYxData
		if err == nil {
			if len(yxList) > 0 {
				for i := range yxList { //遍历结果，重新赋值实体，然后返回
					var tmpYxData ReturnModel.AuxYxData
					tmpMap := nameMap[strconv.Itoa(yxList[i].Code)] //
					if tmpMap != (models.EmDeviceModelCmdParam{}) { //如果存在对应的键值对 测点重新命名
						tmpYxData.Name = tmpMap.Label
						tmpYxData.Unit = tmpMap.Unit
					} else {
						tmpYxData.Name = yxList[i].Name
					}
					tmpYxData.Code = yxList[i].Code
					tmpYxData.Value = yxList[i].Value
					tmpYxData.Ts = yxList[i].Ts
					yxData = append(yxData, tmpYxData)
				}
			}
		}
		result.YxData = yxData
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取数据成功！",
		Data:    result,
	})
	return
}
