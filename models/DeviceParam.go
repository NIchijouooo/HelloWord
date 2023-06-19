package models

// 设备参数相关
type CtrlInfo struct {
	DeviceId int         `json:"deviceId"`
	Code     int         `json:"code"`
	Value    interface{} `json:"value"`
}

type DeviceParam struct {
	DeviceType string `json:"deviceType"`
}
