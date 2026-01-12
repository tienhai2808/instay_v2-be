package router

import (
	"net/http"

	"github.com/InstayPMS/backend/internal/container"
	"github.com/InstayPMS/backend/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine *gin.Engine
}

func NewRouter(r *gin.Engine) *Router {
	return &Router{r}
}

func (r *Router) Setup(cfg config.ServerConfig, ctn *container.Container) {
	v2 := r.Engine.Group(cfg.APIPrefix)

	v2.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	r.setupFileRoutes(v2, ctn.FileHdl)

	r.setupAuthRoutes(v2, ctn.AuthMid, ctn.AuthHdl)

	r.setupUserRoutes(v2, ctn.AuthMid, ctn.UserHdl)

	r.setupDepartmentRoutes(v2, ctn.AuthMid, ctn.DepartmentHdl)
}
