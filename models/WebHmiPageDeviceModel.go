package models

type WebHmiPageDeviceModel struct {
	Id             int    `json:"id" gorm:"primary_key"`
	DeviceId       int    `json:"deviceId"`
	WebHmiPageCode string `json:"webHmiPageCode"`
	CreateTime     string `json:"createTime"`
}

func (u *WebHmiPageDeviceModel) TableName() string {
	return "web_hmi_page_device"
}
