package domain

import "time"

type AssessmentEvaluationStatus string

const (
	AssessmentEvaluationPending   AssessmentEvaluationStatus = "pending"
	AssessmentEvaluationEvaluated AssessmentEvaluationStatus = "evaluated"
)

type AssessmentResult struct {
	AssessmentID string
	Status       AssessmentEvaluationStatus
	Outcome      *AssessmentOutcome
	Feedback     *AssessmentFeedback
	SubmittedAt  *time.Time
	EvaluatedAt  *time.Time
}

type AssessmentOutcome struct {
	Passed bool
}

type AssessmentFeedback struct {
	Text                *string
	Strengths           []string
	AreasForImprovement []string
	AIGenerated         bool
	DisclosedAt         time.Time
}
