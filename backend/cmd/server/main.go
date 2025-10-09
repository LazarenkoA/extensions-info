package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
	onec "your-app/internal/1c"
	"your-app/internal/app"
	"your-app/internal/config"
	"your-app/internal/repository"
	"your-app/internal/usecase/app_settings"
	"your-app/internal/usecase/configuration"
	"your-app/internal/usecase/databases"
	"your-app/internal/usecase/jobs"
	ws_conn "your-app/internal/ws"
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

	log := analyzer.RunAnalyzing(ctx, 1) //todo
	for l := range log {
		_ = l
	}

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
