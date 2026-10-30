package domain

import "errors"

var (
	ErrInvalidMFACode    = errors.New("invalid MFA code")
	ErrMFARequired       = errors.New("MFA is required")
	ErrMFAAlreadyEnabled = errors.New("MFA is already enabled")
)

type MFAEnrollment struct {
	AccountID       string
	AccountType     AccountType
	Email           string
	ProvisioningURI string
}

func MFARequired(accountType AccountType) bool {
	return accountType == AccountTypeInterviewer || accountType == AccountTypeCompany
}
