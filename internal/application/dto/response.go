package dto

import (
	"time"

	"github.com/InstayPMS/backend/internal/domain/model"
)

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type UploadPresignedURLResponse struct {
	Url string `json:"url"`
	Key string `json:"key"`
}

type ViewPresignedURLResponse struct {
	Url string `json:"url"`
}

type UserResponse struct {
	ID         int64                    `json:"id"`
	Username   string                   `json:"username"`
	Email      string                   `json:"email"`
	Phone      string                   `json:"phone"`
	Role       model.UserRole           `json:"role"`
	IsActive   bool                     `json:"is_active"`
	FirstName  string                   `json:"first_name"`
	LastName   string                   `json:"last_name"`
	CreatedAt  time.Time                `json:"created_at"`
	Department *BasicDepartmentResponse `json:"department"`
}

type UserDetailsResponse struct {
	ID         int64                    `json:"id"`
	Username   string                   `json:"username"`
	Email      string                   `json:"email"`
	Phone      string                   `json:"phone"`
	Role       model.UserRole           `json:"role"`
	IsActive   bool                     `json:"is_active"`
	FirstName  string                   `json:"first_name"`
	LastName   string                   `json:"last_name"`
	CreatedAt  time.Time                `json:"created_at"`
	UpdatedAt  time.Time                `json:"updated_at"`
	Department *BasicDepartmentResponse `json:"department"`
	CreatedBy  *BasicUserResponse       `json:"created_by"`
	UpdatedBy  *BasicUserResponse       `json:"updated_by"`
}

type BasicUserResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type BasicDepartmentResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
