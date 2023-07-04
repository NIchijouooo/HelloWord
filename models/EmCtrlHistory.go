package models

type EmCtrlHistory struct {
	Id           int    `json:"id" gorm:"primary_key"`
	DeviceId     int    `json:"deviceId"`
	ParamId      int    `json:"paramId"`
	Value        string `json:"value"`
	CtrlUserName string `json:"ctrlUserName"`
	CtrlStatus   int    `json:"ctrlStatus"`
	CreateTime   string `json:"createTime"`
	UpdateTime   string `json:"updateTime"`
}

type AddEmCtrlHistory struct {
	DeviceId     int    `json:"deviceId"`
	ParamId      int    `json:"paramId"`
	CtrlUserName string `json:"ctrlUserName"`
	CtrlStatus   int    `json:"ctrlStatus"`
	CreateTime   string `json:"createTime"`
	UpdateTime   string `json:"updateTime"`
}

func (u *EmCtrlHistory) TableName() string {
	return "em_ctrl_history"
}
