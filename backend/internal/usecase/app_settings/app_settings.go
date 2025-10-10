package app_settings

import (
	"context"
	"github.com/LazarenkoA/extensions-info/internal/models"
	"github.com/LazarenkoA/extensions-info/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type repo interface {
	GetAppSettings(ctx context.Context) (*models.AppSettings, error)
	SetAppSettings(ctx context.Context, id int32, settings models.AppSettings) error
}

type AppSettings struct {
	repo repo
}

func New(repo repo) *AppSettings {
	return &AppSettings{
		repo: repo,
	}
}

func (a *AppSettings) Register(route *gin.RouterGroup) {
	route.GET("/getAppSettings", a.getAppSettings)
	route.POST("/appSettings/:id", a.appSettings)
}

func (a *AppSettings) appSettings(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	var data models.AppSettings
	if err := ctx.ShouldBindJSON(&data); err != nil {
		usecase.ResponseError(ctx, err)
		return
	}

	err := a.repo.SetAppSettings(ctx, int32(id), data)
	if err != nil {
		usecase.ResponseError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (a *AppSettings) getAppSettings(ctx *gin.Context) {
	data, err := a.repo.GetAppSettings(ctx)
	if err != nil {
		usecase.ResponseError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
