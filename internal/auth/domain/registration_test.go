package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestRegistrationValidate(t *testing.T) {
	tests := []struct {
		name         string
		registration Registration
		wantErr      bool
	}{
		{
			name: "valid candidate",
			registration: Registration{
				AccountType: AccountTypeCandidate, Email: "candidate@example.com",
				Password: "correct horse", FirstName: "Ada", LastName: "Lovelace",
			},
		},
		{
			name: "valid interviewer",
			registration: Registration{
				AccountType: AccountTypeInterviewer, Email: "interviewer@example.com",
				Password: "correct horse", FirstName: "Alan", LastName: "Turing", CompanySecret: "company-secret",
			},
		},
		{name: "invalid account type", registration: Registration{AccountType: "admin"}, wantErr: true},
		{name: "invalid email", registration: Registration{AccountType: AccountTypeCandidate, Email: "not-an-email"}, wantErr: true},
		{name: "short password", registration: Registration{AccountType: AccountTypeCandidate, Email: "a@b.com", Password: "short", FirstName: "A", LastName: "B"}, wantErr: true},
		{name: "long password", registration: Registration{AccountType: AccountTypeCandidate, Email: "a@b.com", Password: strings.Repeat("x", 129), FirstName: "A", LastName: "B"}, wantErr: true},
		{name: "missing first name", registration: Registration{AccountType: AccountTypeCandidate, Email: "a@b.com", Password: "long enough", LastName: "B"}, wantErr: true},
		{name: "long first name", registration: Registration{AccountType: AccountTypeCandidate, Email: "a@b.com", Password: "long enough", FirstName: strings.Repeat("x", 101), LastName: "B"}, wantErr: true},
		{name: "missing last name", registration: Registration{AccountType: AccountTypeCandidate, Email: "a@b.com", Password: "long enough", FirstName: "A"}, wantErr: true},
		{name: "long last name", registration: Registration{AccountType: AccountTypeCandidate, Email: "a@b.com", Password: "long enough", FirstName: "A", LastName: strings.Repeat("x", 101)}, wantErr: true},
		{name: "interviewer missing company", registration: Registration{AccountType: AccountTypeInterviewer, Email: "a@b.com", Password: "long enough", FirstName: "A", LastName: "B"}, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.registration.Validate()
			if test.wantErr != (err != nil) {
				t.Fatalf("Validate() error = %v, wantErr %v", err, test.wantErr)
			}
			if test.wantErr && !errors.Is(err, ErrInvalidRegistration) {
				t.Fatalf("expected ErrInvalidRegistration, got %v", err)
			}
		})
	}
}

func TestRegistrationNormalized(t *testing.T) {
	registration := Registration{
		Email: "  USER@Example.COM ", FirstName: " Ada ", LastName: " Lovelace ", CompanySecret: " secret ",
	}
	got := registration.Normalized()
	if got.Email != "user@example.com" || got.FirstName != "Ada" || got.LastName != "Lovelace" || got.CompanySecret != "secret" {
		t.Fatalf("unexpected normalized registration: %+v", got)
	}
}
