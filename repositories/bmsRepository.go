package repositories

import (
	"gateway/models"
	"gateway/models/ReturnModel"

	"gorm.io/gorm"
)

type BmsRepository struct {
	db *gorm.DB
}

func NewBmsRepository() *BmsRepository {
	return &BmsRepository{db: models.DB}
}

// 获取设备类型下的所有设备数据
func (r *BmsRepository) GetBmsDeviceList(deviceType string) ([]*ReturnModel.TreeDevice, error) {
	var deviceList []*ReturnModel.TreeDevice
	err := r.db.Where("device_type=?", deviceType).Find(&deviceList).Error
	return deviceList, err
}
