package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidProfile  = errors.New("invalid profile")
	ErrProfileNotFound = errors.New("profile not found")
)

type CandidateProfile struct {
	ID                string
	Name              string
	Email             string
	CreatedAt         time.Time
	UpdatedAt         *time.Time
	PhoneNumber       string
	GithubProfile     string
	LinkedinProfile   string
	PortfolioLink     string
	Gender            string
	BirthDate         *time.Time
	Nationality       string
	Address           string
	Biography         string
	YearsOfExperience *int
	Available         *bool
	Amarat            []string
}

type CandidateProfilePatch struct {
	Name              *string
	PhoneNumber       *string
	GithubProfile     *string
	LinkedinProfile   *string
	PortfolioLink     *string
	Gender            *string
	BirthDate         *time.Time
	Nationality       *string
	Address           *string
	Biography         *string
	YearsOfExperience *int32
	Available         *bool
}

type CompanyProfile struct {
	ID              string
	LegalName       string
	DisplayName     *string
	Description     *string
	WebsiteLink     *string
	LogoURL         *string
	ContactPhone    *string
	ContactEmail    *string
	LinkedinProfile *string
	Approved        bool
}

type CompanyProfilePatch struct {
	LegalName            *string
	DisplayName          *string
	ClearDisplayName     bool
	Description          *string
	ClearDescription     bool
	WebsiteLink          *string
	ClearWebsiteLink     bool
	LogoURL              *string
	ClearLogoURL         bool
	ContactPhone         *string
	ClearContactPhone    bool
	ContactEmail         *string
	ClearContactEmail    bool
	LinkedinProfile      *string
	ClearLinkedinProfile bool
}

type InterviewerProfile struct {
	ID            string
	CompanyID     string
	Name          string
	Email         string
	AccountStatus string
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

func ValidateProfileText(value string, max int) bool {
	return strings.TrimSpace(value) != "" && len([]rune(value)) <= max
}
