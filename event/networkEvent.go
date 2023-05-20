package event

type NetworkEventTemplate struct {
	ID    int    //
	Type  string //事件类型
	Value string //事件值
	Time  uint64 //事件时间戳
}

func (d *NetworkEventTemplate) AddEvent() {

}

func (d *NetworkEventTemplate) ModifyEvent() {

}

func (d *NetworkEventTemplate) DeleteEvents() {

}

func (d *NetworkEventTemplate) GetEvents() {

}
