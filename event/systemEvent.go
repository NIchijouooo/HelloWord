package event

type SystemEventTemplate struct {
	ID    int    //
	Type  string //事件类型
	Value string //事件值
	Time  uint64 //事件时间戳
}

func (d *SystemEventTemplate) AddEvent() {

}

func (d *SystemEventTemplate) ModifyEvent() {

}

func (d *SystemEventTemplate) DeleteEvents() {

}

func (d *SystemEventTemplate) GetEvents() {

}
