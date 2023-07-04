package query

type PageBase struct {
	PageSize int `json:"pageSize"`
	PageNum  int `json:"pageNum"`
}
