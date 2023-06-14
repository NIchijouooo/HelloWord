package repositories

import (
	"database/sql"
	"fmt"
	"gateway/models"
	"log"

	"gorm.io/gorm"
)

// 定义字典类型管理的存储库
type DevicePointRepository struct {
	sqldb  *gorm.DB
	taosDb *sql.DB
}

func NewDevicePointRepository() *DevicePointRepository {
	return &DevicePointRepository{sqldb: models.DB, taosDb: models.TaosDB}
}

/*
*

	根据设备id，点位类型获取命令参数属性
	设备去查模型，模型查命令，有多个命令，命令查参数既点位，一个命令数据对应多条参数表数据
	pointType:yx/yc/setting/all
*/
func (r *DevicePointRepository) GetPointsByDeviceId(pointType string, deviceId, code int) []*models.EmDeviceModelCmdParam {
	//var pointParams []*models.PointParam
	var pointParams []*models.EmDeviceModelCmdParam
	//.Joins("left join em_device_model on em_device_model.id = em_device.model_id")
	if len(pointType) == 0 {
		r.sqldb.Joins("LEFT JOIN em_device_model_cmd ON em_device_model_cmd.id = em_device_model_cmd_param.device_model_cmd_id").Joins("LEFT JOIN em_device ON em_device.id = em_device_model_cmd.device_model_id").Where("em_device.id = ?", deviceId).Find(&pointParams).Statement.SQL.String()
	} else {
		r.sqldb.Joins("LEFT JOIN em_device_model_cmd ON em_device_model_cmd.id = em_device_model_cmd_param.device_model_cmd_id").Joins("LEFT JOIN em_device ON em_device.id = em_device_model_cmd.device_model_id").Where("em_device.id = ?", deviceId).Where("em_device_model_cmd_param.iot_data_type = ?", pointType).Find(&pointParams).Statement.SQL.String()
	}
	return pointParams
}

/*
*

	根据设备lable查询设备信息
*/
func (r *DevicePointRepository) GetDeviceByDevLabel(label string) []*models.EmDevice {
	var emDevice []*models.EmDevice
	r.sqldb.Where("em_device.label = ?", label).Find(&emDevice)
	return emDevice
}

/*
*

	根据设备id查询设备信息
*/
func (r *DevicePointRepository) GetDeviceByDeviceId(deviceId int) *models.EmDevice {
	var emDevice *models.EmDevice
	r.sqldb.Where("em_device.id = ?", deviceId).Find(&emDevice)
	return emDevice
}

/*
*
获取yc
*/
func (r *DevicePointRepository) GetYcById(deviceId, code int) (models.YcData, error) {
	var realtime models.YcData
	//tableName := fmt.Sprintf("%v%d%v%d", "realtimedatatest.yc_", deviceId, "_", code)
	sql := fmt.Sprint("select Last(ts), val, device_id, code from realtimedatatest.yc where device_id =", deviceId, "and code =", code)
	rows, err := r.taosDb.Query(sql)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&realtime.Ts, &realtime.Value, &realtime.DeviceId, &realtime.Code)
		if err != nil {
			log.Printf("Request params:%v", err)
		}
	}
	return realtime, err
}

/*
*
获取setting
*/
func (r *DevicePointRepository) GetSettingById(deviceId, code int) (models.SettingData, error) {
	var realtime models.SettingData
	//tableName := fmt.Sprintf("%v%d%v%d", "realtimedatatest.setting_", deviceId, "_", code)
	sql := fmt.Sprint("select Last(ts), val, device_id, code from realtimedatatest.setting where device_id =", deviceId, "and code =", code)
	rows, err := r.taosDb.Query(sql)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&realtime.Ts, &realtime.Value, &realtime.DeviceId, &realtime.Code)
		if err != nil {
			log.Printf("Request params:%v", err)
		}
	}
	return realtime, err
}
