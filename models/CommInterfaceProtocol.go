package models

type CommInterfaceProtocol struct {
	Id    int    `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	Label string `json:"label"`
}

func (u *CommInterfaceProtocol) TableName() string {
	return "comm_interface_protocol"
}
