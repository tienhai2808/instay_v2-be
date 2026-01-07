package di

import (
	"github.com/InstayPMS/backend/internal/infrastructure/api/http/middleware"
	"github.com/InstayPMS/backend/internal/infrastructure/config"
	"github.com/InstayPMS/backend/internal/infrastructure/initialization"
)

func (c *Container) initInfrastructure(cfg *config.Config) error {
	log, err := initialization.InitZap(cfg)
	if err != nil {
		return err
	}
	c.Log = log

	db, err := initialization.InitDatabase(cfg)
	if err != nil {
		return err
	}
	c.Database = db

	stor, err := initialization.InitMinIO(cfg)
	if err != nil {
		return err
	}
	c.Storage = stor

	ctxMid := middleware.NewContextMiddleware(log)
	c.ContextMiddleware = ctxMid

	return nil
}
