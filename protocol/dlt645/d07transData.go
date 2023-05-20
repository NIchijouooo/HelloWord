package dlt645

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/* 数据转换方向 */
const (
	ED07TransU2f int = iota // 数据格式从用户侧转换到帧侧
	ED07TransF2u            // 数据格式从帧侧转换到用户侧
)

type TransD07DataTemplate struct {
	TransDir int
	User     interface{}
	Frame    interface{}
}

// UnpackD07ByFormat 根据数据格式返回执行函数得到的数据
func UnpackD07ByFormat(format string, d *TransD07DataTemplate) (TransD07DataTemplate, error) {
	switch format {
	case "X.XXX":
		return d.TransD07DataFloatTemplate(2, 1, 3)
	case "XX.XXXX":
		return d.TransD07DataFloatTemplate(3, 2, 4)
	case "XXX.XXX":
		return d.TransD07DataFloatTemplate(3, 3, 3)
	case "XXX.X":
		return d.TransD07DataFloatTemplate(2, 3, 1)
	case "XXXXXX.XX":
		return d.TransD07DataFloatTemplate(4, 6, 2)
	case "YYMMDDhhmm":
		return d.TransD07DataTimeTemplate(5)
	default:

	}

	return TransD07DataTemplate{}, errors.New("不符合的数据格式")
}

//TransD07DataFloatTemplate  XX.XXXX   dataLength = 3 、 integerLength = 2 、 decimalsLength = 4
func (d *TransD07DataTemplate) TransD07DataFloatTemplate(dataLength int, integerLength int, decimalsLength int) (TransD07DataTemplate, error) {
	/*  检查接口类型 */
	inUser, userOK := d.User.(float64)
	inFrame, frameOK := d.Frame.([]byte)
	if !userOK || !frameOK {
		return TransD07DataTemplate{}, errors.New("传入的接口数据类型不对")
	}

	if integerLength+decimalsLength > dataLength*2 {
		return TransD07DataTemplate{}, errors.New("TransD07DataFloatTemplate传入的参数异常")
	}
	/* 定义返回值  */
	reTransD07Data := TransD07DataTemplate{d.TransDir, 0, nil} //返回值初始化

	/*  数据转换 */
	if ED07TransF2u == d.TransDir { //把帧数据转换成用户则应用数据
		var temFrame = make([]byte, 0)
		for k, _ := range inFrame {
			if inFrame[k] < 0x33 {
				return TransD07DataTemplate{}, errors.New("数据中有小于0X33的数据")
			}
			temFrame = append(temFrame, inFrame[k]-0x33) //帧数据要减去0x33
		}
		str, ok := D07BCD2Str(temFrame, dataLength) //把BCD转成字符串，从frameData转换N个值，如果frameData长度大于N也只转N个
		if ok != nil {
			return TransD07DataTemplate{}, errors.New("float类型数据D07BCD2Str转换异常")
		}
		//str = strings.TrimLeft(str, "0") //把左边的0全部去掉

		if len(str) < integerLength {
			return TransD07DataTemplate{}, errors.New("float类型数据转换字符长度异常")
		}

		str = str[:integerLength] + "." + str[integerLength:] //在对应位置添加小数点
		if num, err := strconv.ParseFloat(str, 64); nil == err {
			reTransD07Data.User = num      //把字符串转成float64
			reTransD07Data.Frame = inFrame //帧数据不变
			return reTransD07Data, nil
		}
	} else { //把用户则应用数据转换成帧数据
		var numberStr string
		format := "%" + fmt.Sprintf("%d.%df", integerLength, decimalsLength)
		numberStr = fmt.Sprintf(format, inUser)             //把float32的数转成对应的小数字符串
		numberStr = strings.TrimRight(numberStr, "0")       //把右边的0全部去掉
		numberStr = strings.Replace(numberStr, ".", "", -1) //把字符串中所有的点号用空字符替换

		if buf, eer := D07Str2BCD(numberStr, dataLength); eer == nil {
			for k := range buf {
				buf[k] = buf[k] + 0x33 //帧数据要加上0x33
			}
			var f = append(make([]byte, 0), buf...)
			reTransD07Data.Frame = f
			reTransD07Data.User = inUser //用户则数据不变
			return reTransD07Data, nil
		}
	}

	return TransD07DataTemplate{}, errors.New("float类型数据转换异常")
}

//TransD07DataTimeTemplate  dataLength--要转换的日期时间长度
func (d *TransD07DataTemplate) TransD07DataTimeTemplate(dataLength int) (TransD07DataTemplate, error) {
	/*  检查接口类型 */
	inUser, userOK := d.User.(string)
	inFrame, frameOK := d.Frame.([]byte)
	if !userOK || !frameOK {
		return TransD07DataTemplate{}, errors.New("传入的接口数据类型不对")
	}

	/* 定义返回值  */
	reTransD07Data := TransD07DataTemplate{d.TransDir, 0, nil} //返回值初始化

	/*  数据转换 */
	if ED07TransF2u == d.TransDir { //把帧数据转换成用户则应用数据
		var temFrame = make([]byte, 0)
		for k, _ := range inFrame {
			if inFrame[k] < 0x33 {
				return TransD07DataTemplate{}, errors.New("数据中有小于0X33的数据")
			}
			temFrame = append(temFrame, inFrame[k]-0x33) //帧数据要减去0x33
		}
		str, err := D07BCD2Str(temFrame, dataLength) //把BCD转成字符串，从TimeData转换N个值，如果TimeData长度大于2也只转N个
		if nil == err {
			reTransD07Data.User = str
			reTransD07Data.Frame = inFrame //帧数据不变
			return reTransD07Data, nil
		}
	} else { //把用户则应用数据转换成帧数据
		if buf, eer := D07Str2BCD(inUser, dataLength); eer == nil {
			for k := range buf {
				buf[k] = buf[k] + 0x33 //帧数据要加上0x33
			}
			var f = append(make([]byte, 0), buf...)
			reTransD07Data.Frame = f
			reTransD07Data.User = inUser //用户则数据不变
			return reTransD07Data, nil
		}
	}
	return TransD07DataTemplate{}, errors.New("time类型数据转换异常")
}
