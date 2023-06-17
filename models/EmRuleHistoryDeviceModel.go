package models

type EmRuleHistoryDeviceModel struct {
	Id            int    `json:"id" gorm:"primary_key"`
	RuleHistoryId int    `json:"ruleHistoryId"`
	DeviceId      int    `json:"deviceId"`
	PropertyCode  int    `json:"propertyCode"`
	CreateTime    string `json:"createTime"`
}

func (u *EmRuleHistoryDeviceModel) TableName() string {
	return "rule_history_device"
}
