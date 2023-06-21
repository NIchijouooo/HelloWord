package repositories

import (
	"encoding/json"
	"fmt"
	"gateway/models"

	"gorm.io/gorm"
)

type CentralizedRepository struct {
	db *gorm.DB
}

func NewCentralizedRepository() *CentralizedRepository {
	return &CentralizedRepository{db: models.DB}
}

// 获取策略数据
func (r *CentralizedRepository) GetPolicyList() ([]models.EmStrategy, error) {
	var policyList []models.EmStrategy
	err := r.db.Find(&policyList).Error
	return policyList, err
}

// 新增策略
func (r *CentralizedRepository) CreatePolicy(policyData *models.EmStrategy) error {
	return r.db.Create(policyData).Error
}

// 修改策略
func (r *CentralizedRepository) UpdatePolicy(policyData *models.EmStrategy) error {
	var emPolicyData models.EmStrategy
	id := policyData.Id
	r.db.Where("id = ?", id).First(&emPolicyData)
	policyData.Id = emPolicyData.Id
	return r.db.Save(policyData).Error
}

// 删除策略
func (r *CentralizedRepository) DeletePolicy(Id int) error {
	return r.db.Delete(&models.EmStrategy{}, Id).Error
}

// 删除策略
func (r *CentralizedRepository) DeletePolicy2(Id int) error {
	return r.db.Delete(&models.EmStrategy{}, Id).Error
}

type YkYcList struct {
	Id          int    `json:"id" gorm:"primary_key"`
	DeviceName  string `json:"deviceName"`
	ParamName   string `json:"paramName"`
	IotDataType string `json:"iotDataType"`
	Data        string `json:"data"`
	Unit        string `json:"unit"`
}
type YkYcData struct {
	Yx []YkYcList `json:"yx"`
	Yc []YkYcList `json:"yc"`
}

// 查询设备遥控遥调点位数据
func (r *CentralizedRepository) GetDeviceYkYtList() (YkYcData, error) {
	var ykYtList []YkYcList
	err := r.db.Table("em_device_model_cmd_param as param").
		Select("param.id as id, param.name as param_name, param.data, param.iot_data_type, device.name as device_name").
		Joins("LEFT JOIN em_device_model_cmd as cmd on param.device_model_cmd_id = cmd.id").
		Joins("LEFT JOIN em_device as device on cmd.device_model_id = device.model_id").
		Where("(iot_data_type = 'yc' or iot_data_type = 'yx') and device_name not NULL").
		Find(&ykYtList).
		Error

	var result YkYcData

	// 设备模型的读写属性 0读 1写 2读写
	type DeviceModel struct {
		AccessMode int `json:"accessMode"`
	}

	for _, item := range ykYtList {
		var deviceModel DeviceModel

		err := json.Unmarshal([]byte(item.Data), &deviceModel)
		if err != nil {
			fmt.Println("解析 JSON 失败:", err)
			return result, err
		}

		if item.IotDataType == "yx" && deviceModel.AccessMode != 0 {
			result.Yx = append(result.Yx, item)
		} else if item.IotDataType == "yc" && deviceModel.AccessMode != 0 {
			result.Yc = append(result.Yc, item)
		}
	}

	return result, err
}
