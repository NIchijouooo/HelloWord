package controllers

/**
20230605
*/
import (
	"gateway/httpServer/model"
	"gateway/models"
	repositories "gateway/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProjectInfoController struct {
	repo *repositories.ProjectInfoRepository
}

func NewProjectInfoController() *ProjectInfoController {
	return &ProjectInfoController{
		repo: repositories.NewProjectInfoRepository(),
	}
}

func (ctrl *ProjectInfoController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/procjet/createProcjet", ctrl.CreateProcjet)
	router.GET("/api/v2/procjet/getProcjetByID", ctrl.GetProcjetByID)
	router.GET("/api/v2/procjets/getAllProcjets", ctrl.GetAllProcjets)
	router.POST("/api/v2/procjet/updateProcjet", ctrl.UpdateProcjet)
	router.DELETE("/api/v2/procjet/deleteProcjet", ctrl.DeleteProcjet)
	// 注册其他路由...
}

func (c *ProjectInfoController) CreateProcjet(ctx *gin.Context) {
	var procjet models.ProjectInfo
	if err := ctx.ShouldBindJSON(&procjet); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"0",
			"error" + err.Error(),
			"",
		})
		return
	}
	if err := c.repo.Create(&procjet); err != nil {
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
		procjet,
	})
}

func (c *ProjectInfoController) GetProcjetByID(ctx *gin.Context) {
	var procjet models.ProjectInfo
	if err := ctx.Bind(&procjet); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	procjet, err := c.repo.GetById(procjet.ID)
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
		procjet,
	})
}

func (c *ProjectInfoController) GetAllProcjets(ctx *gin.Context) {
	projectName := ctx.Query("projectName")
	createTimeStart := ctx.Query("createTimeStart")
	createTimeEnd := ctx.Query("createTimeEnd")

	procjets, err := c.repo.GetAll(projectName, createTimeStart, createTimeEnd)
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
		procjets,
	})
}

func (c *ProjectInfoController) UpdateProcjet(ctx *gin.Context) {
	var procjet models.ProjectInfo
	if err := ctx.ShouldBindJSON(&procjet); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"0",
			"error" + err.Error(),
			"",
		})
		return
	}
	if err := c.repo.Update(&procjet); err != nil {
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
		procjet,
	})
}

func (c *ProjectInfoController) DeleteProcjet(ctx *gin.Context) {
	var procjet models.ProjectInfo
	if err := ctx.Bind(&procjet); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	if err := c.repo.Delete(procjet.ID); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"0",
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
