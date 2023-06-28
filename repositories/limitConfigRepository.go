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

//根据设备类型和codes获取越限配置信息
func (r *LimitConfigRepository) GetLimitConfigListByDeviceTypeAndCodes(deviceType string, codes []string) ([]models.LimitConfigVo, error) {
	var limitConfigList []models.LimitConfigVo
	if err := r.db.Where("del_flag = ?", 0).Where("device_type = ?", deviceType).Where("property_code in ?", codes).Find(&limitConfigList).Error; err != nil {
		return nil, err
	}
	return limitConfigList, nil
}

//
func (r *LimitConfigRepository) GetLimitConfigListCheckById(id int, propertyCode string) ([]models.LimitConfigVo, error) {
	var (
		limitConfigVoList []models.LimitConfigVo
	)
	var err error
	query := r.db.Model(&models.LimitConfigVo{})
	if propertyCode == "" || id <= 0{
		return nil, err
	}
	query = query.Where("id <> ?", id).Where("property_code = ?", propertyCode)
	if err := query.Find(&limitConfigVoList).Error; err != nil {
		return nil, err
	}

	return limitConfigVoList, nil
}

//
func (r *LimitConfigRepository) GetLimitConfigListList(deviceType, propertyCode string) ([]models.LimitConfigVo, error) {
	var (
		configurationList []models.LimitConfigVo
	)
	query := r.db.Model(&models.LimitConfigVo{})
	if deviceType != "" {
		query = query.Where("device_type = ?", deviceType)
	}
	if propertyCode != "" {
		query = query.Where("property_code = ?", propertyCode)
	}
	if err := query.Find(&configurationList).Error; err != nil {
		return nil, err
	}

	return configurationList, nil
}
