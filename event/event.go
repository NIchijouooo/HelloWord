package event

type EventTemplate interface {
	AddEvent()
	ModifyEvent()
	DeleteEvents()
	GetEvents()
}

var EventMap = make([]EventTemplate, 0)
