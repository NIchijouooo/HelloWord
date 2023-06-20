package dataProcess

import (
	"gateway/models"
	"gateway/repositories"
	"gateway/rule/operation"
	"log"
	"strings"
	"time"
)

func ProcessRecover(action []string, rule models.EmRuleVo, condition string) bool {

	// 处理事件恢复
	var ruleHistory *models.EmRuleHistoryModel
	content := rule.Content
	if strings.Contains(content, "product.${") {
		realTimeDataJson := rule.RealTimeDataJson
		if realTimeDataJson != nil && realTimeDataJson.DeviceId != 0 {
			history, err := repositories.NewRuleHistoryRepository().GetLastRuleHistory(rule.Id, realTimeDataJson.DeviceId)
			if err != nil {
				log.Fatalln("processRecover err : ", err)
				return true
			}
			ruleHistory = history
		}
	}
	if ruleHistory != nil && ruleHistory.Tag == 0 && ruleHistory.Id > 0 {
		if len(action) > 0 {
			for _, s := range action {
				trim := strings.TrimSpace(s)
				if strings.Contains(trim, "record.recover.notMeet") {
					// 不满足恢复
					result := operation.GetResult(rule, condition)
					if result != "" && result != "true" {
						ruleHistory.Tag = 1
						date := time.Now().Format(time.DateTime)
						ruleHistory.UpdateTime = date
						ruleHistory.RecoveryTime = date
						repositories.NewRuleHistoryRepository().UpdateRuleHistoryTag(ruleHistory)
					}
				}
				if strings.Contains(trim, "record.recover.condition") {
					// 条件恢复
					result := operation.GetResult(rule, trim[25:])
					if result == "true" {
						ruleHistory.Tag = 1
						repositories.NewRuleHistoryRepository().UpdateRuleHistoryTag(ruleHistory)
						return false
					}
				}
			}
		}
		return true
	}

	return false
}
