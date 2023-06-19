package repositories

import (
	"database/sql"
	"fmt"
	"gateway/models"
	"gorm.io/gorm"
	"log"
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
	sql := fmt.Sprintf("create database if not exists realtimedatatest cachemodel 'both';")
	_, err := r.taosDb.Exec(sql)
	return err
}

/*
*
添加yx表
*/
func (r *RealtimeDataRepository) CreateYxTable() error {
	// 定义查询参数
	sql := fmt.Sprintf("create table if not exists realtimedatatest.yx (ts timestamp, val int) tags (device_id int, code int, identifier NCHAR);")
	_, err := r.taosDb.Exec(sql)
	return err
}

/*
*
添加yc表
*/
func (r *RealtimeDataRepository) CreateYcTable() error {
	// 定义查询参数
	sql := fmt.Sprintf("create table if not exists realtimedatatest.yc (ts timestamp, val double) tags (device_id int, code int, identifier NCHAR);")
	_, err := r.taosDb.Exec(sql)
	return err
}

/*
*
添加setting表
*/
func (r *RealtimeDataRepository) CreateSettingTable() error {
	// 定义查询参数
	sql := fmt.Sprintf("create table if not exists realtimedatatest.setting (ts timestamp, val NCHAR(16)) tags (device_id int, code int, identifier NCHAR);")
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
	var err error
	for _, realtime := range realtimeList {
		tableName := fmt.Sprintf("%v%d%v%d", "realtimedatatest.yx_", realtime.DeviceId, "_", realtime.Code)
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
func (r *RealtimeDataRepository) BatchCreateYc(realtimeList []*models.YcData) error {
	// 定义查询参数
	var err error
	for _, realtime := range realtimeList {
		tableName := fmt.Sprintf("%v%d%v%d", "realtimedatatest.yc_", realtime.DeviceId, "_", realtime.Code)
		sql := fmt.Sprintf("INSERT INTO %s (ts, val) using realtimedatatest.yc tags(?, ?) VALUES (?, ?)", tableName)
		_, err := r.taosDb.Exec(sql, realtime.DeviceId, realtime.Code, realtime.Ts, realtime.Value)
		err = err
	}
	return err
}

/*
*
批量添加充放电和收益
*/
func (r *RealtimeDataRepository) BatchCreateEsChargeDischargeModel(realtimeList []*models.EsChargeDischargeModel) error {
	var err error
	for _, realtime := range realtimeList {
		tableName := fmt.Sprintf("%v", "realtimedatatest.charge_discharge_", realtime.DeviceId)
		sql := fmt.Sprintf("INSERT INTO %s (ts,charge_capacity,discharge_capacity,profit) using realtimedatatest.charge_discharge tags(?) VALUES (?, ?, ?, ?)", tableName)
		_, err := r.taosDb.Exec(sql, realtime.Ts, realtime.ChargeCapacity, realtime.DischargeCapacity, realtime.Profit)
		err = err
	}
	return err
}

/*
*
批量添加setting
*/
func (r *RealtimeDataRepository) BatchCreateSetting(realtimeList []*models.SettingData) error {
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
获取yx集合
*/
func (r *RealtimeDataRepository) GetYxListByDevIdsAndCodes(deviceIds, codes string) ([]models.YxData, error) {
	var realtime []models.YxData
	sqlStr := fmt.Sprintf("select Last(ts), val, device_id, code,name from yx where device_id in (%s) and code in (%s)", deviceIds, codes)
	fmt.Println(sqlStr)
	rows, err := r.taosDb.Query(sqlStr)
	defer rows.Close()

	for rows.Next() {
		realtimeData := models.YxData{}
		err := rows.Scan(&realtimeData.Ts, &realtimeData.Value, &realtimeData.DeviceId, &realtimeData.Code, &realtimeData.Name)
		if err != nil {
			log.Printf("Request params:%v", err)
		}
		realtime = append(realtime, realtimeData)
	}
	return realtime, err
}
/*
*
获取yx
*/
func (r *RealtimeDataRepository) GetYxListById(deviceId int) ([]models.YxData, error) {
	rowsYx, err := r.taosDb.Query("SELECT last(code), last(name), last(val), last(ts) FROM yx where device_id = ? group by code order by code", deviceId)
	if err != nil {
		return nil, err
	}
	var yxList []models.YxData
	for rowsYx.Next() {
		yx := models.YxData{}
		err := rowsYx.Scan(&yx.Code, &yx.Name, &yx.Value, &yx.Ts)
		if err != nil {
			return nil, err
		}
		yxList = append(yxList, yx)
	}
	return yxList, err
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
获取yc集合
*/
func (r *RealtimeDataRepository) GetYcListByDevIdsAndCodes(deviceIds, codes string) ([]models.YcData, error) {
	var realtime []models.YcData
	sqlStr := fmt.Sprintf("select Last(ts), val, device_id, code,name from yc where device_id in (%s) and code in (%s) group by device_id,code", deviceIds, codes)
	fmt.Println(sqlStr)
	rows, err := r.taosDb.Query(sqlStr)
	defer rows.Close()

	for rows.Next() {
		realtimeData := models.YcData{}
		err := rows.Scan(&realtimeData.Ts, &realtimeData.Value, &realtimeData.DeviceId, &realtimeData.Code, &realtimeData.Name)
		if err != nil {
			log.Printf("Request params:%v", err)
		}
		realtime = append(realtime, realtimeData)
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

func (r *RealtimeDataRepository) GetYcListById(deviceId int) ([]models.YcData, error) {
	rowsYx, err := r.taosDb.Query("SELECT last(code), last(name), last(val),last(ts) FROM yc where device_id = ? group by code order by code", deviceId)
	if err != nil {
		return nil, err
	}
	var list []models.YcData
	for rowsYx.Next() {
		yc := models.YcData{}
		err := rowsYx.Scan(&yc.Code, &yc.Name, &yc.Value, &yc.Ts)
		if err != nil {
			return nil, err
		}
		list = append(list, yc)
	}
	return list, err
}


/*
*
获取setting集合
*/
func (r *RealtimeDataRepository) GetSettingListByDevIdsAndCodes(deviceIds, codes string) ([]models.SettingData, error) {
	var realtime []models.SettingData
	sqlStr := fmt.Sprintf("select Last(ts), val, device_id, code,name from realtimedatatest.setting where device_id in (%s) and code in (%s) group by device_id,code", deviceIds, codes)
	fmt.Println(sqlStr)
	rows, err := r.taosDb.Query(sqlStr)
	defer rows.Close()

	for rows.Next() {
		realtimeData := models.SettingData{}
		err := rows.Scan(&realtimeData.Ts, &realtimeData.Value, &realtimeData.DeviceId, &realtimeData.Code, &realtimeData.Name)
		if err != nil {
			log.Printf("Request params:%v", err)
		}
		realtime = append(realtime, realtimeData)
	}
	return realtime, err
}

/*
*
获取多个yc，好像写重复了
*/
func (r *RealtimeDataRepository) GetLastYcHistoryByDeviceIdListAndCodeList(deviceIds, codes, interval string, beginDt int64, endDt int64) ([]models.YcData, error) {
	var realtimeData []models.YcData
	//tableName := fmt.Sprintf("%v%d%v%d", "realtimedatatest.setting_", deviceId, "_", code)
	sql := fmt.Sprintf("select Last(ts), val, device_id, code,name from realtimedatatest.yc where device_id in (%s) and code in (%s) group by device_id,code", deviceIds, codes)
	rows, err := r.taosDb.Query(sql)
	defer rows.Close()

	for rows.Next() {
		realtime := models.YcData{}
		err := rows.Scan(&realtime.Ts, &realtime.Value, &realtime.DeviceId, &realtime.Code)
		if err != nil {
			log.Printf("Request params:%v", err)
		}
	}
	return realtimeData, err
}
