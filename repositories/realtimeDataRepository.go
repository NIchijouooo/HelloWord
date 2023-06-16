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
func (r *RealtimeDataRepository) BatchCreateYc(realtimeList []*models.YcData) error {
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

/*
*
	根据网关设备的addr查设备，然后更新登录状态
*/
//func (r *RealtimeDataRepository) UpdateGatewayDeviceConnetStatus(gw mqttFeisjy.ReportServiceGWParamFeisjyTemplate) {
//	//var emDev models.EmDevice
//	if err := r.db.Model(&models.ProjectInfo{}).Where("addr = ?", gw.Param.DeviceID).Updates(models.EmDevice{Connectstatus: gw.ReportStatus}).Error; err != nil {
//		log.Printf("Request params:%v", err)
//	}
//}
//
///*
//*
//	根据设备的coll,name查设备，然后更新登录状态
//*/
//func (r *RealtimeDataRepository) UpdateDeviceConnetStatus(node mqttFeisjy.ReportServiceNodeParamFeisjyTemplate) {
//	//var emDev models.EmDevice
//	if err := r.db.Model(&models.ProjectInfo{}).Where("name = ? and coll_interface_id = ?", node.Name, node.CollInterfaceName).Updates(models.EmDevice{Connectstatus: node.ReportStatus}).Error; err != nil {
//		log.Printf("Request params:%v", err)
//	}
//}
//
///*
//*
//
//	更新命令参数属性的实时值到taos
//*/
//func (r *RealtimeDataRepository) SaveRealtimeDataList(devName, collName string, ycPropertyPostParam mqttFeisjy.MQTTFeisjyReportYcTemplate) {
//	var emDev models.EmDevice
//	var str = r.db.Joins("inner JOIN em_coll_interface ON em_coll_interface.id = em_device.coll_interface_id").Where("em_device.name = ? and em_coll_interface.name = ?", devName, collName).First(&emDev).Statement.SQL.String()
//	fmt.Println(str)
//	//if err := r.db.Joins("inner JOIN em_coll_interface ON em_coll_interface.id = em_device.coll_interface_id").Where("em_device.name = ? and em_coll_interface.name = ?", devName, collName).First(&emDev).Error; err != nil {
//	//	log.Printf("Request params:%v", err)
//	//}
//	var pointList []*models.EmDeviceModelCmdParam
//	repo := &devicePointRepository.DevicePointRepository{}
//	pointList = repo.GetPointsByDeviceId("all", emDev.Id, 0)
//
//	var pointListYxList []*models.YxData
//	var pointListYcList []*models.YcData
//	var pointListSettingList []*models.SettingData
//
//	var wg sync.WaitGroup
//	//pointList为查出来的点位list
//	for _, v := range pointList {
//		wg.Add(1)
//		//ycPropertyPostParam.YcList为上报实时数据中的参数list
//		go func() {
//			for _, ycParam := range ycPropertyPostParam.YcList {
//				if numName, err := strconv.Atoi(ycParam.Name); err == nil {
//					if v.Id == numName {
//						t, _ := time.Parse("2021-09-15 14:30:00", ycPropertyPostParam.Time)
//						if v.IotDataType == "yx" {
//							pointListYxList = append(pointListYxList, &models.YxData{
//								DeviceId: emDev.Id,
//								Code:     numName,
//								Value:    ycParam.Value.(int),
//								Ts:       t,
//							})
//						} else if v.IotDataType == "yc" {
//							pointListYcList = append(pointListYcList, &models.YcData{
//								DeviceId: emDev.Id,
//								Code:     numName,
//								Value:    ycParam.Value.(float64),
//								Ts:       t,
//							})
//						} else if v.IotDataType == "setting" {
//							pointListSettingList = append(pointListSettingList, &models.SettingData{
//								DeviceId: emDev.Id,
//								Code:     numName,
//								Value:    ycParam.Value.(string),
//								Ts:       t,
//							})
//						}
//					}
//				}
//			}
//		}()
//	}
//	wg.Wait()
//	/**
//	转出yx。yc。setting分别存入taos
//	*/
//	go r.BatchCreateYx(pointListYxList)
//	go r.BatchCreateYc(pointListYcList)
//	go r.BatchCreateSetting(pointListSettingList)
//}
