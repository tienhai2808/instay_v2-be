package repository

import (
	"context"

	"github.com/InstayPMS/backend/internal/domain/model"
)

type UserRepository interface {
	FindByUsernameWithOutletAndDepartment(ctx context.Context, username string) (*model.User, error)

	FindByIDWithOutletAndDepartment(ctx context.Context, id int64) (*model.User, error)
}
