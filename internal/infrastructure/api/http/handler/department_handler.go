package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/InstayPMS/backend/internal/application/dto"
	departmentUC "github.com/InstayPMS/backend/internal/application/usecase/department"
	"github.com/InstayPMS/backend/internal/infrastructure/api/http/middleware"
	"github.com/InstayPMS/backend/pkg/constants"
	"github.com/InstayPMS/backend/pkg/errors"
	"github.com/InstayPMS/backend/pkg/utils"
	"github.com/InstayPMS/backend/pkg/validator"
	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	departmentUC departmentUC.DepartmentUseCase
}

func NewDepartmentHandler(departmentUC departmentUC.DepartmentUseCase) *DepartmentHandler {
	return &DepartmentHandler{departmentUC}
}

func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		c.Error(errors.ErrUnAuth)
		return
	}

	var req dto.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		c.Error(errors.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	id, err := h.departmentUC.CreateDepartment(ctx, userID, req)
	if err != nil {
		c.Error(err)
		return
	}

	utils.APIResponse(c, http.StatusCreated, constants.CodeCreateDepartmentSuccess, "Department created successfully", gin.H{
		"department_id": id,
	})
}
