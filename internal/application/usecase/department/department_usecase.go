package department

import (
	"context"

	"github.com/InstayPMS/backend/internal/application/dto"
)

type DepartmentUseCase interface {
	CreateDepartment(ctx context.Context, userID int64, req dto.CreateDepartmentRequest) (int64, error)
}
