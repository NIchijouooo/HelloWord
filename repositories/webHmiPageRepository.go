package repositories

import (
	"fmt"
	"gateway/models"
	"gorm.io/gorm"
	"time"
)

type WebHmiPageRepository struct {
	db *gorm.DB
}

func NewWebHmiPageRepository() *WebHmiPageRepository {
	return &WebHmiPageRepository{
		db: models.DB,
	}
}

// GetWebHmiPageDeviceInfo 获取配置
func (r *WebHmiPageRepository) GetWebHmiPageDeviceInfo(deviceId int) ([]models.WebHmiPageDeviceModel, error) {
	var webHmiPageDeviceModelList []models.WebHmiPageDeviceModel
	db := r.db
	fmt.Println("deviceId ： ", deviceId)
	if deviceId > 0 {
		db = db.Where("device_id = ?", deviceId)
	}
	if err := db.Find(&webHmiPageDeviceModelList).Error; err != nil {
		return nil, err
	}
	return webHmiPageDeviceModelList, nil
}

// SaveWebHmiPageDeviceInfo 保存设置
func (r *WebHmiPageRepository) SaveWebHmiPageDeviceInfo(webHmiPageDeviceModelList []*models.WebHmiPageDeviceModel) (int, error) {
	r.db.Where("id > ?", 0).Delete(&models.WebHmiPageDeviceModel{})
	if len(webHmiPageDeviceModelList) > 0 {
		createTime := time.Now().Format(time.DateTime)
		for _, webHmiPageDeviceModel := range webHmiPageDeviceModelList {
			webHmiPageDeviceModel.CreateTime = createTime
		}
		if err := r.db.Create(webHmiPageDeviceModelList).Error; err != nil {
			return 0, err
		}
	}
	return len(webHmiPageDeviceModelList), nil
}
