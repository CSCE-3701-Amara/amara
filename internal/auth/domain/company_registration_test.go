package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestCompanyRegistrationValidate(t *testing.T) {
	valid := CompanyRegistration{
		CompanyName: "Amara Ltd", RegistrationNumber: "REG-123",
		Email: "owner@example.com", Password: "safe password",
	}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	tests := []CompanyRegistration{
		{RegistrationNumber: "REG", Email: "owner@example.com", Password: "safe password"},
		{CompanyName: "Company", Email: "owner@example.com", Password: "safe password"},
		{CompanyName: "Company", RegistrationNumber: "REG", Email: "invalid", Password: "safe password"},
		{CompanyName: "Company", RegistrationNumber: "REG", Email: "owner@example.com", Password: "short"},
		{CompanyName: strings.Repeat("x", 201), RegistrationNumber: "REG", Email: "owner@example.com", Password: "safe password"},
		{CompanyName: "Company", RegistrationNumber: strings.Repeat("x", 101), Email: "owner@example.com", Password: "safe password"},
	}
	for _, registration := range tests {
		if err := registration.Validate(); !errors.Is(err, ErrInvalidRegistration) {
			t.Fatalf("expected invalid registration, got %v", err)
		}
	}
}

func TestCompanyRegistrationNormalized(t *testing.T) {
	got := (CompanyRegistration{
		CompanyName: " Company ", RegistrationNumber: " REG-123 ",
		Email: " OWNER@Example.COM ",
	}).Normalized()
	if got.CompanyName != "Company" || got.RegistrationNumber != "REG-123" || got.Email != "owner@example.com" {
		t.Fatalf("unexpected normalized registration: %+v", got)
	}
}
