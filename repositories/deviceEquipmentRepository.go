package repositories

import (
	"gateway/models"
	"gorm.io/gorm"
)

/*
*
设备台账库
*/
type DeviceEquipmentRepository struct {
	db *gorm.DB
}

func NewDeviceEquipmentRepository() *DeviceEquipmentRepository {
	return &DeviceEquipmentRepository{db: models.DB}
}

/*
*
根据设备id获取设备台账信息
*/
func (r DeviceEquipmentRepository) GetEquipmentInfoByDevId(deviceId int) (models.DeviceEquipmentAccountInfo, error) {
	var result models.DeviceEquipmentAccountInfo
	err := r.db.Where("device_id = ?", deviceId).Find(&result).Error
	return result, err
}
