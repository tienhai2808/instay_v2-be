package utils

import (
	"net/http"

	"github.com/InstayPMS/backend/internal/application/dto"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func APIResponse(c *gin.Context, status, code int, slug, message string, data any) {
	c.JSON(status, dto.APIResponse{
		Code:    code,
		Slug:    slug,
		Message: message,
		Data:    data,
	})
}

func ISEResponse(c *gin.Context) {
	APIResponse(c, http.StatusInternalServerError, 9000, "INTERNAL_SERVER_ERROR", "Internal server error", nil)
}

func GenerateSlug(str string) string {
	return slug.Make(str)
}
