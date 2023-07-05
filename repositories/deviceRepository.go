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

type CtrlHistory struct {
	DeviceName string `json:"deviceName"`
	ParamName  string `json:"paramName"`
	models.EmCtrlHistory
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
	var err error
	if param.DeviceType != "" {
		err = r.db.Where("device_type = ?", param.DeviceType).Find(&result).Error
	} else {
		err = r.db.Find(&result).Error
	}
	return result, err
}

/*
根据设备类型获取设备Id列表
*/
func (r DeviceRepository) GetDeviceIdListByDeviceType(deviceType string) ([]int, error) {
	var result []int
	err := r.db.Table("em_device").Select("id").Where("device_type = ?", deviceType).Find(&result).Error
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
		Joins("left join em_device_model_cmd cmd ON model_id = cmd.device_model_id").
		Joins("left join em_coll_interface eci ON ed.coll_interface_id = eci.id").
		Joins("left join em_device_model_cmd_param cmdp ON cmd.id = cmdp.device_model_cmd_id").
		Where("cmdp.id =?", paramId).
		Scan(&res)
	return res
}

// GetCtrlHistoryList 获取控制历史列表
func (r *DeviceRepository) GetCtrlHistoryList(page, pageSize int) ([]CtrlHistory, int64, error) {
	var (
		res   []CtrlHistory
		total int64
	)
	rows := r.db.Table("em_ctrl_history ech").
		Select("ech.id AS id, ech.value AS value, ech.ctrl_user_name AS ctrl_user_name," +
			" ech.ctrl_status AS ctrl_status, ech.create_time AS create_time," +
			" ech.update_time AS update_time, ed.name AS device_name, cmdp.name AS param_name").
		Joins("left join em_device ed ON ech.device_id = ed.id").
		Joins("left join em_device_model_cmd_param cmdp ON ech.param_id = cmdp.id").
		Order("ech.id desc").
		Scan(&res)

	if err := rows.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := rows.Offset((page - 1) * pageSize).Limit(pageSize).Find(&res).Error; err != nil {
		return nil, 0, err
	}
	return res, total, nil
}

func (r *DeviceRepository) AddCtrlHistory(emCtrlHistory *models.EmCtrlHistory) error {
	return r.db.Create(emCtrlHistory).Error
}
