package errors

import "fmt"

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Errores comunes predefinidos
var (
	ErrInternalServer   = NewAppError(500, "internal server error", nil)
	ErrBadRequest       = NewAppError(400, "bad request", nil)
	ErrUnauthorized     = NewAppError(401, "unauthorized", nil)
	ErrForbidden        = NewAppError(403, "forbidden", nil)
	ErrNotFound         = NewAppError(404, "not found", nil)
	ErrValidationFailed = NewAppError(422, "validation failed", nil)
)
