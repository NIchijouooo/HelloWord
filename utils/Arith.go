package utils

import (
	"gateway/models"
	"math/big"
	"reflect"
	"strconv"
)

func YcValueSum(ycModelList interface{}) float64 {
	result := 0.0

	// 检查 ycModelList 是否为空
	if ycModelList == nil || reflect.ValueOf(ycModelList).Len() == 0 {
		return result
	}

	zero := big.NewFloat(0.0)

	// 使用反射获取 ycModelList 的值并进行遍历
	listValue := reflect.ValueOf(ycModelList)
	for i := 0; i < listValue.Len(); i++ {
		obj := listValue.Index(i).Interface()

		switch obj := obj.(type) {
		case models.YcData:
			value, err := strconv.ParseFloat(obj.Value, 64)
			if err == nil && value != 0 {
				zero.Add(zero, big.NewFloat(value))
			}
		}
	}
	// 保留三位小数点
	resultFloat, _ := zero.SetPrec(3).Float64()
	return resultFloat
}
