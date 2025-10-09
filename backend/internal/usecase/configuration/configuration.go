package configuration

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"your-app/internal/models"
	"your-app/internal/usecase"
)

type repo interface {
	GetConfigurationInfo(ctx context.Context, dbID int32) (*models.ConfigurationInfo, error)
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

	// заполняем мок данными
	info.MetadataTree = &models.MetadataInfo{
		ObjectName: "Конфигурация такая-то",
		Path:       "",
		Type:       models.ObjectTypeConf,
		Children: []*models.MetadataInfo{
			{
				ObjectName: "Документ1",
				Type:       models.ObjectTypeDocument,
			},
			{
				ObjectName:   "Документ2",
				Type:         models.ObjectTypeDocument,
				ExtensionIDs: []int32{1, 2},
			},
			{
				ObjectName:   "Физ лица",
				Type:         models.ObjectTypeCatalog,
				ExtensionIDs: []int32{1},
			},
			{
				ObjectName: "ОбщегоНазначения",
				Type:       models.ObjectTypeCommonModule,
				Funcs: []models.FuncInfo{
					{
						RedefinitionMethod: models.RedefinitionAfter,
						Name:               "Тест",
						Code:               "// тут код",
						Type:               models.ObjectTypeFunction,
						ExtensionIDs:       []int32{1, 2},
					},
					{
						RedefinitionMethod: models.RedefinitionChangeControl,
						Name:               "Тест2",
						Code:               "// тут код2",
						Type:               models.ObjectTypeFunction,
						ExtensionIDs:       []int32{1},
					},
				},
			},
		},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": info,
	})
}
