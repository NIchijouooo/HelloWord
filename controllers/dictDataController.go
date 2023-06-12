package controllers

/**
20230605
*/
import (
	"gateway/httpServer/model"
	"gateway/models"
	repositories "gateway/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义字典数据管理的控制器
type DictDataController struct {
	repo *repositories.DictDataRepository
}

func NewDictDataController() *DictDataController {
	return &DictDataController{repo: repositories.NewDictDataRepository()}
}

func (ctrl *DictDataController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/dictData/createDictData", ctrl.CreateDictData)
	router.POST("/api/v2/dictData/updateDictData", ctrl.UpdateDictData)
	router.POST("/api/v2/dictData/deleteDictData", ctrl.DeleteDictData)
	router.POST("/api/v2/dictData/getDictDataList", ctrl.GetDictDataList)
	router.POST("/api/v2/dictData/getDictDataByID", ctrl.GetDictDataByID)
	// 注册其他路由...
}

type ParamData struct {
	DictLabel string `form:"dictLabel"`
	DictType  string `form:"dictType"`
	DictCode  int    `form:"dictCode"`
	PageNum   int    `form:"pageNum"`
	PageSize  int    `form:"pageSize"`
}

// 新增字典数据
func (c *DictDataController) CreateDictData(ctx *gin.Context) {
	var dictData models.DictData
	if err := ctx.ShouldBindJSON(&dictData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	if err := c.repo.Create(&dictData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"ok",
		dictData,
	})
}

// 修改字典数据
func (c *DictDataController) UpdateDictData(ctx *gin.Context) {
	var dictData models.DictData
	if err := ctx.ShouldBindJSON(&dictData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	if err := c.repo.Update(&dictData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"ok",
		dictData,
	})
}

// 删除字典数据
func (c *DictDataController) DeleteDictData(ctx *gin.Context) {
	var paramData ParamData
	if err := ctx.Bind(&paramData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	if err := c.repo.Delete(paramData.DictCode); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"ok",
		"",
	})
}

// 获取所有字典数据
func (c *DictDataController) GetDictDataList(ctx *gin.Context) {
	var paramData ParamData
	if err := ctx.Bind(&paramData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}

	dictDataList, total, err := c.repo.GetAll(paramData.DictLabel, paramData.DictType, paramData.PageNum, paramData.PageSize)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		gin.H{
			"data":  dictDataList,
			"total": total,
		},
	})
}

// 获取单个字典数据
func (c *DictDataController) GetDictDataByID(ctx *gin.Context) {
	var paramData ParamData
	if err := ctx.Bind(&paramData); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	dictData, err := c.repo.GetById(paramData.DictCode)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"0",
			"error" + err.Error(),
			"",
		})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		dictData,
	})
}
