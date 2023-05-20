package event

import (
	"errors"
	"time"
)

type DeviceEventTemplate struct {
	ID         uint32
	Type       string //事件类型
	DeviceName string //设备名称
	Value      string //事件值
	Time       int64  //事件时间戳
}

var deviceEventID uint32 = 1
var DeviceEvents = make(map[uint32]DeviceEventTemplate)

func AddDeviceEvent(deviceName string, value string) {

	event := DeviceEventTemplate{
		ID:         deviceEventID,
		DeviceName: deviceName,
		Value:      value,
		Time:       time.Now().Unix(),
	}
	DeviceEvents[deviceEventID] = event
	deviceEventID += 1
}

func ModifyEvent() {

}

func DeleteDeviceEvents(id uint32) error {

	_, ok := DeviceEvents[id]
	if !ok {
		return errors.New("DeviceEventID不存在")
	}
	delete(DeviceEvents, id)

	return nil
}

func GetDeviceEvents() []DeviceEventTemplate {

	events := make([]DeviceEventTemplate, 0)
	for _, v := range DeviceEvents {
		events = append(events, v)
	}

	return events
}
