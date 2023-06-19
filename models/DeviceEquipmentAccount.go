package models

/**
设备台账信息类
*/

type DeviceEquipmentAccountInfo struct {
	ID                    int     `json:"id" gorm:"primary_key"`
	Manufacturer          string  `json:"manufacturer" description:"生产厂家"`
	FactoryModel          string  `json:"factoryModel" description:"出厂型号"`
	BatteryCluster        int     `json:"batteryCluster" description:"电池簇"`
	BatteryNumber         int     `json:"batteryNumber" description:"电池数量"`
	BatterySpecifications float64 `json:"batterySpecifications" description:"电池规格"`
	RatedPower            float64 `json:"ratedPower" description:"额定功率"`
	DeviceParameters      string  `json:"deviceParameters" description:"设备参数,json字符串"`
	DeviceId              int     `json:"deviceId" description:"设备id"`
}

func (u *DeviceEquipmentAccountInfo) TableName() string {
	return "device_equipment_account"
}
