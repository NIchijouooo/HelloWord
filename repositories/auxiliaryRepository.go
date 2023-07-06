package repositories

import (
	"gateway/models"

	"gorm.io/gorm"
)

type AuxiliaryRepository struct {
	db *gorm.DB
}

func NewAuxiliaryRepository() *AuxiliaryRepository {
	return &AuxiliaryRepository{db: models.DB}
}

// 获取设备类型下的所有设备数据
func (r *AuxiliaryRepository) GetAuxiliaryDevice(deviceType string) ([]models.EmDevice, error) {
	var deviceList []models.EmDevice
	if deviceType == "" {
		err := r.db.Find(&deviceList).Error
		return deviceList, err
	}
	err := r.db.Where("device_type=?", deviceType).Find(&deviceList).Error
	return deviceList, err
}

// 获取设备的所有标签
func (r *AuxiliaryRepository) GetAuxiliaryDeviceType() ([]string, error) {
	var deviceTypeList []string
	err := r.db.Model(&models.EmDevice{}).Select("label").Group("label").Find(&deviceTypeList).Error
	return deviceTypeList, err
}
