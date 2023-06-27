package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gateway/models"
	"gorm.io/gorm"
	"io"
	"net/http"
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
func (r *WebHmiPageRepository) GetIotWebHmiPageInfo(deviceId int) (int, string, error) {
	var deviceEquipmentAccountInfo models.DeviceEquipmentAccountInfo
	if err := r.db.Where("device_id = ?", deviceId).Find(&deviceEquipmentAccountInfo).Error; err != nil {
		return 0, "", err
	}
	webHmiPageCode := deviceEquipmentAccountInfo.WebHmiPageCode

	webHmiPageId, iotToken := GetWebHmiPageInfo(webHmiPageCode)

	return webHmiPageId, iotToken, nil
}

func GetWebHmiPageInfo(webHmiPageCode string) (id int, iotToken string) {
	// 准备POST请求的数据
	requestData := "{\"code\": \"" + webHmiPageCode + "\"}"

	// 发送POST请求
	req, err := http.NewRequest("POST", "https://interface.feisjy.com/qianhai/hmiPage/getHmiPageInfoByCode", bytes.NewBuffer([]byte(requestData)))
	if err != nil {
		panic(err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		panic("请求失败")
	}

	// 读取响应体的内容
	responseData, err := io.ReadAll(resp.Body)
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
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		panic(err)
	}

	if response.Code == 401 {
		fmt.Println("code is 401 ")
		isLogin := iotLogin()
		if !isLogin {
			return 0, ""
		}
		return GetWebHmiPageInfo(webHmiPageCode)
	} else if response.Code == 200 {
		return response.Data.Id, token
	}
	return 0, ""
}

func iotLogin() (isLogin bool) {
	// 准备POST请求的数据
	requestData := "{\"checkCode\": false,\"code\": false,\"userAccountNum\": \"admin\",\"password\": \"Feisjy@2016\",\"domain\": \"iot.feisjy.com\",\"isHmiLogin\": true,\"isEncryption\": \"1\",\"scenesId\": \"18\"}"
	// 发送POST请求
	resp, err := http.Post("https://interface.feisjy.com/auth/m2mLogin", "application/json", bytes.NewBuffer([]byte(requestData)))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		panic("请求失败")
	}

	// 读取响应体的内容
	responseData, err := io.ReadAll(resp.Body)
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
	err = json.Unmarshal(responseData, &response)
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
