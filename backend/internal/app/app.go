package app

import (
	"context"
	"github.com/LazarenkoA/extensions-info/internal/config"
	"github.com/LazarenkoA/extensions-info/internal/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type UseCases interface {
	Register(route *gin.RouterGroup)
}

type ExtensionsInfo struct {
	cfg *config.Config
	srv *http.Server
}

func NewExtensionsInfo(cfg *config.Config) *ExtensionsInfo {
	return &ExtensionsInfo{
		cfg: cfg,
	}
}

func (e *ExtensionsInfo) Run(ctx context.Context, handlers ...UseCases) {
	r := gin.Default()

	// middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	api := r.Group("/api/v1")
	for _, h := range handlers {
		h.Register(api)
	}

	//api.GET("/getBaseSettings", handlers.GetUsers)
	////api.GET("/users/:id", handlers.GetUserByID)
	////api.POST("/users", handlers.CreateUser)
	////api.PUT("/users/:id", handlers.UpdateUser)
	////api.DELETE("/users/:id", handlers.DeleteUser)
	//
	//api.GET("/health", handlers.HealthCheck)

	e.srv = &http.Server{
		Addr:           ":" + e.cfg.Port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Сервер запущен на порту %s", e.cfg.Port)

	go func() {
		<-ctx.Done()
		_ = e.srv.Shutdown(context.Background())
	}()

	log.Fatal(e.srv.ListenAndServe())
}
