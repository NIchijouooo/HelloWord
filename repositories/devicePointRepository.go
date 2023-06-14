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
	sqldb     *gorm.DB
	taosDb *sql.DB
}

func NewDevicePointRepository() *DevicePointRepository {
	return &DevicePointRepository{sqldb: models.DB, taosDb: models.TaosDB}
}

/**
	设备去查模型，模型查命令，有多个命令，命令查参数既点位，一个命令数据对应一条参数表数据？还是多条
 */
func (r *DevicePointRepository) GetYxById(deviceId, code int) (models.YxData, error) {
	var realtime models.YxData
	//tableName := fmt.Sprintf("%v%d%v%d", "realtimedatatest.yx_", deviceId, "_", code)
	//sql := fmt.Sprintf("select last(*) from ? ", tableName)
	sql := fmt.Sprint("select Last(ts), val, device_id, code from realtimedatatest.yx where device_id =", deviceId, "and code =", code)
	fmt.Println(sql)
	rows, err := r.sqldb.First(sql)
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
