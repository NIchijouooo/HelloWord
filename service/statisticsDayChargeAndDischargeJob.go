package service

import (
	"fmt"
	"gateway/models"
	"gateway/repositories"
	"github.com/jasonlvhit/gocron"

	"gorm.io/gorm"
)

// 创建一个 job 结构体，包含需要执行的方法
type StatisticsDayChargeAndDischargeJob struct{
	db     *gorm.DB
	repoDictData *repositories.DictDataRepository
}


func NewStatisticsDayChargeAndDischargeJob() *StatisticsDayChargeAndDischargeJob {
	return &StatisticsDayChargeAndDischargeJob{
		db: models.DB,
		repoDictData: repositories.NewDictDataRepository()}
}


func (job *StatisticsDayChargeAndDischargeJob) Run() {
	// 需要执行的任务
	fmt.Println("Running StatisticsDayChargeAndDischargeJob job...")
}
/**
定时任务调用接口，
是查什么类型的设备的什么点位，既在字典定义好,taos和iot已经有完整接口示例，搬过来，每次都是传年月日的ts更新，避免数据太多。
code：pcs_charge_discharge的设备。
code：5为充电，6为放电的点位。
创建结构体pcsChargeDischarge和taos表，
用定时任务每小时一次，去taos统计这些deviceId+code点位，昨天的充，放电量，存入taos的pcs_charge_discharge。
*/
func (r *StatisticsDayChargeAndDischargeJob) GetPcsChargeDischargeForType(deviceIds, interval string, startTime, endTime int64){
	dictDataList, _, _ := r.repoDictData.GetAll("", "energy_product_code_setting", 1, 50)

	// 将切片转换为 map
	dictDataMap := make(map[string]string)
	for _, dictData := range dictDataList {
		dictDataMap[dictData.DictLabel] = dictData.DictValue
	}
	//能源交流电表产品CODE数组
	acMeterCategoryId := dictDataMap["energy_storage_ac_meter_category"]
	if acMeterCategoryId != ""{


	}
	//能源辅助电表产品CODE数组
	auxiliaryMeterCategoryId := dictDataMap["energy_storage_auxiliary_meter_category"]
	if auxiliaryMeterCategoryId != ""{


	}



}

func (s *StatisticsDayChargeAndDischargeJob) Start() {
	// 创建一个新的定时任务
	job := &StatisticsDayChargeAndDischargeJob{}
	cron := gocron.NewScheduler()
	// 定义任务执行的时间间隔，例如每分钟执行一次
	cron.Every(1).Hour().Do(job)
	// 启动定时任务
	cron.Start()
}
