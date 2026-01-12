package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/InstayPMS/backend/internal/application/port"
	"github.com/InstayPMS/backend/internal/domain/model"
	"github.com/InstayPMS/backend/internal/infrastructure/config"
	"github.com/InstayPMS/backend/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Role         model.UserRole `json:"role"`
	TokenVersion int            `json:"token_version"`
}

type jwtProviderImpl struct {
	cfg config.JWTConfig
}

func NewJWTProvider(cfg config.JWTConfig) port.JWTProvider {
	return &jwtProviderImpl{cfg}
}

func (p *jwtProviderImpl) GenerateToken(userID int64, role model.UserRole, tokenVersion int, ttl time.Duration) (string, error) {
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role:         role,
		TokenVersion: tokenVersion,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.cfg.SecretKey))
}

func (p *jwtProviderImpl) ParseToken(tokenStr string) (int64, model.UserRole, int, time.Duration, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", t.Header["alg"])
		}
		return []byte(p.cfg.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return 0, "", 0, 0, errors.ErrInvalidToken
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, "", 0, 0, errors.ErrInvalidToken
	}

	ttl := time.Until(claims.ExpiresAt.Time)

	return userID, claims.Role, claims.TokenVersion, ttl, nil
}
