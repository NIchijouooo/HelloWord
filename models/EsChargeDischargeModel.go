package models

import (
	"github.com/shopspring/decimal"
)

/**
存taos
*/

type EsChargeDischargeModel struct {
	Ts                      int64           `json:"ts" gorm:"comment:'时间戳'"`
	ChargeCapacity          float64         `json:"chargeCapacity" gorm:"comment:'充电量'"`
	DischargeCapacity       float64         `json:"dischargeCapacity" gorm:"comment:'放电量'"`
	Profit                  decimal.Decimal `json:"profit" gorm:"comment:'收益'"`
	TopChargeCapacity       float64         `json:"topChargeCapacity" gorm:"comment:'尖充电量'"`
	TopDischargeCapacity    float64         `json:"topDischargeCapacity" gorm:"comment:'尖放电量'"`
	TopProfit               decimal.Decimal `json:"topProfit" gorm:"comment:'尖收益'"`
	PeakChargeCapacity      float64         `json:"peakChargeCapacity" gorm:"comment:'峰充电量'"`
	PeakDischargeCapacity   float64         `json:"peakDischargeCapacity" gorm:"comment:'峰放电量'"`
	PeakProfit              decimal.Decimal `json:"peakProfit" gorm:"comment:'峰收益'"`
	FlatChargeCapacity      float64         `json:"flatChargeCapacity" gorm:"comment:'平充电量'"`
	FlatDischargeCapacity   float64         `json:"flatDischargeCapacity" gorm:"comment:'平放电量'"`
	FlatProfit              decimal.Decimal `json:"flatProfit" gorm:"comment:'平收益'"`
	ValleyChargeCapacity    float64         `json:"valleyChargeCapacity" gorm:"comment:'谷充电量'"`
	ValleyDischargeCapacity float64         `json:"valleyDischargeCapacity" gorm:"comment:'谷放电量'"`
	ValleyProfit            decimal.Decimal `json:"valleyProfit" gorm:"comment:'谷收益'"`
	DeviceId                int             `json:"deviceId" gorm:"comment:'设备id'"`
}

func (u *EsChargeDischargeModel) TableName() string {
	return "charge_discharge"
}
