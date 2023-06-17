package models

type EmRuleHistoryModel struct {
	Id                 int    `json:"id" gorm:"primary_key"`
	RuleId             int    `json:"ruleId"`
	RuleName           string `json:"ruleName"`
	Description        string `json:"description"`
	Level              int    `json:"level"`
	TypeClassification int    `json:"typeClassification"`
	ProduceTime        string `json:"produceTime"`
	RecoveryTime       string `json:"recoveryTime"`
	Tag                int    `json:"tag"`
	CreateTime         string `json:"createTime"`
	UpdateTime         string `json:"updateTime"`
}

func (u *EmRuleHistoryModel) TableName() string {
	return "rule_history"
}
