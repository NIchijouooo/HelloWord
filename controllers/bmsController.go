package controllers

import (
	"gateway/httpServer/model"
	repositories "gateway/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BmsController struct {
	hisRepo *repositories.HistoryDataRepository
}

func NewBmsController() *BmsController {
	return &BmsController{hisRepo: repositories.NewHistoryDataRepository()}
}
func (ctrl *BmsController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/bms/getYcLastByDeviceIdAndDict", ctrl.GetYcLastByDeviceIdAndDict)
	router.POST("/api/v2/bms/getYcLogById", ctrl.GetYcLogById)
}

// 获取设备点位最新的一条非空数据
func (c *BmsController) GetYcLastByDeviceIdAndDict(ctx *gin.Context) {
	//1.根据设备id去查询字典
	//2.将查询出来的所有code拼接成字符串
	deviceId := ctx.Value("deviceId")
	print(deviceId)
	codes := "1,2"
	//3.然后查询涛思数据库
	ycList, err := c.hisRepo.GetLastYcListByCode(codes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"获取信息成功！",
		ycList,
	})
}

//批量获取遥信息GetYcLogById,用于·历史数据
func (c *BmsController) GetYcLogById(ctx *gin.Context) {
	type ycQeury struct {
		DeviceId  int    `form:"DeviceId"`
		Codes     string `form:"codes"`
		StartTime int64  `form:"StartTime"`
		EndTime   int64  `form:"EndTime"`
	}
	var ycQuery ycQeury
	//将传过来的请求体解析到ycQuery中
	if err := ctx.ShouldBindJSON(&ycQuery); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	//查询历史数据
	ycLog, err := c.hisRepo.GetYcLogById(ycQuery.DeviceId, ycQuery.Codes, ycQuery.StartTime, ycQuery.EndTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"获取信息成功！",
		ycLog,
	})
}
