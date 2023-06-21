package models

// 设备参数相关
type CtrlInfo struct {
	DeviceId int         `json:"deviceId"`
	Code     int         `json:"code"`
	ParamId  int         `json:"paramId"`
	Value    interface{} `json:"value"`
}

type DeviceParam struct {
	DeviceType string `json:"deviceType"`
}

type GetDeviceModelCmdParam struct {
	DeviceName string `json:"deviceName"`
	CollName   string `json:"collName"`
	CmdName    string `json:"cmdName"`
	ParamName  string `json:"paramName"`
}
