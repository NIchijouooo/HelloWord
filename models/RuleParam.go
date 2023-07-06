package models

// 设备参数相关
type RuleHistoryParam struct {
	DeviceIds      []int    `json:"deviceIds"`
	Codes          []int    `json:"codes"`
	Level          string   `json:"level"`
	Tag            string   `json:"tag"`
	PageNum        int64    `json:"pageNum"`
	PageSize       int64    `json:"pageSize"`
	StartTime      string   `json:"startTime"`
	EndTime        string   `json:"endTime"`
	Description    string   `json:"description"`
	DeviceName     string   `json:"deviceName"`
	DeviceTypeList []string `json:"deviceTypeList"`
	StatisticType  int      `json:"statisticType"` // 统计类型;1-统计告警等级,2-统计告警状态
}
