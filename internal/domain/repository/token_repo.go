package repository

import (
	"context"

	"github.com/InstayPMS/backend/internal/domain/model"
)

type TokenRepository interface {
	Create(ctx context.Context, token *model.Token) error

	UpdateByToken(ctx context.Context, token string, updateData map[string]any) error

	FindByToken(ctx context.Context, token string) (*model.Token, error)
}
