package repositories

import (
	"database/sql"
	"gateway/models"
	"gorm.io/gorm"
)

type DeviceRepository struct {
	db     *gorm.DB
	taosDb *sql.DB
}

type DevicePoint struct {
	DeviceId int    `json:"deviceId"`
	Code     string `json:"code"`
}

func NewDeviceRepository() *DeviceRepository {
	return &DeviceRepository{db: models.DB, taosDb: models.TaosDB}
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

func (r *DeviceRepository) GetDevicePoint(deviceType string, pointDictType string) ([]DevicePoint, error) {
	var param models.DeviceParam
	var dictData models.DictData
	var res []DevicePoint
	// 找设备
	param.DeviceType = deviceType
	deviceList, err := r.GetDeviceListByType(param)
	err = r.db.Where("dict_type = ?", pointDictType).First(&dictData).Error
	// 组装数据
	for _, device := range deviceList {
		var v DevicePoint
		v.DeviceId = device.Id
		v.Code = dictData.DictValue
		res = append(res, v)
	}
	return res, err
}
