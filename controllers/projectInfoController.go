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
	router.POST("/api/v2/project/createProject", ctrl.CreateProject)
	router.POST("/api/v2/project/getProjectByID", ctrl.GetProjectByID)
	router.GET("/api/v2/project/getAllProjects", ctrl.GetAllProjects)
	router.POST("/api/v2/project/updateProject", ctrl.UpdateProject)
	router.DELETE("/api/v2/project/deleteProject", ctrl.DeleteProject)
	// 注册其他路由...
}

func (c *ProjectInfoController) CreateProject(ctx *gin.Context) {
	var project models.ProjectInfo
	if err := ctx.ShouldBindJSON(&project); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"0",
			"error" + err.Error(),
			"",
		})
		return
	}
	if len(project.Name) > 300 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"项目名输入长度不能大于100",
			"",
		})
		return
	}

	if len(project.Address) > 600 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"地址输入长度不能大于200",
			"",
		})
		return
	}

	if err := c.repo.Create(&project); err != nil {
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
		project,
	})
}

func (c *ProjectInfoController) GetProjectByID(ctx *gin.Context) {
	var project models.ProjectInfo
	if err := ctx.Bind(&project); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	project, err := c.repo.GetById(project.ID)
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
		project,
	})
}

func (c *ProjectInfoController) GetAllProjects(ctx *gin.Context) {
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

func (c *ProjectInfoController) UpdateProject(ctx *gin.Context) {
	var project models.ProjectInfo
	if err := ctx.ShouldBindJSON(&project); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"0",
			"error" + err.Error(),
			"",
		})
		return
	}
	if len(project.Name) > 300 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"项目名输入长度不能大于100",
			"",
		})
		return
	}

	if len(project.Address) > 600 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"地址输入长度不能大于200",
			"",
		})
		return
	}

	if err := c.repo.Update(&project); err != nil {
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
		project,
	})
}

func (c *ProjectInfoController) DeleteProject(ctx *gin.Context) {
	var project models.ProjectInfo
	if err := ctx.Bind(&project); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	if err := c.repo.Delete(project.ID); err != nil {
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
