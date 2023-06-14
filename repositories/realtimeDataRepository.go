package repositories

import (
	"database/sql"
	"fmt"
	"gateway/models"
	"log"

	"gorm.io/gorm"
)

// 定义字典类型管理的存储库
type RealtimeDataRepository struct {
	db     *gorm.DB
	taosDb *sql.DB
}

func NewRealtimeDataRepository() *RealtimeDataRepository {
	return &RealtimeDataRepository{db: models.DB, taosDb: models.TaosDB}
}

/*
*
添加db
*/
func (r *RealtimeDataRepository) CreateDB() error {
	sql := fmt.Sprintf("create database if not exists realtimedatatest;")
	_, err := r.taosDb.Exec(sql)
	return err
}

/*
*
添加yx表
*/
func (r *RealtimeDataRepository) CreateYxTable() error {
	// 定义查询参数
	sql := fmt.Sprintf("create table if not exists realtimedatatest.yx (ts timestamp, val int) tags (device_id int, code int);")
	_, err := r.taosDb.Exec(sql)
	return err
}

/*
*
添加yc表
*/
func (r *RealtimeDataRepository) CreateYcTable() error {
	// 定义查询参数
	sql := fmt.Sprintf("create table if not exists realtimedatatest.yc (ts timestamp, val double) tags (device_id int, code int);")
	_, err := r.taosDb.Exec(sql)
	return err
}

/*
*
添加setting表
*/
func (r *RealtimeDataRepository) CreateSettingTable() error {
	// 定义查询参数
	sql := fmt.Sprintf("create table if not exists realtimedatatest.setting (ts timestamp, val NCHAR(16)) tags (device_id int, code int);")
	_, err := r.taosDb.Exec(sql)
	return err
}

/*
*
添加yx
*/
func (r *RealtimeDataRepository) CreateYx(realtime *models.YxData) error {
	// 定义查询参数
	tableName := fmt.Sprintf("%v%d%v%d", "realtimedatatest.yx_", realtime.DeviceId, "_", realtime.Code)
	sql := fmt.Sprintf("INSERT INTO %s (ts, val) using realtimedatatest.yx tags(%d, %d) VALUES (%v, %d)", tableName, realtime.DeviceId, realtime.Code, realtime.Ts, realtime.Value)
	_, err := r.taosDb.Exec(sql)
	return err
}

/*
*
添加yc
*/
func (r *RealtimeDataRepository) CreateYc(realtime *models.YcData) error {
	// 定义查询参数
	tableName := fmt.Sprintf("%v%d%v%d", "realtimedatatest.yc_", realtime.DeviceId, "_", realtime.Code)
	sql := fmt.Sprintf("INSERT INTO %s (ts, val) using realtimedatatest.yc tags(%d, %d) VALUES (%v, %d)", tableName, realtime.DeviceId, realtime.Code, realtime.Ts, realtime.Value)
	_, err := r.taosDb.Exec(sql)
	return err
}

/*
*
添加setting
*/
func (r *RealtimeDataRepository) CreateSetting(realtime *models.SettingData) error {
	// 定义查询参数
	tableName := fmt.Sprintf("%v%d%v%d", "realtimedatatest.setting_", realtime.DeviceId, "_", realtime.Code)
	sql := fmt.Sprintf("INSERT INTO %s (ts, val) using realtimedatatest.setting tags(%d, %d) VALUES (%v, %d)", tableName, realtime.DeviceId, realtime.Code, realtime.Ts, realtime.Value)
	_, err := r.taosDb.Exec(sql)
	return err
}

/*
*
批量添加yx
*/
func (r *RealtimeDataRepository) BatchCreateYx(realtimeList []*models.YxData) error {
	// 定义查询参数
	//tableName := fmt.Sprintf("%d%d", realtime.DeviceId, realtime.Code)
	//sql := fmt.Sprintf("INSERT INTO %s (ts, val) using realtimedatatest.yx tags(?, ?) VALUES (?, ?)", tableName)
	//result := r.taosDb.Exec(sql, realtime.DeviceId, realtime.Code, realtime.Ts, realtime.Value)
	var err error
	for _, realtime := range realtimeList {
		tableName := fmt.Sprintf("%v%d%v%d", "realtimedatatest.yx_", realtime.DeviceId, "_", realtime.Code)
		//r.taosDb.Table(tableName).Create(&realtime)
		sql := fmt.Sprintf("INSERT INTO %s (ts, val) using realtimedatatest.yx tags(?, ?) VALUES (?, ?)", tableName)
		_, err := r.taosDb.Exec(sql, realtime.DeviceId, realtime.Code, realtime.Ts, realtime.Value)
		err = err
	}
	return err
}

/*
*
批量添加yc
*/
func (r *RealtimeDataRepository) BatchCreateYc(realtimeList []*models.YxData) error {
	// 定义查询参数
	//tableName := fmt.Sprintf("%d%d", realtime.DeviceId, realtime.Code)
	//sql := fmt.Sprintf("INSERT INTO %s (ts, val) using realtimedatatest.yx tags(?, ?) VALUES (?, ?)", tableName)
	//result := r.taosDb.Exec(sql, realtime.DeviceId, realtime.Code, realtime.Ts, realtime.Value)
	var err error
	for _, realtime := range realtimeList {
		tableName := fmt.Sprintf("%v%d%v%d", "realtimedatatest.yc_", realtime.DeviceId, "_", realtime.Code)
		//r.taosDb.Table(tableName).Create(&realtime)
		sql := fmt.Sprintf("INSERT INTO %s (ts, val) using realtimedatatest.yc tags(?, ?) VALUES (?, ?)", tableName)
		_, err := r.taosDb.Exec(sql, realtime.DeviceId, realtime.Code, realtime.Ts, realtime.Value)
		err = err
	}
	return err
}

/*
*
批量添加setting
*/
func (r *RealtimeDataRepository) BatchCreateSetting(realtimeList []*models.YxData) error {
	// 定义查询参数
	//tableName := fmt.Sprintf("%d%d", realtime.DeviceId, realtime.Code)
	//sql := fmt.Sprintf("INSERT INTO %s (ts, val) using realtimedatatest.yx tags(?, ?) VALUES (?, ?)", tableName)
	//result := r.taosDb.Exec(sql, realtime.DeviceId, realtime.Code, realtime.Ts, realtime.Value)
	var err error
	for _, realtime := range realtimeList {
		tableName := fmt.Sprintf("%v%d%v%d", "realtimedatatest.setting_", realtime.DeviceId, "_", realtime.Code)
		//r.taosDb.Table(tableName).Create(&realtime)
		sql := fmt.Sprintf("INSERT INTO %s (ts, val) using realtimedatatest.setting tags(?, ?) VALUES (?, ?)", tableName)
		_, err := r.taosDb.Exec(sql, realtime.DeviceId, realtime.Code, realtime.Ts, realtime.Value)
		err = err
	}
	return err
}

//func (r *RealtimeDataRepository) UpdateYx(realtime *models.YxData) error {
//	return r.taosDb.Save(realtime).Error
//}
//func (r *RealtimeDataRepository) UpdateYc(realtime *models.YcData) error {
//	return r.taosDb.Save(realtime).Error
//}
//func (r *RealtimeDataRepository) UpdateSetting(realtime *models.SettingData) error {
//	return r.taosDb.Save(realtime).Error
//}

//func (r *RealtimeDataRepository) DeleteYx(deviceId, code int) error {
//	return r.taosDb.Delete(&models.YxData{}, deviceId, code).Error
//}
//func (r *RealtimeDataRepository) DeleteYc(deviceId, code int) error {
//	return r.taosDb.Delete(&models.YcData{}, deviceId, code).Error
//}
//func (r *RealtimeDataRepository) DeleteSetting(deviceId, code int) error {
//	return r.taosDb.Delete(&models.SettingData{}, deviceId, code).Error
//}
/**
获取yx
*/
func (r *RealtimeDataRepository) GetYxById(deviceId, code int) (models.YxData, error) {
	var realtime models.YxData
	//tableName := fmt.Sprintf("%v%d%v%d", "realtimedatatest.yx_", deviceId, "_", code)
	//sql := fmt.Sprintf("select last(*) from ? ", tableName)
	sql := fmt.Sprint("select Last(ts), val, device_id, code from realtimedatatest.yx where device_id =", deviceId, "and code =", code)
	fmt.Println(sql)
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
获取yc
*/
func (r *RealtimeDataRepository) GetYcById(deviceId, code int) (models.YcData, error) {
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
func (r *RealtimeDataRepository) GetSettingById(deviceId, code int) (models.SettingData, error) {
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
