package ReturnModel

import "gateway/models"

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
	Code     int              `json:"code"`     //测点编码
	DeviceId int              `json:"deviceId"` //设备id
	Value    float64          `json:"val"`      //值
	Name     string           `json:"name"`     //名称
	Type     string           `json:"type"`     //类型
	Ts       models.LocalTime `json:"ts"`       //时间
	Alias    string           `json:"alias"`    //别名
	Sort     int              `json:"sort"`     //排序
	Uint     string           `json:"uint"`     //单位
	Status   int              `json:"status"`   //状态
}
