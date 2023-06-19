package models

import (
	"github.com/shopspring/decimal"
)
/**
	存taos
 */

type EsChargeDischargeModel struct {
	Ts              		int64    `json:"ts" gorm:"comment:'时间戳'"`
	ChargeCapacity          float64 `json:"chargeCapacity" gorm:"comment:'充电功率'"`
	DischargeCapacity       float64 `json:"dischargeCapacity" gorm:"comment:'放电功率'"`
	Profit         			decimal.Decimal    `json:"profit" gorm:"comment:'收益'"`
	DeviceId 				int    `json:"deviceId" gorm:"comment:'设备id'"`
}

func (u *EsChargeDischargeModel) TableName() string {
	return "charge_discharge"
}
