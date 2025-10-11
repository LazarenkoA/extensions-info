package configuration

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LazarenkoA/extensions-info/internal/models"
	"github.com/LazarenkoA/extensions-info/internal/usecase"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type repo interface {
	GetConfigurationInfo(ctx context.Context, dbID int32) (*models.ConfigurationInfo, error)
	GetMetadata(ctx context.Context, confID int32) ([]byte, error)
	GetCode(ctx context.Context, extID int32, key string) (string, error)
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
	route.GET("/getSourceCode", c.getSourceCode)
}

func (c *Configuration) getSourceCode(ctx *gin.Context) {
	extIDstr := ctx.Query("extid")
	modulekey := ctx.Query("modulekey")
	extID, err := strconv.Atoi(extIDstr)
	if err != nil {
		usecase.ResponseError(ctx, fmt.Errorf("bad id %q. %w", extIDstr, err))
		return
	}

	code, err := c.repo.GetCode(ctx, int32(extID), modulekey)
	if err != nil {
		usecase.ResponseError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": code,
	})
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

	data, err := c.repo.GetMetadata(ctx, info.ID)
	if err != nil {
		log.Println(err)
	}

	info.MetadataTree = &models.MetadataInfo{
		ObjectName: info.Name,
		Type:       "Configuration",
	}
	_ = json.Unmarshal(data, &info.MetadataTree.Children)

	ctx.JSON(http.StatusOK, gin.H{
		"data": info,
	})
}
