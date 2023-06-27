package repositories

import (
	"gateway/models"

	"gorm.io/gorm"
)

type ConfigurationCenterRepository struct {
	db *gorm.DB
}

func NewConfigurationCenterRepository() *ConfigurationCenterRepository {
	return &ConfigurationCenterRepository{db: models.DB}
}

// 新增电费配置
func (r *ConfigurationCenterRepository) AddConfiguration(configuration *models.EmConfiguration) error {
	return r.db.Create(configuration).Error
}

// 更新电费配置
func (r *ConfigurationCenterRepository) UpdateConfiguration(configuration *models.EmConfiguration) error {
	return r.db.Save(configuration).Error
}

// 获取电费配置列表
func (r *ConfigurationCenterRepository) GetConfigurationList(province string, month string, pageNum int, pageSize int) ([]models.EmConfiguration, int64, error) {
	var (
		configurationList []models.EmConfiguration
		total             int64
	)

	query := r.db.Model(&models.EmConfiguration{})
	if province != "" {
		query = query.Where("province = ?", province)
	}
	if month != "" {
		query = query.Where("month = ?", month)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&configurationList).Error; err != nil {
		return nil, 0, err
	}

	return configurationList, total, nil
}

// 删除电费配置
func (r *ConfigurationCenterRepository) DeleteConfiguration(id int) error {
	return r.db.Delete(&models.EmConfiguration{}, id).Error
}

// GetConfigurationByProvince 通过省份和月份获取电费配置
func (r *ConfigurationCenterRepository) GetConfigurationByProvince(province string, month string) (models.EmConfiguration, error) {

	var configuration models.EmConfiguration

	query := r.db.Model(&models.EmConfiguration{})
	if province != "" {
		query = query.Where("province = ?", province)
	}
	if month != "" {
		query = query.Where("month = ?", month)
	}

	if err := query.Find(&configuration).Error; err != nil {
		return configuration, err
	}

	return configuration, nil
}
