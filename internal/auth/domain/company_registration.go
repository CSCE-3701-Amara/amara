package domain

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"
)

const PendingReview = "pending_review"

type CompanyRegistration struct {
	CompanyName        string
	RegistrationNumber string
	Email              string
	Password           string
}

func (r CompanyRegistration) Validate() error {
	if strings.TrimSpace(r.CompanyName) == "" || utf8.RuneCountInString(r.CompanyName) > 200 {
		return fmt.Errorf("%w: companyName is invalid", ErrInvalidRegistration)
	}
	if strings.TrimSpace(r.RegistrationNumber) == "" || utf8.RuneCountInString(r.RegistrationNumber) > 100 {
		return fmt.Errorf("%w: registrationNumber is invalid", ErrInvalidRegistration)
	}
	address, err := mail.ParseAddress(r.Email)
	if err != nil || address.Address != r.Email || strings.ContainsAny(r.Email, " \r\n") {
		return fmt.Errorf("%w: email is invalid", ErrInvalidRegistration)
	}
	if utf8.RuneCountInString(r.Password) < 8 || utf8.RuneCountInString(r.Password) > 128 {
		return fmt.Errorf("%w: password must be between 8 and 128 characters", ErrInvalidRegistration)
	}
	return nil
}

func (r CompanyRegistration) Normalized() CompanyRegistration {
	r.CompanyName = strings.TrimSpace(r.CompanyName)
	r.RegistrationNumber = strings.TrimSpace(r.RegistrationNumber)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	return r
}
