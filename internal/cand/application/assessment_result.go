package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/CSCE-3701-Amara/amara/internal/cand/domain"
	"github.com/CSCE-3701-Amara/amara/internal/cand/ports"
)

var ErrInvalidAssessmentResultQuery = errors.New("invalid assessment result query")

type AssessmentResultItem struct {
	AssessmentID string                            `json:"assessment_id"`
	Status       domain.AssessmentEvaluationStatus `json:"status"`
	Outcome      *AssessmentOutcome                `json:"outcome,omitempty"`
	Feedback     *AssessmentFeedback               `json:"feedback,omitempty"`
}

type AssessmentOutcome struct {
	Passed bool `json:"passed"`
}

type AssessmentFeedback struct {
	Text                *string  `json:"text,omitempty"`
	Strengths           []string `json:"strengths"`
	AreasForImprovement []string `json:"areas_for_improvement"`
	AIGenerated         bool     `json:"ai_generated"`
}

type AssessmentResultPage struct {
	Items []AssessmentResultItem `json:"items"`
}

type AssessmentResults struct {
	repository ports.AssessmentResultRepository
}

func NewAssessmentResults(repository ports.AssessmentResultRepository) *AssessmentResults {
	return &AssessmentResults{repository: repository}
}

func (s *AssessmentResults) List(ctx context.Context, candidateID, applicationID string) (AssessmentResultPage, error) {
	if s == nil || s.repository == nil {
		return AssessmentResultPage{}, fmt.Errorf("assessment result dependencies are not configured")
	}
	if !validUUID(candidateID) || !validUUID(applicationID) {
		return AssessmentResultPage{}, fmt.Errorf("%w: candidate and application IDs must be UUIDs", ErrInvalidAssessmentResultQuery)
	}
	results, applicationFound, err := s.repository.ListByCandidateApplication(ctx, candidateID, applicationID)
	if err != nil {
		return AssessmentResultPage{}, fmt.Errorf("list assessment results: %w", err)
	}
	if !applicationFound {
		return AssessmentResultPage{}, domain.ErrAssessmentNotFound
	}

	page := AssessmentResultPage{Items: make([]AssessmentResultItem, 0, len(results))}
	for _, result := range results {
		page.Items = append(page.Items, mapAssessmentResult(result))
	}
	return page, nil
}

func mapAssessmentResult(result domain.AssessmentResult) AssessmentResultItem {
	item := AssessmentResultItem{AssessmentID: result.AssessmentID, Status: result.Status}
	if result.Outcome != nil {
		item.Outcome = &AssessmentOutcome{Passed: result.Outcome.Passed}
	}
	if result.Feedback != nil {
		item.Feedback = &AssessmentFeedback{
			Text: result.Feedback.Text, Strengths: nonNilStrings(result.Feedback.Strengths),
			AreasForImprovement: nonNilStrings(result.Feedback.AreasForImprovement),
			AIGenerated:         result.Feedback.AIGenerated,
		}
	}
	return item
}

func nonNilStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}
