package models

type PcsDeviceInfo struct {
	DeviceStatus  int                        `json:"deviceStatus" description:"设备状态;0-运行异常;1-运行正常"`
	EquipmentInfo DeviceEquipmentAccountInfo `json:"equipmentInfo" description:"台账信息"`
}
