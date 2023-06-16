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
func (r *RuleHistoryRepository) GetLastRuleHistory(eventId int, deviceId int) (*models.EmRuleHistoryModel, error) {
	var emRuleHistoryModel *models.EmRuleHistoryModel
	if err := r.db.Joins("join rule_history_device on event_history_id = rule_history.id").Where("tag = ?", 0).Where("del_flag = ?", 0).Where("event_id", eventId).Where("device_id", deviceId).Order("rule_history.id desc").Limit(1).Find(&emRuleHistoryModel).Error; err != nil {
		return nil, err
	}
	return emRuleHistoryModel, nil
}

// updateRuleHistoryTag 修改告警恢复
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

// InsertRuleHistory 新增关联设备
func (r *RuleHistoryRepository) InsertRuleHistoryDevice(emRuleHistoryDeviceModelList *[]models.EmRuleHistoryDeviceModel) (int, error) {
	if err := r.db.Create(emRuleHistoryDeviceModelList).Error; err != nil {
		return 0, err
	}
	return 1, nil
}
