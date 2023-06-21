package mqttFeisjy

import (
	"database/sql"
	"fmt"
	"gateway/models"
	"gateway/repositories"
	"gateway/rule"
	"gateway/setting"
	"log"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"
)

// 定义字典类型管理的存储库
type RealtimeDataRepository struct {
	db        *gorm.DB
	taosDb    *sql.DB
	repo      *repositories.RealtimeDataRepository
	repoPoint *repositories.DevicePointRepository
}

func NewRealtimeDataRepository() *RealtimeDataRepository {
	return &RealtimeDataRepository{db: models.DB,
		taosDb:    models.TaosDB,
		repo:      repositories.NewRealtimeDataRepository(),
		repoPoint: repositories.NewDevicePointRepository(),
	}
}

/*
*
	根据网关设备的addr查设备，然后更新登录状态
*/
func (r *RealtimeDataRepository) UpdateGatewayDeviceConnetStatus(gw ReportServiceGWParamFeisjyTemplate) {
	//var emDev models.EmDevice
	setting.ZAPS.Debugf("UpdateGatewayDeviceConnetStatus start......")
	if err := r.db.Model(&models.ProjectInfo{}).Where("addr = ?", gw.Param.DeviceID).Updates(models.EmDevice{ConnectStatus: gw.ReportStatus}).Error; err != nil {
		log.Printf("Request params:%v", err)
	}
}

/*
*
	根据设备的coll,name查设备，然后更新登录状态
*/
func (r *RealtimeDataRepository) UpdateDeviceConnetStatus(node ReportServiceNodeParamFeisjyTemplate) {
	//var emDev models.EmDevice
	setting.ZAPS.Debugf("UpdateGatewayDeviceConnetStatus start......")
	if err := r.db.Model(&models.ProjectInfo{}).Where("name = ? and coll_interface_id = ?", node.Name, node.CollInterfaceName).Updates(models.EmDevice{ConnectStatus: node.ReportStatus}).Error; err != nil {
		log.Printf("Request params:%v", err)
	}
}

/*
*
	更新命令参数属性的实时值到taos
*/
func (r *RealtimeDataRepository) SaveRealtimeDataList(devName, collName string, ycPropertyPostParam MQTTFeisjyReportYcTemplate) {
	var emDev models.EmDevice
	var str = r.db.Joins("inner JOIN em_coll_interface ON em_coll_interface.id = em_device.coll_interface_id").Where("em_device.name = ? and em_coll_interface.name = ?", devName, collName).First(&emDev).Statement.SQL.String()
	fmt.Println(str)
	//if err := r.db.Joins("inner JOIN em_coll_interface ON em_coll_interface.id = em_device.coll_interface_id").Where("em_device.name = ? and em_coll_interface.name = ?", devName, collName).First(&emDev).Error; err != nil {
	//	log.Printf("Request params:%v", err)
	//}
	var pointList []*models.EmDeviceModelCmdParam
	pointList = r.repoPoint.GetPointsByDeviceId("all", emDev.Id, 0)

	var pointListYxList []*models.YxData
	var pointListYcList []*models.YcData
	var pointListSettingList []*models.SettingData

	var wg sync.WaitGroup

	//20230615发给吴总的rule
	realTimeDataJsonList := make([]models.RealTimeDataJson, 0)

	//pointList为查出来的点位list
	for _, v := range pointList {
		wg.Add(1)
		//ycPropertyPostParam.YcList为上报实时数据中的参数list
		go func() {
			for _, ycParam := range ycPropertyPostParam.YcList {
				if numName, err := strconv.Atoi(ycParam.Name); err == nil {
					if v.Id == numName {
						t, _ := time.Parse("2021-09-15 14:30:00", ycPropertyPostParam.Time)
						var pType int
						if v.IotDataType == "yx" {
							pointListYxList = append(pointListYxList, &models.YxData{
								DeviceId: emDev.Id,
								Code:     numName,
								Value:    ycParam.Value.(int),
								Ts:       models.LocalTime{Time: t},
							})
							pType = 0
						} else if v.IotDataType == "yc" {
							pointListYcList = append(pointListYcList, &models.YcData{
								DeviceId: emDev.Id,
								Code:     numName,
								Value:    ycParam.Value.(float64),
								Ts:       models.LocalTime{Time: t},
							})
							pType = 1
						} else if v.IotDataType == "setting" {
							pointListSettingList = append(pointListSettingList, &models.SettingData{
								DeviceId: emDev.Id,
								Code:     numName,
								Value:    ycParam.Value.(string),
								Ts:       models.LocalTime{Time: t},
							})
							pType = 2
						}
						realTimeDataJson := models.RealTimeDataJson{
							Type:     pType,   //遥信-0；遥测-1
							Code:     numName, //遥信/遥测CODE
							Value:    ycParam.Value.(string),
							DeviceId: emDev.Id,
						}
						realTimeDataJsonList = append(realTimeDataJsonList, realTimeDataJson)
					}
				}
			}
		}()
	}
	wg.Wait()
	/**
	转出yx。yc。setting分别存入taos
	*/
	go r.repo.BatchCreateYx(pointListYxList)
	go r.repo.BatchCreateYc(pointListYcList)
	go r.repo.BatchCreateSetting(pointListSettingList)

	rule.ProcessingSignalMsg(realTimeDataJsonList)
}
