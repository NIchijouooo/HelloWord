package event

import (
	"errors"
	"time"
)

type ReportEventTemplate struct {
	ID          uint32 //
	ServiceName string //上报服务名称
	Value       string //事件值
	Time        int64  //事件时间戳
}

var reportEventID uint32 = 1
var ReportEvents = make(map[uint32]ReportEventTemplate)

func AddReportEvent(serviceName string, value string) {
	event := ReportEventTemplate{
		ID:          reportEventID,
		ServiceName: serviceName,
		Value:       value,
		Time:        time.Now().Unix(),
	}
	ReportEvents[reportEventID] = event
	reportEventID += 1
}

func ModifyReportEvent() {

}

func DeleteReportEvents(id uint32) error {
	_, ok := ReportEvents[id]
	if !ok {
		return errors.New("ReportEventID不存在")
	}
	delete(ReportEvents, id)

	return nil
}

func GetReportEvents() []ReportEventTemplate {
	events := make([]ReportEventTemplate, 0)
	for _, v := range ReportEvents {
		events = append(events, v)
	}

	return events
}
