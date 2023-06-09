package models

type EmCollInterface struct {
	Id              int    `json:"id" gorm:"primary_key"`
	Name            string `json:"name"`
	CommInterfaceId int    `json:"commInterfaceId"`
	OfflinePeriod   int    `json:"offlinePeriod"`
	PollPeriod      int    `json:"pollPeriod"`
	Data            string `json:"data"`
}

type AddEmCollInterface struct {
	CollInterfaceName string `json:"collInterfaceName"`
	CommInterfaceName string `json:"commInterfaceName"`
	PollPeriod        int    `json:"pollPeriod"`
	OfflinePeriod     int    `json:"offlinePeriod"`
}

func (u *EmCollInterface) TableName() string {
	return "em_coll_interface"
}
