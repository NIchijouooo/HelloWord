package controllers

import (
	"gateway/httpServer/model"
	"gateway/models/query"
	"gateway/repositories"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/xuri/excelize/v2"
	"net/http"
	"strconv"
)

type LoadTrackingController struct {
	repo *repositories.LoadTrackingRepository
}

func NewLoadTrackingController() *LoadTrackingController {
	return &LoadTrackingController{
		repo: repositories.NewLoadTrackingRepository(),
	}
}

func (c *LoadTrackingController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/loadTracking/GetLoadTrackingData", c.GetLoadTrackingData)
	router.POST("/api/v2/loadTracking/ExportData", c.ExportData)
}

func (c *LoadTrackingController) ExportData(ctx *gin.Context) {
	ycQuery := query.QueryTaoData{}

	if err := ctx.ShouldBindBodyWith(&ycQuery, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ycQuery.StartTime == 0 && ycQuery.EndTime == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "参数错误",
			Data:    "",
		})
		return
	}

	if ycQuery.Interval == 0 {
		ycQuery.Interval = 30
	}
	if ycQuery.IntervalType == 0 {
		ycQuery.IntervalType = 2
	}

	result := repositories.GetLoadTrackingData(ycQuery)

	onGridData := result.OnGridData
	meterData := result.MeterData
	loadData := result.LoadData
	xAxisList := result.XAxisList

	f := excelize.NewFile() // 创建一个新的Excel文件

	f.SetCellValue("Sheet1", "A1", "时间")
	f.SetCellValue("Sheet1", "B1", "并网点功率")
	f.SetCellValue("Sheet1", "C1", "关口表功率")
	f.SetCellValue("Sheet1", "D1", "负载功率")
	for i, date := range xAxisList {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), date)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), onGridData[i])
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), meterData[i])
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), loadData[i])
	}

	// 将Excel文件写入响应体
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=example.xlsx")
	err := f.Write(ctx.Writer)
	if err != nil {
		return
	}
}

func (c *LoadTrackingController) GetLoadTrackingData(ctx *gin.Context) {

	ycQuery := query.QueryTaoData{}

	if err := ctx.ShouldBindBodyWith(&ycQuery, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ycQuery.StartTime == 0 && ycQuery.EndTime == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "参数错误",
			Data:    "",
		})
		return
	}
	if ycQuery.Interval == 0 {
		ycQuery.Interval = 30
	}
	if ycQuery.IntervalType == 0 {
		ycQuery.IntervalType = 2
	}

	result := repositories.GetLoadTrackingData(ycQuery)

	ctx.JSON(http.StatusOK, model.ResponseData{
		Code: "0",
		Data: result,
	})
}
