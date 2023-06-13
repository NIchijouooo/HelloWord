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
