package repositories

import (
	"encoding/json"
	"fmt"
	"gateway/models"
	"gateway/utils"
	"gorm.io/gorm"
)

var token string

type WebHmiPageRepository struct {
	db *gorm.DB
}

func NewWebHmiPageRepository() *WebHmiPageRepository {
	return &WebHmiPageRepository{
		db: models.DB,
	}
}

// GetIotWebHmiPageInfo 获取IOT组态信息
func (r *WebHmiPageRepository) GetIotWebHmiPageInfo(webHmiPageDeviceModel models.WebHmiPageDeviceModel) (int, string, error) {

	webHmiPageCode := webHmiPageDeviceModel.WebHmiPageCode
	if webHmiPageDeviceModel.WebHmiPageCode == "" {
		var deviceEquipmentAccountInfo models.DeviceEquipmentAccountInfo
		if err := r.db.Where("device_id = ?", webHmiPageDeviceModel.DeviceId).Find(&deviceEquipmentAccountInfo).Error; err != nil {
			return 0, "", err
		}
		webHmiPageCode = deviceEquipmentAccountInfo.WebHmiPageCode
	}

	webHmiPageId, iotToken := getWebHmiPageInfo(webHmiPageCode)

	return webHmiPageId, iotToken, nil
}

func getWebHmiPageInfo(webHmiPageCode string) (id int, iotToken string) {
	// 准备POST请求的数据
	requestData := "{\"code\": \"" + webHmiPageCode + "\"}"
	url := "https://interface.feisjy.com/qianhai/hmiPage/getHmiPageInfoByCode"
	header := make(map[string]string)
	header["Authorization"] = token
	result, err := utils.SendPostTls(url, requestData, header)
	if err != nil {
		panic(err)
	}

	type dataJson struct {
		Id int `json:"id"`
	}

	type responseDataJson struct {
		Code int      `json:"code"`
		Data dataJson `json:"data"`
		Msg  string   `json:"msg"`
	}

	// 将响应体内容转换为JSON
	var response responseDataJson
	err = json.Unmarshal([]byte(result), &response)
	if err != nil {
		panic(err)
	}

	if response.Code == 401 {
		fmt.Println("code is 401 ")
		isLogin := iotLogin()
		if !isLogin {
			return 0, ""
		}
		return getWebHmiPageInfo(webHmiPageCode)
	} else if response.Code == 200 {
		return response.Data.Id, token
	}
	return 0, ""
}

func iotLogin() (isLogin bool) {

	// 准备POST请求的数据
	requestData := "{\"checkCode\": false,\"code\": false,\"userAccountNum\": \"admin\",\"password\": \"Feisjy@2016\",\"domain\": \"iot.feisjy.com\",\"isHmiLogin\": true,\"isEncryption\": \"1\",\"scenesId\": \"18\"}"
	url := "https://interface.feisjy.com/auth/m2mLogin"
	result, err := utils.SendPostTls(url, requestData, nil)
	if err != nil {
		panic(err)
	}

	type dataJson struct {
		AccessToken string `json:"access_token"`
	}

	type responseDataJson struct {
		Code int      `json:"code"`
		Data dataJson `json:"data"`
		Msg  string   `json:"msg"`
	}

	// 将响应体内容转换为JSON
	var response responseDataJson
	err = json.Unmarshal([]byte(result), &response)
	if err != nil {
		panic(err)
	}

	if response.Code == 200 {
		// 处理响应
		token = response.Data.AccessToken
		return true
	}

	return false
}
