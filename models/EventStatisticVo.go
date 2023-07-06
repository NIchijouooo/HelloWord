package models

type EventStatisticVo struct {
	Total              int                  `json:"total" description:"告警总数"`
	EventStatisticList []EventStatisticData `json:"eventStatisticList" description:"统计集合"`
}
