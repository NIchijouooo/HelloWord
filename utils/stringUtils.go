package utils

import (
	"strconv"
	"strings"
)

// intArray要处理的数组，separator分隔符，例如数组[1,2,3] 生成str:="1,2,3"
func IntArrayToString(intArray []int, separator string) string {
	strArray := make([]string, len(intArray))
	for i, num := range intArray {
		strArray[i] = strconv.Itoa(num)
	}
	return strings.Join(strArray, separator)
}

// strArray字符串数组转int数组，要处理的数组，返回int数组
func StringArrayToIntArray(strArray []string) ([]int, error) {
	intArray := make([]int, len(strArray))
	for i, str := range strArray {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		intArray[i] = num
	}
	return intArray, nil
}
