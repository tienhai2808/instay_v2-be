package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/InstayPMS/backend/internal/application/dto"
	"github.com/InstayPMS/backend/internal/application/port"
	"github.com/InstayPMS/backend/internal/domain/model"
	"github.com/InstayPMS/backend/internal/infrastructure/config"
	"github.com/InstayPMS/backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	CtxUserID       = "user_id"
	CtxAccessToken  = "access_token"
	CtxRefreshToken = "refresh_token"
	CtxAccessTTL    = "access_ttl"
	CtxRole         = "role"
)

type AuthMiddleware struct {
	cfg      config.JWTConfig
	log      *zap.Logger
	jwtPro   port.JWTProvider
	cachePro port.CacheProvider
}

func NewAuthMiddleware(
	cfg config.JWTConfig,
	log *zap.Logger,
	jwtPro port.JWTProvider,
	cachePro port.CacheProvider,
) *AuthMiddleware {
	return &AuthMiddleware{
		cfg,
		log,
		jwtPro,
		cachePro,
	}
}

func (m *AuthMiddleware) IsAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie(m.cfg.AccessName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.APIResponse{
				Code:    errors.ErrUnAuth.Code,
				Message: errors.ErrUnAuth.Message,
			})
			return
		}

		userID, role, tokenVersion, ttl, err := m.jwtPro.ParseToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.APIResponse{
				Code:    errors.ErrUnAuth.Code,
				Message: errors.ErrUnAuth.Message,
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		blackListKey := fmt.Sprintf("black_list:%s", accessToken)
		str, err := m.cachePro.GetString(ctx, blackListKey)
		if err != nil {
			m.log.Error("get black list failed", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.APIResponse{
				Code:    errors.ErrUnAuth.Code,
				Message: errors.ErrUnAuth.Message,
			})
			return
		}

		if str != "" {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.APIResponse{
				Code:    errors.ErrInvalidUser.Code,
				Message: errors.ErrInvalidUser.Message,
			})
			return
		}

		userVersionKey := fmt.Sprintf("user_version:%d", userID)
		currentTokenVersion, err := m.cachePro.GetInt(ctx, userVersionKey)
		if err != nil {
			m.log.Error("get token version failed", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.APIResponse{
				Code:    errors.ErrUnAuth.Code,
				Message: errors.ErrUnAuth.Message,
			})
			return
		}

		if tokenVersion != currentTokenVersion {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.APIResponse{
				Code:    errors.ErrUnAuth.Code,
				Message: errors.ErrUnAuth.Message,
			})
			return
		}

		c.Set(CtxUserID, userID)
		c.Set(CtxRole, string(role))
		c.Set(CtxAccessTTL, ttl)

		c.Next()
	}
}

func (m *AuthMiddleware) AttachTokens() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie(m.cfg.AccessName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.APIResponse{
				Code:    errors.ErrUnAuth.Code,
				Message: errors.ErrUnAuth.Message,
			})
			return
		}

		refreshToken, err := c.Cookie(m.cfg.RefreshName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.APIResponse{
				Code:    errors.ErrInvalidUser.Code,
				Message: errors.ErrInvalidUser.Message,
			})
			return
		}

		c.Set(CtxAccessToken, accessToken)
		c.Set(CtxRefreshToken, refreshToken)

		c.Next()
	}
}

func (m *AuthMiddleware) HasRole(allowedRole model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleStr := c.GetString(CtxRole)
		role := model.UserRole(roleStr)

		if roleStr == "" || !model.IsValidRole(role) || role != allowedRole {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.APIResponse{
				Code:    errors.ErrForbidden.Code,
				Message: errors.ErrForbidden.Message,
			})
			return
		}

		c.Next()
	}
}
