package repositories

import (
	"fmt"
	"gateway/models"
	"github.com/taosdata/driver-go/v3/errors"
	"gorm.io/gorm"
	"strconv"
)

type RuleHistoryRepository struct {
	db *gorm.DB
}

func NewRuleHistoryRepository() *RuleHistoryRepository {
	return &RuleHistoryRepository{
		db: models.DB,
	}
}

// GetLastRuleHistory 获取最新一条告警
func (r *RuleHistoryRepository) GetLastRuleHistory(ruleId int, deviceId int) (*models.EmRuleHistoryModel, error) {
	var emRuleHistoryModel *models.EmRuleHistoryModel
	if err := r.db.Joins("join rule_history_device on rule_history_id = rule_history.id").Where("tag = ?", 0).Where("del_flag = ?", 0).Where("rule_id", ruleId).Where("device_id", deviceId).Order("rule_history.id desc").Limit(1).Find(&emRuleHistoryModel).Error; err != nil {
		return nil, err
	}
	return emRuleHistoryModel, nil
}

// UpdateRuleHistoryTag 修改告警恢复
func (r *RuleHistoryRepository) UpdateRuleHistoryTag(emRuleHistoryModel *models.EmRuleHistoryModel) (int, error) {
	if err := r.db.Save(emRuleHistoryModel).Error; err != nil {
		return 0, err
	}
	return 1, nil
}

// InsertRuleHistory 新增告警
func (r *RuleHistoryRepository) InsertRuleHistory(emRuleHistoryModel *models.EmRuleHistoryModel) (int, error) {
	if err := r.db.Create(emRuleHistoryModel).Error; err != nil {
		return 0, err
	}
	return 1, nil
}

// InsertRuleHistoryDevice 新增关联设备
func (r *RuleHistoryRepository) InsertRuleHistoryDevice(emRuleHistoryDeviceModelList *[]models.EmRuleHistoryDeviceModel) (int, error) {
	if err := r.db.Create(emRuleHistoryDeviceModelList).Error; err != nil {
		return 0, err
	}
	return 1, nil
}

/*
*
获取告警历史集合
*/
func (r *RuleHistoryRepository) GetRuleHistoryList(param models.RuleHistoryParam) ([]models.EmRuleHistoryVo, int64, error) {
	pageNum := param.PageNum
	pageSize := param.PageSize
	var (
		historyList []models.EmRuleHistoryVo
		total       int64
	)
	query := r.db.Table("rule_history ruleHis").
		Joins("join rule_history_device ruleHisDev").
		Joins("join em_device dev").
		Where("ruleHis.id = ruleHisDev.rule_history_id").
		Where("dev.id = ruleHisDev.device_id")
	// 设备id
	deviceIds := param.DeviceIds
	if len(deviceIds) > 0 {
		query.Where("ruleHisDev.device_id in ?", deviceIds)
	}
	// 设备类型集合
	deviceTypeList := param.DeviceTypeList
	if len(deviceTypeList) > 0 {
		query.Where("dev.device_type in ?", deviceTypeList)
	}
	// 按点位查询
	codes := param.Codes
	if len(codes) > 0 {
		query.Where("ruleHisDev.property_code in ?", codes)
	}
	// 按产生时间查询
	startTime := param.StartTime
	if len(startTime) > 0 {
		query.Where("ruleHis.produce_time >= ?", startTime)
	}
	endTime := param.EndTime
	if len(endTime) > 0 {
		query.Where("ruleHis.produce_time <= ?", endTime)
	}
	// 事件等级(0-通知；1-次要；2-告警；3-故障)
	level := param.Level
	if len(level) > 0 {
		query.Where("ruleHis.level = ?", level)
	}
	// 恢复标记：0-未确认，1-自动恢复 2-手动恢复
	tag := param.Tag
	if len(tag) > 0 {
		query.Where("ruleHis.tag = ?", tag)
	}
	// 告警详情
	description := param.Description
	if len(description) > 0 {
		description = fmt.Sprint("%", description, "%")
		query.Where("ruleHis.description like ?", description)
	}
	deviceName := param.DeviceName
	if len(deviceName) > 0 {
		deviceName = fmt.Sprint("%", deviceName, "%")
		query.Where("dev.name like ?", deviceName)
	}
	// 分页查询
	if pageNum > 0 && pageSize > 0 {
		countErr := query.Count(&total).Error
		if countErr != nil || total == 0 {
			return []models.EmRuleHistoryVo{}, 0, countErr
		}
		// 计算页数
		pages := total / pageSize
		if total%pageSize > 0 {
			pages += 1
		}
		// 调整当前页
		if pageNum > pages {
			pageNum = pages
		}
		// 计算当前页的索引
		offsetIndex := (pageNum - 1) * pageSize
		query.Offset(int(offsetIndex)).Limit(int(offsetIndex + pageSize))
	}
	err := query.Select("ruleHis.*,ruleHisDev.device_id,ruleHisDev.property_code,dev.name as device_name").Order("ruleHis.produce_time desc").Find(&historyList).Error
	if err != nil {
		return []models.EmRuleHistoryVo{}, 0, err
	}
	if pageNum == 0 || pageSize == 0 {
		total = int64(len(historyList))
	}
	return historyList, total, nil
}

func (r *RuleHistoryRepository) GetRuleHistoryStatistic(param models.RuleHistoryParam) (models.EventStatisticVo, error) {
	// 查告警列表
	historyList, _, err := r.GetRuleHistoryList(param)
	if err != nil {
		return models.EventStatisticVo{}, err
	}
	var result models.EventStatisticVo
	// 统计类型;1-统计告警等级,2-统计告警状态
	statisticType := param.StatisticType
	if statisticType == 0 {
		return models.EventStatisticVo{}, errors.NewError(1, "统计类型错误")
	}
	// 告警等级map,key=等级,value=等级对应的历史告警集合
	statisticMap := make(map[int][]models.EmRuleHistoryVo)
	if len(historyList) > 0 {
		for _, event := range historyList {
			level := event.Level
			tag := event.Tag
			var list []models.EmRuleHistoryVo
			var ok bool
			if statisticType == 1 {
				list, ok = statisticMap[level]
			} else {
				list, ok = statisticMap[tag]
			}
			if !ok {
				list = []models.EmRuleHistoryVo{}
			}
			list = append(list, event)
			if statisticType == 1 {
				statisticMap[level] = list
			} else {
				statisticMap[tag] = list
			}
		}
	}
	dictType := "event_level_list"
	if statisticType == 2 {
		dictType = "event_tag_list"
	}
	// 查告警等级字典
	dictList, err := NewDictDataRepository().GetDictDataByDictType(dictType)
	if err != nil {
		return models.EventStatisticVo{}, err
	}
	var eventLevelStatisticList []models.EventStatisticData
	total := 0
	if len(dictList) > 0 {
		// 封装数据
		for _, data := range dictList {
			dictValInt, err := strconv.Atoi(data.DictValue)
			if err != nil {
				continue
			}
			// 按等级获取历史告警集合
			list := statisticMap[dictValInt]
			size := len(list)
			total += size
			statistic := models.EventStatisticData{
				Total: size,
				Code:  data.DictValue,
				Name:  data.DictLabel,
			}
			eventLevelStatisticList = append(eventLevelStatisticList, statistic)
		}
	}
	result.EventStatisticList = eventLevelStatisticList
	result.Total = total
	return result, err
}
