package dataProcess

import (
	"gateway/models"
	"gateway/repositories"
	"gateway/rule/utils"
	"log"
	"strings"
	"time"
)

func ProcessAction(actionList []string, rule models.EmRuleVo, ruleType string) {
	if actionList == nil || len(actionList) == 0 {
		log.Println("actionList is null")
		return
	}

	var linkage, record, debance, assign []string

	for _, action := range actionList {
		substring := action[:strings.Index(action, ".")]
		switch substring {
		case "linkage":
			linkage = append(linkage, action)
		case "record":
			record = append(record, action)
		case "debance":
			debance = append(debance, action)
		case "assign":
			assign = append(assign, action)
		default:
			break
		}
	}

	if len(record) > 0 {
		// 生成历史、报警
		insertHistory(record, rule)
	}
}

func insertHistory(actionList []string, rule models.EmRuleVo) {
	ruleHistory := models.EmRuleHistoryModel{
		EventId:            rule.Id,
		EventName:          rule.Name,
		Description:        "[" + rule.Name + "]规则触发",
		Level:              rule.Level,
		TypeClassification: rule.TypeClassification,
	}

	outputList := make([]string, 0)

	for _, action := range actionList {
		if strings.Contains(action, "record.name") {
			split := strings.Split(action, "=")
			if len(split) >= 2 {
				ruleHistory.EventName = split[1]
			}
		}
		if strings.Contains(action, "record.content") {
			split := strings.Split(action, "=")
			if len(split) >= 2 {
				ruleHistory.Description = split[1]
			}
		}
		if strings.Contains(action, "record.output") {
			outputList = append(outputList, action)
		}
	}

	date := time.Now().Format("2006-01-02 15:04:05")
	ruleHistory.ProduceTime = date
	ruleHistory.CreateTime = date
	ruleHistory.UpdateTime = date

	insert, err := repositories.NewRuleHistoryRepository().InsertRuleHistory(&ruleHistory)
	if err != nil {
		log.Fatalln("getPropertyValueData err : ", err)
		return
	}
	if insert > 0 {
		content := rule.Content
		split := strings.Split(content, "\\n")
		condition := strings.Builder{}
		previous := ""
		typeStr := ""
		for _, s := range split {
			str := strings.TrimSpace(s)
			if strings.Contains(str, "type=") {
				typeStr = str[5:10]
				continue
			}
			if str == "if" || str == "then" {
				previous = str
				continue
			}
			if previous == "if" && typeStr == "logic" {
				condition.WriteString(str)
			}
		}

		realTimeDataJson := rule.RealTimeDataJson
		if condition.Len() > 0 {
			deviceMap := make(map[int]int, 0)
			operators := utils.GetOperator(condition.String())
			length := len(operators)
			for i := 0; i < length; i++ {
				operator := operators[i]
				var deviceId int
				if strings.Contains(operator, "product.${") {
					if realTimeDataJson != nil {
						deviceId = realTimeDataJson.DeviceId
						deviceMap[deviceId] = realTimeDataJson.Code
					}
				}
			}
			if len(deviceMap) > 0 {
				ruleHistoryDevice := make([]models.EmRuleHistoryDeviceModel, 0)
				for key, value := range deviceMap {
					ruleHistoryDeviceItem := models.EmRuleHistoryDeviceModel{
						EventHistoryId: ruleHistory.Id,
						DeviceId:       key,
						PropertyCode:   value,
						CreateTime:     time.DateTime,
					}
					ruleHistoryDevice = append(ruleHistoryDevice, ruleHistoryDeviceItem)
				}
				repositories.NewRuleHistoryRepository().InsertRuleHistoryDevice(&ruleHistoryDevice)
			}
		}
	}
}
