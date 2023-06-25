package repositories

import (
	"database/sql"
	"fmt"
	"gateway/models"
	"gateway/utils"
	"gorm.io/gorm"
	"log"
)

// 定义字典类型管理的存储库
type RealtimeDataRepository struct {
	db     *gorm.DB
	taosDb *sql.DB
}

type Res struct {
	Ts  string  `json:"ts"`
	Val float32 `json:"val"`
}

func NewRealtimeDataRepository() *RealtimeDataRepository {
	return &RealtimeDataRepository{db: models.DB, taosDb: models.TaosDB}
}

/*
*
添加db
*/
func (r *RealtimeDataRepository) CreateDB() error {
	sql := fmt.Sprintf("create database if not exists realtimedata cachemodel 'both';")
	_, err := r.taosDb.Exec(sql)
	return err
}

//realtimedata.charge_discharge_${item.deviceId} (ts,charge_capacity,discharge_capacity,profit) using realtimedata.charge_discharge
//tags(${item.deviceId})
//values (#{item.ts}, #{item.chargeCapacity}, #{item.dischargeCapacity}, #{item.profit})

/*
*
添加charge_discharge表
*/
func (r *RealtimeDataRepository) CreateChargeDischargeTable() error {
	// 定义查询参数
	sql := fmt.Sprintf("create table if not exists realtimedata.charge_discharge (ts timestamp, charge_capacity double, discharge_capacity double, profit double) tags (device_id int);")
	_, err := r.taosDb.Exec(sql)
	return err
}

/*
*
添加charge_discharge_hour表
*/
func (r *RealtimeDataRepository) CreateChargeDischargeHourTable() error {
	// 定义查询参数
	sql := fmt.Sprintf("create table if not exists realtimedata.charge_discharge_hour (ts timestamp, charge_capacity double, discharge_capacity double, profit double) tags (device_id int);")
	_, err := r.taosDb.Exec(sql)
	return err
}

/*
*
添加yx表
*/
func (r *RealtimeDataRepository) CreateYxTable() error {
	// 定义查询参数
	sql := fmt.Sprintf("create table if not exists realtimedata.yx (ts timestamp, val int) tags (device_id int, code int, identifier NCHAR);")
	_, err := r.taosDb.Exec(sql)
	return err
}

/*
*
添加yc表
*/
func (r *RealtimeDataRepository) CreateYcTable() error {
	// 定义查询参数
	sql := fmt.Sprintf("create table if not exists realtimedata.yc (ts timestamp, val double) tags (device_id int, code int, identifier NCHAR);")
	_, err := r.taosDb.Exec(sql)
	return err
}

/*
*
添加setting表
*/
func (r *RealtimeDataRepository) CreateSettingTable() error {
	// 定义查询参数
	sql := fmt.Sprintf("create table if not exists realtimedata.setting (ts timestamp, val NCHAR(16)) tags (device_id int, code int, identifier NCHAR);")
	_, err := r.taosDb.Exec(sql)
	return err
}

/*
*
添加yx
*/
func (r *RealtimeDataRepository) CreateYx(realtime *models.YxData) error {
	// 定义查询参数
	tableName := fmt.Sprintf("%v%d%v%d", "realtimedata.yx_", realtime.DeviceId, "_", realtime.Code)
	sql := fmt.Sprintf("INSERT INTO %s (ts, val) using realtimedata.yx tags(%d, %d) VALUES (%v, %d)", tableName, realtime.DeviceId, realtime.Code, realtime.Ts, realtime.Value)
	_, err := r.taosDb.Exec(sql)
	return err
}

/*
*
添加yc
*/
func (r *RealtimeDataRepository) CreateYc(realtime *models.YcData) error {
	// 定义查询参数
	tableName := fmt.Sprintf("%v%d%v%d", "realtimedata.yc_", realtime.DeviceId, "_", realtime.Code)
	sql := fmt.Sprintf("INSERT INTO %s (ts, val) using realtimedata.yc tags(%d, %d) VALUES (%v, %d)", tableName, realtime.DeviceId, realtime.Code, realtime.Ts, realtime.Value)
	_, err := r.taosDb.Exec(sql)
	return err
}

/*
*
添加setting
*/
func (r *RealtimeDataRepository) CreateSetting(realtime *models.SettingData) error {
	// 定义查询参数
	tableName := fmt.Sprintf("%v%d%v%d", "realtimedata.setting_", realtime.DeviceId, "_", realtime.Code)
	sql := fmt.Sprintf("INSERT INTO %s (ts, val) using realtimedata.setting tags(%d, %d) VALUES (%v, %d)", tableName, realtime.DeviceId, realtime.Code, realtime.Ts, realtime.Value)
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
		tableName := fmt.Sprintf("%v%d%v%d", "realtimedata.yx_", realtime.DeviceId, "_", realtime.Code)
		sql := fmt.Sprintf("INSERT INTO %s (ts, val) using realtimedata.yx tags(?, ?) VALUES (?, ?)", tableName)
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
		tableName := fmt.Sprintf("%v%d%v%d", "realtimedata.yc_", realtime.DeviceId, "_", realtime.Code)
		sql := fmt.Sprintf("INSERT INTO %s (ts, val) using realtimedata.yc tags(?, ?) VALUES (?, ?)", tableName)
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
		tableName := fmt.Sprintf("%v", "realtimedata.charge_discharge_", realtime.DeviceId)
		sql := fmt.Sprintf("INSERT INTO %s (ts,charge_capacity,discharge_capacity,profit) using realtimedata.charge_discharge tags(?) VALUES (?, ?, ?, ?)", tableName)
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
	//sql := fmt.Sprintf("INSERT INTO %s (ts, val) using realtimedata.yx tags(?, ?) VALUES (?, ?)", tableName)
	//result := r.taosDb.Exec(sql, realtime.DeviceId, realtime.Code, realtime.Ts, realtime.Value)
	var err error
	for _, realtime := range realtimeList {
		tableName := fmt.Sprintf("%v%d%v%d", "realtimedata.setting_", realtime.DeviceId, "_", realtime.Code)
		//r.taosDb.Table(tableName).Create(&realtime)
		sql := fmt.Sprintf("INSERT INTO %s (ts, val) using realtimedata.setting tags(?, ?) VALUES (?, ?)", tableName)
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
	//tableName := fmt.Sprintf("%v%d%v%d", "realtimedata.yx_", deviceId, "_", code)
	//sql := fmt.Sprintf("select last(*) from ? ", tableName)
	sql := fmt.Sprint("select Last(ts), val, device_id, code from realtimedata.yx where device_id =", deviceId, "and code =", code)
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
	//tableName := fmt.Sprintf("%v%d%v%d", "realtimedata.yc_", deviceId, "_", code)
	sql := fmt.Sprint("select Last(ts), val, device_id, code from realtimedata.yc where device_id =", deviceId, "and code =", code)
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
	//tableName := fmt.Sprintf("%v%d%v%d", "realtimedata.setting_", deviceId, "_", code)
	sql := fmt.Sprint("select Last(ts), val, device_id, code from realtimedata.setting where device_id =", deviceId, "and code =", code)
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

func (r *RealtimeDataRepository) GetSettingPointParamListById(deviceId int) ([]models.PointParam, error) {
	rowsYx, err := r.taosDb.Query("SELECT last(code), last(name), last(val),last(ts) FROM setting where device_id = ? group by code order by code", deviceId)
	if err != nil {
		return nil, err
	}
	var list []models.PointParam
	for rowsYx.Next() {
		pointParam := models.PointParam{}
		err := rowsYx.Scan(&pointParam.Code, &pointParam.Name, &pointParam.Value, &pointParam.Ts)
		if err != nil {
			return nil, err
		}
		list = append(list, pointParam)
	}
	return list, err
}
func (r *RealtimeDataRepository) GetYxPointParamListById(deviceId int) ([]models.PointParam, error) {
	rowsYx, err := r.taosDb.Query("SELECT last(code), last(name), last(val),last(ts) FROM yx where device_id = ? group by code order by code", deviceId)
	if err != nil {
		return nil, err
	}
	var list []models.PointParam
	for rowsYx.Next() {
		pointParam := models.PointParam{}
		err := rowsYx.Scan(&pointParam.Code, &pointParam.Name, &pointParam.Value, &pointParam.Ts)
		if err != nil {
			return nil, err
		}
		list = append(list, pointParam)
	}
	return list, err
}
func (r *RealtimeDataRepository) GetYcPointParamListById(deviceId int) ([]models.PointParam, error) {
	rowsYx, err := r.taosDb.Query("SELECT last(code), last(name), last(val),last(ts) FROM yc where device_id = ? group by code order by code", deviceId)
	if err != nil {
		return nil, err
	}
	var list []models.PointParam
	for rowsYx.Next() {
		pointParam := models.PointParam{}
		err := rowsYx.Scan(&pointParam.Code, &pointParam.Name, &pointParam.Value, &pointParam.Ts)
		if err != nil {
			return nil, err
		}
		list = append(list, pointParam)
	}
	return list, err
}

/*
*
获取setting集合
*/
func (r *RealtimeDataRepository) GetSettingListByDevIdsAndCodes(deviceIds, codes string) ([]models.SettingData, error) {
	var realtime []models.SettingData
	sqlStr := fmt.Sprintf("select Last(ts), val, device_id, code,name from realtimedata.setting where device_id in (%s) and code in (%s) group by device_id,code", deviceIds, codes)
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
	//tableName := fmt.Sprintf("%v%d%v%d", "realtimedata.setting_", deviceId, "_", code)
	sql := fmt.Sprintf("select Last(ts), val, device_id, code,name from realtimedata.yc where device_id in (%s) and code in (%s) group by device_id,code", deviceIds, codes)
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

func (r *RealtimeDataRepository) GetChartByDeviceIdAndCode(deviceId int, code string) ([]Res, error) {
	sql := fmt.Sprint("SELECT _WSTART AS ts,LAST(VAL) AS val FROM yc_100_66 WHERE ts>= NOW-2d and ts<=NOW+1d INTERVAL(1h) FILL(VALUE,0)")
	rows, err := r.taosDb.Query(sql)
	if err != nil {
		return nil, err
	}
	var list []Res
	for rows.Next() {
		yc := Res{}
		err := rows.Scan(&yc.Ts, &yc.Val)
		if err != nil {
			return nil, err
		}
		list = append(list, yc)
	}
	return list, err
}

// GetGenerateElectricityChartByDeviceIds 获取前一天每小时的充电电量
func (r *RealtimeDataRepository) GetGenerateElectricityChartByDeviceIds(deviceIds []int, fieldName string, intervalType string) ([]Res, error) {
	ids := utils.IntArrayToString(deviceIds, ",")
	sql := fmt.Sprintf("SELECT _WSTART AS ts,SUM(%v) AS charge_capacity FROM charge_discharge_hour WHERE device_id IN (%s) AND ts>= NOW-1d and ts<=NOW INTERVAL(1%s) FILL(VALUE,0);", ids, fieldName, intervalType)
	rows, err := r.taosDb.Query(sql)
	if err != nil {
		return nil, err
	}
	var list []Res
	for rows.Next() {
		yc := Res{}
		err := rows.Scan(&yc.Ts, &yc.Val)
		if err != nil {
			return nil, err
		}
		list = append(list, yc)
	}
	return list, err
}

// GetGenerateElectricitySumByDeviceIds 统计累计可充电电量信息
func (r *RealtimeDataRepository) GetGenerateElectricitySumByDeviceIds(deviceIds []int) Res {
	ids := utils.IntArrayToString(deviceIds, ",")
	var res Res
	sql := fmt.Sprintf("SELECT SUM(charge_capacity) AS val FROM charge_discharge WHERE device_id IN (%s)", ids)
	err := r.taosDb.QueryRow(sql).Scan(&res.Val)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

// GetReleaseElectricitySumByDeviceIds 统计累计可放电电量信息
func (r *RealtimeDataRepository) GetReleaseElectricitySumByDeviceIds(deviceIds []int) Res {
	ids := utils.IntArrayToString(deviceIds, ",")
	var res Res
	sql := fmt.Sprintf("SELECT SUM(discharge_capacity) AS val FROM charge_discharge WHERE device_id IN (%s)", ids)
	err := r.taosDb.QueryRow(sql).Scan(&res.Val)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

// GetGenerateElectricitySumByDeviceIds 统计累计每日收益信息
func (r *RealtimeDataRepository) GetProfitChartByDeviceIds(deviceIds []int, startTime int64, endTime int64) ([]Res, error) {
	ids := utils.IntArrayToString(deviceIds, ",")
	sql := fmt.Sprintf("SELECT _WSTART AS ts,SUM(profit) AS profit FROM charge_discharge WHERE device_id IN (%s) AND ts>= %v and ts<=%v INTERVAL(1d) FILL(VALUE,0)", ids, startTime, endTime)
	rows, err := r.taosDb.Query(sql)
	if err != nil {
		return nil, err
	}
	var list []Res
	for rows.Next() {
		yc := Res{}
		err := rows.Scan(&yc.Ts, &yc.Val)
		if err != nil {
			return nil, err
		}
		list = append(list, yc)
	}
	return list, err
}

// GetProfitSumByDeviceIds 统计累计收益信息
func (r *RealtimeDataRepository) GetProfitSumByDeviceIds(deviceIds []int, startTime int64, endTime int64) Res {
	ids := utils.IntArrayToString(deviceIds, ",")
	var res Res
	sql := fmt.Sprintf("SELECT SUM(profit) AS val FROM charge_discharge WHERE device_id IN (%s) AND ts >= %v AND ts <= %v", ids, startTime, endTime)
	err := r.taosDb.QueryRow(sql).Scan(&res.Val)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

// GetGenerateElectricityChartByDeviceIds 获取充放电量信息
func (r *RealtimeDataRepository) GetElectricityChartByDeviceIds(deviceIds []int, fieldName string, startTime, endTime int64, intervalType string, tableName string) ([]Res, error) {
	ids := utils.IntArrayToString(deviceIds, ",")
	sql := fmt.Sprintf("SELECT _WSTART AS ts,SUM(%[2]v) AS %[2]v FROM %[6]v WHERE device_id IN (%[1]s) AND ts>=%[3]d and ts<=%[4]d INTERVAL(1%[5]s) FILL(VALUE,0);", ids, fieldName, startTime, endTime, intervalType, tableName)
	rows, err := r.taosDb.Query(sql)
	if err != nil {
		return nil, err
	}
	var list []Res
	for rows.Next() {
		yc := Res{}
		err := rows.Scan(&yc.Ts, &yc.Val)
		if err != nil {
			return nil, err
		}
		list = append(list, yc)
	}
	return list, err
}
