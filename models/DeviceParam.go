package models

// 设备参数相关
type CtrlInfo struct {
	DeviceId     int         `json:"deviceId"`
	Code         int         `json:"code"`
	CtrlUserName string      `json:"ctrlUserName"`
	ParamId      int         `json:"paramId"`
	Value        interface{} `json:"value"`
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

type DevicePageParam struct {
	KeyWord     string `json:"keyWord"`
	DeviceName  string `json:"deviceName"`
	DeviceLabel string `json:"deviceLabel"`
	PageNum     int    `form:"pageNum"`
	PageSize    int    `form:"pageSize"`
}

type EmDeviceParamVO struct {
	Id              int    `json:"id" gorm:"primary_key"`
	Name            string `json:"name"`
	Label           string `json:"label"`
	DeviceType      string `json:"deviceType"`
	DeviceTypeKey   string `json:"deviceTypeKey"`
	ModelId         int    `json:"modelId"`
	CollInterfaceId int    `json:"collInterfaceId"`
	Addr            string `json:"addr"`
	Data            string `json:"data"`
	ConnectStatus   string `json:"connectStatus"`
	Manufacturer    string `json:"manufacturer" description:"生产厂家"`
	Polarity        int    `json:"polarity" description:"极性"`
	FactoryModel    string `json:"factoryModel" description:"出厂型号"`
}
