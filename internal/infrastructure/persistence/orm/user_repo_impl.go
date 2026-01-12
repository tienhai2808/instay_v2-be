package orm

import (
	"context"
	"errors"

	"github.com/InstayPMS/backend/internal/domain/model"
	"github.com/InstayPMS/backend/internal/domain/repository"
	customErr "github.com/InstayPMS/backend/pkg/errors"
	"gorm.io/gorm"
)

type Preload struct {
	Relation string
	Scope    func(*gorm.DB) *gorm.DB
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepositoryImpl) FindByUsernameWithDepartment(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).
		Preload("Department").
		Where("username = ?", username).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) FindByIDWithDepartment(ctx context.Context, id int64) (*model.User, error) {
	return findByIDBase(r.db.WithContext(ctx), id, Preload{Relation: "Department"})
}

func (r *userRepositoryImpl) FindByIDWithDetails(ctx context.Context, id int64) (*model.User, error) {
	return findByIDBase(r.db.WithContext(ctx), id,
		Preload{Relation: "Department"},
		Preload{Relation: "CreatedBy"},
		Preload{Relation: "UpdatedBy"},
	)
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.User, error) {
	return findByIDBase(r.db.WithContext(ctx), id)
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) UpdateTx(tx *gorm.DB, id int64, updateData map[string]any) error {
	result := tx.Model(&model.User{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return customErr.ErrUserNotFound
	}

	return nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, id int64, updateData map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return customErr.ErrUserNotFound
	}

	return nil
}

func (r *userRepositoryImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func findByIDBase(tx *gorm.DB, id int64, preloads ...Preload) (*model.User, error) {
	var user model.User

	for _, preload := range preloads {
		if preload.Scope != nil {
			tx = tx.Preload(preload.Relation, preload.Scope)
		} else {
			tx = tx.Preload(preload.Relation)
		}
	}

	if err := tx.Where("id = ?", id).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
