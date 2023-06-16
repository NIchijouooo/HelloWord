package models

type EmDeviceModelCmdParam struct {
	Id               int    `json:"id" gorm:"primary_key"`
	DeviceModelCmdId int    `json:"deviceModelCmdId"`
	Name             string `json:"name"`
	Label            string `json:"label"`
	Data             string `json:"data"`
	IotDataType      string `json:"iotDataType"`
}

type AddEmDeviceModelCmdParam struct {
	TslName     string `json:"tslName"`
	CmdName     string `json:"cmdName"`
	Name        string `json:"name"`
	Label       string `json:"label"`
	AccessMode  int    `json:"accessMode"`
	Type        int    `json:"type"`
	Decimals    int    `json:"decimals"`
	Unit        string `json:"unit"`
	RegCnt      int    `json:"regCnt"`
	RegAddr     int    `json:"regAddr"`
	RuleType    string `json:"ruleType"`
	Formula     string `json:"formula"`
	IotDataType string `json:"IotDataType"`
}

func (u *EmDeviceModelCmdParam) TableName() string {
	return "em_device_model_cmd_param"
}
