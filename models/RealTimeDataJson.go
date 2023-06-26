package models

type RealTimeDataJson struct {
	Id       int    `json:"id"`
	Type     int    `json:"type"`
	Code     int    `json:"code"`
	Value    string `json:"value"`
	DeviceId int    `json:"deviceId"`
}
