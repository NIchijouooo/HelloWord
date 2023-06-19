package controllers

import (
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
}

func NewPcsController() *PcsController {
	return &PcsController{
		deviceEquipmentRepo: repositories.NewDeviceEquipmentRepository(),
		dictDataRepo:        repositories.NewDictDataRepository(),
		realtimeRepo:        repositories.NewRealtimeDataRepository(),
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
	dictData, _ := ctrl.dictDataRepo.SelectDictValue("em_dict_info", "pcs_device_status_yx_code")
	if len(dictData.DictValue) == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "设备状态字典获取失败",
			Data:    "",
		})
		return
	}
	info := models.PcsDeviceInfo{}
	// 设备状态遥信编码
	deviceStatusCode, _ := strconv.Atoi(dictData.DictValue)
	yxData, _ := ctrl.realtimeRepo.GetYxById(deviceId, deviceStatusCode)
	info.DeviceStatus = yxData.Value
	equipmentInfo, _ := ctrl.deviceEquipmentRepo.GetEquipmentInfoByDevId(deviceId)
	info.EquipmentInfo = equipmentInfo
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "",
		Data:    info,
	})
}
