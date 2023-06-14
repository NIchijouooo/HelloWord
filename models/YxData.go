package models

import "time"

// 定义字典数据表的模型
type YxData struct {
	Code     int       `json:"code"`
	DeviceId int       `json:"deviceId"`
	Value    int       `json:"val"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Ts       time.Time `json:"ts"`
}

func (u *YxData) TableName() string {
	return "yx"
}
