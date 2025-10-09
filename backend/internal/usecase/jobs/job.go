package jobs

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"net/http"
	"strconv"
	"your-app/internal/models"
	"your-app/internal/usecase"
)

type repo interface {
	GetCronSettings(ctx context.Context) (*models.CRONInfo, error)
	SetSchedule(ctx context.Context, bdID int32, schedule string) error
	DeleteSchedule(ctx context.Context, bdID int32) error
}

type IAnalyzer interface {
	RunAnalyzing(ctx context.Context, dbID int32) chan string
}

type IWS interface {
	Write(msg string) error
}

type Job struct {
	repo     repo
	analyzer IAnalyzer
	ws       IWS
}

const (
	JobStateNew        = "new"
	JobStateInProgress = "in_progress"
	JobStateDone       = "done"
	JobStateError      = "error"
)

func New(repo repo, analyzer IAnalyzer, ws IWS) *Job {
	return &Job{
		repo:     repo,
		analyzer: analyzer,
		ws:       ws,
	}
}

func (j *Job) Register(route *gin.RouterGroup) {
	route.GET("/getCronSettings", j.getCronSettings)
	route.POST("/setCronSettings/:db_id", j.setCronSettings)
	route.DELETE("/deleteCronSettings/:db_id", j.deleteCronSettings)
	route.POST("/startManualAnalysis/:db_id", j.startAnalysis)

	go j.runJobs()
}

func (j *Job) runJobs() {

}

func (j *Job) getCronSettings(ctx *gin.Context) {
	data, err := j.repo.GetCronSettings(ctx)
	if err != nil {
		usecase.ResponseError(ctx, err)
		return
	}

	data.NextCheckAsString = lo.If(data.NextCheck.IsZero(), "").Else(data.NextCheck.Format("02-01-2006 15:04"))

	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (j *Job) setCronSettings(ctx *gin.Context) {
	idStr := ctx.Param("db_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		usecase.ResponseError(ctx, errors.Wrap(err, "bad db_id"))
		return
	}

	var data models.CRONInfo
	if err := ctx.ShouldBindJSON(&data); err != nil {
		usecase.ResponseError(ctx, err)
		return
	}

	err = j.repo.SetSchedule(ctx, int32(id), data.Schedule)
	if err != nil {
		usecase.ResponseError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (j *Job) deleteCronSettings(ctx *gin.Context) {
	idStr := ctx.Param("db_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		usecase.ResponseError(ctx, errors.Wrap(err, "bad db_id"))
		return
	}

	err = j.repo.DeleteSchedule(ctx, int32(id))
	if err != nil {
		usecase.ResponseError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (j *Job) startAnalysis(ctx *gin.Context) {
	idStr := ctx.Param("db_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		usecase.ResponseError(ctx, errors.Wrap(err, "bad db_id"))
		return
	}

	logs := j.analyzer.RunAnalyzing(context.Background(), int32(id))

	go func() {
		for log := range logs {
			j.ws.Write(log)
		}
	}()

	ctx.JSON(http.StatusOK, gin.H{})
}
