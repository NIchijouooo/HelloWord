package rule

import (
	repositories "gateway/repositories"
	"gateway/rule/dataProcess"
	"log"
)

func Init() {
	ruleList, err := repositories.NewRuleRepository().GetAllRule()
	if err != nil {
		log.Fatalln("GetAllRule err : ", err)
		return
	}
	if ruleList == nil || len(ruleList) == 0 {
		log.Fatalln("ruleList len is 0 ")
		return
	}
	for _, rule := range ruleList {
		//fmt.Println(rule)
		dataProcess.ProcessRuleVariable(rule)
	}

}
