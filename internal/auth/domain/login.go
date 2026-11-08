package domain

import "errors"

var (
	ErrEmailMFARequired   = errors.New("enter the six-digit sign-in code sent to your email")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrLoginRateLimited   = errors.New("login attempts rate limited")
	ErrInvalidSession     = errors.New("invalid session")
	ErrAccountInactive    = errors.New("account is not active")
	ErrAccountForbidden   = errors.New("account is forbidden")
)

type Login struct {
	Email    string
	Password string
	MFACode  string
}

type Session struct {
	Token       string
	AccountID   string
	AccountType AccountType
	Email       string
	ExpiresAt   string
}
