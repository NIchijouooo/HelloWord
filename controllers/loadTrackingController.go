package controllers

import (
	"fmt"
	"gateway/httpServer/model"
	"gateway/models"
	"gateway/models/ReturnModel"
	"gateway/models/query"
	"gateway/repositories"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
	"time"
)

type LoadTrackingController struct {
	repo *repositories.LoadTrackingRepository
}

func NewLoadTrackingController() *LoadTrackingController {
	return &LoadTrackingController{
		repo: repositories.NewLoadTrackingRepository(),
	}
}

func (c *LoadTrackingController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/loadTracking/GetLoadTrackingData", c.GetLoadTrackingData)
}

type resultData struct {
	XAxisList    []string      `json:"xAxisList"`
	OnGridData   []interface{} `json:"onGridData"`
	MeterData    []interface{} `json:"meterData"`
	LoadData     []interface{} `json:"loadData"`
	TopPeriod    string        `json:"topPeriod"`
	PeakPeriod   string        `json:"peakPeriod"`
	FlatPeriod   string        `json:"flatPeriod"`
	ValleyPeriod string        `json:"valleyPeriod"`
}

func (c *LoadTrackingController) GetLoadTrackingData(ctx *gin.Context) {

	ycQuery := query.QueryTaoData{}

	if err := ctx.ShouldBindBodyWith(&ycQuery, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ycQuery.StartTime == 0 && ycQuery.EndTime == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "参数错误",
			Data:    "",
		})
		return
	}
	ycQuery.Interval = 30
	ycQuery.IntervalType = 2

	historyDataRepository := repositories.NewHistoryDataRepository()
	dictDataRepository := repositories.NewDictDataRepository()
	deviceRepository := repositories.NewDeviceRepository()
	pcsDeviceType, _ := dictDataRepository.GetDictValueByDictTypeAndDictLabel("device_type", "PCS")
	pcsPowerCode, _ := dictDataRepository.GetDictValueByDictTypeAndDictLabel("energy_product_code_setting", "energy_storage_pcs_active_power_code")
	var pcsData ReturnModel.CharData
	if pcsDeviceType != "" && pcsPowerCode != "" {
		deviceIdList, _ := deviceRepository.GetDeviceIdListByDeviceType(pcsDeviceType)
		ycQuery.DeviceIds = deviceIdList
		intPcsPowerCode, _ := strconv.Atoi(pcsPowerCode)
		ycQuery.CodeList = []int{intPcsPowerCode}
		pcsDataResult, _ := historyDataRepository.GetCharData(ycQuery)
		pcsData = pcsDataResult
	}

	meterPowerCode, _ := dictDataRepository.GetDictValueByDictTypeAndDictLabel("energy_product_code_setting", "energy_storage_ac_meter_yggl")
	meterDeviceType, _ := dictDataRepository.GetDictValueByDictTypeAndDictLabel("device_type", "电能表")
	var meterData ReturnModel.CharData
	if meterPowerCode != "" && meterDeviceType != "" {
		deviceIdList, _ := deviceRepository.GetDeviceIdListByDeviceType(meterDeviceType)
		ycQuery.DeviceIds = deviceIdList
		intMeterPowerCode, _ := strconv.Atoi(meterPowerCode)
		ycQuery.CodeList = []int{intMeterPowerCode}
		meterDataResult, _ := historyDataRepository.GetCharData(ycQuery)
		meterData = meterDataResult
	}

	xaxisList := pcsData.XAxisList
	intPcsPowerCode, _ := strconv.Atoi(pcsPowerCode)
	pcsValueData := pcsData.DataMap[intPcsPowerCode]
	intMeterPowerCode, _ := strconv.Atoi(meterPowerCode)
	meterValueData := meterData.DataMap[intMeterPowerCode]

	result := resultData{
		XAxisList: xaxisList,
	}
	result.OnGridData = pcsValueData
	result.MeterData = meterValueData
	loadData := make([]interface{}, 0)
	for i, _ := range xaxisList {
		var loadValueI interface{} = nil
		var loadValue float64 = 0
		pcsValue := pcsValueData[i]
		var isHaveValue = false
		if pcsValue != nil {
			pcsValueF, ok := pcsValue.(float64)
			if !ok {
				fmt.Println("pcsValue to float64 failed")
				return
			}
			loadValue = pcsValueF
			isHaveValue = true
		}

		meterValue := meterValueData[i]
		if meterValue != nil {
			meterValueF, ok := meterValue.(float64)
			if !ok {
				fmt.Println("meterValue to float64 failed")
				return
			}
			loadValue = meterValueF + loadValue
			isHaveValue = true
		}

		if isHaveValue {
			loadValueI = loadValue
		}
		loadData = append(loadData, loadValueI)
	}

	result.LoadData = loadData

	projectRepository := repositories.NewProjectInfoRepository()
	projectList, _ := projectRepository.GetAll("", "", "")
	priceConfig := models.EmConfiguration{}
	if len(projectList) > 0 {
		configurationCenterRepository := repositories.NewConfigurationCenterRepository()
		priceConfigData, _ := configurationCenterRepository.GetConfigurationByProvince(projectList[0].Province, time.Now().AddDate(0, 0, -1).Format("2006-01"))
		priceConfig = priceConfigData
	}
	result.FlatPeriod = priceConfig.FlatPeriod
	result.PeakPeriod = priceConfig.PeakPeriod
	result.TopPeriod = priceConfig.TopPeriod
	result.ValleyPeriod = priceConfig.ValleyPeriod
	//if err != nil {
	//	ctx.JSON(http.StatusOK, model.ResponseData{
	//		Code:    "1",
	//		Message: "error" + err.Error(),
	//		Data:    "",
	//	})
	//	return
	//}
	//result := map[string]interface{}{
	//	"webHmiPageId": webHmiPageId,
	//	"token":        token,
	//}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: result,
	})
}
