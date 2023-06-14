package dataProcess

import (
	"fmt"
	"gateway/models"
	"gateway/rule/utils"
	"strings"
)

var DeviceRuleMap = make(map[string]models.EmRuleModel)
var RuleMap map[int]models.EmRuleModel
var PropertyRuleMap map[string][]int
var SystemEventRuleMap map[string][]int

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

	fmt.Println("condition : ", condition.String())
	processLogicRule(rule, condition.String(), action)
	RuleMap[rule.Id] = rule
}

func processLogicRule(rule models.EmRuleModel, condition string, action []string) {
	ruleId := rule.Id
	if len(condition) > 0 {
		operators := utils.GetOperator(condition)
		fmt.Println("operators : ", operators)
		length := len(operators)
		for i := 0; i < length; i++ {
			operator := operators[i]
			if strings.Contains(operator, "device.${") || strings.Contains(operator, "event.${") {
				processDeviceAndEventCondition(ruleId, operator)
			} else if strings.Contains(operator, "systemEvent.${") {
				processSystemEventCondition(ruleId, operator)
			} else if strings.Contains(operator, "product.${") {
				processProductCondition(rule, operator)
			}
		}
	}

	if len(action) > 0 {
		//processRuleAction(ruleId, action)
	}
}

func processProductCondition(rule models.EmRuleModel, operator string) {
	//ruleId := rule.Id
	propertyName := utils.GetConditionVariablePropertyId(operator)
	fmt.Println("propertyName ： ", propertyName)
	deviceLabel := utils.GetConditionVariableObjectId(operator)
	fmt.Println("deviceLabel ： ", deviceLabel)
	//label := repositories.NewDevicePointRepository().GetDeviceByDevLabel(deviceLabel)
	//for _, deviceId := range deviceIdList {
	//	key := "device." + strconv.Itoa(deviceId) + "_" + propertyName
	//	PropertyRuleMap[key] = append(PropertyRuleMap[key], ruleId)
	//}
}

func processSystemEventCondition(ruleId int, operator string) {
	identifier := utils.GetConditionVariableEventIdentifier(operator)
	key := "systemEvent." + identifier
	SystemEventRuleMap[key] = append(SystemEventRuleMap[key], ruleId)
}

func processDeviceAndEventCondition(ruleId int, operator string) {
	propertyId := utils.GetConditionVariablePropertyId(operator)
	var key string
	if propertyId == "status" {
		objectId := utils.GetConditionVariableObjectId(operator)
		key = ""
		if strings.Contains(operator, "device.${") {
			key = "device"
		} else if strings.Contains(operator, "event.${") {
			key = "event"
		}
		key += "." + objectId + "." + propertyId
	} else {
		key = ""
		if strings.Contains(operator, "device.${") {
			key = "device"
		} else if strings.Contains(operator, "event.${") {
			key = "event"
		}
		key += "." + propertyId
	}

	PropertyRuleMap[key] = append(PropertyRuleMap[key], ruleId)
}

func processRuleAction(ruleId int, action []string) {
	for _, str := range action {
		operatorSplit := strings.Split(str, "}=")
		if len(operatorSplit) < 2 {
			continue
		}
		variable := strings.TrimSpace(operatorSplit[1])
		if variable == "" {
			continue
		}
		operators := utils.GetOperator(variable)
		length := len(operators)
		for i := 0; i < length; i++ {
			operator := operators[i]
			if strings.Contains(operator, "device.${") {
				propertyId := utils.GetConditionVariablePropertyId(operator)
				key := "device." + propertyId
				PropertyRuleMap[key] = append(PropertyRuleMap[key], ruleId)
			}
		}
	}
}
