package contorl

import (
	"net/http"
	"gateway/httpServer/model"

	"gateway/setting"

	"github.com/gin-gonic/gin"
)

func ApiGetSerial(context *gin.Context) {
	type SerialPortNameTemplate struct {
		Name string `json:"Name"`
	}

	data := make([]SerialPortNameTemplate, 0)

	SerialPortName := SerialPortNameTemplate{}
	for _, v := range setting.SerialPortNameTemplateMap.Name {
		SerialPortName.Name = v
		data = append(data, SerialPortName)
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "",
		Data:    data,
	})
}
