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

func (r *DeviceRepository) GetEmDevice() {
	type Device struct {
		id    int
		name  string
		label string
	}
	rows, err := r.db.Table("em_device").Select("*").Joins("left join em_device_model_cmd on em_device.model_id = em_device_model_cmd.device_model_id").Rows()
	if err != nil {
	}

	var device Device
	for rows.Next() {
		err := rows.Scan(&device.id, &device.name, &device.label)
		if err != nil {
			return
		}
	}
	//var emDeviceModelCmd []models.EmDeviceModelCmd
	//if err := r.db.Where("device_model_id = ?", modelId).Find(&emDeviceModelCmd).Error; err != nil {
	//	return nil, err
	//}
	//return emDeviceModelCmd, nil
}
