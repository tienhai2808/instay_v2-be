package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/InstayPMS/backend/internal/application/dto"
	fileUC "github.com/InstayPMS/backend/internal/application/usecase/file"
	"github.com/InstayPMS/backend/pkg/utils"
	"github.com/InstayPMS/backend/pkg/validator"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileUC fileUC.FileUseCase
}

func NewFileHandler(fileUC fileUC.FileUseCase) *FileHandler {
	return &FileHandler{fileUC}
}

func (h *FileHandler) UploadPresignedURLs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req dto.UploadPresignedURLsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := validator.HandleRequestError(err)
		utils.APIResponse(c, http.StatusBadRequest, 4000, "INVALID_REQUEST_ERROR", mess, nil)
		return
	}

	presignedURLs, err := h.fileUC.CreateUploadURLs(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	utils.APIResponse(c, http.StatusOK, 1000, "SUCCESS", "Generate upload presigned urls successfully", gin.H{
		"presigned_urls": presignedURLs,
	})
}

func (h *FileHandler) ViewPresignedURLs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req dto.ViewPresignedURLsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := validator.HandleRequestError(err)
		utils.APIResponse(c, http.StatusBadRequest, 4000, "INVALID_REQUEST_ERROR", mess, nil)
		return
	}

	presignedURLs, err := h.fileUC.CreateViewURLs(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	utils.APIResponse(c, http.StatusOK, 1000, "SUCCESS", "Generate view presigned urls successfully", gin.H{
		"presigned_urls": presignedURLs,
	})
}
