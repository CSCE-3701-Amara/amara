package domain

import (
	"errors"
	"time"
)

const (
	ApplicationAdmitted    = "admitted"
	ApplicationScreenedOut = "screened_out"
)

var (
	ErrApplicationPostingUnavailable = errors.New("posting is not accepting applications")
	ErrApplicationLimitReached       = errors.New("posting application limit reached")
	ErrDuplicateApplication          = errors.New("candidate has already applied to this posting")
	ErrCandidateUnavailable          = errors.New("candidate is not eligible to apply")
)

// Application is JOB's submission and initial Amarat screening result. Later pipeline state is owned by INVW.
type Application struct {
	ID                   string    `json:"id"`
	ReferenceNumber      int32     `json:"reference_number"`
	PostingID            string    `json:"posting_id"`
	Status               string    `json:"status"`
	SubmittedAt          time.Time `json:"submitted_at"`
	ScreenedOut          bool      `json:"screened_out"`
	ScreeningCompletedAt time.Time `json:"screening_completed_at"`
	FailedAmarat         []string  `json:"failed_amarat"`
}
