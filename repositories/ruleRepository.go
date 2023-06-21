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

// DeleteRuleList 批量删除规则
func (r *RuleRepository) DeleteRuleList(ruleIdList []int) (int, error) {
	if err := r.db.Table("rule").Where("id in ?", ruleIdList).Where("del_flag = ?", 0).Update("del_flag", "1").Error; err != nil {
		return 0, err
	}
	return len(ruleIdList), nil
}

// InsertRuleList 批量新增规则
func (r *RuleRepository) InsertRuleList(ruleList []*models.EmRuleModel) (int, error) {
	if err := r.db.Create(ruleList).Error; err != nil {
		return 0, err
	}
	return len(ruleList), nil
}

// UpdateRuleList 批量修改规则公式
func (r *RuleRepository) UpdateRuleList(ruleList []*models.EmRuleModel) (int, error) {
	for _, rule := range ruleList {
		r.db.Table("rule").Where("id = ?", rule.Id).Update("content", rule.Content).Update("update_time", rule.UpdateTime)
	}
	return len(ruleList), nil
}
