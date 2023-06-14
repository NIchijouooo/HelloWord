package repositories

import (
	sqlt "database/sql"
	"fmt"
	"gateway/models"
	"gorm.io/gorm"
)

// 定义字典类型管理的存储库
type HistoryDataRepository struct {
	db     *gorm.DB
	taosDb *sqlt.DB
}

func NewHistoryDataRepository() *HistoryDataRepository {
	return &HistoryDataRepository{db: models.DB, taosDb: models.TaosDB}
}

/*
*

	批量获取yx
	必传deviceIds直接给加逗号拼接好的字符串“121,111,131”
	必传codes直接给加逗号拼接好的字符串“1,2,3”
	startTime，endTime为时间戳
	interval按语法传参
*/
func (r *HistoryDataRepository) GetYxLogByDeviceIdsCodes(deviceIds, codes, interval string, startTime, endTime int64) ([]*models.PointParam, error) {
	var realtimeList []*models.PointParam
	tableName := "realtimedatatest.yx"
	var err error
	var sql = ""
	if startTime > 0 && endTime > 0 && interval == "" {
		sql = fmt.Sprintf("select * from %s Where device_id in (%s) and code in (%s) and ts >=%v and ts <%v", tableName, deviceIds, codes, startTime, endTime)
		if len(sql) <= 0 {
			return nil, err
		}

		rows, err := r.taosDb.Query(sql)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			realtime := &models.PointParam{}
			var value sqlt.NullString
			err := rows.Scan(&realtime.Ts, &value, &realtime.DeviceId, &realtime.Code)
			if err != nil {
				return nil, err
			}
			realtime.Value = &value.String
			realtimeList = append(realtimeList, realtime)
		}
	} else if startTime > 0 && endTime > 0 && interval != "" {
		sql = fmt.Sprintf("select sum(val) as val, device_id, code from %s Where device_id in (%s) and code in (%s) and ts >=%v and ts <%v partition by device_id,code interval(%s) FILL(NULL)", tableName, deviceIds, codes, startTime, endTime, interval)
		if len(sql) <= 0 {
			return nil, err
		}

		rows, err := r.taosDb.Query(sql)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			realtime := &models.PointParam{}
			var value sqlt.NullString
			err := rows.Scan(&value, &realtime.DeviceId, &realtime.Code)
			if err != nil {
				return nil, err
			}
			realtime.Value = &value.String
			realtimeList = append(realtimeList, realtime)
		}
	}
	return realtimeList, err
}

/*
*

	批量获取yc
	必传deviceIds直接给加逗号拼接好的字符串“121,111,131”
	必传codes直接给加逗号拼接好的字符串“1,2,3”
	startTime，endTime为时间戳
	interval按语法传参
*/
func (r *HistoryDataRepository) GetYcLogByDeviceIdsCodes(deviceIds, codes, interval string, startTime, endTime int64) ([]*models.PointParam, error) {
	var realtimeList []*models.PointParam
	tableName := "realtimedatatest.yc"
	var err error
	var sql = ""
	if startTime > 0 && endTime > 0 && interval == "" {
		sql = fmt.Sprintf("select * from %s Where device_id in (%s) and code in (%s) and ts >=%v and ts <%v", tableName, deviceIds, codes, startTime, endTime)
		if len(sql) <= 0 {
			return nil, err
		}
		rows, err := r.taosDb.Query(sql)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			realtime := &models.PointParam{}
			var value sqlt.NullString
			err := rows.Scan(&realtime.Ts, &value, &realtime.DeviceId, &realtime.Code)
			if err != nil {
				return nil, err
			}
			realtime.Value = &value.String
			realtimeList = append(realtimeList, realtime)
		}
	} else if startTime > 0 && endTime > 0 && interval != "" {
		sql = fmt.Sprintf("select sum(val) as val, device_id, code from %s Where device_id in (%s) and code in (%s) and ts >=%v and ts <%v partition by device_id,code interval(%s) FILL(NULL)", tableName, deviceIds, codes, startTime, endTime, interval)
		if len(sql) <= 0 {
			return nil, err
		}
		rows, err := r.taosDb.Query(sql)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			realtime := &models.PointParam{}
			var value sqlt.NullString
			err := rows.Scan(&value, &realtime.DeviceId, &realtime.Code)
			if err != nil {
				return nil, err
			}
			realtime.Value = &value.String
			realtimeList = append(realtimeList, realtime)
		}
	}
	return realtimeList, err
}

/*
*

	批量获取setting
	必传deviceIds直接给加逗号拼接好的字符串“121,111,131”
	必传codes直接给加逗号拼接好的字符串“1,2,3”
	startTime，endTime为时间戳
	interval按语法传参
*/
func (r *HistoryDataRepository) GetSettingLogByDeviceIdsCodes(deviceIds, codes, interval string, startTime, endTime int64) ([]*models.SettingData, error) {
	var realtimeList []*models.SettingData
	tableName := "realtimedatatest.setting"
	var err error

	var sql = ""
	if startTime > 0 && endTime > 0 && interval == "" {
		sql = fmt.Sprintf("select * from %s Where device_id in (%s) and code in (%s) and ts >=%v and ts <%v", tableName, deviceIds, codes, startTime, endTime)
	}
	//else if startTime > 0 && endTime > 0 && interval != ""{
	//	sql = fmt.Sprintf("select sum(val) as val, device_id, code from %s Where device_id in (%s) and code in (%s) and ts >=%v and ts <%v partition by device_id,code interval(%s) FILL(NULL)", tableName, deviceIds, codes, startTime, endTime, interval)
	//}
	if len(sql) <= 0 {
		return nil, err
	}
	rows, err := r.taosDb.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		realtime := &models.SettingData{}
		err := rows.Scan(&realtime.Ts, &realtime.Value, &realtime.DeviceId, &realtime.Code)
		if err != nil {
			return nil, err
		}
		realtimeList = append(realtimeList, realtime)
	}
	return realtimeList, err
}

///**
//	deviceIds
// */
//func (r *HistoryDataRepository) GetYxLogBydeviceIds(code int, deviceIds, interval string, startTime, endTime int64) ([]*models.YxData, error) {
//	var realtimeList []*models.YxData
//	tableName := "realtimedatatest.yc"
//	var err error
//	var sql = ""
//	if startTime > 0 && endTime > 0 && interval == "" {
//		sql = fmt.Sprintf("select * from %s Where code = %d and device_id in (%s) and ts >=%v and ts <%v", tableName, code, deviceIds, startTime, endTime)
//	} else if startTime > 0 && endTime > 0 && interval != ""{
//		sql = fmt.Sprintf("select * from %s Where code = %d and device_id in (%s) and ts >=%v and ts <%v partition by device_id,code interval(%s) FILL(NULL)", tableName, code, deviceIds, startTime, endTime, interval)
//	}
//	if len(sql) <= 0{return nil, err}
//
//	rows, err := r.taosDb.Query(sql)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	for rows.Next() {
//		realtime := &models.YxData{}
//		err := rows.Scan(&realtime.Ts, &realtime.Value, &realtime.DeviceID, &realtime.Code)
//		if err != nil {
//			return nil, err
//		}
//		realtimeList = append(realtimeList, realtime)
//	}
//	return realtimeList, err
//}
//
//func (r *HistoryDataRepository) GetYcLogBydeviceIds(code int, deviceIds, interval string, startTime, endTime int64) ([]*models.YcData, error) {
//	var realtimeList []*models.YcData
//	tableName := "realtimedatatest.yc"
//	var err error
//	var sql = ""
//	if startTime > 0 && endTime > 0 && interval == "" {
//		sql = fmt.Sprintf("select * from %s Where code = %d and device_id in (%s) and ts >=%v and ts <%v", tableName, code, deviceIds, startTime, endTime)
//	} else if startTime > 0 && endTime > 0 && interval != ""{
//		sql = fmt.Sprintf("select * from %s Where code = %d and device_id in (%s) and ts >=%v and ts <%v partition by device_id,code interval(%s) FILL(NULL)", tableName, code, deviceIds, startTime, endTime, interval)
//	}
//	if len(sql) <= 0{return nil, err}
//
//	rows, err := r.taosDb.Query(sql)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	for rows.Next() {
//		realtime := &models.YcData{}
//		err := rows.Scan(&realtime.Ts, &realtime.Value, &realtime.DeviceID, &realtime.Code)
//		if err != nil {
//			return nil, err
//		}
//		realtimeList = append(realtimeList, realtime)
//	}
//	return realtimeList, err
//}
//
//func (r *HistoryDataRepository) GetSettingLogBydeviceIds(code int, deviceIds, interval string, startTime, endTime int64) ([]*models.SettingData, error) {
//	var realtimeList []*models.SettingData
//	tableName := "realtimedatatest.yc"
//	var err error
//	var sql = ""
//	if startTime > 0 && endTime > 0 && interval == "" {
//		sql = fmt.Sprintf("select * from %s Where code = %d and device_id in (%s) and ts >=%v and ts <%v", tableName, code, deviceIds, startTime, endTime)
//	} else if startTime > 0 && endTime > 0 && interval != ""{
//		sql = fmt.Sprintf("select * from %s Where code = %d and device_id in (%s) and ts >=%v and ts <%v partition by device_id,code interval(%s) FILL(NULL)", tableName, code, deviceIds, startTime, endTime, interval)
//	}
//	if len(sql) <= 0{return nil, err}
//
//	rows, err := r.taosDb.Query(sql)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	for rows.Next() {
//		realtime := &models.SettingData{}
//		err := rows.Scan(&realtime.Ts, &realtime.Value, &realtime.DeviceID, &realtime.Code)
//		if err != nil {
//			return nil, err
//		}
//		realtimeList = append(realtimeList, realtime)
//	}
//	return realtimeList, err
//}
