package models

type EmRuleHistoryModel struct {
	Id           int    `json:"id" gorm:"primary_key"`
	EventId      int    `json:"eventId"`
	EventName    string `json:"eventName"`
	Description  string `json:"description"`
	Level        int    `json:"level"`
	ProduceTime  string `json:"produceTime"`
	RecoveryTime string `json:"recoveryTime"`
	Tag          int    `json:"tag"`
}

func (u *EmRuleHistoryModel) TableName() string {
	return "rule_history"
}
