package onec

import (
	"context"
	"encoding/json"
	"github.com/LazarenkoA/extensions-info/internal/models"
	"github.com/LazarenkoA/extensions-info/internal/usecase/databases"
	"github.com/LazarenkoA/extensions-info/internal/utils"
	"github.com/pkg/errors"
	"os"
	"time"
)

type logWriter func(msgType string, msg interface{})

type repo interface {
	SetBDState(ctx context.Context, bdID int32, newState string, lastCheck time.Time) error
	GetDataBaseByID(ctx context.Context, id int32) (*models.DatabaseSettings, error)
	GetAppSettings(ctx context.Context) (*models.AppSettings, error)
	StoreConfigurationInfo(ctx context.Context, dbID int32, confInfo *ConfigurationInfo) (int32, error)
	StoreExtensionsInfo(ctx context.Context, confID int32, confInfo []ConfigurationInfo) error
	GetExtensionsInfo(ctx context.Context, confID int32) ([]ConfigurationInfo, error)
	SetMetadata(ctx context.Context, confID int32, value []byte) error
	GetMetadata(ctx context.Context, confID int32) ([]byte, error)
	GetChildObjectsConf(ctx context.Context, confID int32) (*ConfigurationStruct, error)
	SetCode(ctx context.Context, extID int32, key, code string) error
}

type Analyzer1C struct {
	repo repo
}

func NewAnalyzer1C(repo repo) *Analyzer1C {
	return &Analyzer1C{
		repo: repo,
	}
}

func (a *Analyzer1C) RunAnalyzing(ctx context.Context, dbID int32) chan string {
	logC := make(chan string)
	log := formatLog(int(dbID), logC)

	step := &steps{log: log, state: make(map[string]any)}

	go func() {
		defer close(logC)

		err := a.stateAnalyzing(ctx, dbID, log)
		if err != nil {
			log("error", err.Error())
			return
		}

		step.add("Получение информации о конфигурации", func() error {
			confID, extDir, err := a.confAnalyzing(ctx, dbID)
			step.state["extDir"] = extDir
			step.state["confID"] = confID
			step.destruct(func() {
				os.RemoveAll(extDir)
			})

			return err
		})
		step.add("Анализ метаданных расширений", func() error {
			extDir, ok := step.state["extDir"]
			if !ok {
				return errors.New("extensions dir is not defined")
			}
			confID, ok := step.state["confID"]
			if !ok {
				return errors.New("confID is not defined")
			}

			return a.metadataAnalyzing(extDir.(string), confID.(int32))
		})
		step.add("Анализ кода расширений", func() error {
			confID, ok := step.state["confID"]
			if !ok {
				return errors.New("confID is not defined")
			}

			return a.codeAnalyzing(confID.(int32))
		})
		err = step.run()

		if err != nil {
			log("log", errors.Wrap(err, "ERROR").Error())
			_ = a.stateError(ctx, dbID, log)
		} else {
			_ = a.stateDone(ctx, dbID, log)
		}
	}()

	return logC
}

func (a *Analyzer1C) confAnalyzing(ctx context.Context, dbID int32) (int32, string, error) {
	info, err := a.repo.GetDataBaseByID(ctx, dbID)
	if err != nil {
		return 0, "", errors.Wrap(err, "get dataBase by id")
	}

	appSettings, err := a.repo.GetAppSettings(ctx)
	if err != nil {
		return 0, "", errors.Wrap(err, "get application settings")
	}

	if appSettings.PlatformPath == "" {
		return 0, "", errors.New("platform bin path is not defined")
	}

	confInfo, err := loadConfigurationInfo(appSettings.PlatformPath, info.ConnectionString, utils.PtrToVal(info.Username), utils.PtrToVal(info.Password))
	if err != nil {
		return 0, "", err
	}

	confID, err := a.repo.StoreConfigurationInfo(ctx, dbID, confInfo)
	if err != nil {
		return 0, "", errors.Wrap(err, "store error")
	}

	extDir, extInfo, err := loadExtensionsInfo(appSettings.PlatformPath, info.ConnectionString, utils.PtrToVal(info.Username), utils.PtrToVal(info.Password))
	if err != nil {
		return 0, "", errors.Wrap(err, "store error")
	}

	err = a.repo.StoreExtensionsInfo(ctx, confID, extInfo)
	if err != nil {
		return 0, "", errors.Wrap(err, "store error")
	}

	return confID, extDir, nil
}

func (a *Analyzer1C) stateDone(ctx context.Context, dbID int32, log logWriter) (err error) {
	if err = a.repo.SetBDState(ctx, dbID, databases.DbStateDone, time.Now()); err == nil {
		log("new_state", databases.DbStateDone)
	}
	return err
}

func (a *Analyzer1C) stateAnalyzing(ctx context.Context, dbID int32, log logWriter) (err error) {
	if err = a.repo.SetBDState(ctx, dbID, databases.DbStateAnalyzing, time.Time{}); err == nil {
		log("new_state", databases.DbStateAnalyzing)
	}
	return err
}

func (a *Analyzer1C) stateError(ctx context.Context, dbID int32, log logWriter) (err error) {
	if err = a.repo.SetBDState(ctx, dbID, databases.DbStateError, time.Time{}); err == nil {
		log("new_state", databases.DbStateError)
	}
	return err
}

func formatLog(dbID int, log chan string) logWriter {
	return func(msgType string, msg interface{}) {
		tmp := map[string]interface{}{
			"db_id": dbID,
			"type":  msgType,
			"msg":   msg,
			"time":  time.Now().Format(time.TimeOnly),
		}

		bdata, _ := json.Marshal(tmp)
		log <- string(bdata)
	}
}
