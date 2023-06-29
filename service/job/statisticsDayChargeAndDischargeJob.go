package job

import (
	"fmt"
	"gateway/models"
	"gateway/repositories"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
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
	Power    float64
	Profit   float64
	DeviceId int
	Ts       int64
}

// StatisticsHourChargingAndDischarging 小时充放电量收益统计
func StatisticsHourChargingAndDischarging() {
	dictDataRepository := repositories.NewDictDataRepository()
	deviceRepository := repositories.NewDevicePointRepository()
	deviceEquipmentRepository := repositories.NewDeviceEquipmentRepository()
	configurationCenterRepository := repositories.NewConfigurationCenterRepository()
	projectRepository := repositories.NewProjectInfoRepository()

	dictDataList, _ := dictDataRepository.GetDictDataByDictType("energy_product_code_setting")

	// 将切片转换为 map
	energyProductCodeSettingMap := make(map[string]string)
	for _, dictData := range dictDataList {
		energyProductCodeSettingMap[dictData.DictLabel] = dictData.DictValue
	}

	deviceType, _ := dictDataRepository.SelectDictValue("device_type", "电能表")
	if deviceType.DictValue == "" {
		fmt.Println("statisticsDayChargingAndDischarging deviceType is null")
		return
	}

	deviceIdList := deviceRepository.GetDeviceIdListByDevType(deviceType.DictValue)
	if deviceIdList == nil || len(deviceIdList) == 0 {
		fmt.Println("statisticsDayChargingAndDischarging deviceIdList : ", deviceIdList)
		return
	}

	stringDevIds := make([]string, len(deviceIdList))
	for i, v := range deviceIdList {
		stringDevIds[i] = strconv.Itoa(v)
	}
	strDeviceIdList := strings.Join(stringDevIds, ",")
	deviceEquipmentList, _ := deviceEquipmentRepository.GetEquipmentInfoByDevIdList(strDeviceIdList)

	// 将切片转换为 map
	deviceMeterMagnificationMap := make(map[int]int)
	deviceMeterReadFlipMap := make(map[int]int)
	polarityMap := make(map[int]int)
	for _, deviceEquipment := range deviceEquipmentList {
		deviceMeterMagnificationMap[deviceEquipment.DeviceId] = deviceEquipment.MeterMagnification
		deviceMeterReadFlipMap[deviceEquipment.DeviceId] = deviceEquipment.MeterReadFlip
		polarityMap[deviceEquipment.DeviceId] = deviceEquipment.Polarity
	}

	projectList, _ := projectRepository.GetAll("", "", "")
	priceConfig := models.EmConfiguration{}
	if len(projectList) > 0 {
		priceConfigData, _ := configurationCenterRepository.GetConfigurationByProvince(projectList[0].Province, time.Now().AddDate(0, 0, -1).Format("2006-01"))
		priceConfig = priceConfigData
	}
	deviceSettingMap := make(map[string]map[int]int)
	deviceSettingMap["deviceMeterMagnificationMap"] = deviceMeterMagnificationMap
	deviceSettingMap["deviceMeterReadFlipMap"] = deviceMeterReadFlipMap

	currentTime := time.Now()
	year, month, day := currentTime.Date()
	hour, _, _ := currentTime.Clock()
	zeroTime := time.Date(year, month, day, hour, 0, 0, 0, currentTime.Location())
	previousHourStartDate := zeroTime.Add(-time.Hour).UnixMilli()
	previousHourEndDate := previousHourStartDate + 3599999

	//当前正向有功尖
	dqzxygjCode := energyProductCodeSettingMap["energy_storage_ac_meter_dqzxygj"]
	dqzxygjDataList := getHourPower(deviceIdList, deviceSettingMap, priceConfig.TopPrice, dqzxygjCode,
		previousHourStartDate, previousHourEndDate)
	var rdjzxygjMap = make(map[int]PvPowerGenerationModel)
	for _, rdjzxygjData := range dqzxygjDataList {
		deviceId := rdjzxygjData.DeviceId
		rdjzxygjMap[deviceId] = rdjzxygjData
	}

	//当前正向有功峰
	dqzxygfCode := energyProductCodeSettingMap["energy_storage_ac_meter_dqzxygf"]
	dqzxygfDataList := getHourPower(deviceIdList, deviceSettingMap, priceConfig.PeakPrice, dqzxygfCode,
		previousHourStartDate, previousHourEndDate)
	var rdjzxygfMap = make(map[int]PvPowerGenerationModel)
	for _, rdjzxygfData := range dqzxygfDataList {
		deviceId := rdjzxygfData.DeviceId
		rdjzxygfMap[deviceId] = rdjzxygfData
	}

	//当前正向有功平
	dqzxygpCode := energyProductCodeSettingMap["energy_storage_ac_meter_dqzxygp"]
	dqzxygpDataList := getHourPower(deviceIdList, deviceSettingMap, priceConfig.FlatPrice, dqzxygpCode,
		previousHourStartDate, previousHourEndDate)
	var rdjzxygpMap = make(map[int]PvPowerGenerationModel)
	for _, rdjzxygpData := range dqzxygpDataList {
		deviceId := rdjzxygpData.DeviceId
		rdjzxygpMap[deviceId] = rdjzxygpData
	}

	//当前正向有功谷
	dqzxyggCode := energyProductCodeSettingMap["energy_storage_ac_meter_dqzxygg"]
	dqzxyggDataList := getHourPower(deviceIdList, deviceSettingMap, priceConfig.ValleyPrice, dqzxyggCode,
		previousHourStartDate, previousHourEndDate)
	var rdjzxyggMap = make(map[int]PvPowerGenerationModel)
	for _, rdjzxyggData := range dqzxyggDataList {
		deviceId := rdjzxyggData.DeviceId
		rdjzxyggMap[deviceId] = rdjzxyggData
	}

	//当前反向有功尖
	dqfxygjCode := energyProductCodeSettingMap["energy_storage_ac_meter_dqfxygj"]
	dqfxygjDataList := getHourPower(deviceIdList, deviceSettingMap, priceConfig.TopPrice, dqfxygjCode,
		previousHourStartDate, previousHourEndDate)
	var rdjfxygjMap = make(map[int]PvPowerGenerationModel)
	for _, rdjfxygjData := range dqfxygjDataList {
		deviceId := rdjfxygjData.DeviceId
		rdjfxygjMap[deviceId] = rdjfxygjData
	}

	//当前反向有功峰
	dqfxygfCode := energyProductCodeSettingMap["energy_storage_ac_meter_dqfxygf"]
	dqfxygfDataList := getHourPower(deviceIdList, deviceSettingMap, priceConfig.PeakPrice, dqfxygfCode,
		previousHourStartDate, previousHourEndDate)
	var rdjfxygfMap = make(map[int]PvPowerGenerationModel)
	for _, rdjfxygfData := range dqfxygfDataList {
		deviceId := rdjfxygfData.DeviceId
		rdjfxygfMap[deviceId] = rdjfxygfData
	}

	//当前反向有功平
	dqfxygpCode := energyProductCodeSettingMap["energy_storage_ac_meter_dqfxygp"]
	dqfxygpDataList := getHourPower(deviceIdList, deviceSettingMap, priceConfig.FlatPrice, dqfxygpCode,
		previousHourStartDate, previousHourEndDate)
	var rdjfxygpMap = make(map[int]PvPowerGenerationModel)
	for _, rdjfxygpData := range dqfxygpDataList {
		deviceId := rdjfxygpData.DeviceId
		rdjfxygpMap[deviceId] = rdjfxygpData
	}

	//当前反向有功谷
	dqfxyggCode := energyProductCodeSettingMap["energy_storage_ac_meter_dqfxygg"]
	dqfxyggDataList := getHourPower(deviceIdList, deviceSettingMap, priceConfig.ValleyPrice, dqfxyggCode,
		previousHourStartDate, previousHourEndDate)
	var rdjfxyggMap = make(map[int]PvPowerGenerationModel)
	for _, rdjfxyggData := range dqfxyggDataList {
		deviceId := rdjfxyggData.DeviceId
		rdjfxyggMap[deviceId] = rdjfxyggData
	}

	dataMap := make(map[string]map[int]PvPowerGenerationModel)
	dataMap["zxygjMap"] = rdjzxygjMap
	dataMap["zxygfMap"] = rdjzxygfMap
	dataMap["zxygpMap"] = rdjzxygpMap
	dataMap["zxyggMap"] = rdjzxyggMap
	dataMap["fxygjMap"] = rdjfxygjMap
	dataMap["fxygfMap"] = rdjfxygfMap
	dataMap["fxygpMap"] = rdjfxygpMap
	dataMap["fxyggMap"] = rdjfxyggMap

	insertList := statisticsChargingAndDischarging(deviceIdList, polarityMap, dataMap, previousHourStartDate)

	if len(insertList) > 0 {
		fmt.Println("insertList")
		realtimeDataRepository := repositories.NewRealtimeDataRepository()
		realtimeDataRepository.BatchCreateEsHourChargeDischargeModel(insertList)
	}
}

// StatisticsDayChargingAndDischarging 日充放电量收益统计
func StatisticsDayChargingAndDischarging() {
	dictDataRepository := repositories.NewDictDataRepository()
	deviceRepository := repositories.NewDevicePointRepository()
	deviceEquipmentRepository := repositories.NewDeviceEquipmentRepository()
	configurationCenterRepository := repositories.NewConfigurationCenterRepository()
	projectRepository := repositories.NewProjectInfoRepository()

	dictDataList, _ := dictDataRepository.GetDictDataByDictType("energy_product_code_setting")

	// 将切片转换为 map
	energyProductCodeSettingMap := make(map[string]string)
	for _, dictData := range dictDataList {
		energyProductCodeSettingMap[dictData.DictLabel] = dictData.DictValue
	}

	deviceType, _ := dictDataRepository.SelectDictValue("device_type", "电能表")
	if deviceType.DictValue == "" {
		fmt.Println("statisticsDayChargingAndDischarging deviceType is null")
		return
	}

	deviceIdList := deviceRepository.GetDeviceIdListByDevType(deviceType.DictValue)
	if deviceIdList == nil || len(deviceIdList) == 0 {
		fmt.Println("statisticsDayChargingAndDischarging deviceIdList : ", deviceIdList)
		return
	}

	stringDevIds := make([]string, len(deviceIdList))
	for i, v := range deviceIdList {
		stringDevIds[i] = strconv.Itoa(v)
	}
	strDeviceIdList := strings.Join(stringDevIds, ",")
	deviceEquipmentList, _ := deviceEquipmentRepository.GetEquipmentInfoByDevIdList(strDeviceIdList)

	// 将切片转换为 map
	deviceMeterMagnificationMap := make(map[int]int)
	deviceMeterReadFlipMap := make(map[int]int)
	polarityMap := make(map[int]int)
	for _, deviceEquipment := range deviceEquipmentList {
		deviceMeterMagnificationMap[deviceEquipment.DeviceId] = deviceEquipment.MeterMagnification
		deviceMeterReadFlipMap[deviceEquipment.DeviceId] = deviceEquipment.MeterReadFlip
		polarityMap[deviceEquipment.DeviceId] = deviceEquipment.Polarity
	}

	projectList, _ := projectRepository.GetAll("", "", "")
	priceConfig := models.EmConfiguration{}
	if len(projectList) > 0 {
		priceConfigData, _ := configurationCenterRepository.GetConfigurationByProvince(projectList[0].Province, time.Now().AddDate(0, 0, -1).Format("2006-01"))
		priceConfig = priceConfigData
	}
	deviceSettingMap := make(map[string]map[int]int)
	deviceSettingMap["deviceMeterMagnificationMap"] = deviceMeterMagnificationMap
	deviceSettingMap["deviceMeterReadFlipMap"] = deviceMeterReadFlipMap

	currentTime := time.Now()
	year, month, day := currentTime.Date()
	zeroTime := time.Date(year, month, day, 0, 0, 0, 0, currentTime.Location())
	todayStartDate := zeroTime.UnixMilli()
	todayEndDate := todayStartDate + 86399999
	yesterdayStartDate := zeroTime.AddDate(0, 0, -1).UnixMilli()
	yesterdayEndDate := yesterdayStartDate + 86399999

	//日冻结正向有功尖
	rdjzxygjCode := energyProductCodeSettingMap["energy_storage_ac_meter_rdjzxygj"]
	rdjzxygjDataList := getDayPower(deviceIdList, deviceSettingMap, priceConfig.TopPrice, rdjzxygjCode,
		todayStartDate, todayEndDate, yesterdayStartDate, yesterdayEndDate)
	var rdjzxygjMap = make(map[int]PvPowerGenerationModel)
	for _, rdjzxygjData := range rdjzxygjDataList {
		deviceId := rdjzxygjData.DeviceId
		rdjzxygjMap[deviceId] = rdjzxygjData
	}

	//日冻结正向有功峰
	rdjzxygfCode := energyProductCodeSettingMap["energy_storage_ac_meter_rdjzxygf"]
	rdjzxygfDataList := getDayPower(deviceIdList, deviceSettingMap, priceConfig.PeakPrice, rdjzxygfCode,
		todayStartDate, todayEndDate, yesterdayStartDate, yesterdayEndDate)
	var rdjzxygfMap = make(map[int]PvPowerGenerationModel)
	for _, rdjzxygfData := range rdjzxygfDataList {
		deviceId := rdjzxygfData.DeviceId
		rdjzxygfMap[deviceId] = rdjzxygfData
	}

	//日冻结正向有功平
	rdjzxygpCode := energyProductCodeSettingMap["energy_storage_ac_meter_rdjzxygp"]
	rdjzxygpDataList := getDayPower(deviceIdList, deviceSettingMap, priceConfig.FlatPrice, rdjzxygpCode,
		todayStartDate, todayEndDate, yesterdayStartDate, yesterdayEndDate)
	var rdjzxygpMap = make(map[int]PvPowerGenerationModel)
	for _, rdjzxygpData := range rdjzxygpDataList {
		deviceId := rdjzxygpData.DeviceId
		rdjzxygpMap[deviceId] = rdjzxygpData
	}

	//日冻结正向有功谷
	rdjzxyggCode := energyProductCodeSettingMap["energy_storage_ac_meter_rdjzxygg"]
	rdjzxyggDataList := getDayPower(deviceIdList, deviceSettingMap, priceConfig.ValleyPrice, rdjzxyggCode,
		todayStartDate, todayEndDate, yesterdayStartDate, yesterdayEndDate)
	var rdjzxyggMap = make(map[int]PvPowerGenerationModel)
	for _, rdjzxyggData := range rdjzxyggDataList {
		deviceId := rdjzxyggData.DeviceId
		rdjzxyggMap[deviceId] = rdjzxyggData
	}

	//日冻结反向有功尖
	rdjfxygjCode := energyProductCodeSettingMap["energy_storage_ac_meter_rdjfxygj"]
	rdjfxygjDataList := getDayPower(deviceIdList, deviceSettingMap, priceConfig.TopPrice, rdjfxygjCode,
		todayStartDate, todayEndDate, yesterdayStartDate, yesterdayEndDate)
	var rdjfxygjMap = make(map[int]PvPowerGenerationModel)
	for _, rdjfxygjData := range rdjfxygjDataList {
		deviceId := rdjfxygjData.DeviceId
		rdjfxygjMap[deviceId] = rdjfxygjData
	}

	//日冻结反向有功峰
	rdjfxygfCode := energyProductCodeSettingMap["energy_storage_ac_meter_rdjfxygf"]
	rdjfxygfDataList := getDayPower(deviceIdList, deviceSettingMap, priceConfig.PeakPrice, rdjfxygfCode,
		todayStartDate, todayEndDate, yesterdayStartDate, yesterdayEndDate)
	var rdjfxygfMap = make(map[int]PvPowerGenerationModel)
	for _, rdjfxygfData := range rdjfxygfDataList {
		deviceId := rdjfxygfData.DeviceId
		rdjfxygfMap[deviceId] = rdjfxygfData
	}

	//日冻结反向有功平
	rdjfxygpCode := energyProductCodeSettingMap["energy_storage_ac_meter_rdjfxygp"]
	rdjfxygpDataList := getDayPower(deviceIdList, deviceSettingMap, priceConfig.FlatPrice, rdjfxygpCode,
		todayStartDate, todayEndDate, yesterdayStartDate, yesterdayEndDate)
	var rdjfxygpMap = make(map[int]PvPowerGenerationModel)
	for _, rdjfxygpData := range rdjfxygpDataList {
		deviceId := rdjfxygpData.DeviceId
		rdjfxygpMap[deviceId] = rdjfxygpData
	}

	//日冻结反向有功谷
	rdjfxyggCode := energyProductCodeSettingMap["energy_storage_ac_meter_rdjfxygg"]
	rdjfxyggDataList := getDayPower(deviceIdList, deviceSettingMap, priceConfig.ValleyPrice, rdjfxyggCode,
		todayStartDate, todayEndDate, yesterdayStartDate, yesterdayEndDate)
	var rdjfxyggMap = make(map[int]PvPowerGenerationModel)
	for _, rdjfxyggData := range rdjfxyggDataList {
		deviceId := rdjfxyggData.DeviceId
		rdjfxyggMap[deviceId] = rdjfxyggData
	}

	dataMap := make(map[string]map[int]PvPowerGenerationModel)
	dataMap["zxygjMap"] = rdjzxygjMap
	dataMap["zxygfMap"] = rdjzxygfMap
	dataMap["zxygpMap"] = rdjzxygpMap
	dataMap["zxyggMap"] = rdjzxyggMap
	dataMap["fxygjMap"] = rdjfxygjMap
	dataMap["fxygfMap"] = rdjfxygfMap
	dataMap["fxygpMap"] = rdjfxygpMap
	dataMap["fxyggMap"] = rdjfxyggMap

	yesterday := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, 0, 0, 0, 0, time.Local).UnixMilli()
	insertList := statisticsChargingAndDischarging(deviceIdList, polarityMap, dataMap, yesterday)

	if len(insertList) > 0 {
		realtimeDataRepository := repositories.NewRealtimeDataRepository()
		realtimeDataRepository.BatchCreateEsChargeDischargeModel(insertList)
	}
}

func statisticsChargingAndDischarging(deviceIdList []int, polarityMap map[int]int, dataMap map[string]map[int]PvPowerGenerationModel, dataTime int64) []*models.EsChargeDischargeModel {

	zxygjMap := dataMap["zxygjMap"]
	zxygfMap := dataMap["zxygfMap"]
	zxygpMap := dataMap["zxygpMap"]
	zxyggMap := dataMap["zxyggMap"]
	fxygjMap := dataMap["fxygjMap"]
	fxygfMap := dataMap["fxygfMap"]
	fxygpMap := dataMap["fxygpMap"]
	fxyggMap := dataMap["fxyggMap"]

	insertList := make([]*models.EsChargeDischargeModel, 0)
	for _, deviceId := range deviceIdList {

		esChargeDischargeModel := &models.EsChargeDischargeModel{
			DeviceId: deviceId,
			Ts:       dataTime,
		}
		totalCharge := float64(0)
		totalDischarge := float64(0)
		totalPrice := decimal.NewFromFloat(0)
		//尖
		topChargeCapacity, topDischargeCapacity, topPrice := statisticalQuantityOfElectricity(polarityMap, deviceId, zxygjMap, fxygjMap)
		if topChargeCapacity != nil || topDischargeCapacity != nil {
			esChargeDischargeModel.TopChargeCapacity = *topChargeCapacity
			esChargeDischargeModel.TopDischargeCapacity = *topDischargeCapacity
			esChargeDischargeModel.TopProfit = topPrice
			totalPrice = totalPrice.Add(topPrice)
			totalCharge = totalCharge + *topChargeCapacity
			totalDischarge = totalDischarge + *topDischargeCapacity
		}

		//峰
		peakChargeCapacity, peakDischargeCapacity, peakPrice := statisticalQuantityOfElectricity(polarityMap, deviceId, zxygfMap, fxygfMap)
		if peakChargeCapacity != nil || peakDischargeCapacity != nil {
			esChargeDischargeModel.PeakChargeCapacity = *peakChargeCapacity
			esChargeDischargeModel.PeakDischargeCapacity = *peakDischargeCapacity
			esChargeDischargeModel.PeakProfit = peakPrice
			totalPrice = totalPrice.Add(peakPrice)
			totalCharge = totalCharge + *peakChargeCapacity
			totalDischarge = totalDischarge + *peakDischargeCapacity
		}

		//平
		flatChargeCapacity, flatDischargeCapacity, flatPrice := statisticalQuantityOfElectricity(polarityMap, deviceId, zxygpMap, fxygpMap)
		if flatChargeCapacity != nil || flatDischargeCapacity != nil {
			esChargeDischargeModel.FlatChargeCapacity = *flatChargeCapacity
			esChargeDischargeModel.FlatDischargeCapacity = *flatDischargeCapacity
			esChargeDischargeModel.FlatProfit = flatPrice
			totalPrice = totalPrice.Add(flatPrice)
			totalCharge = totalCharge + *flatChargeCapacity
			totalDischarge = totalDischarge + *flatDischargeCapacity
		}

		//谷
		valleyChargeCapacity, valleyDischargeCapacity, valleyPrice := statisticalQuantityOfElectricity(polarityMap, deviceId, zxyggMap, fxyggMap)
		if valleyChargeCapacity != nil || valleyDischargeCapacity != nil {
			esChargeDischargeModel.ValleyChargeCapacity = *valleyChargeCapacity
			esChargeDischargeModel.ValleyDischargeCapacity = *valleyDischargeCapacity
			esChargeDischargeModel.ValleyProfit = valleyPrice
			totalPrice = totalPrice.Add(valleyPrice)
			totalCharge = totalCharge + *valleyChargeCapacity
			totalDischarge = totalDischarge + *valleyDischargeCapacity
		}
		esChargeDischargeModel.ChargeCapacity = totalCharge
		esChargeDischargeModel.DischargeCapacity = totalDischarge
		esChargeDischargeModel.Profit = totalPrice
		insertList = append(insertList, esChargeDischargeModel)
	}

	return insertList
}

// 统计充放电量以及收益
func statisticalQuantityOfElectricity(polarityMap map[int]int, deviceId int, rdjzxMap map[int]PvPowerGenerationModel, rdjfxMap map[int]PvPowerGenerationModel) (*float64, *float64, decimal.Decimal) {
	var zxPower *float64
	var zxPrice *float64
	var fxPower *float64
	var fxPrice *float64
	if rdjzxData, ok := rdjzxMap[deviceId]; ok {
		zxPower = &rdjzxData.Power
		zxPrice = &rdjzxData.Profit
	}
	if rdjfxData, ok := rdjfxMap[deviceId]; ok {
		fxPower = &rdjfxData.Power
		fxPrice = &rdjfxData.Profit
	}

	if zxPower == nil || fxPower == nil {
		return nil, nil, decimal.NewFromFloat(0)
	}

	if polarityMap == nil || polarityMap[deviceId] != 1 {
		chargeCapacity := zxPower
		dischargeCapacity := fxPower
		price := decimal.NewFromFloat(*fxPrice).Sub(decimal.NewFromFloat(*zxPrice))
		return chargeCapacity, dischargeCapacity, price
	} else {
		chargeCapacity := fxPower
		dischargeCapacity := zxPower
		price := decimal.NewFromFloat(*zxPrice).Sub(decimal.NewFromFloat(*fxPrice))
		return chargeCapacity, dischargeCapacity, price
	}
}

func getHourPower(deviceIdList []int, deviceSettingMap map[string]map[int]int, price decimal.Decimal, code string,
	startTime int64, endTime int64) []PvPowerGenerationModel {
	intCode, _ := strconv.Atoi(code)
	lastValueResultMap := getLastValueMap(deviceIdList, intCode, startTime, endTime)
	firstValueResultMap := getFirstValueMap(deviceIdList, intCode, startTime, endTime)
	return statisticalPower(deviceIdList, deviceSettingMap, price, lastValueResultMap, firstValueResultMap)
}

func getDayPower(deviceIdList []int, deviceSettingMap map[string]map[int]int, price decimal.Decimal, code string,
	currentValueStartTime int64, currentValueEndTime int64, previousValueStartTime int64, previousValueEndTime int64) []PvPowerGenerationModel {
	intCode, _ := strconv.Atoi(code)
	currentValueResultMap := getLastValueMap(deviceIdList, intCode, currentValueStartTime, currentValueEndTime)
	previousValueResultMap := getLastValueMap(deviceIdList, intCode, previousValueStartTime, previousValueEndTime)
	return statisticalPower(deviceIdList, deviceSettingMap, price, currentValueResultMap, previousValueResultMap)
}

// 获取点位值
func statisticalPower(deviceIdList []int, deviceSettingMap map[string]map[int]int, price decimal.Decimal, currentValueResultMap map[int]float64, previousValueResultMap map[int]float64) []PvPowerGenerationModel {

	var dailyStatistics []PvPowerGenerationModel
	if currentValueResultMap != nil && previousValueResultMap != nil {
		dailyStatistics = processData(deviceIdList, deviceSettingMap, price, previousValueResultMap, currentValueResultMap)
	}
	return dailyStatistics
}

// 封装点位第一条数据
func getFirstValueMap(DeviceIdList []int, code int, startTime int64, endTime int64) map[int]float64 {
	valueList := getFirstYcHistoryByDeviceIdListAndCodeList(DeviceIdList, []int{code}, startTime, endTime)
	if valueList == nil || len(valueList) == 0 {
		return nil
	}

	result := make(map[int]float64)
	for _, ycModel := range valueList {
		result[ycModel.DeviceId] = ycModel.Value
	}

	return result
}

// 封装点位最后一条数据
func getLastValueMap(DeviceIdList []int, code int, startTime int64, endTime int64) map[int]float64 {
	valueList := getLastYcHistoryByDeviceIdListAndCodeList(DeviceIdList, []int{code}, startTime, endTime)
	if valueList == nil || len(valueList) == 0 {
		return nil
	}

	result := make(map[int]float64)
	for _, ycModel := range valueList {
		result[ycModel.DeviceId] = ycModel.Value
	}

	return result
}

// 统计电量
func processData(DeviceIdList []int, deviceSettingMap map[string]map[int]int, price decimal.Decimal, yesterdayValueMap map[int]float64, todayValueMap map[int]float64) []PvPowerGenerationModel {
	if len(yesterdayValueMap) == 0 || len(todayValueMap) == 0 {
		return nil
	}

	magnMap := deviceSettingMap["deviceMeterMagnificationMap"]
	fzMap := deviceSettingMap["deviceMeterReadFlipMap"]

	insertList := []PvPowerGenerationModel{}
	yesterday := time.Now().AddDate(0, 0, -1).Unix()
	for _, DeviceId := range DeviceIdList {
		deviceYesterdayValue := yesterdayValueMap[DeviceId]
		deviceTodayValue := todayValueMap[DeviceId]
		if deviceYesterdayValue == 0 || deviceTodayValue == 0 {
			fmt.Println("deviceYesterdayValue == 0 || deviceTodayValue == 0. DeviceId :", DeviceId)
			continue
		}

		yesterdayValueBd := decimal.NewFromFloat(float64(deviceYesterdayValue))
		todayValueBd := decimal.NewFromFloat(float64(deviceTodayValue))

		// 获取倍率
		iMagn := decimal.NewFromInt(int64(1))
		if magnMap[DeviceId] != 0 {
			iMagn, _ = decimal.NewFromString(string(magnMap[DeviceId]))
		}

		// 获取示数翻转值
		iFz, _ := strconv.Atoi(string(fzMap[DeviceId]))

		// 计算电量，包含翻转逻辑
		diffBd := getElectricityCalculation(todayValueBd, yesterdayValueBd, iFz)
		devicePower := diffBd.Mul(iMagn)

		// 计算收益
		deviceProfit := decimal.NewFromInt(int64(0))
		fmt.Println()
		if !price.Equal(decimal.Decimal{}) {
			deviceProfit = devicePower.Mul(price)
		}

		pvPowerGenerationModel := PvPowerGenerationModel{
			Power:    devicePower.InexactFloat64(),
			Profit:   deviceProfit.InexactFloat64(),
			DeviceId: DeviceId,
			Ts:       yesterday,
		}
		insertList = append(insertList, pvPowerGenerationModel)
	}
	if len(insertList) == 0 {
		return nil
	}
	return insertList
}

// 获取遥测时间区间内最后一条数据
func getFirstYcHistoryByDeviceIdListAndCodeList(devIDList []int, codeList []int, beginDt int64, endDt int64) []models.YcData {
	if len(devIDList) == 0 || len(codeList) == 0 || beginDt == 0 || endDt == 0 {
		return []models.YcData{}
	}

	idSet := make(map[int]bool)
	for _, id := range devIDList {
		idSet[id] = true
	}

	stringCodes := make([]string, len(codeList))
	for i, v := range codeList {
		stringCodes[i] = strconv.Itoa(v)
	}
	codes := strings.Join(stringCodes, ",")

	result := []models.YcData{}
	realtimeDataRepository := repositories.NewRealtimeDataRepository()
	modelList, _ := realtimeDataRepository.GetFirstYcHistoryByDeviceIdListAndCodeList(codes, beginDt, endDt)
	if len(modelList) != 0 {
		for _, ycModel := range modelList {
			DeviceId := ycModel.DeviceId
			if idSet[DeviceId] {
				result = append(result, ycModel)
			}
		}
	}
	return result
}

// 获取遥测时间区间内最后一条数据
func getLastYcHistoryByDeviceIdListAndCodeList(devIDList []int, codeList []int, beginDt int64, endDt int64) []models.YcData {
	if len(devIDList) == 0 || len(codeList) == 0 || beginDt == 0 || endDt == 0 {
		return []models.YcData{}
	}

	idSet := make(map[int]bool)
	for _, id := range devIDList {
		idSet[id] = true
	}

	stringCodes := make([]string, len(codeList))
	for i, v := range codeList {
		stringCodes[i] = strconv.Itoa(v)
	}
	codes := strings.Join(stringCodes, ",")

	result := []models.YcData{}
	realtimeDataRepository := repositories.NewRealtimeDataRepository()
	modelList, _ := realtimeDataRepository.GetLastYcHistoryByDeviceIdListAndCodeList(codes, beginDt, endDt)
	if len(modelList) != 0 {
		for _, ycModel := range modelList {
			DeviceId := ycModel.DeviceId
			if idSet[DeviceId] {
				result = append(result, ycModel)
			}
		}
	}
	return result
}

// 电量计算，包含示数反转逻辑
func getElectricityCalculation(currentValue, lastValue decimal.Decimal, flipIndicator int) decimal.Decimal {
	subtract := currentValue.Sub(lastValue)
	if flipIndicator <= 0 {
		return subtract
	}

	flipIndicatorBd := decimal.NewFromFloat(float64(flipIndicator))
	zero := decimal.NewFromFloat(float64(0))
	onePercent := decimal.NewFromFloat(float64(0.01))

	// 差值小于0 并且 翻转示数减上一次的值小于翻转值的1%，则判断为翻转
	if subtract.Cmp(zero) < 0 && flipIndicatorBd.Sub(lastValue).Cmp(flipIndicatorBd.Mul(onePercent)) < 0 {
		return flipIndicatorBd.Sub(lastValue).Add(currentValue)
	} else {
		return subtract
	}
}
