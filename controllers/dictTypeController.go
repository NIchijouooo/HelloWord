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

// 定义字典类型管理的控制器
type DictTypeController struct {
	repoType *repositories.DictTypeRepository
	repoData *repositories.DictDataRepository
}

func NewDictTypeController() *DictTypeController {
	return &DictTypeController{repoType: repositories.NewDictTypeRepository(), repoData: repositories.NewDictDataRepository()}
}

func (ctrl *DictTypeController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/dictType/createDictType", ctrl.CreateDictType)
	router.POST("/api/v2/dictType/updateDictType", ctrl.UpdateDictType)
	router.POST("/api/v2/dictType/deleteDictType", ctrl.DeleteDictType)
	router.POST("/api/v2/dictType/getDictTypeList", ctrl.GetDictTypeList)
	router.POST("/api/v2/dictType/getDictTypeByID", ctrl.GetDictTypeByID)
	router.POST("/api/v2/dictData/getDictTypeListByDictTypeId", ctrl.GetDictTypeListByDictTypeId)
	router.POST("/api/v2/dictData/getDictDataByTypeAndLabel", ctrl.GetDictDataByTypeAndLabel)
	// 注册其他路由...
}

type ParamType struct {
	DictId          int    `form:"dictId"`
	DictType        string `form:"dictType"`
	DictCode        int    `form:"dictCode"`
	PageNum         int    `form:"pageNum"`
	PageSize        int    `form:"pageSize"`
	DictName        string `form:"dictName"`
	CreateTimeStart string `form:"createTimeStart"`
	CreateTimeEnd   string `form:"createTimeEnd"`
}

// 新增字典类型
func (c *DictTypeController) CreateDictType(ctx *gin.Context) {
	var dictType models.DictType
	if err := ctx.ShouldBindJSON(&dictType); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	if err := c.repoType.Create(&dictType); err != nil {
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
		dictType,
	})

}

// 修改字典类型
func (c *DictTypeController) UpdateDictType(ctx *gin.Context) {
	var dictType models.DictType
	if err := ctx.ShouldBindJSON(&dictType); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	if err := c.repoType.Update(&dictType); err != nil {
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
		dictType,
	})
}

// 删除字典类型
func (c *DictTypeController) DeleteDictType(ctx *gin.Context) {
	var paramType ParamType
	if err := ctx.Bind(&paramType); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	if err := c.repoType.Delete(paramType.DictId); err != nil {
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

// 获取所有字典类型
func (c *DictTypeController) GetDictTypeList(ctx *gin.Context) {
	var (
		dictTypeList []models.DictType
	)
	var paramType ParamType
	if err := ctx.Bind(&paramType); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}

	dictTypeList, total, err := c.repoType.GetAll(paramType.DictName, paramType.CreateTimeStart, paramType.CreateTimeEnd, paramType.PageNum, paramType.PageSize)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}

	// 将查询结果返回给客户端
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		gin.H{
			"data":  dictTypeList,
			"total": total,
		},
	})

}

// 获取单个字典类型
func (c *DictTypeController) GetDictTypeByID(ctx *gin.Context) {
	var paramType ParamType
	if err := ctx.Bind(&paramType); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	dictType, err := c.repoType.GetById(paramType.DictId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		dictType,
	})
}

// 获取字典类型下的所有字典数据
func (c *DictTypeController) GetDictTypeListByDictTypeId(ctx *gin.Context) {
	var paramType ParamType
	if err := ctx.Bind(&paramType); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	dictDataList, err := c.repoType.GetDictDataListByDictTypeId(paramType.DictType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		dictDataList,
	})
}

// 获取字典类型下的所有字典数据
func (c *DictTypeController) GetDictDataByTypeAndLabel(ctx *gin.Context) {
	var paramType models.DictData
	if err := ctx.Bind(&paramType); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	dictData, err := c.repoData.SelectDictValue(paramType.DictType, paramType.DictLabel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		dictData,
	})
}
