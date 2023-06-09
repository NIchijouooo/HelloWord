package models

type EmDeviceModel struct {
	Id    int    `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	Label string `json:"label"`
	Type  int    `json:"type"`
	Data  string `json:"data"`
}

type AddEmDeviceModel struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Param string `json:"param"`
	Type  int    `json:"type"`
}

func (u *EmDeviceModel) TableName() string {
	return "em_device_model"
}
