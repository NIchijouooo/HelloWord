package repositories

import (
	"gateway/models"
	"gorm.io/gorm"
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
func (r *RuleHistoryRepository) GetRuleHistoryList(param models.RuleHistoryParam) ([]models.EmRuleHistoryModel, int64, error) {
	pageNum := param.PageNum
	pageSize := param.PageSize
	var (
		historyList []models.EmRuleHistoryModel
		total       int64
	)
	query := r.db.Table("rule_history ruleHis").
		Joins("join rule_history_device ruleHisDev").
		Where("ruleHis.id = ruleHisDev.rule_history_id").
		Where("ruleHisDev.device_id in ?", param.DeviceIds)
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
	// 分页查询
	if pageNum > 0 && pageSize > 0 {
		countErr := query.Count(&total).Error
		if countErr != nil || total == 0 {
			return []models.EmRuleHistoryModel{}, 0, countErr
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
	err := query.Select("ruleHis.*,ruleHisDev.device_id,ruleHisDev.property_code").Order("ruleHis.produce_time desc").Find(&historyList).Error
	if err != nil {
		return []models.EmRuleHistoryModel{}, 0, err
	}
	if pageNum == 0 || pageSize == 0 {
		total = int64(len(historyList))
	}
	return historyList, total, nil
}
