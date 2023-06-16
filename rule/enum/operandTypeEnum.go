package enum

type OperandTypeEnum int

const (
	NUM OperandTypeEnum = iota
	DATE
	STR
)

var operandTypeNames = [...]string{
	"数字",
	"日期",
	"字符串",
}

func (o OperandTypeEnum) String() string {
	return operandTypeNames[o]
}

func GetOperandTypeEnumByName(name string) OperandTypeEnum {
	for i, n := range operandTypeNames {
		if n == name {
			return OperandTypeEnum(i)
		}
	}
	return -1
}
