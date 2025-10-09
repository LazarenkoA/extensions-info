package configuration

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"your-app/internal/models"
	"your-app/internal/usecase"
)

type repo interface {
	GetConfigurationInfo(ctx context.Context, dbID int32) (*models.ConfigurationInfo, error)
	GetMetadata(ctx context.Context, confID int32) ([]byte, error)
}

type Configuration struct {
	repo repo
}

func New(repo repo) *Configuration {
	return &Configuration{
		repo: repo,
	}
}

func (c *Configuration) Register(route *gin.RouterGroup) {
	route.GET("/getConfigurationInfo", c.getConfigurationInfo)
}

func (c *Configuration) getConfigurationInfo(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		usecase.ResponseError(ctx, fmt.Errorf("bad id %q. %w", idStr, err))
		return
	}

	info, err := c.repo.GetConfigurationInfo(ctx, int32(id))
	if err != nil {
		usecase.ResponseError(ctx, err)
		return
	}

	data, err := c.repo.GetMetadata(ctx, int32(id))
	if err != nil {
		log.Println(err)
	}

	info.MetadataTree = &models.MetadataInfo{
		ObjectName: info.Name,
		Type:       models.ObjectTypeConf,
	}
	_ = json.Unmarshal(data, &info.MetadataTree.Children)

	ctx.JSON(http.StatusOK, gin.H{
		"data": info,
	})
}
