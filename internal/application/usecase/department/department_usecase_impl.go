package department

import (
	"context"

	"github.com/InstayPMS/backend/internal/application/dto"
	"github.com/InstayPMS/backend/internal/domain/model"
	"github.com/InstayPMS/backend/internal/domain/repository"
	customErr "github.com/InstayPMS/backend/pkg/errors"
	"github.com/InstayPMS/backend/pkg/utils"
	"github.com/sony/sonyflake/v2"
	"go.uber.org/zap"
)

type departmentUseCaseImpl struct {
	log            *zap.Logger
	idGen          *sonyflake.Sonyflake
	departmentRepo repository.DepartmentRepository
}

func NewDepartmentUseCase(
	log *zap.Logger,
	idGen *sonyflake.Sonyflake,
	departmentRepo repository.DepartmentRepository,
) DepartmentUseCase {
	return &departmentUseCaseImpl{
		log,
		idGen,
		departmentRepo,
	}
}

func (u *departmentUseCaseImpl) CreateDepartment(ctx context.Context, userID int64, req dto.CreateDepartmentRequest) (int64, error) {
	id, err := u.idGen.NextID()
	if err != nil {
		u.log.Error("generate department id failed", zap.Error(err))
		return 0, err
	}

	dept := &model.Department{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
		CreatedByID: &userID,
		UpdatedByID: &userID,
	}

	if err = u.departmentRepo.Create(ctx, dept); err != nil {
		if ok, constraint := utils.IsUniqueViolation(err); ok {
			switch constraint {
			case "departments_name_key":
				return 0, customErr.ErrNameAlreadyExists
			case "departments_phone_key":
				return 0, customErr.ErrPhoneAlreadyExists
			}
		}
		u.log.Error("create department failed", zap.Error(err))
		return 0, err
	}

	return id, nil
}
