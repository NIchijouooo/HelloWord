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

/*
*
根据设备类型获取设备列表
*/
func (r DeviceRepository) GetDeviceListByCollAndName(collInterfaceId, name string) (models.EmDevice, error) {
	var result models.EmDevice
	r.db.Model(&models.EmDevice{})
	err := r.db.Where("coll_interface_id = ? and name = ?", collInterfaceId, name).Find(&result).Error
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

func (r *DeviceRepository) GetDeviceModelCmdParam(paramId int) models.GetDeviceModelCmdParam {
	var res models.GetDeviceModelCmdParam
	r.db.Table("em_device ed").
		Select("ed.name AS device_name,cmd.name AS cmd_name,eci.name AS coll_name,cmdp.name AS param_name").
		Joins("left join em_device_model_cmd cmd").
		Joins("left join em_coll_interface eci").
		Joins("left join em_device_model_cmd_param cmdp").
		Where("model_id = cmd.device_model_id").
		Where("ed.coll_interface_id = eci.id").
		Where("cmd.id = cmdp.device_model_cmd_id").
		Where("cmdp.id =?", paramId).
		Scan(&res)
	return res
}
