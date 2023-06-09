package models

type CommInterface struct {
	Id                      int    `json:"id" gorm:"primary_key"`
	CommInterfaceProtocolId int    `json:"commInterfaceProtocolId"`
	Name                    string `json:"name"`
	CommInterfaceProtocol   string `json:"commInterfaceProtocol"`
	Data                    string `json:"data"`
}

func (u *CommInterface) TableName() string {
	return "em_comm_interface"
}
