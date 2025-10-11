package main

import (
	"context"
	onec "github.com/LazarenkoA/extensions-info/internal/1c"
	"github.com/LazarenkoA/extensions-info/internal/app"
	"github.com/LazarenkoA/extensions-info/internal/config"
	"github.com/LazarenkoA/extensions-info/internal/repository"
	"github.com/LazarenkoA/extensions-info/internal/usecase/app_settings"
	"github.com/LazarenkoA/extensions-info/internal/usecase/configuration"
	"github.com/LazarenkoA/extensions-info/internal/usecase/databases"
	"github.com/LazarenkoA/extensions-info/internal/usecase/jobs"
	ws_conn "github.com/LazarenkoA/extensions-info/internal/ws"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("can't load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	repo, err := repository.NewPG(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("postgree connect error: %v", err)
	}

	ws := ws_conn.New()
	baseSettings := databases.New(repo)
	conf := configuration.New(repo)
	appSettings := app_settings.New(repo)
	analyzer := onec.NewAnalyzer1C(repo)
	job := jobs.New(repo, analyzer, ws)

	mainApp := app.NewExtensionsInfo(cfg)

	go shutdown(cancel)
	mainApp.Run(ctx, baseSettings, conf, appSettings, job, ws)
}

func shutdown(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("shutting down")
	cancel()
}
