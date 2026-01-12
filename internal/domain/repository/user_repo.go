package repository

import (
	"context"

	"github.com/InstayPMS/backend/internal/domain/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	
	FindByUsernameWithDepartment(ctx context.Context, username string) (*model.User, error)

	FindByIDWithDepartment(ctx context.Context, id int64) (*model.User, error)

	FindByIDWithDetails(ctx context.Context, id int64) (*model.User, error)

	FindByID(ctx context.Context, id int64) (*model.User, error)

	UpdateTx(tx *gorm.DB, id int64, updateData map[string]any) error

	ExistsByEmail(ctx context.Context, email string) (bool, error)

	FindByEmail(ctx context.Context, email string) (*model.User, error)

	Update(ctx context.Context, id int64, updateData map[string]any) error
}
