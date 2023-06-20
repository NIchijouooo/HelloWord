package models

// 定义字典数据表的模型
type SettingData struct {
	Code     int       `json:"code"`
	DeviceId int       `json:"deviceId"`
	Value    string    `json:"val"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Ts       LocalTime `json:"ts"`
}

func (u *SettingData) TableName() string {
	return "setting"
}
