package job

import (
	"fmt"
	"gateway/models"
	"gateway/repositories"
	"strconv"
	"time"
)

// 创建一个 job 结构体，包含需要执行的方法
//type StatisticsDayChargeAndDischargeJob struct{
//	db     *gorm.DB
//}


//func NewStatisticsDayChargeAndDischargeJob() *StatisticsDayChargeAndDischargeJob {
//	return &StatisticsDayChargeAndDischargeJob{
//		db: models.DB,
//		repoDictData: repositories.NewDictDataRepository(),
//		repoRealTime: repositories.NewRealtimeDataRepository()}
//}

//func (s *StatisticsDayChargeAndDischargeJob) Start() {
//	// 创建一个新的定时任务
//	job := &StatisticsDayChargeAndDischargeJob{}
//	cron := gocron.NewScheduler()
//	// 定义任务执行的时间间隔，例如每分钟执行一次
//	cron.Every(1).Hour().Do(job)
//	// 启动定时任务
//	cron.Start()
//}

//func Run() {
//	// 需要执行的任务
//	setting.ZAPS.Debug("注册StatisticsDayChargeAndDischargeJob定时任务到 GoCron")
//
//}

type PvPowerGenerationModel struct {
	DeviceId int
	Power    float64
}

/**
定时任务调用接口，
是查什么类型的设备的什么点位，既在字典定义好,taos和iot已经有完整接口示例，搬过来，每次都是传年月日的ts更新，避免数据太多。
code：pcs_charge_discharge的设备。
code：5为充电，6为放电的点位。
创建结构体pcsChargeDischarge和taos表，
用定时任务每小时一次，去taos统计这些deviceId+code点位，昨天的充，放电量，存入taos的pcs_charge_discharge。
*/
func StatisticsDayChargingAndDischarging() {
	fmt.Println("Running StatisticsDayChargeAndDischargeJob job...")
	dictDataRepo := repositories.DictDataRepository{}
	dictDataList, _, _ := dictDataRepo.GetAll("", "energy_product_code_setting", 1, 50)

	deviceIdList := make([]int, 0)
	// 将切片转换为 map
	energyProductCodeSettingMap := make(map[string]string)
	for _, dictData := range dictDataList {
		energyProductCodeSettingMap[dictData.DictLabel] = dictData.DictValue
	}


	acMeterCategoryId := energyProductCodeSettingMap["energy_storage_ac_meter_category"]
	if acMeterCategoryId != "" {
		acMeterCategoryIdInt, _ := strconv.Atoi(acMeterCategoryId)
		acMeterProductCodeList := getDeviceTypeCodeByCategoryId(nil, acMeterCategoryIdInt)
		if acMeterProductCodeList != nil && len(acMeterProductCodeList) > 0 {
			acDeviceIdList := getDeviceIdListByProjectIdList(nil, acMeterProductCodeList)
			if acDeviceIdList != nil && len(acDeviceIdList) > 0 {
				deviceIdList = append(deviceIdList, acDeviceIdList...)
			}
		}
	}

	auxiliaryMeterCategoryId := energyProductCodeSettingMap["energy_storage_auxiliary_meter_category"]
	var auxiliaryMeterDeviceIdList []int
	if auxiliaryMeterCategoryId != "" {
		auxiliaryMeterCategoryIdInt, _ := strconv.Atoi(auxiliaryMeterCategoryId)
		auxiliaryMeterProductCodeList := getDeviceTypeCodeByCategoryId(nil, auxiliaryMeterCategoryIdInt)
		if auxiliaryMeterProductCodeList != nil && len(auxiliaryMeterProductCodeList) > 0 {
			auxiliaryMeterDeviceIdList = getDeviceIdListByProjectIdList(nil, auxiliaryMeterProductCodeList)
			if auxiliaryMeterDeviceIdList != nil && len(auxiliaryMeterDeviceIdList) > 0 {
				deviceIdList = append(deviceIdList, auxiliaryMeterDeviceIdList...)
			}
		}
	}

	if len(deviceIdList) == 0 {
		xxlJobHelperLog("statisticsDayChargingAndDischarging deviceIdList : ", deviceIdList)
		return
	}

	zxCode := energyProductCodeSettingMap["energy_storage_ac_meter_rdjzxygz"]
	zxCodeInt, _ := strconv.Atoi(zxCode)
	zxTodayValueResultMap := getTodayValueMap(deviceIdList, zxCodeInt)
	zxYesterdayValueResultMap := getYesterdayValueMap(deviceIdList, zxCodeInt)
	var zxDeviceDailyStatistics []PvPowerGenerationModel
	if zxTodayValueResultMap != nil && zxYesterdayValueResultMap != nil {
		zxDeviceDailyStatistics = processData(deviceIdList, zxYesterdayValueResultMap, zxTodayValueResultMap, nil)
	}
	zxDeviceDailyStatisticsMap := make(map[int]PvPowerGenerationModel)
	if zxDeviceDailyStatistics != nil {
		for _, vo := range zxDeviceDailyStatistics {
			zxDeviceDailyStatisticsMap[vo.DeviceId] = vo
		}
	}

	fxCode := energyProductCodeSettingMap["energy_storage_ac_meter_rdjfxygz"]
	fxCodeInt, _ := strconv.Atoi(fxCode)
	fxTodayValueResultMap := getTodayValueMap(deviceIdList, fxCodeInt)
	fxYesterdayValueResultMap := getYesterdayValueMap(deviceIdList, fxCodeInt)
	var fxDeviceDailyStatistics []PvPowerGenerationModel
	if fxTodayValueResultMap != nil && fxYesterdayValueResultMap != nil {
		fxDeviceDailyStatistics = processData(deviceIdList, fxYesterdayValueResultMap, fxTodayValueResultMap, nil)
	}
	fxDeviceDailyStatisticsMap := make(map[int]PvPowerGenerationModel)
	if fxDeviceDailyStatistics != nil {
		for _, vo := range fxDeviceDailyStatistics {
			fxDeviceDailyStatisticsMap[vo.DeviceId] = vo
		}
	}

	if fxDeviceDailyStatistics == nil && zxDeviceDailyStatisticsMap == nil {
		xxlJobHelperLog("statisticsDayChargingAndDischarging fxDeviceDailyStatistics == nil && zxDeviceDailyStatisticsMap == nil")
		return
	}

	polarity := energyProductCodeSettingMap["pv_polarity_code"]
	polarityInt, _ := strconv.Atoi(polarity)
	polarityValueList := getSettingValueListByDeviceIdListAndYcCode(deviceIdList, polarityInt)
	polarityMap := make(map[int]string)
	if polarityValueList != nil {
		for _, po := range polarityValueList {
			polarityMap[po.DeviceId] = po.Value
		}
	}

	yesterday := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, 0, 0, 0, 0, time.Local).Unix()
	insertList := make([]*models.EsChargeDischargeModel, 0)
	for _, deviceId := range deviceIdList {
		var zxPower *float64
		var fxPower *float64
		if zxPvPowerGenerationModel, ok := zxDeviceDailyStatisticsMap[deviceId]; ok {
			zxPower = &zxPvPowerGenerationModel.Power
		}
		if fxPvPowerGenerationModel, ok := fxDeviceDailyStatisticsMap[deviceId]; ok {
			fxPower = &fxPvPowerGenerationModel.Power
		}

		if zxPower == nil && fxPower == nil {
			continue
		}

		esChargeDischargeModel := &models.EsChargeDischargeModel{
			DeviceId: deviceId,
			Ts:       yesterday,
		}

		var chargeCapacity *float64
		var dischargeCapacity *float64
		if polarityMap == nil || polarityMap[deviceId] == "0" || polarityMap[deviceId] != "1" {
			chargeCapacity = zxPower
			dischargeCapacity = fxPower
		} else {
			chargeCapacity = zxPower
			dischargeCapacity = fxPower
		}

		chargeCapacityL := *chargeCapacity
		chargeCapacityDisL := *dischargeCapacity

		esChargeDischargeModel.ChargeCapacity = chargeCapacityL
		esChargeDischargeModel.DischargeCapacity = chargeCapacityDisL

		insertList = append(insertList, esChargeDischargeModel)
	}

	if len(insertList) > 0 {
		realtimeDataRepo := repositories.RealtimeDataRepository{}
		realtimeDataRepo.BatchCreateEsChargeDischargeModel(insertList)
	}
}

func getDeviceTypeCodeByCategoryId(categoryId []int, acMeterCategoryId int) []int {
	// TODO: Implement getDeviceTypeCodeByCategoryId
	return nil
}

func getDeviceIdListByProjectIdList(projectIdList []int, acMeterProductCodeList []int) []int {
	// TODO: Implement getDeviceIdListByProjectIdList
	return nil
}

func getTodayValueMap(deviceIdList []int, code int) map[int]float64 {
	// TODO: Implement getTodayValueMap
	return nil
}

func getYesterdayValueMap(deviceIdList []int, code int) map[int]float64 {
	// TODO: Implement getYesterdayValueMap
	return nil
}

func processData(deviceIdList []int, yesterdayValueResultMap map[int]float64, todayValueResultMap map[int]float64, param interface{}) []PvPowerGenerationModel {
	// TODO: Implement processData
	return nil
}

func getSettingValueListByDeviceIdListAndYcCode(deviceIdList []int, code int) []models.SettingData {
	// TODO: Implement getSettingValueListByDeviceIdListAndYcCode
	return nil
}

func xxlJobHelperLog(message ...interface{}) {
	// TODO: Implement XxlJobHelperLog
}

