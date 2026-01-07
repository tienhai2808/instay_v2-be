package di

import (
	fileUC "github.com/InstayPMS/backend/internal/application/usecase/file"
	"github.com/InstayPMS/backend/internal/infrastructure/api/http/handler"
	"github.com/InstayPMS/backend/internal/infrastructure/api/http/middleware"
	"github.com/InstayPMS/backend/internal/infrastructure/config"
	"github.com/InstayPMS/backend/internal/infrastructure/initialization"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type Container struct {
	Config            *config.Config
	Database          *initialization.Database
	Log               *zap.Logger
	Storage           *minio.Client
	FileUseCase       fileUC.FileUseCase
	FileHandler       *handler.FileHandler
	ContextMiddleware *middleware.ContextMiddleware
}

func NewContainer(cfg *config.Config) (*Container, error) {
	c := &Container{
		Config: cfg,
	}

	if err := c.initInfrastructure(cfg); err != nil {
		return nil, err
	}

	c.initUseCases()

	c.initHandlers()

	return c, nil
}

func (c *Container) Cleanup() {
	if c.Database != nil {
		c.Database.Close()
	}
}
