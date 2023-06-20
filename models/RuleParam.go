package models

// 设备参数相关
type RuleHistoryParam struct {
	DeviceIds string `json:"deviceIds"`
	Codes     string `json:"codes"`
	Level     string `json:"level"`
	Tag       string `json:"tag"`
	PageNum   int64  `json:"pageNum"`
	PageSize  int64  `json:"pageSize"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
