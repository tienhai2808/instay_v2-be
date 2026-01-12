package mapper

import (
	"github.com/InstayPMS/backend/internal/application/dto"
	"github.com/InstayPMS/backend/internal/domain/model"
)

func ToBasicDepartmentResponse(dept *model.Department) *dto.BasicDepartmentResponse {
	if dept == nil {
		return nil
	}

	return &dto.BasicDepartmentResponse{
		ID:   dept.ID,
		Name: dept.Name,
	}
}

func ToBasicUserResponse(usr *model.User) *dto.BasicUserResponse {
	if usr == nil {
		return nil
	}

	return &dto.BasicUserResponse{
		ID:        usr.ID,
		Username:  usr.Username,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
	}
}

func ToUserResponse(usr *model.User) *dto.UserResponse {
	if usr == nil {
		return nil
	}

	return &dto.UserResponse{
		ID:         usr.ID,
		Email:      usr.Email,
		Phone:      usr.Phone,
		Username:   usr.Username,
		FirstName:  usr.FirstName,
		LastName:   usr.LastName,
		Role:       usr.Role,
		IsActive:   usr.IsActive,
		CreatedAt:  usr.CreatedAt,
		Department: ToBasicDepartmentResponse(usr.Department),
	}
}

func ToUserDetailsResponse(usr *model.User) *dto.UserDetailsResponse {
	if usr == nil {
		return nil
	}

	return &dto.UserDetailsResponse{
		ID:         usr.ID,
		Email:      usr.Email,
		Phone:      usr.Phone,
		Username:   usr.Username,
		FirstName:  usr.FirstName,
		LastName:   usr.LastName,
		Role:       usr.Role,
		IsActive:   usr.IsActive,
		CreatedAt:  usr.CreatedAt,
		UpdatedAt:  usr.UpdatedAt,
		Department: ToBasicDepartmentResponse(usr.Department),
		CreatedBy:  ToBasicUserResponse(usr.CreatedBy),
		UpdatedBy:  ToBasicUserResponse(usr.UpdatedBy),
	}
}
