package models

type EmDevice struct {
	Id              int    `json:"id" gorm:"primary_key"`
	Name            string `json:"name"`
	Label           string `json:"label"`
	ModelId         int    `json:"modelId"`
	CollInterfaceId int    `json:"collInterfaceId"`
	Addr            string `json:"addr"`
	Data            string `json:"data"`
}

type AddEmDevice struct {
	InterfaceName string `json:"interfaceName"`
	Name          string `json:"name"`
	Label         string `json:"label"`
	Addr          string `json:"addr"`
	Tsl           string `json:"tsl"`
}

func (u *EmDevice) TableName() string {
	return "em_device"
}
