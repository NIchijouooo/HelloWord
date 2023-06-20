package models

type EmStrategy struct {
	Id            int     `json:"id" gorm:"primary_key"`
	StartDate     string  `json:"startDate"`
	EndDate       string  `json:"endDate"`
	StartTime     string  `json:"startTime"`
	EndTime       string  `json:"endTime"`
	ActivePower   float64 `json:"activePower"`
	ReactivePower float64 `json:"reactivePower"`
	Status        int     `json:"status"`
}

func (u *EmStrategy) TableName() string {
	return "em_strategy"
}
