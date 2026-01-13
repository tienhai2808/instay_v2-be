package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
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

	FindAllWithDepartmentPaginated(ctx context.Context, query dto.UserPaginationQuery) ([]*model.User, int64, error)

	ExistsActiveAdminExceptID(ctx context.Context, id int64) (bool, error)

	DeleteTx(tx *gorm.DB, id int64) error

	DeleteAllByIDsTx(tx *gorm.DB, ids []int64) (int64, error)

	ExistsActiveAdmin(ctx context.Context) (bool, error)
}
