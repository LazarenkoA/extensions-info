package main

import (
	"context"
	"flag"
	onec "github.com/LazarenkoA/extensions-info/internal/1c"
	"github.com/LazarenkoA/extensions-info/internal/app"
	"github.com/LazarenkoA/extensions-info/internal/config"
	"github.com/LazarenkoA/extensions-info/internal/repository"
	"github.com/LazarenkoA/extensions-info/internal/repository/redis"
	"github.com/LazarenkoA/extensions-info/internal/usecase/app_settings"
	"github.com/LazarenkoA/extensions-info/internal/usecase/configuration"
	"github.com/LazarenkoA/extensions-info/internal/usecase/databases"
	"github.com/LazarenkoA/extensions-info/internal/usecase/health"
	"github.com/LazarenkoA/extensions-info/internal/usecase/jobs"
	"github.com/LazarenkoA/extensions-info/internal/usecase/notify"
	ws_conn "github.com/LazarenkoA/extensions-info/internal/ws"
	"github.com/joho/godotenv"
	"github.com/samber/lo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "", "Порт для прослушивания")
	flag.Parse()

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

	r, err := redis.New(ctx, cfg.RedisHost)
	if err != nil {
		log.Fatalf("redis connect error: %v", err)
	}

	ws := ws_conn.New()
	notify.Listener(ctx, r, ws)

	baseSettings := databases.New(repo)
	conf := configuration.New(repo)
	appSettings := app_settings.New(repo)
	analyzer := onec.NewAnalyzer1C(repo)
	job := jobs.New(repo, analyzer, r)
	h := health.New()

	cfg.Port = lo.If(port != "", port).Else(cfg.Port) // если задан порт как параметр cli то он "перекрывает" значение из env
	mainApp := app.NewExtensionsInfo(cfg)

	go shutdown(cancel)
	mainApp.Run(ctx, baseSettings, conf, appSettings, job, ws, h)
}

func shutdown(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("shutting down")
	cancel()
}
