package apperrors

import (
	"errors"
	"fmt"
	"testing"
)

func TestAppError(t *testing.T) {
	internalErr := errors.New("db query failed: connection closed")
	appErr := Wrap(internalErr, ErrNotFound, "The requested resource was not found.", "job_not_found")

	if !errors.Is(appErr, ErrNotFound) {
		t.Fatal("expected errors.Is(appErr, ErrNotFound) to be true")
	}
	if !errors.Is(appErr, internalErr) {
		t.Fatal("expected errors.Is(appErr, internalErr) to be true")
	}

	wrapped := fmt.Errorf("application layer: %w", appErr)
	if !errors.Is(wrapped, ErrNotFound) {
		t.Fatal("expected errors.Is(wrapped, ErrNotFound) to be true")
	}

	var target *AppError
	if !errors.As(wrapped, &target) {
		t.Fatal("expected errors.As(wrapped, &target) to match *AppError")
	}
	if target.Code != "job_not_found" {
		t.Fatalf("expected code %q, got %q", "job_not_found", target.Code)
	}
	if target.Message != "The requested resource was not found." {
		t.Fatalf("expected message %q, got %q", "The requested resource was not found.", target.Message)
	}
}

func TestValidationError(t *testing.T) {
	valErr := Validation("Invalid fields",
		FieldError{Field: "email", Message: "Email format is invalid."},
		FieldError{Field: "password", Message: "Password is too short."},
	)

	if !errors.Is(valErr, ErrValidation) {
		t.Fatal("expected errors.Is(valErr, ErrValidation) to be true")
	}

	wrapped := fmt.Errorf("service: %w", valErr)
	if !errors.Is(wrapped, ErrValidation) {
		t.Fatal("expected errors.Is(wrapped, ErrValidation) to be true")
	}

	var target *ValidationError
	if !errors.As(wrapped, &target) {
		t.Fatal("expected errors.As to match *ValidationError")
	}
	if len(target.Errors) != 2 {
		t.Fatalf("expected 2 field errors, got %d", len(target.Errors))
	}
	if target.Errors[0].Field != "email" || target.Errors[1].Field != "password" {
		t.Fatalf("unexpected field errors: %+v", target.Errors)
	}
}

func TestAppErrorConstructors(t *testing.T) {
	tests := []struct {
		name     string
		err      *AppError
		expected error
	}{
		{"NotFound", NotFound("not found", "code_404"), ErrNotFound},
		{"Conflict", Conflict("conflict", "code_409"), ErrConflict},
		{"Forbidden", Forbidden("forbidden", "code_403"), ErrForbidden},
		{"Unauthenticated", Unauthenticated("unauth", "code_401"), ErrUnauthenticated},
		{"InvalidInput", InvalidInput("bad input", "code_400"), ErrInvalidInput},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !errors.Is(tt.err, tt.expected) {
				t.Fatalf("expected errors.Is(err, %v) to be true", tt.expected)
			}
		})
	}
}
