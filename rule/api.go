package rule

import (
	"fmt"
	"gateway/models"
	"gateway/rule/dataProcess"
	"log"
)

func ProcessingSignalMsg(realTimeDataJsonList []models.RealTimeDataJson) {
	if len(realTimeDataJsonList) == 0 {
		log.Println("processingSignalMsg  signalInfos size = 0")
		return
	}

	var ruleList []models.EmRuleVo

	for _, realTimeDataJson := range realTimeDataJsonList {
		typeValue := realTimeDataJson.Type
		code := realTimeDataJson.Code
		deviceId := realTimeDataJson.DeviceId
		strType := "4"
		fmt.Println("typeValue : ", typeValue)
		if typeValue != 0 {
			strType = "0"
		}
		mapKey := fmt.Sprintf("device.%d_%s_%d", deviceId, strType, code)
		ruleIdList := dataProcess.PropertyRuleMap[mapKey]
		if ruleIdList == nil {
			log.Println("ProcessingSignalMsg  ruleIdList = nil ; realTimeDataJson:", realTimeDataJson, ", mapKey:", mapKey)
			continue
		}
		for _, ruleId := range ruleIdList {
			rule := dataProcess.RuleMap[ruleId]
			if &rule == nil {
				log.Println("ProcessingSignalMsg  rule = nil ; realTimeDataJson:", realTimeDataJson)
				continue
			}

			ruleVo := models.EmRuleVo{
				Id:                 rule.Id,
				Name:               rule.Name,
				Content:            rule.Content,
				Description:        rule.Description,
				Level:              rule.Level,
				EnableFlag:         rule.EnableFlag,
				TypeClassification: rule.TypeClassification,
			}
			ruleVo.RealTimeDataJson = &models.RealTimeDataJson{
				Type:     realTimeDataJson.Type,
				Code:     realTimeDataJson.Code,
				Value:    realTimeDataJson.Value,
				DeviceId: realTimeDataJson.DeviceId,
			}
			if !containsRule(ruleList, ruleVo) {
				ruleList = append(ruleList, ruleVo)
			}
		}
	}
	if len(ruleList) > 0 {
		for _, vo := range ruleList {
			dataProcess.ProcessRule(vo)
		}
	}
}

// 判断 ruleList 是否包含 rule
func containsRule(ruleList []models.EmRuleVo, rule models.EmRuleVo) bool {
	for _, r := range ruleList {
		if r == rule {
			return true
		}
	}
	return false
}
