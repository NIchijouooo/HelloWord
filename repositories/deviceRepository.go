package repositories

import (
	"fmt"
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

func (r *DeviceRepository) GetEmDeviceInfo(deviceType string) (interface{}, error) {
	type Param struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	type Res struct {
		Id     int     `json:"id"`
		Name   string  `json:"name"`
		Label  string  `json:"label"`
		Params []Param `json:"params"`
	}

	var res []Res
	var dictData []models.DictData
	var deviceParam models.DeviceParam
	deviceParam.DeviceType = deviceType
	// 获取设备列表
	deviceList, err := r.GetDeviceListByType(deviceParam)
	// 获取字典
	switch deviceType {
	case "pcs":
		fmt.Print(deviceList)
		err = r.db.Where("dict_type = ?", deviceType).Find(&dictData).Error
	case "bms":

	case "辅控":
	}
	// 构建数据
	for _, device := range deviceList {
		var v Res
		v.Id = device.Id
		v.Name = device.Name
		v.Label = device.Label
		// 查数据
		for _, dict := range dictData {
			var p Param
			p.Name = dict.DictValue
			v.Params = append(v.Params, p)
		}
		res = append(res, v)
	}
	return res, err

}
