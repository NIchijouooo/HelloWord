package dataProcess

import (
	"gateway/models"
	"gateway/repositories"
	"gateway/rule/utils"
	"strconv"
	"strings"
)

var DeviceRuleMap = make(map[string]models.EmRuleModel)
var RuleMap = make(map[int]models.EmRuleModel)
var PropertyRuleMap = make(map[string][]int)
var SystemEventRuleMap = make(map[string][]int)

func ProcessRuleVariable(rule models.EmRuleModel) {
	content := rule.Content
	split := strings.Split(content, "\\n")
	action := make([]string, 0)
	condition := strings.Builder{}
	previous := ""
	ruleType := ""

	for _, s := range split {
		str := strings.TrimSpace(s)

		if strings.Contains(str, "type=") {
			ruleType = str[5:]
			continue
		}

		if str == "if" || str == "then" {
			previous = str
			continue
		}
		if previous == "if" {
			condition.WriteString(str)
		}

		if previous == "then" && ruleType == "calculate" && strings.Contains(str, "assign.${") {
			action = append(action, str)
		}
	}

	processLogicRule(rule, condition.String())
	RuleMap[rule.Id] = rule
}

func processLogicRule(rule models.EmRuleModel, condition string) {
	if len(condition) > 0 {
		operators := utils.GetOperator(condition)
		length := len(operators)
		for i := 0; i < length; i++ {
			operator := operators[i]
			processProductCondition(rule, operator)
		}
	}

}

func processProductCondition(rule models.EmRuleModel, operator string) {
	ruleId := rule.Id
	propertyName := utils.GetConditionVariablePropertyId(operator)
	deviceLabel := utils.GetConditionVariableObjectId(operator)
	deviceList := repositories.NewDevicePointRepository().GetDeviceByDevLabel(deviceLabel)
	s := strings.Split(propertyName, "_")
	propertyType := s[0]
	propertyCode := s[1]
	for _, device := range deviceList {
		key := "device." + strconv.Itoa(device.Id) + "_" + propertyType + "_" + propertyCode
		PropertyRuleMap[key] = append(PropertyRuleMap[key], ruleId)
	}
}
