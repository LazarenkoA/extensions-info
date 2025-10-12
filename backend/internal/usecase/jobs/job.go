package jobs

import (
	"context"
	"fmt"
	"github.com/LazarenkoA/extensions-info/internal/models"
	"github.com/LazarenkoA/extensions-info/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
	"net/http"
	"strconv"
	"sync"
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
	repo       repo
	analyzer   IAnalyzer
	ws         IWS
	cronJobs   *cron.Cron
	cronJobsID map[int]cron.EntryID
	mx         sync.Mutex
}

const (
	JobStateNew        = "new"
	JobStateInProgress = "in_progress"
	JobStateDone       = "done"
	JobStateError      = "error"
)

func New(repo repo, analyzer IAnalyzer, ws IWS) *Job {
	return &Job{
		repo:       repo,
		analyzer:   analyzer,
		ws:         ws,
		cronJobs:   cron.New(),
		cronJobsID: make(map[int]cron.EntryID),
	}
}

func (j *Job) Register(route *gin.RouterGroup) {
	route.GET("/getCronSettings", j.getCronSettings)
	route.POST("/setCronSettings/:db_id", j.setCronSettings)
	route.DELETE("/deleteCronSettings/:db_id", j.deleteCronSettings)
	route.POST("/startManualAnalysis/:db_id", j.startManualAnalysis)

	j.cronJobs.Start()
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

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	if _, err := parser.Parse(data.Schedule); err != nil {
		usecase.ResponseError(ctx, fmt.Errorf("%q bad cron format", data.Schedule))
		return
	}

	err = j.repo.SetSchedule(ctx, int32(id), data.Schedule)
	if err != nil {
		usecase.ResponseError(ctx, err)
		return
	}

	j.Start(id, data.Schedule)
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

	j.Stop(id) // todo запрос на отключение может попасть не на тот ПОД на котором работает джоба, поэтому нужно оповестить всех
	ctx.JSON(http.StatusOK, gin.H{})
}

func (j *Job) startManualAnalysis(ctx *gin.Context) {
	idStr := ctx.Param("db_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		usecase.ResponseError(ctx, errors.Wrap(err, "bad db_id"))
		return
	}

	j.startAnalysis(id)
	ctx.JSON(http.StatusOK, gin.H{})
}

func (j *Job) startAnalysis(dbID int) {
	logs := j.analyzer.RunAnalyzing(context.Background(), int32(dbID))

	go func() {
		for log := range logs {
			j.ws.Write(log)
		}
	}()
}

func (j *Job) Stop(dbID int) {
	j.mx.Lock()
	defer j.mx.Unlock()

	if jobID, ok := j.cronJobsID[dbID]; ok {
		j.cronJobs.Remove(jobID)
		delete(j.cronJobsID, dbID)
	}
}

func (j *Job) Start(dbID int, schedule string) {
	j.mx.Lock()
	defer j.mx.Unlock()

	id, _ := j.cronJobs.AddFunc(schedule, func() {
		j.startAnalysis(dbID)
	})

	j.cronJobsID[dbID] = id
}
