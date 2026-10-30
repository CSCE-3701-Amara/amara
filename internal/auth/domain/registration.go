package domain

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"
)

type AccountType string

const (
	AccountTypeCandidate   AccountType = "candidate"
	AccountTypeInterviewer AccountType = "interviewer"
	AccountTypeCompany     AccountType = "company"
)

const PendingVerification = "pending_verification"

var (
	ErrInvalidRegistration       = errors.New("invalid registration")
	ErrEmailAlreadyExists        = errors.New("email already registered")
	ErrCompanyNotFound           = errors.New("approved company association not found")
	ErrCompanyRegistrationExists = errors.New("company registration already exists")
)

type Registration struct {
	AccountType   AccountType
	Email         string
	Password      string
	FirstName     string
	LastName      string
	CompanySecret string
}

func (r Registration) Validate() error {
	if r.AccountType != AccountTypeCandidate && r.AccountType != AccountTypeInterviewer {
		return fmt.Errorf("%w: account type is invalid", ErrInvalidRegistration)
	}
	if !validEmail(r.Email) {
		return fmt.Errorf("%w: email is invalid", ErrInvalidRegistration)
	}
	if utf8.RuneCountInString(r.Password) < 8 || utf8.RuneCountInString(r.Password) > 128 {
		return fmt.Errorf("%w: password must be between 8 and 128 characters", ErrInvalidRegistration)
	}
	if strings.TrimSpace(r.FirstName) == "" || utf8.RuneCountInString(r.FirstName) > 100 {
		return fmt.Errorf("%w: firstName is invalid", ErrInvalidRegistration)
	}
	if strings.TrimSpace(r.LastName) == "" || utf8.RuneCountInString(r.LastName) > 100 {
		return fmt.Errorf("%w: lastName is invalid", ErrInvalidRegistration)
	}
	if r.AccountType == AccountTypeInterviewer && strings.TrimSpace(r.CompanySecret) == "" {
		return fmt.Errorf("%w: companySecret is required", ErrInvalidRegistration)
	}
	return nil
}

func (r Registration) Normalized() Registration {
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.CompanySecret = strings.TrimSpace(r.CompanySecret)
	return r
}

func validEmail(value string) bool {
	address, err := mail.ParseAddress(value)
	return err == nil && address.Address == value && !strings.ContainsAny(value, " \r\n")
}

type Account struct {
	ID        string
	Type      AccountType
	Email     string
	Status    string
	CreatedAt string
}
