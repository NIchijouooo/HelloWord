package ReturnModel

type CharData struct {
	XAxisList []string `json:"xAxisList"`
	//DataMap   map[string][]float64 `json:"dataMap"`
	DataList []ResYcData `json:"dataList"`
}
type ResYcData struct {
	Name string    `json:"name"`
	Data []float64 `json:"data"`
}
