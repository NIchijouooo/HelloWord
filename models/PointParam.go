package models

type PointParam struct {
	DeviceId int       `json:"deviceId"    description:"设备id"`
	Code     int       `json:"code"    description:"遥信/遥测/参数编码"`
	Name     string    `json:"name"    description:"遥信/遥测/参数名称"`
	Label    string    `json:"label"	description:"别名"`
	Value    *string   `json:"val"    description:"遥信/遥测/参数实时值，遥信遥测为数字类型，参数为字符串类型"`
	Ts       LocalTime `json:"ts"`
}
