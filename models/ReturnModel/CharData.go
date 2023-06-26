package ReturnModel

import (
	"gateway/models"
)

type CharData struct {
	XAxisList []string              `json:"xAxisList"`
	DataMap   map[int][]interface{} `json:"dataMap"`
	//DataList []ResYcData `json:"dataList"`
}

//	type ResYcData struct {
//		Name string   `json:"name"`
//		Data []string `json:"data"`
//	}
type YcData struct {
	Code       int              `json:"code"`       //测点编码
	DeviceId   int              `json:"deviceId"`   //设备id
	Value      float64          `json:"val"`        //值
	Name       string           `json:"name"`       //名称
	Type       string           `json:"type"`       //类型
	Ts         models.LocalTime `json:"ts"`         //时间
	Alias      string           `json:"alias"`      //别名
	Sort       int              `json:"sort"`       //排序
	Unit       string           `json:"unit"`       //单位
	Level      int              `json:"level"`      //告警状态
	LimitScope LimitScope       `json:"limitScope"` //告警范围
}
type AuxYcData struct {
	Code     int              `json:"code"`     //测点编码
	DeviceId int              `json:"deviceId"` //设备id
	Value    float64          `json:"val"`      //值
	Name     string           `json:"name"`     //名称
	Type     string           `json:"type"`     //类型
	Ts       models.LocalTime `json:"ts"`       //时间
	Unit     string           `json:"unit"`     //单位
}
type AuxYxData struct {
	Code     int              `json:"code"`     //测点编码
	DeviceId int              `json:"deviceId"` //设备id
	Value    int              `json:"val"`      //值
	Name     string           `json:"name"`     //名称
	Type     string           `json:"type"`     //类型
	Ts       models.LocalTime `json:"ts"`       //时间
	Unit     string           `json:"unit"`     //单位
}

type LimitScope struct {
	DeviceType   string `json:"deviceType"`   //设备类型
	PropertyCode string `json:"propertyCode"` //测点编码
	NotifyMin    string `json:"notifyMin"`
	NotifyMax    string `json:"notifyMax"`
	SecondaryMin string `json:"secondaryMin"`
	SecondaryMax string `json:"secondaryMax"`
	SeriousMin   string `json:"seriousMin"`
	SeriousMax   string `json:"seriousMax"`
	UrgentMin    string `json:"urgentMin"`
	UrgentMax    string `json:"urgentMax"`
}
