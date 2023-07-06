package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var Location, _ = time.LoadLocation("Asia/Shanghai")

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
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, Location).AddDate(0, 0, offset)
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
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, Location)
}

// GetLastDateOfMonth  获取本月最后一天
func GetLastDateOfMonth(t time.Time) time.Time {
	nextMonth := t.AddDate(0, 1, 0)
	firstDayOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, Location)

	// 减去一天，即为本月最后一天的日期
	return firstDayOfNextMonth.AddDate(0, 0, -1)
}

// GetLastMonthFirstDate 获取上个月的第一天
func GetLastMonthFirstDate(t time.Time) time.Time {
	lastMonth := t.AddDate(0, -1, 0)
	return time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, Location)
}

// GetLastDateOfLastMonth 获取上个月的最后一天
func GetLastDateOfLastMonth(t time.Time) time.Time {
	lastMonth := t.AddDate(0, -1, 0)
	// 获取上个月的下一个月的第一天
	nextMonthFirstDay := time.Date(lastMonth.Year(), lastMonth.Month()+1, 1, 0, 0, 0, 0, Location)
	// 上个月的最后一天即为下个月的第一天的前一天
	return nextMonthFirstDay.AddDate(0, 0, -1)
}

func GetFirstDateOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, Location)
}

// GetLastDateOfYear 获取今年的最后一天
func GetLastDateOfYear(t time.Time) time.Time {
	nextYearFirstDay := time.Date(t.Year()+1, time.January, 1, 0, 0, 0, 0, Location)
	// 今年的最后一天即为明年的第一天的前一天
	return nextYearFirstDay.AddDate(0, 0, -1)
}

// GetFirstDateOfFirstYear 获取去年的第一天
func GetFirstDateOfFirstYear(t time.Time) time.Time {
	return time.Date(t.Year()-1, time.January, 1, 0, 0, 0, 0, Location)
}

// GetLastDateOfLastYear 获取去年的最后一天
func GetLastDateOfLastYear(t time.Time) time.Time {
	// 获取今年的第一天
	thisYearFirstDay := time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, Location)
	// 去年的最后一天即为今年的第一天的前一天
	return thisYearFirstDay.AddDate(0, 0, -1)
}

// 获取今天之前7天每一天的日期，格式为dd日，返回一个日期数组
func GetLast7Days() []string {
	var dates []string
	// 获取当前日期
	today := time.Now()
	// 构建日期数组
	for i := 7; i >= 1; i-- {
		date := today.AddDate(0, 0, -i)
		dates = append(dates, date.Format("1月2日"))
	}
	return dates
}

// 获取当月每一天的日期，格式为dd日，返回一个日期数组
func GetCurrentMonthDays() []string {
	var dates []string
	// 获取当月第一天的日期
	firstDay := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	// 获取下个月第一天的日期
	nextMonth := firstDay.AddDate(0, 1, 0)
	// 计算当月的天数
	numDays := nextMonth.Sub(firstDay).Hours() / 24
	// 构建日期数组
	for i := 0; i < int(numDays); i++ {
		date := firstDay.AddDate(0, 0, i)
		dates = append(dates, date.Format("2日"))
	}
	return dates
}

// 获取一年中12个月，格式为MM月，返回一个月份数组
func GetAllMonths() []string {
	var months []string
	// 构建月份数组
	for i := 1; i <= 12; i++ {
		month := fmt.Sprintf("%2d月", i)
		months = append(months, month)
	}
	return months
}

// 获取今天之前7天的开始时间戳和结束时间戳
func GetLast7DaysTimestamps() (int64, int64) {
	// 获取当前时间
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// 计算开始时间
	startTime := startOfDay.AddDate(0, 0, -7).UnixMilli()
	// 计算结束时间
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	endTime := endOfDay.UnixMilli()
	return startTime, endTime
}

// 获取当月的开始时间戳和结束时间戳
func GetCurrentMonthTimestamps() (int64, int64) {
	// 获取当前时间
	now := time.Now()
	// 获取当月第一天的日期
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	// 获取下个月第一天的日期
	nextMonth := firstDay.AddDate(0, 1, 0)
	// 计算开始时间
	startTime := firstDay.UnixMilli()
	// 计算结束时间
	endTime := nextMonth.UnixMilli()
	return startTime, endTime
}

// 获取当年的开始时间戳和结束时间戳
func GetCurrentYearTimestamps() (int64, int64) {
	// 获取当前时间
	now := time.Now()
	// 获取当年第一天的日期
	firstDay := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	// 获取下一年第一天的日期
	nextYear := time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, now.Location())
	// 计算开始时间
	startTime := firstDay.UnixMilli()
	// 计算结束时间
	endTime := nextYear.UnixMilli()
	return startTime, endTime
}
