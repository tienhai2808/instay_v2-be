package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/InstayPMS/backend/internal/application/dto"
	"github.com/InstayPMS/backend/internal/application/port"
	"github.com/InstayPMS/backend/internal/domain/model"
	"github.com/InstayPMS/backend/internal/domain/repository"
	"github.com/InstayPMS/backend/internal/infrastructure/config"
	"github.com/InstayPMS/backend/pkg/constants"
	customErr "github.com/InstayPMS/backend/pkg/errors"
	"github.com/InstayPMS/backend/pkg/utils"
	"github.com/google/uuid"
	"github.com/sony/sonyflake/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type authUseCaseImpl struct {
	cfg       config.JWTConfig
	db        *gorm.DB
	log       *zap.Logger
	idGen     *sonyflake.Sonyflake
	jwtPro    port.JWTProvider
	cachePro  port.CacheProvider
	mqPro     port.MessageQueueProvider
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewAuthUseCase(
	cfg config.JWTConfig,
	db *gorm.DB,
	log *zap.Logger,
	idGen *sonyflake.Sonyflake,
	jwtPro port.JWTProvider,
	cachePro port.CacheProvider,
	mqPro port.MessageQueueProvider,
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
) AuthUseCase {
	return &authUseCaseImpl{
		cfg,
		db,
		log,
		idGen,
		jwtPro,
		cachePro,
		mqPro,
		userRepo,
		tokenRepo,
	}
}

func (u *authUseCaseImpl) Login(ctx context.Context, ua string, req dto.LoginRequest) (*model.User, string, string, error) {
	user, err := u.userRepo.FindByUsernameWithDepartment(ctx, req.Username)
	if err != nil {
		u.log.Error("find user by username failed", zap.String("username", req.Username), zap.Error(err))
		return nil, "", "", err
	}
	if user == nil || !user.IsActive {
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

	accessToken, err := u.jwtPro.GenerateToken(user.ID, user.Role, tokenVersion, u.cfg.AccessExpiresIn)
	if err != nil {
		u.log.Error("generate access token failed", zap.Error(err))
		return nil, "", "", err
	}

	refreshToken, err := generateRefreshToken()
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
		u.log.Error("update token by token failed", zap.Error(err))
		return err
	}

	redisKey := fmt.Sprintf("black_list:%s", accessToken)
	if err := u.cachePro.SetString(ctx, redisKey, "1", accessTTL); err != nil {
		u.log.Error("save black list failed", zap.Error(err))
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

	user, err := u.userRepo.FindByID(ctx, token.UserID)
	if err != nil {
		u.log.Error("find user by id failed", zap.Error(err))
		return "", "", nil
	}
	if user == nil || !user.IsActive {
		return "", "", customErr.ErrInvalidUser
	}

	redisKey := fmt.Sprintf("user_version:%d", user.ID)
	tokenVersion, err := u.cachePro.GetInt(ctx, redisKey)
	if err != nil {
		u.log.Error("get token version failed", zap.Error(err))
		return "", "", err
	}
	if tokenVersion == 0 {
		return "", "", customErr.ErrInvalidUser
	}

	newAccessToken, err := u.jwtPro.GenerateToken(user.ID, user.Role, tokenVersion, u.cfg.AccessExpiresIn)
	if err != nil {
		u.log.Error("generate access token failed", zap.Error(err))
		return "", "", err
	}

	newRefreshToken, err := generateRefreshToken()
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
		UserID:    user.ID,
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
	user, err := u.userRepo.FindByIDWithDepartment(ctx, userID)
	if err != nil {
		u.log.Error("find user by id failed", zap.Int64("id", userID), zap.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, customErr.ErrInvalidUser
	}

	return user, nil
}

func (u *authUseCaseImpl) ChangePassword(ctx context.Context, userID int64, req dto.ChangePasswordRequest) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		u.log.Error("find user by id failed", zap.Int64("id", userID), zap.Error(err))
		return err
	}
	if user == nil {
		return customErr.ErrInvalidUser
	}

	if err = utils.VerifyPassword(req.OldPassword, user.Password); err != nil {
		return customErr.ErrInvalidPassword
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		u.log.Error("hash password failed", zap.Error(err))
		return err
	}

	if err = u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = u.userRepo.UpdateTx(tx, userID, map[string]any{"password": hashedPassword}); err != nil {
			if errors.Is(err, customErr.ErrUserNotFound) {
				return customErr.ErrInvalidUser
			}
			u.log.Error("update password failed", zap.Error(err))
			return err
		}

		if err := u.tokenRepo.UpdateByUserIDTx(tx, userID, map[string]any{"revoked_at": time.Now()}); err != nil {
			u.log.Error("update token by user id failed", zap.Error(err))
			return err
		}
		return nil
	}); err != nil {
		return nil
	}

	redisKey := fmt.Sprintf("user_version:%d", user.ID)
	if err = u.cachePro.Increment(ctx, redisKey); err != nil {
		u.log.Error("increase token version failed", zap.Error(err))
	}

	return nil
}

func (u *authUseCaseImpl) ForgotPassword(ctx context.Context, email string) (string, error) {
	exists, err := u.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		u.log.Error("check user by email failed", zap.String("email", email), zap.Error(err))
		return "", err
	}
	if !exists {
		return "", customErr.ErrEmailDoesNotExist
	}

	otp := utils.GenerateOTP(6)
	forgotPasswordToken := uuid.NewString()

	forgData := dto.ForgotPasswordData{
		Email:    email,
		Otp:      otp,
		Attempts: 0,
	}

	bytes, err := json.Marshal(forgData)
	if err != nil {
		u.log.Error("json marshal forgot password data failed", zap.Error(err))
		return "", err
	}

	redisKey := fmt.Sprintf("forgot_password:%s", forgotPasswordToken)
	if err = u.cachePro.SetObject(ctx, redisKey, bytes, 3*time.Minute); err != nil {
		u.log.Error("save forgot password data failed", zap.Error(err))
		return "", err
	}

	emailMsg := dto.AuthEmailMessage{
		To:      email,
		Subject: "Xác thực quên mật khẩu tại Instay",
		Otp:     otp,
	}

	go func(msg dto.AuthEmailMessage) {
		body, err := json.Marshal(msg)
		if err != nil {
			u.log.Error("json marshal failed", zap.Error(err))
		}

		if u.mqPro.PublishMessage(constants.ExchangeEmail, constants.RoutingKeyAuthEmail, body); err != nil {
			u.log.Error("publish auth email message failed", zap.String("email", email), zap.Error(err))
		}
	}(emailMsg)

	return forgotPasswordToken, nil
}

func (u *authUseCaseImpl) VerifyForgotPassword(ctx context.Context, req dto.VerifyForgotPasswordRequest) (string, error) {
	redisKey := fmt.Sprintf("forgot_password:%s", req.ForgotPasswordToken)
	bytes, err := u.cachePro.GetObject(ctx, redisKey)
	if err != nil {
		u.log.Error("get forgot password data failed", zap.Error(err))
		return "", err
	}
	if bytes == nil {
		return "", customErr.ErrInvalidToken
	}

	var forgData dto.ForgotPasswordData
	if err = json.Unmarshal(bytes, &forgData); err != nil {
		u.log.Error("json unmarshal forgot password data failed", zap.Error(err))
		return "", nil
	}

	if forgData.Attempts >= 3 {
		if err = u.cachePro.Del(ctx, redisKey); err != nil {
			u.log.Error("delete forgot password data failed", zap.Error(err))
			return "", err
		}
		return "", customErr.ErrTooManyAttempts
	}
	if forgData.Otp != req.Otp {
		return "", customErr.ErrInvalidOTP
	}

	resetPasswordToken := uuid.NewString()
	key := fmt.Sprintf("reset_password:%s", resetPasswordToken)

	if err = u.cachePro.SetString(ctx, key, forgData.Email, 3*time.Minute); err != nil {
		u.log.Error("save email reset password failed", zap.Error(err))
		return "", err
	}

	if err = u.cachePro.Del(ctx, redisKey); err != nil {
		u.log.Error("delete forgot password data failed", zap.Error(err))
	}

	return resetPasswordToken, nil
}

func (u *authUseCaseImpl) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	redisKey := fmt.Sprintf("reset_password:%s", req.ResetPasswordToken)
	email, err := u.cachePro.GetString(ctx, redisKey)
	if err != nil {
		u.log.Error("get email reset password failed", zap.Error(err))
		return err
	}
	if email == "" {
		return customErr.ErrInvalidToken
	}

	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		u.log.Error("find user by email failed", zap.String("email", email), zap.Error(err))
		return err
	}
	if user == nil {
		return customErr.ErrUnAuth
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		u.log.Error("hash password failed", zap.Error(err))
		return err
	}

	if err = u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = u.userRepo.UpdateTx(tx, user.ID, map[string]any{"password": hashedPassword}); err != nil {
			if errors.Is(err, customErr.ErrUserNotFound) {
				return customErr.ErrUnAuth
			}
			u.log.Error("update password failed", zap.Error(err))
			return err
		}

		if err := u.tokenRepo.UpdateByUserIDTx(tx, user.ID, map[string]any{"revoked_at": time.Now()}); err != nil {
			u.log.Error("update token by user id failed", zap.Error(err))
			return err
		}
		return nil
	}); err != nil {
		return nil
	}

	if err = u.cachePro.Del(ctx, redisKey); err != nil {
		u.log.Error("delete reset password failed", zap.Error(err))
	}

	redisKey = fmt.Sprintf("user_version:%d", user.ID)
	if err = u.cachePro.Increment(ctx, redisKey); err != nil {
		u.log.Error("increase token version failed", zap.Error(err))
	}

	return nil
}

func (u *authUseCaseImpl) UpdateInfo(ctx context.Context, userID int64, req dto.UpdateInfoRequest) (*model.User, error) {
	updateData := map[string]any{
		"email":      req.Email,
		"phone":      req.Phone,
		"first_name": req.FirstName,
		"last_name":  req.LastName,
	}

	if err := u.userRepo.Update(ctx, userID, updateData); err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			return nil, customErr.ErrInvalidUser
		}
		ok, constraint := utils.IsUniqueViolation(err)
		if ok {
			switch constraint {
			case "users_email_key":
				return nil, customErr.ErrEmailAlreadyExists
			case "users_phone_key":
				return nil, customErr.ErrPhoneAlreadyExists
			}
		}
		u.log.Error("update user failed", zap.Int64("id", userID), zap.Error(err))
		return nil, err
	}

	user, err := u.userRepo.FindByIDWithDepartment(ctx, userID)
	if err != nil {
		u.log.Error("find user by id failed", zap.Int64("id", userID), zap.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, customErr.ErrInvalidUser
	}

	return user, nil
}

func generateRefreshToken() (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	rawToken := base64.RawURLEncoding.EncodeToString(randomBytes)

	return rawToken, nil
}
