package dataProcess

import (
	"fmt"
	"gateway/models"
	"gateway/repositories"
	"gateway/rule/operation"
	"log"
	"strings"
)

func ProcessRecover(action []string, rule models.EmRuleModel, condition string) bool {

	fmt.Println("processRecover start ; action :", action, " , rule : ", rule, " , condition : ", condition)

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
			fmt.Println("getLastRuleHistoryByRuleIdAndDeviceId start ; rule.Id :", rule.Id, " , realTimeDataJson.DeviceId : ", realTimeDataJson.DeviceId, " , ruleHistory : ", ruleHistory)
			ruleHistory = history
		}
	}

	if ruleHistory != nil && ruleHistory.Tag == 0 {
		if len(action) > 0 {
			for _, s := range action {
				trim := strings.TrimSpace(s)
				if strings.Contains(trim, "record.recover.notMeet") {
					fmt.Println("start record.recover.notMeet ruleId :", rule.Id)
					// 不满足恢复
					result := operation.GetResult(rule, condition)
					fmt.Println("start record.recover.notMeet result :", result)
					if result != "" && result != "true" {
						ruleHistory.Tag = 1
						repositories.NewRuleHistoryRepository().UpdateRuleHistoryTag(ruleHistory)
					}
				}
				if strings.Contains(trim, "record.recover.condition") {
					fmt.Println("start record.recover.condition ruleId :", rule.Id)
					// 条件恢复
					result := operation.GetResult(rule, trim[25:])
					fmt.Println("start record.recover.condition result :", result)
					if result == "true" {
						ruleHistory.Tag = 1
						repositories.NewRuleHistoryRepository().UpdateRuleHistoryTag(ruleHistory)
						fmt.Println("processRecover start ; return false , ruleId :", rule.Id)
						return false
					}
				}
			}
		}
		fmt.Println("processRecover start ; return true , ruleId :", rule.Id)
		return true
	}

	fmt.Println("processRecover start ; return false , ruleId :", rule.Id)
	return false
}
