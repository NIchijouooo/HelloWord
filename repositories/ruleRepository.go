package repositories

import (
	"gateway/models"
	"gorm.io/gorm"
)

type RuleRepository struct {
	db *gorm.DB
}

func NewRuleRepository() *RuleRepository {
	return &RuleRepository{
		db: models.DB,
	}
}

// GetAllRule 获取所有规则
func (r *RuleRepository) GetAllRule() ([]models.EmRuleModel, error) {
	var emRuleModel []models.EmRuleModel
	if err := r.db.Where("del_flag = ?", 0).Find(&emRuleModel).Error; err != nil {
		return nil, err
	}
	return emRuleModel, nil
}

// GetRuleByDeviceLabel 通过设备标签获取规则
func (r *RuleRepository) GetRuleByDeviceLabel(deviceLabel string) ([]models.EmRuleModel, error) {
	var emRuleModel []models.EmRuleModel
	tx := r.db.Where("del_flag = ?", 0)
	if len(deviceLabel) > 0 {
		tx = tx.Where("content like ?", "%product.${"+deviceLabel+":\"%")
	}
	if err := tx.Find(&emRuleModel).Error; err != nil {
		return nil, err
	}
	return emRuleModel, nil
}
