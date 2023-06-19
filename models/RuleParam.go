package models

// 设备参数相关
type RuleHistoryParam struct {
	DeviceId  int    `json:"deviceId"`
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
