package operation

import (
	"bytes"
	"gateway/models"
	"gateway/repositories"
	"gateway/rule/enum"
	"log"
	"strconv"
	"strings"
)

// 将关键字替换为实际值
func ProcessKeyword(rule models.EmRuleModel, condition string) string {
	var result bytes.Buffer
	var sb bytes.Buffer

	for i := 0; i < len(condition); i++ {
		// 汉字
		if condition[i] < 0 {
			// 字符串异常截断
			if i == len(condition)-1 {
				// 不要这个字
			} else {
				sb.WriteString(condition[i : i+2])
				i++
			}
			continue
		}

		thisOp := string(condition[i])

		// 直接入栈
		if thisOp == "(" {
			result.WriteString(thisOp)
			continue
		}

		// 当前是否为操作符
		if enum.IsOperator(thisOp) || thisOp == ")" {
			if sb.Len() > 0 {
				s := getPropertyValue(sb.String(), rule)
				result.WriteString(s)
				sb.Reset()
			}
			result.WriteString(thisOp)
		} else {
			sb.WriteString(thisOp)
		}
	}

	if sb.Len() > 0 {
		s := getPropertyValue(sb.String(), rule)
		result.WriteString(s)
	}

	return result.String()
}

// 获取属性值
func getPropertyValue(operand string, rule models.EmRuleModel) string {
	// product.${10016:”逆变器”}.${0_1:”日发电量”}
	split := strings.Split(operand, ".")
	if len(split) != 3 {
		return operand
	}
	productCodeData := split[1]
	productCode := getKeyValue(productCodeData)

	signalCodeData := split[2]
	signalCode := getKeyValue(signalCodeData)

	realTimeDataJson := rule.RealTimeDataJson
	if realTimeDataJson == nil {
		return operand
	}
	deviceId := realTimeDataJson.DeviceId
	device := repositories.NewDevicePointRepository().GetDeviceByDeviceId(deviceId)
	if device == nil {
		return operand
	}

	if productCode != device.Label {
		return operand
	}

	propertyValue := getPropertyValueData(signalCode, deviceId)
	if propertyValue == "" {
		return operand
	}
	return propertyValue
}

// 获取键值
func getKeyValue(valueData string) string {
	start := strings.Index(valueData, "{") + 1
	end := strings.Index(valueData, ":")
	return valueData[start:end]
}

// 获取属性值数据
func getPropertyValueData(propertyCodeData string, deviceId int) string {
	var value string
	s := strings.Split(propertyCodeData, "_")
	propertyType := s[0]
	code, err := strconv.Atoi(s[1])
	if err != nil {
		log.Fatalln("getPropertyValueData err : ", err)
		return value
	}
	switch propertyType {
	case "0", "6", "7", "8":
		// 遥测
		yc, err := repositories.NewRealtimeDataRepository().GetYcById(deviceId, code)
		if err != nil {
			log.Fatalln("getPropertyValueData err : ", err)
			return value
		}
		value = strconv.FormatFloat(yc.Value, 'f', -1, 64)
	case "4":
		// 遥信
		yx, err := repositories.NewRealtimeDataRepository().GetYxById(deviceId, code)
		if err != nil {
			log.Fatalln("getPropertyValueData err : ", err)
			return value
		}
		value = strconv.Itoa(yx.Value)
	default:
		// 参数
	}
	return value
}
