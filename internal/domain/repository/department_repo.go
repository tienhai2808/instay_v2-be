package repository

import (
	"context"

	"github.com/InstayPMS/backend/internal/domain/model"
)

type DepartmentRepository interface {
	FindByID(ctx context.Context, id int64) (*model.Department, error)

	Create(ctx context.Context, dept *model.Department) error
}
