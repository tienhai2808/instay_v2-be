package errors

type APIError struct {
	Status  int
	Code    int
	Slug    string
	Message string
}

func NewAPIError(status, code int, slug, message string) *APIError {
	return &APIError{
		status,
		code,
		slug,
		message,
	}
}

func (e *APIError) Error() string {
	return e.Message
}
