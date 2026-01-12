package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/InstayPMS/backend/internal/application/dto"
	userUC "github.com/InstayPMS/backend/internal/application/usecase/user"
	"github.com/InstayPMS/backend/internal/domain/model"
	"github.com/InstayPMS/backend/internal/infrastructure/api/http/middleware"
	"github.com/InstayPMS/backend/pkg/constants"
	"github.com/InstayPMS/backend/pkg/errors"
	"github.com/InstayPMS/backend/pkg/mapper"
	"github.com/InstayPMS/backend/pkg/utils"
	"github.com/InstayPMS/backend/pkg/validator"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUC userUC.UserUseCase
}

func NewUserHandler(userUC userUC.UserUseCase) *UserHandler {
	return &UserHandler{userUC}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		c.Error(errors.ErrUnAuth)
		return
	}

	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		c.Error(errors.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if req.Role == model.RoleAdmin {
		if req.DepartmentID != nil {
			c.Error(errors.ErrBadRequest.WithData(gin.H{
				"field": "departmentid",
				"tag":   "notrequired",
				"param": "",
			}))
			return
		}
	} else {
		if req.DepartmentID == nil {
			c.Error(errors.ErrBadRequest.WithData(gin.H{
				"field": "departmentid",
				"tag":   "required",
				"param": "",
			}))
			return
		}
	}

	id, err := h.userUC.CreateUser(ctx, userID, req)
	if err != nil {
		c.Error(err)
		return
	}

	utils.APIResponse(c, http.StatusCreated, constants.CodeCreateUserSuccess, "User created successfully", gin.H{
		"user_id": id,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.Error(errors.ErrInvalidID)
		return
	}

	user, err := h.userUC.GetUserByID(ctx, userID)
	if err != nil {
		c.Error(err)
		return
	}

	utils.OKResponse(c, gin.H{
		"user": mapper.ToUserDetailsResponse(user),
	})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	
}
