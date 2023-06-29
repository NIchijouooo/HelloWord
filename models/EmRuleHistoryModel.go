package models

type EmRuleHistoryModel struct {
	Id                 int    `json:"id" gorm:"primary_key"`
	RuleId             int    `json:"ruleId"`
	RuleName           string `json:"ruleName"`
	Description        string `json:"description"`
	Level              int    `json:"level"` // 事件等级(0-通知；1-次要；2-告警；3-故障)
	TypeClassification int    `json:"typeClassification"`
	ProduceTime        string `json:"produceTime"`
	RecoveryTime       string `json:"recoveryTime"`
	Tag                int    `json:"tag"` // 恢复标记：0-未确认，1-自动恢复 2-手动恢复
	CreateTime         string `json:"createTime"`
	UpdateTime         string `json:"updateTime"`
}

func (u *EmRuleHistoryModel) TableName() string {
	return "rule_history"
}
