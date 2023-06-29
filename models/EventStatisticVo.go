package models

type EventStatisticVo struct {
	Total                   int                  `json:"total" description:"告警总数"`
	EventLevelStatisticList []EventStatisticData `json:"eventLevelStatisticList" description:"按事件等级分类统计集合"`
}
