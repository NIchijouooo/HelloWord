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
