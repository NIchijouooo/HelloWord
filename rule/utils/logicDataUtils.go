package utils

import (
	"bytes"
	"gateway/rule/enum"
	"strings"
)

func GetConditionVariableObjectId(condition string) string {
	conditionSplit := strings.Split(condition, ".")
	if len(conditionSplit) != 3 {
		return ""
	}
	propertyData := conditionSplit[1]
	if propertyData == "" {
		return ""
	}
	startIndex := strings.Index(propertyData, "{") + 1
	endIndex := strings.Index(propertyData, ":")
	if startIndex == -1 || endIndex == -1 {
		return ""
	}
	return propertyData[startIndex:endIndex]
}

func GetConditionVariablePropertyId(condition string) string {
	conditionSplit := strings.Split(condition, ".")
	if len(conditionSplit) != 3 {
		return ""
	}
	propertyData := conditionSplit[2]
	if propertyData == "" || propertyData == "status" {
		return propertyData
	}
	startIndex := strings.Index(propertyData, "{") + 1
	endIndex := strings.Index(propertyData, ":")
	if startIndex == -1 || endIndex == -1 {
		return ""
	}
	return propertyData[startIndex:endIndex]
}

func GetConditionVariableEventIdentifier(condition string) string {
	conditionSplit := strings.Split(condition, ".")
	if len(conditionSplit) != 3 {
		return ""
	}
	propertyData := conditionSplit[1]
	if propertyData == "" {
		return propertyData
	}
	startIndex := strings.Index(propertyData, "{") + 1
	endIndex := strings.Index(propertyData, "}")
	if startIndex == -1 || endIndex == -1 {
		return ""
	}
	return propertyData[startIndex:endIndex]
}

func GetOperator(expression string) []string {
	operands := make([]string, 0)
	operate := make([]string, 0)
	var sb bytes.Buffer

	bytes := []byte(expression)
	for i := 0; i < len(bytes); i++ {
		// 汉字
		if bytes[i] < 0 {
			// 字符串异常截断
			if i == len(bytes)-1 {
				// 不要这个字
			} else {
				sb.WriteString(string(bytes[i : i+2]))
				i++
			}
			continue
		}
		thisOp := string(bytes[i])
		// 直接入栈
		if thisOp == "(" {
			pushOperandIntoStack(&operands, &sb)
			operate = append([]string{thisOp}, operate...)
			continue
		}

		// 到“(”之前的全部出栈
		if thisOp == ")" {
			pushOperandIntoStack(&operands, &sb)
			topOp := operate[0]
			operate = operate[1:]
			for topOp != "(" {
				if len(operate) == 0 {
					//括号没匹配上，逻辑表达式出错
					return nil
				}
			}
			continue
		}

		// 当前是否为操作符
		if enum.IsOperator(thisOp) {
			// 负数运算的特殊处理
			if thisOp == "-" {
				// 第一个为-或者-之前是运算符
				if i == 0 || enum.IsOperator(string(bytes[i-1])) || thisOp == "(" || thisOp == ")" {
					sb.WriteString(thisOp)
					continue
				}
			}

			// 是，1.查看之前是否有字符未入栈，先入栈字符
			pushOperandIntoStack(&operands, &sb)
			// 2.查看下一个是否为操作，并且非括号，是合并当前一起入操作栈
			nextOp := string(bytes[i+1])
			// 下个与当前一起组成一个操作
			if nextOp != "(" && nextOp != ")" && nextOp != "-" && enum.IsOperator(nextOp) {
				thisOp += nextOp
				i++
			}
			operate = append([]string{thisOp}, operate...)
		} else {
			sb.WriteString(thisOp)
		}
	}

	if sb.Len() > 0 {
		operands = append([]string{sb.String()}, operands...)
	}

	return operands
}

func pushOperandIntoStack(operands *[]string, sb *bytes.Buffer) {
	if sb.Len() > 0 {
		*operands = append([]string{sb.String()}, *operands...)
		sb.Reset()
	}
}

func getIdValue(valueData string) string {
	startIndex := strings.Index(valueData, "{") + 1
	endIndex := strings.Index(valueData, ":")
	if startIndex == -1 || endIndex == -1 {
		return ""
	}
	return valueData[startIndex:endIndex]
}
