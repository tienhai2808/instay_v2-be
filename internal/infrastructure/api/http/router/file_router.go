package router

import (
	"github.com/InstayPMS/backend/internal/infrastructure/api/http/handler"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupFileRoutes(rg *gin.RouterGroup, hdl *handler.FileHandler) {
	file := rg.Group("/files")
	{
		file.POST("/presigned-urls/uploads", hdl.UploadPresignedURLs)

		file.POST("/presigned-urls/views", hdl.ViewPresignedURLs)
	}
}
