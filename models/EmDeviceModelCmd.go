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

type AddEmDevicePlcModelCmd struct {
	Name     string `json:"name"`
	Property struct {
		Name       string `json:"name"`
		Label      string `json:"label"`
		AccessMode int    `json:"accessMode"`
		Type       int    `json:"type"`
		Decimals   int    `json:"decimals"`
		Unit       string `json:"unit"`
		Params     struct {
			DbNumber    string `json:"dbNumber"`
			StartAddr   string `json:"startAddr"`
			DataType    int    `json:"dataType"`
			IotDataType string `json:"iotDataType"`
		} `json:"params"`
	} `json:"property"`
}

func (u *EmDeviceModelCmd) TableName() string {
	return "em_device_model_cmd"
}
