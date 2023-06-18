package repositories

import (
	"gateway/models"
	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository() *DeviceRepository {
	return &DeviceRepository{db: models.DB}
}

/*
*
根据设备类型获取设备列表
*/
func (r DeviceRepository) GetDeviceListByType(param models.DeviceParam) ([]models.EmDevice, error) {
	var result []models.EmDevice
	r.db.Model(&models.EmDevice{})
	err := r.db.Where("device_type = ?", param.DeviceType).Find(&result).Error
	return result, err
}
