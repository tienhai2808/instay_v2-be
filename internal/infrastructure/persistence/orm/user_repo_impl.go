package orm

import (
	"context"
	"errors"
	"strings"

	"github.com/InstaySystem/is_v2-be/internal/application/dto"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	customErr "github.com/InstaySystem/is_v2-be/pkg/errors"
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
	return r.findByIDBase(r.db.WithContext(ctx), id, Preload{Relation: "Department"})
}

func (r *userRepositoryImpl) FindByIDWithDetails(ctx context.Context, id int64) (*model.User, error) {
	return r.findByIDBase(r.db.WithContext(ctx), id,
		Preload{Relation: "Department"},
		Preload{Relation: "CreatedBy"},
		Preload{Relation: "UpdatedBy"},
	)
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.User, error) {
	return r.findByIDBase(r.db.WithContext(ctx), id)
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
	result := r.db.WithContext(ctx).
		Model(&model.User{}).
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
	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepositoryImpl) ExistsActiveAdminExceptID(ctx context.Context, id int64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("role = ? AND is_active = true AND id <> ?", model.RoleAdmin, id).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepositoryImpl) ExistsActiveAdmin(ctx context.Context) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("role = ? AND is_active = true", model.RoleAdmin).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepositoryImpl) DeleteTx(tx *gorm.DB, id int64) error {
	result := tx.Where("id = ?", id).
		Delete(&model.User{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return customErr.ErrUserNotFound
	}

	return nil
}

func (r *userRepositoryImpl) DeleteAllByIDsTx(tx *gorm.DB, ids []int64) (int64, error) {
	result := tx.Where("id IN ?", ids).
		Delete(&model.User{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (r *userRepositoryImpl) FindAllWithDepartmentPaginated(ctx context.Context, query dto.UserPaginationQuery) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	db := r.db.WithContext(ctx).
		Model(&model.User{})

	db = r.applyFilters(db, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.User{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	db = db.Preload("Department", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	})

	db = r.applySorting(db, query)

	offset := (query.Page - 1) * query.Limit

	if err := db.Select("id", "role", "first_name", "last_name", "is_active", "department_id", "created_at").
		Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepositoryImpl) findByIDBase(tx *gorm.DB, id int64, preloads ...Preload) (*model.User, error) {
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

func (r *userRepositoryImpl) applyFilters(db *gorm.DB, query dto.UserPaginationQuery) *gorm.DB {
	if query.Search != "" {
		term := "%" + query.Search + "%"
		db = db.Where(
			"username ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
			term, term, term,
		)
	}

	if query.Role != "" {
		db = db.Where("role = ?", query.Role)
	}

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	if query.DepartmentID != 0 {
		db = db.Where("department_id = ?", query.DepartmentID)
	}

	return db
}

func (r *userRepositoryImpl) applySorting(db *gorm.DB, query dto.UserPaginationQuery) *gorm.DB {
	allowedSorts := map[string]string{
		"created_at": "created_at",
		"first_name": "first_name",
		"last_name":  "last_name",
	}

	sortField := "created_at"
	if field, ok := allowedSorts[query.Sort]; ok {
		sortField = field
	}

	order := "DESC"
	if strings.ToUpper(query.Order) == "ASC" {
		order = "ASC"
	}

	return db.Order(sortField + " " + order)
}
