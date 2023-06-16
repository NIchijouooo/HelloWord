package utils

import "time"

func GetIntervalDateFormat(intervalType int) string {
	pattern := "2006-01-02 15:04:05"
	// 1-秒;2-分;3-时;4-日;5-月;
	switch intervalType {
	case 2:
		pattern = "2006-01-02 15:04"
	case 3:
		pattern = "2006-01-02 15"
	case 4:
		pattern = "2006-01-02"
	case 5:
		pattern = "2006-01"
	case 6:
		pattern = "2006"
	}
	return pattern
}
func GetIntervalTime(calendar time.Time, intervalType int, interval int) int64 {
	timeInMillis := calendar.UnixNano() / int64(time.Millisecond)
	var intervalLong int64

	if intervalType == 1 {
		calendar = calendar.Add(time.Duration(interval) * time.Second)
		//calendar = time.Date(calendar.Year(), calendar.Month(), calendar.Day(), 0, 0, 0, 0, calendar.Location())
		intervalLong = calendar.UnixNano()/int64(time.Millisecond) - timeInMillis
	} else if intervalType == 2 {
		calendar = calendar.Add(time.Duration(interval) * time.Minute)
		//calendar = time.Date(calendar.Year(), calendar.Month(), calendar.Day(), 0, 0, 0, 0, calendar.Location())
		intervalLong = calendar.UnixNano()/int64(time.Millisecond) - timeInMillis
	} else if intervalType == 3 {
		calendar = calendar.Add(time.Duration(interval) * time.Hour)
		//calendar = time.Date(calendar.Year(), calendar.Month(), calendar.Day(), 0, 0, 0, 0, calendar.Location())
		intervalLong = calendar.UnixNano()/int64(time.Millisecond) - timeInMillis
	} else if intervalType == 4 {
		calendar = calendar.AddDate(0, 0, interval)
		//calendar = time.Date(calendar.Year(), calendar.Month(), calendar.Day(), 0, 0, 0, 0, calendar.Location())
		intervalLong = calendar.UnixNano()/int64(time.Millisecond) - timeInMillis
	} else if intervalType == 5 {
		calendar = calendar.AddDate(0, interval, 0)
		//calendar = time.Date(calendar.Year(), calendar.Month(), 1, 0, 0, 0, 0, calendar.Location())
		intervalLong = calendar.UnixNano()/int64(time.Millisecond) - timeInMillis
	} else if intervalType == 6 {
		calendar = calendar.AddDate(interval, 0, 0)
		//calendar = time.Date(calendar.Year(), 1, 1, 0, 0, 0, 0, calendar.Location())
		intervalLong = calendar.UnixNano()/int64(time.Millisecond) - timeInMillis
	}
	return intervalLong

}
