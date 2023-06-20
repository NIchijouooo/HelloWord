package models

type LimitConfigVo struct {
	Id              int    `json:"id" gorm:"primary_key"`
	DeviceType      string `json:"deviceType"`
	PropertyCode    string `json:"propertyCode"`
	EnableFlag      string `json:"enableFlag"`
	NotifyMin       string `json:"notifyMin"`
	NotifyMax       string `json:"notifyMax"`
	NotifyRuleId    string `json:"notifyRuleId"`
	SecondaryMin    string `json:"secondaryMin"`
	SecondaryMax    string `json:"secondaryMax"`
	SecondaryRuleId string `json:"secondaryRuleId"`
	SeriousMin      string `json:"seriousMin"`
	SeriousMax      string `json:"seriousMax"`
	SeriousRuleId   string `json:"seriousRuleId"`
	UrgentMin       string `json:"urgentMin"`
	UrgentMax       string `json:"urgentMax"`
	UrgentRuleId    string `json:"urgentRuleId"`
	DelFlag         int    `json:"delFlag"`
	CreateTime      string `json:"createTime"`
	UpdateTime      string `json:"updateTime"`
}

func (u *LimitConfigVo) TableName() string {
	return "em_limit_config"
}
