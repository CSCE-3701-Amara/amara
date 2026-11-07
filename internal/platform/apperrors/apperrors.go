package apperrors

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidInput indicates that request syntax is invalid or malformed.
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthenticated indicates missing or invalid authentication credentials.
	ErrUnauthenticated = errors.New("unauthenticated")

	// ErrUnauthorized is an alias for ErrUnauthenticated.
	ErrUnauthorized = ErrUnauthenticated

	// ErrForbidden indicates the authenticated caller lacks permission.
	ErrForbidden = errors.New("forbidden")

	// ErrNotFound indicates the requested resource does not exist.
	ErrNotFound = errors.New("resource not found")

	// ErrConflict indicates a resource conflict or duplicate operation.
	ErrConflict = errors.New("resource conflict")

	// ErrValidation indicates a semantically invalid request.
	ErrValidation = errors.New("validation failed")

	// ErrRateLimitExceeded indicates request rate limit has been exceeded.
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	// ErrInternal indicates an unexpected internal failure.
	ErrInternal = errors.New("internal server error")

	// ErrServiceUnavailable indicates the service is temporarily unavailable.
	ErrServiceUnavailable = errors.New("service unavailable")
)

// FieldError represents an error associated with a specific field.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationError represents semantic validation failures containing field-level errors.
type ValidationError struct {
	Message string       `json:"message"`
	Errors  []FieldError `json:"errors,omitempty"`
}

func (e *ValidationError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "validation failed"
}

func (e *ValidationError) Is(target error) bool {
	return target == ErrValidation
}

// AppError carries an error category, safe client message, optional machine-readable code,
// and optional internal cause.
type AppError struct {
	Category error
	Message  string
	Code     string
	Cause    error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *AppError) Unwrap() []error {
	var errs []error
	if e.Category != nil {
		errs = append(errs, e.Category)
	}
	if e.Cause != nil {
		errs = append(errs, e.Cause)
	}
	return errs
}

func (e *AppError) Is(target error) bool {
	if errors.Is(e.Category, target) {
		return true
	}
	if e.Cause != nil && errors.Is(e.Cause, target) {
		return true
	}
	return false
}

// New creates an AppError with category, safe message, and code.
func New(category error, message string, code string) *AppError {
	return &AppError{
		Category: category,
		Message:  message,
		Code:     code,
	}
}

// Wrap wraps an internal error with an AppError.
func Wrap(cause error, category error, message string, code string) *AppError {
	return &AppError{
		Category: category,
		Message:  message,
		Code:     code,
		Cause:    cause,
	}
}

// NotFound creates a not-found application error.
func NotFound(message, code string) *AppError {
	return New(ErrNotFound, message, code)
}

// Conflict creates a conflict application error.
func Conflict(message, code string) *AppError {
	return New(ErrConflict, message, code)
}

// Forbidden creates a forbidden application error.
func Forbidden(message, code string) *AppError {
	return New(ErrForbidden, message, code)
}

// Unauthenticated creates an unauthenticated application error.
func Unauthenticated(message, code string) *AppError {
	return New(ErrUnauthenticated, message, code)
}

// InvalidInput creates an invalid input application error.
func InvalidInput(message, code string) *AppError {
	return New(ErrInvalidInput, message, code)
}

// Validation creates a validation error with optional field errors.
func Validation(message string, errors ...FieldError) *ValidationError {
	return &ValidationError{
		Message: message,
		Errors:  errors,
	}
}
