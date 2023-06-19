package repositories

import (
	"gateway/models"

	"gorm.io/gorm"
)

type CentralizedRepository struct {
	db *gorm.DB
}

func NewCentralizedRepository() *CentralizedRepository {
	return &CentralizedRepository{db: models.DB}
}

// 获取策略数据
func (r *CentralizedRepository) GetPolicyList() ([]models.EmStrategy, error) {
	var policyList []models.EmStrategy
	err := r.db.Find(&policyList).Error
	return policyList, err
}

// 新增策略
func (r *CentralizedRepository) CreatePolicy(policyData *models.EmStrategy) error {
	return r.db.Create(policyData).Error
}

// 修改策略
func (r *CentralizedRepository) UpdatePolicy(policyData *models.EmStrategy) error {
	var emPolicyData models.EmStrategy
	id := policyData.Id
	r.db.Where("id = ?", id).First(&emPolicyData)
	policyData.Id = emPolicyData.Id
	return r.db.Save(policyData).Error
}

// 删除策略
func (r *CentralizedRepository) DeletePolicy(Id int) error {
	return r.db.Delete(&models.EmStrategy{}, Id).Error
}
