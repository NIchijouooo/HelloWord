package dlt645

import (
	"errors"
	"fmt"
	"strconv"
)

func D07BCD2Str(bcdData []byte, sliceLen int) (string, error) {

	/* 检查长度是否大于切片的长度  */
	if sliceLen > len(bcdData) {
		return "", errors.New("D07BCD2Str转换长度过长")
	}

	var outStr string
	for i := sliceLen - 1; i >= 0; i-- {
		outStr += fmt.Sprintf("%02X", bcdData[i])
	}

	return outStr, nil
}

func D07Str2BCD(strData string, sliceLen int) ([]byte, error) {
	outData := make([]byte, 0)
	sLen := len(strData) //字符串的长度
	j := sLen - sLen%2
	k := sLen/2 + sLen%2 //字符能转换出的最大切片长度

	/* 检查字符串是否能转出需要的切片长度sliceLen  */
	if sliceLen > k {
		return nil, errors.New("D07Str2BCD转换字符长度不够")
	}

	start := 0
	if sLen%2 != 0 {
		str := string(strData[0])
		if num, ok := strconv.ParseUint(str, 16, 8); ok == nil {
			outData = append(outData, byte(num))
			start = 1
		} else {
			return nil, errors.New("字符串中存在非法BCD字符2")
		}
	}

	for i := start; i < j; i += 2 {
		str := string(strData[i]) + string(strData[i+1])
		if num, ok := strconv.ParseUint(str, 16, 8); ok == nil {
			outData = append(outData, byte(num))
		} else {
			return nil, errors.New("字符串中存在非法BCD字符1")
		}
	}

	i, j := 0, len(outData)-1 //把数据反序号
	for i < j {
		outData[i], outData[j] = outData[j], outData[i]
		i++
		j--
	}

	return outData, nil
}

//func Hex2Byte(str string) []byte {
//	slen := len(str)
//	bHex := make([]byte, len(str)/2)
//	ii := 0
//	for i := 0; i < len(str); i = i + 2 {
//		if slen != 1 {
//			ss := string(str[i]) + string(str[i+1])
//			bt, _ := strconv.ParseInt(ss, 16, 32)
//			bHex[ii] = byte(bt)
//			ii = ii + 1
//			slen = slen - 2
//		}
//	}
//	return bHex
//}
//
//func String2bcd(number string) []byte {
//	var rNumber = number
//
//	for i := 0; i < 8-len(number); i++ {
//		rNumber = "f" + rNumber
//	}
//	bcd := Hex2Byte(rNumber)
//	return bcd
//}
//
//func Bcd2String(bcd []byte) string {
//	var number string
//	for _, i := range bcd {
//		number += fmt.Sprintf("%02X", i)
//	}
//	pos := strings.LastIndex(number, "F")
//	if pos == 8 {
//		return "0"
//	}
//	return number[pos+1:]
//}
