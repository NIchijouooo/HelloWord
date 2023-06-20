package repositories

import (
	"gateway/models"
	"gorm.io/gorm"
)

type LimitConfigRepository struct {
	db *gorm.DB
}

func NewLimitConfigRepository() *LimitConfigRepository {
	return &LimitConfigRepository{
		db: models.DB,
	}
}

// GetLimitConfigListByDeviceType GetRuleByDeviceLabel 通过设备类型获取配置
func (r *LimitConfigRepository) GetLimitConfigListByDeviceType(deviceType string) ([]models.LimitConfigVo, error) {
	var limitConfigList []models.LimitConfigVo
	if err := r.db.Where("del_flag = ?", 0).Where("device_type = ?", deviceType).Find(&limitConfigList).Error; err != nil {
		return nil, err
	}
	return limitConfigList, nil
}

// InsertLimitConfig 新增越限配置
func (r *LimitConfigRepository) InsertLimitConfig(limitConfig *models.LimitConfigVo) (int, error) {
	if err := r.db.Create(limitConfig).Error; err != nil {
		return 0, err
	}
	return 1, nil
}

// UpdateLimitConfig 修改越限配置
func (r *LimitConfigRepository) UpdateLimitConfig(limitConfig *models.LimitConfigVo) (int, error) {
	if err := r.db.Save(limitConfig).Error; err != nil {
		return 0, err
	}
	return 1, nil
}
