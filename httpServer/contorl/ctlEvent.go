package contorl

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gateway/event"
	"gateway/httpServer/model"
	"gateway/setting"
	"sort"
)

func ApiGetReportEvent(context *gin.Context) {
	events := event.GetReportEvents()

	//排序，方便前端页面显示
	sort.Slice(events, func(i, j int) bool {
		iName := events[i].ID
		jName := events[j].ID
		return iName < jName
	})

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "",
		Data:    events,
	})
}

func ApiDeleteReportEvent(context *gin.Context) {

	reqParam := struct {
		ID []uint32
	}{}

	err := context.BindJSON(&reqParam)
	if err != nil {
		setting.ZAPS.Errorf("ReportEvent json解析错误 %v", err)
		context.JSON(http.StatusOK, gin.H{
			"Code":    "1",
			"Message": "json unMarshall err",
			"Data":    "",
		})
		return
	}

	for _, r := range reqParam.ID {
		_, ok := event.ReportEvents[r]
		if !ok {
			continue
		}
		_ = event.DeleteReportEvents(r)
	}

	context.JSON(http.StatusOK, gin.H{
		"Code":    "0",
		"Message": "",
		"Data":    "",
	})
}
