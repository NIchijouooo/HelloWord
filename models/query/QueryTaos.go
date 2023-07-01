package query

// 定义涛思查询模型
type QueryTaoData struct {
	DeviceId     int      `form:"deviceId"`
	DeviceIds    []int    `form:"deviceIds"`
	CodeList     []int    `form:"codeList"`
	CodeNameList []string `form:"codeNameList"`
	Codes        string   `form:"codes"`
	StartTime    int64    `form:"startTime"`
	EndTime      int64    `form:"endTime"`
	Interval     int      `form:"interval"`     //间隔时间
	IntervalType int      `form:"intervalType"` //间隔类型
}

type DeviceIdAndCode struct {
	DeviceId int    `json:"deviceId"`
	Code     string `json:"code"`
}
