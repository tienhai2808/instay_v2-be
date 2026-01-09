package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/InstayPMS/backend/internal/application/dto"
	"github.com/InstayPMS/backend/internal/application/port"
	"github.com/InstayPMS/backend/internal/domain/model"
	"github.com/InstayPMS/backend/internal/domain/repository"
	"github.com/InstayPMS/backend/internal/infrastructure/config"
	customErr "github.com/InstayPMS/backend/pkg/errors"
	"github.com/InstayPMS/backend/pkg/utils"
	"github.com/sony/sonyflake/v2"
	"go.uber.org/zap"
)

type authUseCaseImpl struct {
	cfg       config.JWTConfig
	log       *zap.Logger
	idGen     *sonyflake.Sonyflake
	jwtPro    port.JWTProvider
	cachePro  port.CacheProvider
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewAuthUseCase(
	cfg config.JWTConfig,
	log *zap.Logger,
	idGen *sonyflake.Sonyflake,
	jwtPro port.JWTProvider,
	cachePro port.CacheProvider,
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
) AuthUseCase {
	return &authUseCaseImpl{
		cfg,
		log,
		idGen,
		jwtPro,
		cachePro,
		userRepo,
		tokenRepo,
	}
}

func (u *authUseCaseImpl) Login(ctx context.Context, ua string, req dto.LoginRequest) (*model.User, string, string, error) {
	user, err := u.userRepo.FindByUsernameWithOutletAndDepartment(ctx, req.Username)
	if err != nil {
		u.log.Error("find user by username failed", zap.String("username", req.Username), zap.Error(err))
		return nil, "", "", err
	}

	if user == nil {
		return nil, "", "", customErr.ErrLoginFailed
	}

	if !user.IsActive {
		return nil, "", "", customErr.ErrLoginFailed
	}

	if err = utils.VerifyPassword(req.Password, user.Password); err != nil {
		return nil, "", "", customErr.ErrLoginFailed
	}

	redisKey := fmt.Sprintf("user_version:%d", user.ID)
	tokenVersion, err := u.cachePro.GetInt(ctx, redisKey)
	if err != nil {
		u.log.Error("get token version failed", zap.Error(err))
		return nil, "", "", err
	}

	if tokenVersion == 0 {
		if err = u.cachePro.SetString(ctx, redisKey, "1", 0); err != nil {
			u.log.Error("save token version failed", zap.Error(err))
			return nil, "", "", err
		}
		tokenVersion = 1
	}

	accessToken, err := u.jwtPro.GenerateToken(user.ID, tokenVersion, u.cfg.AccessExpiresIn)
	if err != nil {
		u.log.Error("generate access token failed", zap.Error(err))
		return nil, "", "", err
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		u.log.Error("generate refresh token failed", zap.Error(err))
		return nil, "", "", err
	}

	id, err := u.idGen.NextID()
	if err != nil {
		u.log.Error("generate token id failed", zap.Error(err))
		return nil, "", "", err
	}

	token := &model.Token{
		ID:        id,
		UserID:    user.ID,
		Token:     utils.SHA256Hash(refreshToken),
		UserAgent: utils.ConvertUserAgent(ua),
		RevokedAt: nil,
		ExpiresAt: time.Now().Add(u.cfg.RefreshExpiresIn),
	}

	if err := u.tokenRepo.Create(ctx, token); err != nil {
		u.log.Error("create token failed", zap.Error(err))
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (u *authUseCaseImpl) Logout(ctx context.Context, accessToken, refreshToken string, accessTTL time.Duration) error {
	hashedToken := utils.SHA256Hash(refreshToken)

	if err := u.tokenRepo.UpdateByToken(ctx, hashedToken, map[string]any{"revoked_at": time.Now()}); err != nil {
		if errors.Is(err, customErr.ErrInvalidUser) {
			return err
		}
		u.log.Error("update token by user id and token failed", zap.Error(err))
		return err
	}

	redisKey := fmt.Sprintf("black_list:%s", accessToken)
	if err := u.cachePro.SetString(ctx, redisKey, "1", accessTTL); err != nil {
		u.log.Error("save black list failed", zap.Error(err))
		return err
	}

	return nil
}

func (u *authUseCaseImpl) RefreshToken(ctx context.Context, ua, refreshToken string) (string, string, error) {
	hashedToken := utils.SHA256Hash(refreshToken)

	token, err := u.tokenRepo.FindByToken(ctx, hashedToken)
	if err != nil {
		u.log.Error("find token by token failed", zap.Error(err))
		return "", "", nil
	}

	if token == nil || token.RevokedAt != nil || token.ExpiresAt.Before(time.Now()) {
		return "", "", customErr.ErrInvalidUser
	}

	userID := token.UserID

	redisKey := fmt.Sprintf("user_version:%d", userID)
	tokenVersion, err := u.cachePro.GetInt(ctx, redisKey)
	if err != nil {
		u.log.Error("get token version failed", zap.Error(err))
		return "", "", err
	}

	if tokenVersion == 0 {
		return "", "", customErr.ErrInvalidUser
	}

	newAccessToken, err := u.jwtPro.GenerateToken(userID, tokenVersion, u.cfg.AccessExpiresIn)
	if err != nil {
		u.log.Error("generate access token failed", zap.Error(err))
		return "", "", err
	}

	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		u.log.Error("generate refresh token failed", zap.Error(err))
		return "", "", err
	}

	id, err := u.idGen.NextID()
	if err != nil {
		u.log.Error("generate token id failed", zap.Error(err))
		return "", "", err
	}

	newToken := &model.Token{
		ID:        id,
		UserID:    userID,
		Token:     utils.SHA256Hash(newRefreshToken),
		UserAgent: utils.ConvertUserAgent(ua),
		RevokedAt: nil,
		ExpiresAt: time.Now().Add(u.cfg.RefreshExpiresIn),
	}

	if err := u.tokenRepo.Create(ctx, newToken); err != nil {
		u.log.Error("create token failed", zap.Error(err))
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (u *authUseCaseImpl) GetMe(ctx context.Context, userID int64) (*model.User, error) {
	user, err := u.userRepo.FindByIDWithOutletAndDepartment(ctx, userID)
	if err != nil {
		u.log.Error("find user by id failed", zap.Int64("id", userID), zap.Error(err))
		return nil, err
	}

	if user == nil {
		return nil, customErr.ErrUnAuth
	}

	return user, nil
}
