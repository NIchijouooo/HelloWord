package rule

import (
	"bytes"
	"gateway/rule/enum"
	"gateway/utils"
	"strconv"
	"strings"
)

// Parse 运算类
func Parse(expression string) string {
	//操作数栈
	operands := make([]string, 0)
	//操作栈
	operate := make([]string, 0)
	var sb bytes.Buffer
	for i := 0; i < len(expression); i++ {
		ch := expression[i]
		//汉字
		if ch < 0 {
			//字符串异常截断
			if i == len(expression)-1 {
				//不要这个字
			} else {
				sb.WriteString(expression[i : i+2])
				i++
			}
			continue
		}
		thisOp := string(ch)
		//直接入栈
		if thisOp == "(" {
			pushOperandIntoStack(&operands, &sb)
			operate = append([]string{thisOp}, operate...)
			continue
		}
		//到“(”之前的全部出栈
		if thisOp == ")" {
			pushOperandIntoStack(&operands, &sb)
			topOp := operate[0]
			operate = operate[1:]
			for topOp != "(" {
				if len(operate) == 0 {
					//括号没匹配上，逻辑表达式出错
					return ""
				}
				calculate(&operands, topOp)
				topOp = operate[0]
				operate = operate[1:]
			}
			continue
		}
		//当前是否为操作符
		if enum.IsOperator(thisOp) {
			//负数运算的特殊处理
			if thisOp == "-" {
				//第一个为-或者-之前是运算符
				if i == 0 || enum.IsOperator(string(expression[i-1])) ||
					"(" == string(expression[i-1]) || ")" == string(expression[i-1]) {
					sb.WriteString(thisOp)
					continue
				}
			}
			//是，1.查看之前是否有字符未入栈，先入栈字符
			pushOperandIntoStack(&operands, &sb)
			//2.查看下一个是否为操作，并且非括号，是合并当前一起入操作栈
			nextOp := string(expression[i+1])
			//下个与当前一起组成一个操作
			if "(" != nextOp && ")" != nextOp && "-" != nextOp && enum.IsOperator(nextOp) {
				thisOp += nextOp
				i++
			}
			//判断当前操作与栈顶操作优先级
			if len(operate) > 0 {
				topOp := operate[0]
				for topOp != "(" && enum.GetOperatorEnumByName(topOp) <= enum.GetOperatorEnumByName(thisOp) {
					//优先级高，出栈进行计算
					operate = operate[1:]
					calculate(&operands, topOp)
					if len(operate) > 0 {
						topOp = operate[0]
					} else {
						break
					}
				}
			}
			operate = append([]string{thisOp}, operate...)
		} else {
			sb.WriteString(thisOp)
		}
	}
	if sb.Len() > 0 {
		operands = append([]string{sb.String()}, operands...)
	}
	for len(operate) > 0 {
		topOp := operate[0]
		operate = operate[1:]
		calculate(&operands, topOp)
	}
	if len(operands) > 0 {
		return operands[0]
	}
	return ""
}

func pushOperandIntoStack(operands *[]string, sb *bytes.Buffer) {
	if sb.Len() > 0 {
		*operands = append([]string{sb.String()}, *operands...)
		sb.Reset()
	}
}

func calculate(operands *[]string, topOp string) {
	operand2 := strings.TrimSpace((*operands)[0])
	operand1 := strings.TrimSpace((*operands)[1])
	*operands = (*operands)[2:]
	//判断两个操作数类型，不一致不可比较直接返回false
	type1 := judgeType(operand1)
	type2 := judgeType(operand2)

	if type1 == type2 {
		switch type1 {
		case enum.NUM.String():
			if topOp == "+" || topOp == "-" || topOp == "*" || topOp == "/" {
				s := numCalculate1(operand1, operand2, topOp)
				*operands = append([]string{s}, *operands...)
			} else {
				*operands = append([]string{strconv.FormatBool(numCalculate(operand1, operand2, topOp))}, *operands...)
			}
		case enum.DATE.String():
			*operands = append([]string{strconv.FormatBool(dateCalculate(operand1, operand2, topOp))}, *operands...)
		case enum.STR.String():
			*operands = append([]string{strconv.FormatBool(strCalculate(operand1, operand2, topOp))}, *operands...)
		default:
			break
		}
	} else {
		*operands = append([]string{"false"}, *operands...)
	}
}

func judgeType(operands string) string {
	operands = strings.TrimSpace(operands)
	_, err := strconv.ParseFloat(operands, 64)
	if err == nil {
		return enum.NUM.String()
	}
	if utils.VerifyDateLegal(operands) {
		return enum.DATE.String()
	}
	return enum.STR.String()
}

func numCalculate1(operand1 string, operand2 string, operate string) string {
	num1, _ := strconv.ParseFloat(operand1, 64)
	num2, _ := strconv.ParseFloat(operand2, 64)
	switch operate {
	case enum.ADD.String():
		return strconv.FormatFloat(num1+num2, 'f', -1, 64)
	case enum.SUBTRACT.String():
		return strconv.FormatFloat(num1-num2, 'f', -1, 64)
	case enum.MULTIPLY.String():
		return strconv.FormatFloat(num1*num2, 'f', -1, 64)
	case enum.DIVIDE.String():
		return strconv.FormatFloat(num1/num2, 'f', -1, 64)
	default:
		return "0"
	}
}

func numCalculate(operand1 string, operand2 string, operate string) bool {
	num1, _ := strconv.ParseFloat(operand1, 64)
	num2, _ := strconv.ParseFloat(operand2, 64)
	switch operate {
	case enum.LT.String():
		return num1 < num2
	case enum.ELT.String():
		return num1 <= num2
	case enum.GT.String():
		return num1 > num2
	case enum.EGT.String():
		return num1 >= num2
	case enum.EQ.String():
		return num1 == num2
	case enum.NEQ.String():
		return num1 != num2
	default:
		return true
	}
}

func strCalculate(operand1 string, operand2 string, operate string) bool {
	switch operate {
	case enum.EQ.String():
		return operand1 == operand2
	case enum.NEQ.String():
		return operand1 != operand2
	case enum.AND.String():
		return operand1 == "true" && operand2 == "true"
	case enum.OR.String():
		return operand1 == "true" || operand2 == "true"
	default:
		return true
	}
}

func dateCalculate(operand1 string, operand2 string, operate string) bool {
	switch operate {
	case enum.LT.String():
		return utils.CompareDate(operand1, operand2) == -1
	case enum.ELT.String():
		return utils.CompareDate(operand1, operand2) <= 0
	case enum.GT.String():
		return utils.CompareDate(operand1, operand2) == 1
	case enum.EGT.String():
		return utils.CompareDate(operand1, operand2) >= 0
	case enum.EQ.String():
		return utils.CompareDate(operand1, operand2) == 0
	case enum.NEQ.String():
		return utils.CompareDate(operand1, operand2) != 0
	default:
		return true
	}
}
