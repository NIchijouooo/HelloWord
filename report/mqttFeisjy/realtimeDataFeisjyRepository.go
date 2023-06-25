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
	//查出所有相关点位
	var emDev models.EmDevice
	var str = r.db.Joins("inner JOIN em_coll_interface ON em_coll_interface.id = em_device.coll_interface_id").Where("em_device.name = ? and em_coll_interface.name = ?", devName, collName).First(&emDev).Statement.SQL.String()
	fmt.Println(str)
	//if err := r.db.Joins("inner JOIN em_coll_interface ON em_coll_interface.id = em_device.coll_interface_id").Where("em_device.name = ? and em_coll_interface.name = ?", devName, collName).First(&emDev).Error; err != nil {
	//	log.Printf("Request params:%v", err)
	//}
	var pointList []*models.EmDeviceModelCmdParam
	pointList = r.repoPoint.GetPointsByDeviceId("", emDev.Id, 0)

	var pointListYxList []*models.YxData
	var pointListYcList []*models.YcData
	var pointListSettingList []*models.SettingData

	// 转换为map
	dataMap := make(map[int]MQTTFeisjyReportDataTemplate)
	for _, ycParam := range ycPropertyPostParam.YcList {
		dataMap[ycParam.ID] = ycParam
	}

	//var wg sync.WaitGroup

	//20230615发给吴总的rule
	realTimeDataJsonList := make([]models.RealTimeDataJson, 0)

	//pointList为查出来的点位list
	for _, v := range pointList {
		//wg.Add(1)
		//ycPropertyPostParam.YcList为上报实时数据中的参数list
		//go func() {

		// 从map中取数据
		ycParam, exists := dataMap[v.Id]
		if exists {
			fmt.Println(ycParam)
		} else {
			fmt.Println("ycParam not found")
		}
		//for _, ycParam := range ycPropertyPostParam.YcList {
		//	setting.ZAPS.Infof("设备point:[%v]和[%v]和[%v]", v.Id, ycParam.ID, ycParam.Value, ycPropertyPostParam.Time)
			//if numName, err := strconv.Atoi(ycParam.ID); err == nil {
		//if v.Id == numName {
						t, _ := time.Parse("2006-01-02 15:04:05", ycPropertyPostParam.Time)
						var pType int

						var num = GetInterfaceToInt(ycParam.Value)

						if v.IotDataType == "yx" {
							pointListYxList = append(pointListYxList, &models.YxData{
								DeviceId: emDev.Id,
								Code:     ycParam.ID,
								Value:    num,
								Ts:       models.LocalTime{Time: t},
							})
							pType = 0
						} else if v.IotDataType == "yc" {
							pointListYcList = append(pointListYcList, &models.YcData{
								DeviceId: emDev.Id,
								Code:     ycParam.ID,
								Value:    float64(num),
								Ts:       models.LocalTime{Time: t},
							})
							pType = 1
						} else if v.IotDataType == "setting" {
							pointListSettingList = append(pointListSettingList, &models.SettingData{
								DeviceId: emDev.Id,
								Code:     ycParam.ID,
								Value:    strconv.Itoa(num),
								Ts:       models.LocalTime{Time: t},
							})
							pType = 2
						}
						realTimeDataJson := models.RealTimeDataJson{
							Type:     pType,   //遥信-0；遥测-1
							Code:     ycParam.ID, //遥信/遥测CODE
							Value:    strconv.Itoa(num),
							DeviceId: emDev.Id,
						}
						realTimeDataJsonList = append(realTimeDataJsonList, realTimeDataJson)
					//}
				//}
			//}
		//}()
	}
	//wg.Wait()
	/**
	转出yx。yc。setting分别存入taos
	*/
	go r.repo.BatchCreateYx(pointListYxList)
	go r.repo.BatchCreateYc(pointListYcList)
	go r.repo.BatchCreateSetting(pointListSettingList)
	//20230615发给吴总的rule
	rule.ProcessingSignalMsg(realTimeDataJsonList)
}

func GetInterfaceToInt(t1 interface{}) int {
	var t2 int
	switch t1.(type) {
	case uint:
		t2 = int(t1.(uint))
		break
	case int8:
		t2 = int(t1.(int8))
		break
	case uint8:
		t2 = int(t1.(uint8))
		break
	case int16:
		t2 = int(t1.(int16))
		break
	case uint16:
		t2 = int(t1.(uint16))
		break
	case int32:
		t2 = int(t1.(int32))
		break
	case uint32:
		t2 = int(t1.(uint32))
		break
	case int64:
		t2 = int(t1.(int64))
		break
	case uint64:
		t2 = int(t1.(uint64))
		break
	case float32:
		t2 = int(t1.(float32))
		break
	case float64:
		t2 = int(t1.(float64))
		break
	case string:
		t2, _ = strconv.Atoi(t1.(string))
		if t2 == 0 && len(t1.(string)) > 0 {
			f, _ := strconv.ParseFloat(t1.(string), 64)
			t2 = int(f)
		}
		break
	case nil:
		t2 = 0
		break
	default:
		t2 = t1.(int)
		break
	}
	return t2
}
