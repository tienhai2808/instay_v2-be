package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine *gin.Engine
}

func NewRouter(r *gin.Engine) *Router {
	return &Router{
		r,
	}
}

func (r *Router) Setup(apiPrefix string) {
	v2 := r.Engine.Group(apiPrefix)

	v2.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})
}
