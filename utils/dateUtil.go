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

// 格式化日期
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

// 计算根据间隔以及间隔单位得出时间差
func GetIntervalTime(calendar time.Time, intervalType int, interval int) (int64, time.Time) {
	timeInMillis := calendar.UnixMilli()
	var intervalLong int64

	if intervalType == 1 {
		calendar = calendar.Add(time.Duration(interval) * time.Second)
		calendar = time.Date(calendar.Year(), calendar.Month(), calendar.Day(), calendar.Hour(), calendar.Minute(), calendar.Second(), 0, calendar.Location())
		intervalLong = calendar.UnixNano()/int64(time.Millisecond) - timeInMillis
	} else if intervalType == 2 {
		calendar = calendar.Add(time.Duration(interval) * time.Minute)
		calendar = time.Date(calendar.Year(), calendar.Month(), calendar.Day(), calendar.Hour(), calendar.Minute(), 0, 0, calendar.Location())
		intervalLong = calendar.UnixNano()/int64(time.Millisecond) - timeInMillis
	} else if intervalType == 3 {
		calendar = calendar.Add(time.Duration(interval) * time.Hour)
		calendar = time.Date(calendar.Year(), calendar.Month(), calendar.Day(), calendar.Hour(), 0, 0, 0, calendar.Location())
		intervalLong = calendar.UnixNano()/int64(time.Millisecond) - timeInMillis
	} else if intervalType == 4 {
		calendar = calendar.AddDate(0, 0, interval)
		calendar = time.Date(calendar.Year(), calendar.Month(), calendar.Day(), 0, 0, 0, 0, calendar.Location())
		intervalLong = calendar.UnixNano()/int64(time.Millisecond) - timeInMillis
	} else if intervalType == 5 {
		calendar = calendar.AddDate(0, interval, 0)
		calendar = time.Date(calendar.Year(), calendar.Month(), 1, 0, 0, 0, 0, calendar.Location())
		intervalLong = calendar.UnixNano()/int64(time.Millisecond) - timeInMillis
	} else if intervalType == 6 {
		calendar = calendar.AddDate(interval, 0, 0)
		calendar = time.Date(calendar.Year(), 1, 1, 0, 0, 0, 0, calendar.Location())
		intervalLong = calendar.UnixNano()/int64(time.Millisecond) - timeInMillis
	}
	return intervalLong, calendar

}

// GetFirstDateOfWeek 获取本周周一的日期
func GetFirstDateOfWeek(t time.Time) time.Time {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
}

// GetLastDateOfWeek 获取本周周日
func GetLastDateOfWeek(t time.Time) time.Time {
	return GetFirstDateOfWeek(t).AddDate(0, 0, 6)
}

// GetLastWeekFirstDate 获取上周的周一
func GetLastWeekFirstDate(t time.Time) time.Time {
	thisWeekMonday := GetFirstDateOfWeek(t)
	return thisWeekMonday.AddDate(0, 0, -7)
}

// GetNextFirstDateOfWeek 获取下周周一
func GetNextFirstDateOfWeek(t time.Time) time.Time {
	return GetFirstDateOfWeek(t).AddDate(0, 0, 7)
}

// GetLastWeekLastDate 获取下周周日
func GetLastWeekLastDate(t time.Time) time.Time {
	return GetLastDateOfWeek(t).AddDate(0, 0, -7)
}

// GetFirstDateOfMonth 获取本月第一天
func GetFirstDateOfMonth(t time.Time) time.Time {
	return t.AddDate(0, 0, -t.Day()+1)
}

// GetLastDateOfMonth  获取本月最后一天
func GetLastDateOfMonth(t time.Time) time.Time {
	return GetFirstDateOfMonth(t).AddDate(0, 1, -1)
}

// GetLastMonthFirstDate 获取上个月的第一天
func GetLastMonthFirstDate(t time.Time) time.Time {
	lastMonth := t.AddDate(0, -1, 0)
	return time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// GetLastDateOfLastMonth 获取上个月的最后一天
func GetLastDateOfLastMonth(t time.Time) time.Time {
	lastMonth := t.AddDate(0, -1, 0)
	// 获取上个月的下一个月的第一天
	nextMonthFirstDay := time.Date(lastMonth.Year(), lastMonth.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	// 上个月的最后一天即为下个月的第一天的前一天
	return nextMonthFirstDay.AddDate(0, 0, -1)
}

func GetFirstDateOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
}

// GetLastDateOfYear 获取今年的最后一天
func GetLastDateOfYear(t time.Time) time.Time {
	nextYearFirstDay := time.Date(t.Year()+1, time.January, 1, 0, 0, 0, 0, time.UTC)
	// 今年的最后一天即为明年的第一天的前一天
	return nextYearFirstDay.AddDate(0, 0, -1)
}

// GetFirstDateOfFirstYear 获取去年的第一天
func GetFirstDateOfFirstYear(t time.Time) time.Time {
	return time.Date(t.Year()-1, time.January, 1, 0, 0, 0, 0, time.UTC)
}

// GetLastDateOfLastYear 获取去年的最后一天
func GetLastDateOfLastYear(t time.Time) time.Time {
	// 获取今年的第一天
	thisYearFirstDay := time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	// 去年的最后一天即为今年的第一天的前一天
	return thisYearFirstDay.AddDate(0, 0, -1)
}
