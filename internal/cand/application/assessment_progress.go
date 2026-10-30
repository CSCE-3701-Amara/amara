package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/cand/domain"
	"github.com/CSCE-3701-Amara/amara/internal/cand/ports"
)

var ErrInvalidAssessmentProgressQuery = errors.New("invalid assessment progress query")

type AssessmentProgressItem struct {
	ID                   string                  `json:"id"`
	Status               domain.AssessmentStatus `json:"status"`
	SupportsInterruption bool                    `json:"supports_interruption"`
	CanResume            bool                    `json:"can_resume"`
	SavedProgress        json.RawMessage         `json:"saved_progress,omitempty"`
	UpdatedAt            time.Time               `json:"updated_at"`
}

type AssessmentProgressPage struct {
	Items []AssessmentProgressItem `json:"items"`
}

type AssessmentProgress struct {
	repository ports.AssessmentProgressRepository
}

func NewAssessmentProgress(repository ports.AssessmentProgressRepository) *AssessmentProgress {
	return &AssessmentProgress{repository: repository}
}

func (s *AssessmentProgress) List(ctx context.Context, candidateID, applicationID string) (AssessmentProgressPage, error) {
	if s == nil || s.repository == nil {
		return AssessmentProgressPage{}, fmt.Errorf("assessment progress dependencies are not configured")
	}
	if !validUUID(candidateID) || !validUUID(applicationID) {
		return AssessmentProgressPage{}, fmt.Errorf("%w: candidate and application IDs must be UUIDs", ErrInvalidAssessmentProgressQuery)
	}
	assessments, applicationFound, err := s.repository.ListByCandidateApplication(ctx, candidateID, applicationID)
	if err != nil {
		return AssessmentProgressPage{}, fmt.Errorf("list assessment progress: %w", err)
	}
	if !applicationFound {
		return AssessmentProgressPage{}, domain.ErrAssessmentNotFound
	}

	page := AssessmentProgressPage{Items: make([]AssessmentProgressItem, 0, len(assessments))}
	for _, assessment := range assessments {
		page.Items = append(page.Items, mapAssessmentProgress(assessment))
	}
	return page, nil
}

func (s *AssessmentProgress) Save(ctx context.Context, candidateID, applicationID, assessmentID string, progress json.RawMessage) (AssessmentProgressItem, error) {
	if s == nil || s.repository == nil {
		return AssessmentProgressItem{}, fmt.Errorf("assessment progress dependencies are not configured")
	}
	if !validUUID(candidateID) || !validUUID(applicationID) || !validUUID(assessmentID) {
		return AssessmentProgressItem{}, fmt.Errorf("%w: candidate, application, and assessment IDs must be UUIDs", ErrInvalidAssessmentProgressQuery)
	}
	if !validProgressObject(progress) {
		return AssessmentProgressItem{}, fmt.Errorf("%w: saved_progress must be a JSON object", ErrInvalidAssessmentProgressQuery)
	}

	assessment, err := s.repository.SaveProgress(ctx, candidateID, applicationID, assessmentID, progress)
	if err != nil {
		return AssessmentProgressItem{}, fmt.Errorf("save assessment progress: %w", err)
	}
	return mapAssessmentProgress(assessment), nil
}

func validProgressObject(progress json.RawMessage) bool {
	var object map[string]json.RawMessage
	return json.Valid(progress) && json.Unmarshal(progress, &object) == nil && object != nil
}

func mapAssessmentProgress(assessment domain.AssessmentProgress) AssessmentProgressItem {
	return AssessmentProgressItem{
		ID: assessment.ID, Status: assessment.Status,
		SupportsInterruption: assessment.SupportsInterruption,
		CanResume:            assessment.CanResume(), SavedProgress: assessment.SavedProgress,
		UpdatedAt: assessment.UpdatedAt,
	}
}
