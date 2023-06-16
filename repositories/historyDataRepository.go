package models

import (
	"database/sql"
	"fmt"
	"gateway/models"
	"gorm.io/gorm"
)

// 定义字典类型管理的存储库
type HistoryDataRepository struct {
	db     *gorm.DB
	taosDb *sql.DB
}

func NewHistoryDataRepository() *HistoryDataRepository {
	return &HistoryDataRepository{db: models.DB, taosDb: models.TaosDB}
}

/*
*
批量获取yx
*/
func (r *HistoryDataRepository) GetYxLogById(deviceId int, codes string, startTime, endTime int64) ([]*models.YxData, error) {
	var realtimeList []*models.YxData
	tableName := "realtimedatatest.yx"

	sql := fmt.Sprintf("select * from %s Where device_id = %d and code in (%s) and ts >=%v and ts <%v", tableName, deviceId, codes, startTime, endTime)
	rows, err := r.taosDb.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		realtime := &models.YxData{}
		err := rows.Scan(&realtime.Ts, &realtime.Value, &realtime.DeviceID, &realtime.Code)
		if err != nil {
			return nil, err
		}
		realtimeList = append(realtimeList, realtime)
	}
	return realtimeList, err
}

/*
*
批量获取yc
*/
func (r *HistoryDataRepository) GetYcLogById(deviceId int, codes string, startTime, endTime int64) ([]*models.YcData, error) {
	var realtimeList []*models.YcData
	tableName := "realtimedatatest.yc"
	sql := fmt.Sprintf("select * from %s Where device_id = %d and code in (%s) and ts >=%v and ts <%v", tableName, deviceId, codes, startTime, endTime)
	rows, err := r.taosDb.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		realtime := &models.YcData{}
		err := rows.Scan(&realtime.Ts, &realtime.Value, &realtime.DeviceID, &realtime.Code)
		if err != nil {
			return nil, err
		}
		realtimeList = append(realtimeList, realtime)
	}
	return realtimeList, err
}

/*
*
批量获取setting
*/
func (r *HistoryDataRepository) GetSettingLogById(deviceId int, codes string, startTime, endTime int64) ([]*models.SettingData, error) {
	var realtimeList []*models.SettingData
	tableName := "realtimedatatest.setting"
	sql := fmt.Sprintf("select * from %s Where device_id = %d and code in (%s) and ts >=%v and ts <%v", tableName, deviceId, codes, startTime, endTime)
	rows, err := r.taosDb.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		realtime := &models.SettingData{}
		err := rows.Scan(&realtime.Ts, &realtime.Value, &realtime.DeviceID, &realtime.Code)
		if err != nil {
			return nil, err
		}
		realtimeList = append(realtimeList, realtime)
	}
	return realtimeList, err
}

//批量code获取yc的最新一条信息
func (r *HistoryDataRepository) GetLastYcListByCode(codes string) ([]*models.YcData, error) {
	var realtimeList []*models.YcData
	tableName := "realtimedatatest.yc"
	//
	sql := fmt.Sprintf("SELECT last(ts),last(val),last(device_id),last(code) FROM %s  where code in (%s) group by code", tableName, codes)
	rows, err := r.taosDb.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		realtime := &models.YcData{}
		err := rows.Scan(&realtime.Ts, &realtime.Value, &realtime.DeviceID, &realtime.Code)
		if err != nil {
			return nil, err
		}
		realtimeList = append(realtimeList, realtime)
	}
	return realtimeList, err
}

//根据时间段获取历史数据数据
func (r *HistoryDataRepository) GetLastYcHistoryByDeviceIdAndCodeList(deviceId int, codes string, startTime, endTime int64, interval string) ([]*models.YcData, error) {
	var realtimeList []*models.YcData
	tableName := "realtimedata.yc"
	//select  from realtimedatatest.yc Where device_id = 121 and code in (1,2) and ts >=1683689768000 and ts<1684553768000 partition by device_id,code INTERVAL(1d) FILL(null);
	sql := fmt.Sprintf("select FIRST(ts) as ts,val,device_id,code from %s Where device_id = %d and code in (%s) and ts >=%v and ts <%v  partition by code INTERVAL(%s);", tableName, deviceId, codes, startTime, endTime, interval)
	rows, err := r.taosDb.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		realtime := &models.YcData{}
		err := rows.Scan(&realtime.Ts, &realtime.Value, &realtime.DeviceID, &realtime.Code)
		if err != nil {
			return nil, err
		}
		realtimeList = append(realtimeList, realtime)
	}
	return realtimeList, err
}
