package domain

import "fmt"

// AppError is the base application error type used throughout the service layer.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Common error constructors.

// ErrNotFound indicates the requested resource does not exist.
func ErrNotFound(resource string) *AppError {
	return &AppError{Code: 404, Message: fmt.Sprintf("%s not found", resource)}
}

// ErrConflict indicates a unique constraint violation.
func ErrConflict(msg string) *AppError {
	return &AppError{Code: 409, Message: msg}
}

// ErrBadRequest indicates invalid client input.
func ErrBadRequest(msg string) *AppError {
	return &AppError{Code: 400, Message: msg}
}

// ErrUnauthorized indicates missing or invalid credentials.
func ErrUnauthorized(msg string) *AppError {
	return &AppError{Code: 401, Message: msg}
}

// ErrForbidden indicates the user lacks permission.
func ErrForbidden(msg string) *AppError {
	return &AppError{Code: 403, Message: msg}
}

// ErrInternal wraps an unexpected server-side error.
func ErrInternal(err error) *AppError {
	return &AppError{Code: 500, Message: "internal server error", Err: err}
}
