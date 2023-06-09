package models

type EmDeviceModelCmd struct {
	Id            int    `json:"id" gorm:"primary_key"`
	DeviceModelId int    `json:"deviceModelId"`
	Name          string `json:"name"`
	Label         string `json:"label"`
	Data          string `json:"data"`
}

type AddEmDeviceModelCmd struct {
	TslName      string `json:"tslName"`
	Name         string `json:"name"`
	Label        string `json:"label"`
	FunCode      int    `json:"funCode"`
	StartRegAddr int    `json:"startRegAddr"`
	RegCnt       int    `json:"regCnt"`
}

func (u *EmDeviceModelCmd) TableName() string {
	return "em_device_model_cmd"
}
