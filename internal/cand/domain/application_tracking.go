package domain

import (
	"encoding/json"
	"errors"
	"time"
)

var ErrApplicationNotFound = errors.New("application not found")

type ApplicationSummary struct {
	ID                   string     `json:"id"`
	ReferenceNumber      int32      `json:"referenceNumber"`
	PostingID            string     `json:"postingId"`
	JobTitle             string     `json:"jobTitle"`
	CompanyName          string     `json:"companyName"`
	CreatedAt            *time.Time `json:"createdAt"`
	UpdatedAt            *time.Time `json:"updatedAt"`
	SubmittedAt          *time.Time `json:"submittedAt"`
	WithdrawnAt          *time.Time `json:"withdrawnAt"`
	ScreenedOut          *bool      `json:"screenedOut"`
	ScreeningCompletedAt *time.Time `json:"screeningCompletedAt"`
	Status               string     `json:"status"`
	Closed               bool       `json:"closed"`
	IdentityHidden       bool       `json:"identityHidden"`
}

type ApplicationDetail struct {
	ApplicationSummary
	PipelineExecutionID *string         `json:"pipelineExecutionId"`
	FinalDecision       *bool           `json:"finalDecision"`
	FinalFeedback       *string         `json:"finalFeedback,omitempty"`
	RejectionFeedback   json.RawMessage `json:"rejectionFeedback,omitempty"`
}

type ApplicationFilter struct {
	Limit  int
	Offset int
	Status string
}

type ApplicationPage struct {
	Data       []ApplicationSummary  `json:"data"`
	Pagination ApplicationPagination `json:"pagination"`
}

type ApplicationPagination struct {
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total"`
}
