package databases

import (
	"context"
	"github.com/LazarenkoA/extensions-info/internal/models"
	"github.com/LazarenkoA/extensions-info/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"net/http"
	"strconv"
)

type repo interface {
	GetDataBaseSettings(ctx context.Context) ([]models.DatabaseSettings, error)
	AddDataBase(ctx context.Context, data models.DatabaseSettings) error
	DeleteDataBase(ctx context.Context, id int32) error
}

type Settings struct {
	repo repo
}

const (
	DbStateNew       = "new"
	DbStateAnalyzing = "analyzing"
	DbStateDone      = "done"
	DbStateError     = "error"
)

func New(repo repo) *Settings {
	return &Settings{
		repo: repo,
	}
}

func (s *Settings) Register(route *gin.RouterGroup) {
	route.GET("/getBaseSettings", s.getSettings)
	route.POST("/addBaseSettings", s.addDataBase)
	route.DELETE("/deleteBaseSettings/:id", s.deleteDataBase)
}

func (s *Settings) getSettings(ctx *gin.Context) {
	data, err := s.repo.GetDataBaseSettings(ctx)
	if err != nil {
		usecase.ResponseError(ctx, err)
		return
	}

	for i, base := range data {
		if base.LastCheck != nil {
			data[i].LastCheckAsString = lo.If(base.LastCheck.IsZero(), "").Else(base.LastCheck.Format("02-01-2006 15:04"))
		}
		if data[i].Cron != nil {
			data[i].Cron.NextCheckAsString = lo.If(base.Cron.NextCheck.IsZero(), "").Else(base.Cron.NextCheck.Format("02-01-2006 15:04"))
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (s *Settings) addDataBase(ctx *gin.Context) {
	var data models.DatabaseSettings

	if err := ctx.ShouldBindJSON(&data); err != nil {
		usecase.ResponseError(ctx, err)
		return
	}

	err := s.repo.AddDataBase(ctx, data)
	if err != nil {
		usecase.ResponseError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
	//c.Status(http.StatusCreated)
}

func (s *Settings) deleteDataBase(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		usecase.ResponseError(ctx, errors.Wrap(err, "bad id"))
		return
	}

	err = s.repo.DeleteDataBase(ctx, int32(id))
	if err != nil {
		usecase.ResponseError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
