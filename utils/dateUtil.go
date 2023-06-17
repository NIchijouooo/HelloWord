package utils

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 比较日期大小
func CompareDate(date1, date2 string) int {
	d1 := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(date1), "-", ""), ":", ""), "/", "")
	d2 := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(date2), "-", ""), ":", ""), "/", "")
	sb1 := strings.Builder{}
	sb2 := strings.Builder{}
	sb1.WriteString(d1)
	sb2.WriteString(d2)
	for sb1.Len() < 14 {
		sb1.WriteString("0")
	}
	for sb2.Len() < 14 {
		sb2.WriteString("0")
	}
	num1, _ := strconv.ParseInt(sb1.String(), 10, 64)
	num2, _ := strconv.ParseInt(sb2.String(), 10, 64)
	if num1 == num2 {
		return 0
	} else if num1 > num2 {
		return 1
	} else {
		return -1
	}
}

// 验证字符串是否为合法日期 支持2019-03-12 2019/03/12 2019.03.12   HH:mm:ss HH:mm常用格式
func VerifyDateLegal(date string) bool {
	if strings.ContainsAny(date, "-/.") || strings.Count(date, "-") != 2 {
		return false
	}
	date = regexp.MustCompile("[./]").ReplaceAllString(date, "-")
	timeSb := strings.Builder{}
	parts := strings.SplitN(date, " ", 2)
	timeSb.WriteString(parts[0])
	timeSb.WriteString(" ")
	if len(parts) > 1 {
		timeSb.WriteString(parts[1])
	}
	for i := len(parts[1]); i < 8; i++ {
		if i == 2 || i == 5 {
			timeSb.WriteString(":")
		} else {
			timeSb.WriteString("0")
		}
	}
	_, err := time.Parse("2006-01-02 15:04:05", timeSb.String())
	return err == nil
}

//格式化日期
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

//计算根据间隔以及间隔单位得出时间差
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
