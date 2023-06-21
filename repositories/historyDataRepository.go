package repositories

import (
	sqlt "database/sql"
	"fmt"
	"gateway/models"
	"gateway/models/ReturnModel"
	"gateway/models/query"
	"gateway/service"
	"gateway/utils"
	"github.com/taosdata/driver-go/v3/errors"
	"gorm.io/gorm"
	"strconv"
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
	tableName := "realtimedata.yx"
	var err error
	var sql = ""
	if startTime > 0 && endTime > 0 && interval == "" {
		sql = fmt.Sprintf("select ts, val, device_id, code from %s Where device_id in (%s) and code in (%s) and ts >=%v and ts <%v", tableName, deviceIds, codes, startTime, endTime)
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
	tableName := "realtimedata.yc"
	var err error
	var sql = ""
	if startTime > 0 && endTime > 0 && interval == "" {
		sql = fmt.Sprintf("select ts, val, device_id, code from %s Where device_id in (%s) and code in (%s) and ts >=%v and ts <%v", tableName, deviceIds, codes, startTime, endTime)
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
	tableName := "realtimedata.setting"
	var err error

	var sql = ""
	if startTime > 0 && endTime > 0 && interval == "" {
		sql = fmt.Sprintf("select ts, val, device_id, code from %s Where device_id in (%s) and code in (%s) and ts >=%v and ts <%v", tableName, deviceIds, codes, startTime, endTime)
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

// 批量code获取yc的最新一条信息
func (r *HistoryDataRepository) GetLastYcListByCode(deviceIds, codes string) ([]*models.YcData, error) {
	var realtimeList []*models.YcData
	tableName := "realtimedata.yc"
	//
	sql := fmt.Sprintf("SELECT last(ts),val,device_id,name,code FROM %s  where device_id in (%s) and  code in (%s) group by device_id,code", tableName, deviceIds, codes)
	rows, err := r.taosDb.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		realtime := &models.YcData{}
		err := rows.Scan(&realtime.Ts, &realtime.Value, &realtime.DeviceId, &realtime.Name, &realtime.Code)
		if err != nil {
			return nil, err
		}
		realtimeList = append(realtimeList, realtime)
	}
	return realtimeList, err
}

// 根据时间段获取历史数据数据
func (r *HistoryDataRepository) GetLastYcHistoryByDeviceIdAndCodeList(deviceId int, codes string, startTime, endTime int64, interval string) ([]*models.YcData, error) {
	var realtimeList []*models.YcData
	tableName := "realtimedata.yc"
	//select  from realtimedata.yc Where device_id = 121 and code in (1,2) and ts >=1683689768000 and ts<1684553768000 partition by device_id,code INTERVAL(1d) FILL(null);
	sql := fmt.Sprintf("select Last(ts) as ts,val,device_id,code from %s Where device_id = %d and code in (%s) and ts >=%v and ts <%v  partition by code INTERVAL(%s);", tableName, deviceId, codes, startTime, endTime, interval)
	rows, err := r.taosDb.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		realtime := &models.YcData{}
		err := rows.Scan(&realtime.Ts, &realtime.Value, &realtime.DeviceId, &realtime.Code)
		if err != nil {
			return nil, err
		}
		realtimeList = append(realtimeList, realtime)
	}
	return realtimeList, err
}

// 根据设备id和遥测编码集合获取历史数据数据
func (r *HistoryDataRepository) GetLastYcHistoryByDeviceIdListAndCodeList(deviceIdList []int, codes string, startTime, endTime int64, interval string) ([]*models.YcData, error) {
	if len(deviceIdList) == 0 {
		return []*models.YcData{}, nil
	}
	var realtimeList []*models.YcData
	tableName := "realtimedata.yc"
	sql := fmt.Sprintf("select Last(ts) as ts,val,device_id,code from %s Where code in (%s) and ts >=%v and ts <%v  partition by device_id,code INTERVAL(%s);", tableName, codes, startTime, endTime, interval)
	rows, err := r.taosDb.Query(sql)
	if err != nil {
		return nil, err
	}
	// 设备id map,在map里的设备才返回数据
	devIdMap := make(map[int]bool)
	for _, devId := range deviceIdList {
		devIdMap[devId] = true
	}
	defer rows.Close()
	for rows.Next() {
		realtime := &models.YcData{}
		err := rows.Scan(&realtime.Ts, &realtime.Value, &realtime.DeviceId, &realtime.Code)
		if err != nil {
			return nil, err
		}
		if !devIdMap[realtime.DeviceId] {
			continue
		}
		realtimeList = append(realtimeList, realtime)
	}
	return realtimeList, err
}

/**
 * 获取充放电量大于日的降采样
 * @param deviceIdList
 * @param startTime
 * @param endTime
 * @return
 */
func (r *HistoryDataRepository) getDayEsChargeDischarge(deviceIds string, startTime, endTime int64) ([]*models.EsChargeDischargeModel, error) {
	var realtimeList []*models.EsChargeDischargeModel
	tableName := "realtimedata.charge_discharge"

	sql := fmt.Sprintf("select last_row(ts) as ts,charge_capacity as chargeCapacity,discharge_capacity as dischargeCapacity,profit,device_id as deviceId from %s Where device_id in (%s) and ts >=%v and ts <%v  partition by device_id INTERVAL(1d);", tableName, deviceIds, startTime, endTime)
	rows, err := r.taosDb.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		realtime := &models.EsChargeDischargeModel{}
		err := rows.Scan(&realtime.Ts, &realtime.ChargeCapacity, &realtime.DischargeCapacity, &realtime.Profit, &realtime.DeviceId)
		if err != nil {
			return nil, err
		}
		realtimeList = append(realtimeList, realtime)
	}
	return realtimeList, err
}

/**
 * 获取充放电量小时的降采样
 * @param deviceIdList
 * @param startTime
 * @param endTime
 * @return
 */
func (r *HistoryDataRepository) getDayEsChargeDischargeHour(deviceIds string, startTime, endTime int64) ([]*models.EsChargeDischargeModel, error) {
	var realtimeList []*models.EsChargeDischargeModel
	tableName := "realtimedata.charge_discharge_hour"

	sql := fmt.Sprintf("select last_row(ts) as ts,charge_capacity as chargeCapacity,discharge_capacity as dischargeCapacity,profit,device_id as deviceId from %s Where device_id in (%s) and ts >=%v and ts <%v  partition by device_id INTERVAL(1h);", tableName, deviceIds, startTime, endTime)
	rows, err := r.taosDb.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		realtime := &models.EsChargeDischargeModel{}
		err := rows.Scan(&realtime.Ts, &realtime.ChargeCapacity, &realtime.DischargeCapacity, &realtime.Profit, &realtime.DeviceId)
		if err != nil {
			return nil, err
		}
		realtimeList = append(realtimeList, realtime)
	}
	return realtimeList, err
}

/*
*
通用图表接口
*/
func (r *HistoryDataRepository) GetCharData(ycQuery query.QueryTaoData) (ReturnModel.CharData, error) {
	//必要条件校验
	if ycQuery.StartTime == 0 || ycQuery.EndTime == 0 || ycQuery.Interval == 0 || ycQuery.IntervalType == 0 || len(ycQuery.CodeList) == 0 {
		return ReturnModel.CharData{}, errors.NewError(1, "缺少必要参数")
	}
	//间隔时间段
	intervalStr := "s"
	if ycQuery.IntervalType == 2 {
		intervalStr = "m"
	} else if ycQuery.IntervalType == 3 {
		intervalStr = "h"
	} else if ycQuery.IntervalType > 3 {
		intervalStr = "d"
	}
	// 将 int 数组转换为字符串数组
	ycQuery.Codes = utils.IntArrayToString(ycQuery.CodeList, ",")

	//查询历史数据
	ycList, err := r.GetLastYcHistoryByDeviceIdListAndCodeList(ycQuery.DeviceIds, ycQuery.Codes, ycQuery.StartTime, ycQuery.EndTime, strconv.Itoa(ycQuery.Interval)+intervalStr)
	if err != nil {
		return ReturnModel.CharData{}, err
	}
	var xAxisList []string
	// 初始化x轴数据,返回x轴时间对应的历史数据分组,key=x轴,value=x轴对应的历史数据集合
	returnMap := service.GetCharData(xAxisList, ycQuery.StartTime, ycQuery.EndTime, ycQuery.Interval, ycQuery.IntervalType, ycList, ycQuery.CodeList)
	return returnMap, nil
}
