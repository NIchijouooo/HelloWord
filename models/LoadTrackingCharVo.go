package models

type LoadTrackingCharVo struct {
	XAxisList    []string      `json:"xAxisList"`
	OnGridData   []interface{} `json:"onGridData"`
	MeterData    []interface{} `json:"meterData"`
	LoadData     []interface{} `json:"loadData"`
	TopPeriod    string        `json:"topPeriod"`
	PeakPeriod   string        `json:"peakPeriod"`
	FlatPeriod   string        `json:"flatPeriod"`
	ValleyPeriod string        `json:"valleyPeriod"`
}
