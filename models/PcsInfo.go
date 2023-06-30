package models

type PcsDeviceInfo struct {
	DeviceStatus  string                     `json:"deviceStatus" description:"PCS状态详情"`
	EquipmentInfo DeviceEquipmentAccountInfo `json:"equipmentInfo" description:"台账信息"`
}
