package enum

type OperatorEnum int

const (
	NOT OperatorEnum = iota
	LT
	ELT
	GT
	EGT
	EQ
	NEQ
	BAND
	BOR
	AND
	OR
	E
	ADD
	SUBTRACT
	MULTIPLY
	DIVIDE
)

var operatorNames = [...]string{
	"!",
	"<",
	"<=",
	">",
	">=",
	"==",
	"!=",
	"&",
	"|",
	"&&",
	"||",
	"=",
	"+",
	"-",
	"*",
	"/",
}

var operatorPriorities = [...]int{
	0,
	3,
	3,
	3,
	3,
	3,
	3,
	4,
	4,
	4,
	4,
	2,
	1,
	1,
	0,
	0,
}

func (o OperatorEnum) String() string {
	return operatorNames[o]
}

func (o OperatorEnum) Priority() int {
	return operatorPriorities[o]
}

func GetOperatorEnumByName(name string) int {
	for i, n := range operatorNames {
		if n == name {
			return operatorPriorities[i]
		}
	}
	return -1
}

func IsOperator(name string) bool {
	for _, n := range operatorNames {
		if n == name {
			return true
		}
	}
	return false
}
