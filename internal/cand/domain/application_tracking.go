package domain

import (
	"encoding/json"
	"errors"
	"time"
)

var ErrApplicationNotFound = errors.New("application not found")

type ApplicationSummary struct {
	ID              string     `json:"id"`
	ReferenceNumber int32      `json:"reference_number"`
	PostingID       string     `json:"posting_id"`
	JobTitle        string     `json:"job_title"`
	CompanyName     string     `json:"company_name"`
	SubmittedAt     *time.Time `json:"submitted_at"`
	Status          string     `json:"status"`
	Closed          bool       `json:"closed"`
}

type ApplicationDetail struct {
	ApplicationSummary
	PipelineExecutionID *string         `json:"pipeline_execution_id"`
	FinalDecision       *bool           `json:"final_decision"`
	FinalFeedback       *string         `json:"final_feedback"`
	RejectionFeedback   json.RawMessage `json:"rejection_feedback"`
}
