package models

type ProjectInfo struct {
	ID             int    `gorm:"primary_key;auto_increment;comment:'主键'" json:"id"`       // 主键
	Description    string `gorm:"type:varchar(255);comment:'描述'" json:"description"`       // 描述
	Avatar         string `gorm:"type:varchar(255);comment:'头像'" json:"avatar"`            // 头像
	CommissionDate string `gorm:"type:date;comment:'投产时间'" json:"commissionDate"`          // 投产时间
	Address        string `gorm:"type:varchar(255);comment:'地址'" json:"address"`           // 地址
	RatedPower     string `gorm:"type:decimal(10,2);comment:'额定功率'" json:"ratedPower"`     // 额定功率
	InstalledPower string `gorm:"type:decimal(10,2);comment:'装机容量'" json:"installedPower"` // 装机容量
	PCSCount       int    `gorm:"type:varchar(255);comment:'功率转换系统'" json:"pcsCount"`      // 功率转换系统
	BatteryCount   int    `gorm:"type:int;comment:'电池簇数量'" json:"batteryCount"`            // 电池簇数量
	GridVoltage    string `gorm:"type:decimal(10,2);comment:'并网电压等级'" json:"gridVoltage"`  // 并网电压等级
	GridDate       string `gorm:"type:date;comment:'并网日期'" json:"gridDate"`                // 并网日期
	Status         string `gorm:"type:char(1);comment:'状态（0正常 1停用）'" json:"status"`        // 状态
}

func (u *ProjectInfo) TableName() string {
	return "project_info"
}
