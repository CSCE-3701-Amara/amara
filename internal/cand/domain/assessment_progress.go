package domain

import (
	"encoding/json"
	"errors"
	"time"
)

type AssessmentStatus string

const (
	AssessmentStatusNotStarted AssessmentStatus = "not_started"
	AssessmentStatusInProgress AssessmentStatus = "in_progress"
	AssessmentStatusSubmitted  AssessmentStatus = "submitted"
	AssessmentStatusCompleted  AssessmentStatus = "completed"
)

var (
	ErrAssessmentNotFound     = errors.New("assessment not found")
	ErrAssessmentNotResumable = errors.New("assessment cannot be resumed")
)

// AssessmentProgress is an assessment assigned to an application and the
// candidate-owned state that can be restored when interruption is supported.
type AssessmentProgress struct {
	ID                   string
	Status               AssessmentStatus
	SupportsInterruption bool
	SavedProgress        json.RawMessage
	UpdatedAt            time.Time
}

func (a AssessmentProgress) CanResume() bool {
	return a.SupportsInterruption && a.Status == AssessmentStatusInProgress
}
