package dataProcess

import (
	"gateway/models"
	"gateway/rule/operation"
	"strings"
)

func ProcessRule(rule models.EmRuleVo) {
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

		if previous == "then" {
			action = append(action, str)
		}

	}
	if ruleType == "logic" {
		if !ProcessRecover(action, rule, condition.String()) {
			result := operation.GetResult(rule, condition.String())
			if result == "true" {
				ProcessAction(action, rule, condition.String())
			}
		}
	}
}
