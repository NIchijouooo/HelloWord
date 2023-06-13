package models

import "time"

// 定义字典数据表的模型
type SettingData struct {
	Code     int       `json:"code"`
	DeviceID int       `json:"deviceID"`
	Value    string    `json:"val"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Ts       time.Time `json:"ts"`
}

func (u *SettingData) TableName() string {
	return "setting"
}
