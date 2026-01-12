package port

import (
	"time"

	"github.com/InstayPMS/backend/internal/domain/model"
)

type JWTProvider interface {
	GenerateToken(userID int64, role model.UserRole, tokenVersion int, ttl time.Duration) (string, error)

	ParseToken(tokenStr string) (int64, model.UserRole, int, time.Duration, error)
}
