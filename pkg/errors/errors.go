package errors

import (
	"net/http"

	"github.com/InstayPMS/backend/pkg/constants"
)

var (
	ErrLoginFailed = NewAPIError(http.StatusBadRequest, constants.CodeLoginFailed, "Incorrect username or password")

	ErrInvalidToken = NewAPIError(http.StatusBadRequest, constants.CodeInvalidToken, "Invalid or expired token")

	ErrBadRequest = NewAPIError(http.StatusBadRequest, constants.CodeBadRequest, "Invalid data")

	ErrUnAuth = NewAPIError(http.StatusUnauthorized, constants.CodeUnAuth, "Unauthorized")

	ErrForbidden = NewAPIError(http.StatusForbidden, constants.CodeForbidden, "Forbidden")

	ErrInvalidPassword = NewAPIError(http.StatusBadRequest, constants.CodeInvalidPassword, "Incorrect password")

	ErrInvalidUser = NewAPIError(http.StatusForbidden, constants.CodeNoRefreshToken, "Please login again")

	ErrUserNotFound = NewAPIError(http.StatusNotFound, constants.CodeUserNotFound, "User not found")

	ErrEmailDoesNotExist = NewAPIError(http.StatusBadRequest, constants.CodeEmailDoesNotExist, "Email does not exist")

	ErrTooManyAttempts = NewAPIError(http.StatusTooManyRequests, constants.CodeTooManyAttempts, "Too many attempts")

	ErrInvalidOTP = NewAPIError(http.StatusBadRequest, constants.CodeInvalidOTP, "Invalid or expired OTP")

	ErrEmailAlreadyExists = NewAPIError(http.StatusConflict, constants.CodeEmailAlreadyExists, "Email already exists")

	ErrNameAlreadyExists = NewAPIError(http.StatusConflict, constants.CodeNameAlreadyExists, "Name already exists")

	ErrPhoneAlreadyExists = NewAPIError(http.StatusConflict, constants.CodePhoneAlreadyExists, "Phone already exists")

	ErrDepartmentNotFound = NewAPIError(http.StatusNotFound, constants.CodeDepartmentNotFound, "Department not found")

	ErrInvalidID = NewAPIError(http.StatusBadRequest, constants.CodeInvalidID, "Invalid id")

	ErrUsernameAlreadyExists = NewAPIError(http.StatusConflict, constants.CodeUsernameAlreadyExists, "Username already exists")
)

type APIError struct {
	Status  int
	Code    int
	Message string
	Data    any
}

func NewAPIError(status, code int, message string) *APIError {
	return &APIError{
		status,
		code,
		message,
		nil,
	}
}

func (e *APIError) Error() string {
	return e.Message
}

func (e *APIError) WithData(data any) *APIError {
	e.Data = data
	return e
}
