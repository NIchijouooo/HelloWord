package models

import (
	"github.com/shopspring/decimal"
)

type EmConfiguration struct {
	Id           int             `json:"id"`
	Province     string          `json:"province"`
	Month        string          `json:"month"`
	TopPrice     decimal.Decimal `json:"topPrice"`
	TopPeriod    string          `json:"topPeriod"`
	PeakPrice    decimal.Decimal `json:"peakPrice"`
	PeakPeriod   string          `json:"peakPeriod"`
	FlatPrice    decimal.Decimal `json:"flatPrice"`
	FlatPeriod   string          `json:"flatPeriod"`
	ValleyPrice  decimal.Decimal `json:"valleyPrice"`
	ValleyPeriod string          `json:"valleyPeriod"`
}

func (u *EmConfiguration) TableName() string {
	return "em_configuration"
}
