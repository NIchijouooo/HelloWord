package job

import (
	"fmt"
	"gateway/models"
	"gateway/repositories"
	"gateway/setting"
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

/**
定时任务调用接口，
是查什么类型的设备的什么点位，既在字典定义好,taos和iot已经有完整接口示例，搬过来，每次都是传年月日的ts更新，避免数据太多。
code：pcs_charge_discharge的设备。
code：5为充电，6为放电的点位。
创建结构体pcsChargeDischarge和taos表，
用定时任务每小时一次，去taos统计这些DeviceId+code点位，昨天的充，放电量，存入taos的pcs_charge_discharge。
*/
func StatisticsDayChargingAndDischarging() {
	fmt.Println("Running StatisticsDayChargeAndDischargeJob job...")
	dictDataRepo := repositories.DictDataRepository{}
	dictDataList, _, _ := dictDataRepo.GetAll("", "energy_product_code_setting", 1, 500)

	DeviceIdList := make([]int, 0)
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
				DeviceIdList = append(DeviceIdList, acDeviceIdList...)
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
				DeviceIdList = append(DeviceIdList, auxiliaryMeterDeviceIdList...)
			}
		}
	}

	if len(DeviceIdList) == 0 {
		setting.ZAPS.Debug("statisticsDayChargingAndDischarging DeviceIdList : ", DeviceIdList)
		return
	}

	zxCode := energyProductCodeSettingMap["energy_storage_ac_meter_rdjzxygz"]
	zxCodeInt, _ := strconv.Atoi(zxCode)
	zxTodayValueResultMap := GetTodayValueMap(DeviceIdList, zxCodeInt)
	zxYesterdayValueResultMap := GetYesterdayValueMap(DeviceIdList, zxCodeInt)
	var zxDeviceDailyStatistics []PvPowerGenerationModel
	if zxTodayValueResultMap != nil && zxYesterdayValueResultMap != nil {
		zxDeviceDailyStatistics = ProcessData(DeviceIdList, zxYesterdayValueResultMap, zxTodayValueResultMap, nil)
	}
	zxDeviceDailyStatisticsMap := make(map[int]PvPowerGenerationModel)
	if zxDeviceDailyStatistics != nil {
		for _, vo := range zxDeviceDailyStatistics {
			zxDeviceDailyStatisticsMap[vo.DeviceId] = vo
		}
	}

	fxCode := energyProductCodeSettingMap["energy_storage_ac_meter_rdjfxygz"]
	fxCodeInt, _ := strconv.Atoi(fxCode)
	fxTodayValueResultMap := GetTodayValueMap(DeviceIdList, fxCodeInt)
	fxYesterdayValueResultMap := GetYesterdayValueMap(DeviceIdList, fxCodeInt)
	var fxDeviceDailyStatistics []PvPowerGenerationModel
	if fxTodayValueResultMap != nil && fxYesterdayValueResultMap != nil {
		fxDeviceDailyStatistics = ProcessData(DeviceIdList, fxYesterdayValueResultMap, fxTodayValueResultMap, nil)
	}
	fxDeviceDailyStatisticsMap := make(map[int]PvPowerGenerationModel)
	if fxDeviceDailyStatistics != nil {
		for _, vo := range fxDeviceDailyStatistics {
			fxDeviceDailyStatisticsMap[vo.DeviceId] = vo
		}
	}

	if fxDeviceDailyStatistics == nil && zxDeviceDailyStatisticsMap == nil {
		setting.ZAPS.Debug("statisticsDayChargingAndDischarging fxDeviceDailyStatistics == nil && zxDeviceDailyStatisticsMap == nil")
		return
	}

	polarity := energyProductCodeSettingMap["pv_polarity_code"]
	polarityInt, _ := strconv.Atoi(polarity)
	polarityValueList := GetSettingValueListByDeviceIdListAndYcCode(DeviceIdList, polarityInt)
	polarityMap := make(map[int]string)
	if polarityValueList != nil {
		for _, po := range polarityValueList {
			polarityMap[po.DeviceId] = po.Value
		}
	}

	yesterday := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, 0, 0, 0, 0, time.Local).Unix()
	insertList := make([]*models.EsChargeDischargeModel, 0)
	for _, DeviceId := range DeviceIdList {
		var zxPower *float64
		var fxPower *float64
		if zxPvPowerGenerationModel, ok := zxDeviceDailyStatisticsMap[DeviceId]; ok {
			zxPower = &zxPvPowerGenerationModel.Power
		}
		if fxPvPowerGenerationModel, ok := fxDeviceDailyStatisticsMap[DeviceId]; ok {
			fxPower = &fxPvPowerGenerationModel.Power
		}

		if zxPower == nil && fxPower == nil {
			continue
		}

		esChargeDischargeModel := &models.EsChargeDischargeModel{
			DeviceId: DeviceId,
			Ts:       yesterday,
		}

		var chargeCapacity *float64
		var dischargeCapacity *float64
		if polarityMap == nil || polarityMap[DeviceId] == "0" || polarityMap[DeviceId] != "1" {
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

/**
	方法作废，需要根据网关项目改写逻辑
 */
func getDeviceTypeCodeByCategoryId(categoryId []int, acMeterCategoryId int) []int {
	// TODO: Implement getDeviceTypeCodeByCategoryId
	return nil
}
/**
方法作废，需要根据网关项目改写逻辑
*/
func getDeviceIdListByProjectIdList(projectIdList []int, acMeterProductCodeList []int) []int {
	// TODO: Implement getDeviceIdListByProjectIdList
	return nil
}

func GetTodayValueMap(DeviceIdList []int, code int) map[int]float64 {
	todayDate := time.Now().Unix()

	todayValueList := GetLastYcHistoryByDeviceIdListAndCodeList(DeviceIdList, []int{code}, "", todayDate, todayDate + 86399999)
	if todayValueList == nil || len(todayValueList) == 0 {
		return nil
	}

	result := make(map[int]float64)
	for _, ycModel := range todayValueList {
		result[ycModel.DeviceId] = ycModel.Value
	}

	return result
}

func GetYesterdayValueMap(DeviceIdList []int, code int) map[int]float64 {
	yesterdayDate := time.Now().AddDate(0, 0, -1).Unix()

	yesterdayValueList := GetLastYcHistoryByDeviceIdListAndCodeList(DeviceIdList, []int{code}, "", yesterdayDate, yesterdayDate + 86399999)
	if yesterdayValueList == nil || len(yesterdayValueList) == 0 {
		return nil
	}

	result := make(map[int]float64)
	for _, ycModel := range yesterdayValueList {
		result[ycModel.DeviceId] = ycModel.Value
	}

	return result
}

func ProcessData(DeviceIdList []int, yesterdayValueMap map[int]float64, todayValueMap map[int]float64, devicePriceMap map[int]float64) []PvPowerGenerationModel {
	if len(yesterdayValueMap) == 0 || len(todayValueMap) == 0 {
		return nil
	}
	dictDataRepo := repositories.DictDataRepository{}
	dictDataList, _, _ := dictDataRepo.GetAll("", "device_metering_point_magn_setting_code", 1, 500)

	// 将切片转换为 map
	energyProductCodeSettingMap := make(map[string]string)
	for _, dictData := range dictDataList {
		energyProductCodeSettingMap[dictData.DictLabel] = dictData.DictValue
	}

	magn := energyProductCodeSettingMap["1"];
	magnInt, _ := strconv.Atoi(magn)
	magnValueList := GetSettingValueListByDeviceIdListAndYcCode(DeviceIdList, magnInt)
	magnMap := make(map[int]string)
	for _, item := range magnValueList {
		magnMap[item.DeviceId] = item.Value
	}

	fz := energyProductCodeSettingMap["2"];
	fzInt, _ := strconv.Atoi(fz)
	fzValueList := GetSettingValueListByDeviceIdListAndYcCode(DeviceIdList, fzInt)
	fzMap := make(map[int]string)
	for _, item := range fzValueList {
		fzMap[item.DeviceId] = item.Value
	}

	insertList := []PvPowerGenerationModel{}
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05")
	fmt.Println("yesterday :", yesterday)
	for _, DeviceId := range DeviceIdList {
		fmt.Println("DeviceId :", DeviceId)
		deviceYesterdayValue := yesterdayValueMap[DeviceId]
		fmt.Println("deviceYesterdayValue :", deviceYesterdayValue)
		deviceTodayValue := todayValueMap[DeviceId]
		fmt.Println("deviceTodayValue :", deviceTodayValue)
		if deviceYesterdayValue == 0 || deviceTodayValue == 0 {
			fmt.Println("deviceYesterdayValue == 0 || deviceTodayValue == 0. DeviceId :", DeviceId)
			continue
		}

		yesterdayValueBd := decimal.NewFromFloat(float64(deviceYesterdayValue))
		todayValueBd := decimal.NewFromFloat(float64(deviceTodayValue))

		// 获取倍率
		iMagn := decimal.NewFromInt(int64(1))
		if magnMap[DeviceId] != "" {
			iMagn, _ = decimal.NewFromString((magnMap[DeviceId]))
		}

		// 获取示数翻转值
		iFz := 0
		if fzMap[DeviceId] != "" {
			iFz, _ = strconv.Atoi(fzMap[DeviceId])
		}

		// 计算电量，包含翻转逻辑
		diffBd := GetElectricityCalculation(todayValueBd, yesterdayValueBd, iFz)
		devicePower := diffBd.Mul(iMagn)

		// 计算收益
		deviceProfit := decimal.NewFromInt(int64(0))
		if price, ok := devicePriceMap[DeviceId]; ok {
			deviceProfit = devicePower.Mul(decimal.NewFromFloat(float64(price)))
		}

		pvPowerGenerationModel := PvPowerGenerationModel{
			Power:    devicePower.InexactFloat64(),
			Profit:   deviceProfit.InexactFloat64(),
			DeviceId: DeviceId,
			Ts:       time.Now().AddDate(0, 0, -1).Unix(),
		}
		insertList = append(insertList, pvPowerGenerationModel)
	}
	if len(insertList) == 0 {
		return nil
	}
	return insertList
}

func GetSettingValueListByDeviceIdListAndYcCode(devIDList []int, code int) []models.SettingData {
	result := []models.SettingData{}
	//DeviceIdList := append([]int{}, devIDList...)
	//直接查taos，这里没有Redis
	realtimeDataRepo := repositories.RealtimeDataRepository{}

	stringDevIds := make([]string, len(devIDList))
	for i, v := range devIDList {
		stringDevIds[i] = strconv.Itoa(v)
	}
	subDevIds := strings.Join(stringDevIds, ",")
	//直接查taos
	settingValueList, _ := realtimeDataRepo.GetSettingListByDevIdsAndCodes(subDevIds, strconv.Itoa(code))
	if settingValueList != nil {
		for _, settingValuePo := range settingValueList {
			//DeviceId := settingValuePo.DeviceId
			value := settingValuePo.Value
			if value != "" {
				result = append(result, settingValuePo)
				//DeviceIdList = removeDeviceId(DeviceIdList, DeviceId)
			}
		}
	}
	//这里就不用再查taos了
	//if len(DeviceIdList) > 0 {
	//	mysqlYcValuePosValueList := GetValueListByDeviceIdAndCode(DeviceIdList, code)
	//	if mysqlYcValuePosValueList != nil {
	//		result = append(result, mysqlYcValuePosValueList...)
	//	}
	//}
	return result
}

func removeDeviceId(slice []int, DeviceId int) []int {
	for i, id := range slice {
		if id == DeviceId {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func GetLastYcHistoryByDeviceIdListAndCodeList(devIDList []int, codeList []int, interval string, beginDt int64, endDt int64) []models.YcData {
	if len(devIDList) == 0 || len(codeList) == 0 || beginDt == 0 || endDt == 0 {
		return []models.YcData{}
	}

	size := len(devIDList)
	num := 500
	count := size / num
	if size%num != 0 {
		count++
	}

	idSet := make(map[int]bool)
	for _, id := range devIDList {
		idSet[id] = true
	}

	result := []models.YcData{}
	for i := 0; i < count; i++ {
		start := i * num
		end := start + num
		if i == count-1 {
			end = size
		}
		subDevIDList := devIDList[start:end]
		realtimeDataRepo := repositories.RealtimeDataRepository{}

		stringDevIds := make([]string, len(subDevIDList))
		for i, v := range subDevIDList {
			stringDevIds[i] = strconv.Itoa(v)
		}
		subDevIds := strings.Join(stringDevIds, ",")

		stringCodes := make([]string, len(codeList))
		for i, v := range subDevIDList {
			stringCodes[i] = strconv.Itoa(v)
		}
		codes := strings.Join(stringCodes, ",")

		modelList, _ := realtimeDataRepo.GetLastYcHistoryByDeviceIdListAndCodeList(subDevIds, codes, interval, beginDt, endDt)
		if len(modelList) == 0 {
			continue
		}
		for _, ycModel := range modelList {
			DeviceId := ycModel.DeviceId
			if idSet[DeviceId] {
				result = append(result, ycModel)
			}
		}
	}

	return result
}

func GetElectricityCalculation(currentValue, lastValue decimal.Decimal, flipIndicator int) decimal.Decimal {
	subtract := currentValue.Sub(lastValue)
	if flipIndicator <= 0{
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


