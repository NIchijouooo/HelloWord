package models

type EmStrategy struct {
	Id            int     `json:"id" gorm:"primary_key"`
	Date          string  `json:"date"`
	StartTime     string  `json:"startTime"`
	EndTime       string  `json:"endTime"`
	ActivePower   float64 `json:"activePower"`
	ReactivePower float64 `json:"reactivePower"`
	Status        int     `json:"status"`
}

func (u *EmStrategy) TableName() string {
	return "em_strategy"
}
