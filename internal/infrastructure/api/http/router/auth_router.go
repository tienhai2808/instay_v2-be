package router

import (
	"github.com/InstayPMS/backend/internal/infrastructure/api/http/handler"
	"github.com/InstayPMS/backend/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupAuthRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", hdl.Login)

		auth.POST("/logout", authMid.IsAuthentication(), authMid.AttachTokens(), hdl.Logout)

		auth.POST("/refresh-token", hdl.RefreshToken)

		auth.GET("/me", authMid.IsAuthentication(), hdl.GetMe)
	}
}
