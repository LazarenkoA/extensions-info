package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type health struct {
}

func New() *health {
	return &health{}
}

func (h *health) Register(route *gin.RouterGroup) {
	route.GET("/health", h.health)
}

func (h *health) health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})
}
