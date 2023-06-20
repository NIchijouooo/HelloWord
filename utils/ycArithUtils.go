package utils

import (
	"gateway/models"
	"math/big"
	"reflect"
	"time"
)

//遥测累加
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
			if obj.Value != 0 {
				zero.Add(zero, big.NewFloat(obj.Value))
			}
		}
	}
	// 保留三位小数点
	resultFloat, _ := zero.SetPrec(3).Float64()
	return resultFloat
}

//取出最大遥测值
func YcValueMax(ycModelList interface{}) float64 {
	result := 0.0

	// 检查 ycModelList 是否为空
	if ycModelList == nil || reflect.ValueOf(ycModelList).Len() == 0 {
		return result
	}

	var zero float64
	//初始化一个开始时间
	calendar := time.Unix(0, 643000*int64(time.Millisecond))
	// 使用反射获取 ycModelList 的值并进行遍历
	listValue := reflect.ValueOf(ycModelList)
	for i := 0; i < listValue.Len(); i++ {
		obj := listValue.Index(i).Interface()
		switch obj := obj.(type) {
		case models.YcData:
			if obj.Ts.Time.After(calendar) { //如果当日期大于calendar日期就可以重新赋值
				zero = obj.Value //重新赋值zero
				calendar = obj.Ts.Time
			}
		}
	}
	// 保留三位小数点
	//resultFloat, _ := zero.SetPrec(3).Float64()
	//fmt.Println("修改后:", resultFloat)
	return zero

}
