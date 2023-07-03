package repositories

import (
	"fmt"
	"gateway/models"
	"gateway/models/ReturnModel"
	"gateway/models/query"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type LoadTrackingRepository struct {
	db *gorm.DB
}

func NewLoadTrackingRepository() *LoadTrackingRepository {
	return &LoadTrackingRepository{
		db: models.DB,
	}
}

func GetLoadTrackingData(ycQuery query.QueryTaoData) models.LoadTrackingCharVo {
	historyDataRepository := NewHistoryDataRepository()
	dictDataRepository := NewDictDataRepository()
	deviceRepository := NewDeviceRepository()
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

	result := models.LoadTrackingCharVo{
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
				return result
			}
			loadValue = pcsValueF
			isHaveValue = true
		}

		meterValue := meterValueData[i]
		if meterValue != nil {
			meterValueF, ok := meterValue.(float64)
			if !ok {
				fmt.Println("meterValue to float64 failed")
				return result
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

	projectRepository := NewProjectInfoRepository()
	projectList, _ := projectRepository.GetAll("", "", "")
	priceConfig := models.EmConfiguration{}
	if len(projectList) > 0 {
		configurationCenterRepository := NewConfigurationCenterRepository()
		priceConfigData, _ := configurationCenterRepository.GetConfigurationByProvince(projectList[0].Province, time.Now().AddDate(0, 0, -1).Format("2006-01"))
		priceConfig = priceConfigData
	}
	result.FlatPeriod = priceConfig.FlatPeriod
	result.PeakPeriod = priceConfig.PeakPeriod
	result.TopPeriod = priceConfig.TopPeriod
	result.ValleyPeriod = priceConfig.ValleyPeriod

	return result
}
