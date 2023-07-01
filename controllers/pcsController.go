package controllers

import (
	"encoding/json"
	"fmt"
	"gateway/httpServer/model"
	"gateway/models"
	repositories "gateway/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PcsController struct {
	deviceEquipmentRepo *repositories.DeviceEquipmentRepository
	dictDataRepo        *repositories.DictDataRepository
	realtimeRepo        *repositories.RealtimeDataRepository
	deviceRepo          *repositories.EmRepository
}

func NewPcsController() *PcsController {
	return &PcsController{
		deviceEquipmentRepo: repositories.NewDeviceEquipmentRepository(),
		dictDataRepo:        repositories.NewDictDataRepository(),
		realtimeRepo:        repositories.NewRealtimeDataRepository(),
		deviceRepo:          repositories.NewEmRepository(),
	}
}
func (ctrl *PcsController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/pcs/getPcsDeviceInfoById", ctrl.getPcsDeviceInfoById)
}

/*
*
根据设备id获取pcs设备信息
*/
func (ctrl *PcsController) getPcsDeviceInfoById(ctx *gin.Context) {
	var param models.PcsParam
	if err := ctx.Bind(&param); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	deviceId := param.DeviceId
	if deviceId == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "参数错误",
			Data:    "",
		})
		return
	}
	dev, err := ctrl.deviceRepo.GetEmDeviceById(deviceId)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	dictData, _ := ctrl.dictDataRepo.SelectDictValue("energy_product_code_setting", "energy_storage_pcs_running_state")
	pcsRunStatusData, _ := ctrl.dictDataRepo.SelectDictValue("pcs_yc_device_status", dev.DeviceType)
	if len(dictData.DictValue) == 0 || len(pcsRunStatusData.DictValue) == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "设备状态字典获取失败",
			Data:    "",
		})
		return
	}
	info := models.PcsDeviceInfo{}
	// 设备状态遥测编码
	deviceStatusCode, _ := strconv.Atoi(dictData.DictValue)
	ycData, _ := ctrl.realtimeRepo.GetYcById(deviceId, deviceStatusCode)
	status := "--"
	if ycData.DeviceId > 0 {
		status = fmt.Sprintf("%v", ycData.Value)
		statusMap := map[string]string{}
		err := json.Unmarshal([]byte(pcsRunStatusData.DictValue), &statusMap)
		if err == nil {
			statusStr, ok := statusMap[status]
			if ok {
				status = statusStr
			}
		}
	}
	info.DeviceStatus = status
	equipmentInfo, _ := ctrl.deviceEquipmentRepo.GetEquipmentInfoByDevId(deviceId)
	info.EquipmentInfo = equipmentInfo
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "",
		Data:    info,
	})
}
