package ReturnModel

import "time"

type CharData struct {
	XAxisList []string `json:"xAxisList"`
	//DataMap   map[string][]float64 `json:"dataMap"`
	DataList []ResYcData `json:"dataList"`
}
type ResYcData struct {
	Name string   `json:"name"`
	Data []string `json:"data"`
}
type YcData struct {
	Code     int       `json:"code"`
	DeviceId int       `json:"deviceId"`
	Value    float64   `json:"val"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Ts       time.Time `json:"ts"`
	Alias    string    `json:"alias"`
	Sort     int       `json:"sort"`
}
