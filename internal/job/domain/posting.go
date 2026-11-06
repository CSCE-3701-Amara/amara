package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	PostingDraft     = "draft"
	PostingPublished = "published"
	PostingArchived  = "archived"

	EmploymentOnSiteFullTime = "on-site full-time"
	EmploymentOnSitePartTime = "on-site part-time"
	EmploymentRemoteFullTime = "remote full-time"
	EmploymentRemotePartTime = "remote part-time"
)

var ErrPostingNotFound = errors.New("job posting not found")
var ErrInvalidPosting = errors.New("invalid job posting")
var ErrInvalidTransition = errors.New("invalid job posting state transition")

// Posting is the JOB-owned management view of a posting and its screening data.
type Posting struct {
	ID                  string        `json:"id"`
	CompanyID           string        `json:"company_id"`
	Title               string        `json:"title"`
	Description         *string       `json:"description"`
	Language            *string       `json:"language"`
	Currency            *string       `json:"currency"`
	SeniorityLevel      string        `json:"seniority_level"`
	EmploymentType      string        `json:"employment_type"`
	Location            *string       `json:"location"`
	SalaryMin           *int32        `json:"salary_min"`
	SalaryMax           *int32        `json:"salary_max"`
	SalaryVisible       *bool         `json:"salary_visible"`
	OpenPositions       *int32        `json:"open_positions"`
	ApplicationOpening  *time.Time    `json:"application_opening"`
	ApplicationDeadline *time.Time    `json:"application_deadline"`
	MaxApplications     *int32        `json:"max_applications"`
	Status              string        `json:"status"`
	LastAutosavedAt     *time.Time    `json:"last_autosaved_at"`
	PublishedAt         *time.Time    `json:"published_at"`
	ArchivedAt          *time.Time    `json:"archived_at"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
	Requirements        []Requirement `json:"requirements"`
	Amarat              []Amara       `json:"amarat"`
}

type Requirement struct {
	ID   string `json:"id,omitempty"`
	Text string `json:"text"`
	Type string `json:"requirement_type"`
}

type Amara struct {
	ID          string  `json:"id,omitempty"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
}

func (p Posting) Normalized() Posting {
	p.Title = strings.TrimSpace(p.Title)
	p.SeniorityLevel = strings.TrimSpace(p.SeniorityLevel)
	p.EmploymentType = strings.ToLower(strings.TrimSpace(p.EmploymentType))

	for i := range p.Requirements {
		p.Requirements[i].Text = strings.TrimSpace(p.Requirements[i].Text)
		p.Requirements[i].Type = strings.ToLower(strings.TrimSpace(p.Requirements[i].Type))
	}
	for i := range p.Amarat {
		p.Amarat[i].Title = strings.TrimSpace(p.Amarat[i].Title)
	}
	return p
}

// ValidateForSave enforces this feature's required fields without changing the
// nullability of the pre-existing database columns.
func (p Posting) ValidateForSave() error {
	p = p.Normalized()
	if p.Title == "" {
		return fmt.Errorf("%w: title is required", ErrInvalidPosting)
	}
	if p.SeniorityLevel == "" {
		return fmt.Errorf("%w: seniority_level is required", ErrInvalidPosting)
	}
	switch p.EmploymentType {
	case EmploymentOnSiteFullTime, EmploymentOnSitePartTime, EmploymentRemoteFullTime, EmploymentRemotePartTime:
	default:
		return fmt.Errorf("%w: employment_type must be on-site full-time, on-site part-time, remote full-time, or remote part-time", ErrInvalidPosting)
	}
	if len(p.Amarat) == 0 {
		return fmt.Errorf("%w: at least one Amara is required", ErrInvalidPosting)
	}
	if p.SalaryMin != nil && *p.SalaryMin < 0 || p.SalaryMax != nil && *p.SalaryMax < 0 {
		return fmt.Errorf("%w: salary values cannot be negative", ErrInvalidPosting)
	}
	if p.SalaryMin != nil && p.SalaryMax != nil && *p.SalaryMin > *p.SalaryMax {
		return fmt.Errorf("%w: salary_min cannot exceed salary_max", ErrInvalidPosting)
	}
	if p.OpenPositions != nil && *p.OpenPositions < 1 {
		return fmt.Errorf("%w: open_positions must be positive", ErrInvalidPosting)
	}
	if p.MaxApplications != nil && *p.MaxApplications < 1 {
		return fmt.Errorf("%w: max_applications must be positive", ErrInvalidPosting)
	}
	if p.ApplicationOpening != nil && p.ApplicationDeadline != nil && !p.ApplicationDeadline.After(*p.ApplicationOpening) {
		return fmt.Errorf("%w: application_deadline must be after application_opening", ErrInvalidPosting)
	}
	for i := range p.Requirements {
		p.Requirements[i].Text = strings.TrimSpace(p.Requirements[i].Text)
		if p.Requirements[i].Text == "" {
			return fmt.Errorf("%w: requirement text is required", ErrInvalidPosting)
		}
		if p.Requirements[i].Type != "must_have" && p.Requirements[i].Type != "nice_to_have" {
			return fmt.Errorf("%w: requirement_type must be must_have or nice_to_have", ErrInvalidPosting)
		}
	}
	seen := make(map[string]struct{}, len(p.Amarat))
	for i := range p.Amarat {
		p.Amarat[i].Title = strings.TrimSpace(p.Amarat[i].Title)
		if p.Amarat[i].Title == "" {
			return fmt.Errorf("%w: Amara title is required", ErrInvalidPosting)
		}
		key := strings.ToLower(p.Amarat[i].Title)
		if _, exists := seen[key]; exists {
			return fmt.Errorf("%w: Amara titles must be unique within a posting", ErrInvalidPosting)
		}
		seen[key] = struct{}{}
	}
	return nil
}
