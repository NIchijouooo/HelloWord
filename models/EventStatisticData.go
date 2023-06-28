package models

type EventStatisticData struct {
	Total int    `json:"total" description:"告警总数"`
	Code  string `json:"code" description:"编码"`
	Name  string `json:"name" description:"事件类型名称"`
}
