package models

type EmRuleModel struct {
	Id                 int    `json:"id" gorm:"primary_key"`
	Name               string `json:"name"`
	Content            string `json:"content"`
	Description        string `json:"description"`
	Level              int    `json:"level"`
	EnableFlag         int    `json:"enableFlag"`
	TypeClassification int    `json:"typeClassification"`
	ModelJson          string `json:"modelJson"`
}

func (u *EmRuleModel) TableName() string {
	return "rule"
}
