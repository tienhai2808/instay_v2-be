package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/InstayPMS/backend/internal/application/dto"
	authUC "github.com/InstayPMS/backend/internal/application/usecase/auth"
	"github.com/InstayPMS/backend/internal/infrastructure/api/http/middleware"
	"github.com/InstayPMS/backend/internal/infrastructure/config"
	"github.com/InstayPMS/backend/pkg/constants"
	"github.com/InstayPMS/backend/pkg/errors"
	"github.com/InstayPMS/backend/pkg/mapper"
	"github.com/InstayPMS/backend/pkg/utils"
	"github.com/InstayPMS/backend/pkg/validator"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	cfg    *config.Config
	authUC authUC.AuthUseCase
}

func NewAuthHandler(
	cfg *config.Config,
	authUC authUC.AuthUseCase,
) *AuthHandler {
	return &AuthHandler{
		cfg,
		authUC,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		c.Error(errors.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	user, accessToken, refreshToken, err := h.authUC.Login(ctx, c.Request.UserAgent(), req)
	if err != nil {
		c.Error(err)
		return
	}

	h.storeTokenInCookie(c, accessToken, refreshToken, int(h.cfg.JWT.AccessExpiresIn.Seconds()), int(h.cfg.JWT.RefreshExpiresIn.Seconds()))

	utils.APIResponse(c, http.StatusOK, constants.CodeLoginSuccess, "Login successfully", gin.H{
		"user": mapper.ToUserResponse(user),
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	accessTTLAny, ok := c.Get(middleware.CtxAccessTTL)
	if !ok {
		c.Error(errors.ErrUnAuth)
		return
	}

	accessTTL, ok := accessTTLAny.(time.Duration)
	if !ok {
		c.Error(errors.ErrUnAuth)
		return
	}

	accessToken, err := c.Cookie(h.cfg.JWT.AccessName)
	if err != nil || accessToken == "" {
		c.Error(errors.ErrUnAuth)
		return
	}

	refreshToken, err := c.Cookie(h.cfg.JWT.RefreshName)
	if err != nil || refreshToken == "" {
		c.Error(errors.ErrInvalidUser)
		return
	}

	if err := h.authUC.Logout(ctx, accessToken, refreshToken, accessTTL); err != nil {
		c.Error(err)
		return
	}

	h.storeTokenInCookie(c, "", "", -1, -1)

	utils.APIResponse(c, http.StatusOK, constants.CodeLogoutSuccess, "Logout successfully", nil)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	refreshToken, err := c.Cookie(h.cfg.JWT.RefreshName)
	if err != nil || refreshToken == "" {
		c.Error(errors.ErrInvalidUser)
		return
	}

	newAccessToken, newRefreshToken, err := h.authUC.RefreshToken(ctx, c.Request.UserAgent(), refreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	h.storeTokenInCookie(c, newAccessToken, newRefreshToken, int(h.cfg.JWT.AccessExpiresIn.Seconds()), int(h.cfg.JWT.RefreshExpiresIn.Seconds()))

	utils.OKResponse(c, nil)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		c.Error(errors.ErrUnAuth)
		return
	}

	user, err := h.authUC.GetMe(ctx, userID)
	if err != nil {
		c.Error(err)
		return
	}

	utils.OKResponse(c, gin.H{
		"user": mapper.ToUserResponse(user),
	})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		c.Error(errors.ErrUnAuth)
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		c.Error(errors.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.authUC.ChangePassword(ctx, userID, req); err != nil {
		c.Error(err)
		return
	}

	h.storeTokenInCookie(c, "", "", -1, -1)

	utils.APIResponse(c, http.StatusOK, constants.CodeChangePasswordSuccess, "Change password successfully", nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		c.Error(errors.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	forgotPasswordToken, err := h.authUC.ForgotPassword(ctx, req.Email)
	if err != nil {
		c.Error(err)
		return
	}

	utils.APIResponse(c, http.StatusOK, constants.CodeForgotPasswordSuccess, "Forgot password successfully", gin.H{
		"forgot_password_token": forgotPasswordToken,
	})
}

func (h *AuthHandler) VerifyForgotPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req dto.VerifyForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		c.Error(errors.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	resetPasswordToken, err := h.authUC.VerifyForgotPassword(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	utils.APIResponse(c, http.StatusOK, constants.CodeVerifyForgotPasswordSuccess, "Verify forgot password successfully", gin.H{
		"reset_password_token": resetPasswordToken,
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		c.Error(errors.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.authUC.ResetPassword(ctx, req); err != nil {
		c.Error(err)
		return
	}

	utils.APIResponse(c, http.StatusOK, constants.CodeResetPasswordSuccess, "Reset password successfully", nil)
}

func (h *AuthHandler) UpdateInfo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		c.Error(errors.ErrUnAuth)
		return
	}

	var req dto.UpdateInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		c.Error(errors.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	updatedUser, err := h.authUC.UpdateInfo(ctx, userID, req)
	if err != nil {
		c.Error(err)
		return
	}

	utils.APIResponse(c, http.StatusOK, constants.CodeUpdateInfoSuccess, "User updated successfully", gin.H{
		"user": mapper.ToUserResponse(updatedUser),
	})
}

func (h *AuthHandler) storeTokenInCookie(c *gin.Context, accessToken, refreshToken string, accessExpiresIn, refreshExpiresIn int) {
	isSecure := c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"
	domain := utils.ExtractRootDomain(c.Request.Host)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		h.cfg.JWT.AccessName,
		accessToken,
		accessExpiresIn,
		fmt.Sprintf("%s", h.cfg.Server.APIPrefix),
		domain,
		isSecure,
		true,
	)
	c.SetCookie(
		h.cfg.JWT.RefreshName,
		refreshToken,
		refreshExpiresIn,
		fmt.Sprintf("%s/auth", h.cfg.Server.APIPrefix),
		domain,
		isSecure,
		true,
	)
}
