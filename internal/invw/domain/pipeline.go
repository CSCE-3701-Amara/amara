package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const AssessmentStage = "assessment"
const LiveInterviewStage = "live_interview"

var (
	ErrPipelineNotFound = errors.New("pipeline resource not found")
	ErrPipelineExists   = errors.New("the posting already has a pipeline")
	ErrPipelineLocked   = errors.New("only an unpublished draft can be edited")
)

type InvalidPipelineError struct {
	Reason string
}

func (e *InvalidPipelineError) Error() string { return e.Reason }

type PipelineDraft struct {
	ID            string `json:"id"`
	PostingID     string `json:"postingId"`
	VersionID     string `json:"versionId"`
	VersionNumber int    `json:"versionNumber"`
	Published     bool   `json:"published"`
}

type Stage struct {
	Rule           string         `json:"advancementRule,omitempty"`
	ScoreThreshold *int           `json:"scoreThreshold,omitempty"`
	Below          string         `json:"belowThresholdAction,omitempty"`
	X              float64        `json:"x"`
	Y              float64        `json:"y"`
	ID             string         `json:"id,omitempty"`
	Name           string         `json:"name"`
	Type           string         `json:"type"`
	OpensAt        time.Time      `json:"scheduledStartAt"`
	Deadline       time.Time      `json:"scheduledEndAt"`
	ModuleIDs      []string       `json:"moduleIds,omitempty"`
	Interview      *LiveInterview `json:"interview,omitempty"`
}

type LiveInterview struct {
	Behavioral      bool        `json:"behavioral,omitempty"`
	DurationMinutes int         `json:"durationMinutes"`
	InterviewerIDs  []string    `json:"interviewerIds"`
	Threshold       *int        `json:"threshold"`
	Rubric          []Criterion `json:"rubric"`
	Amarat          []Amara     `json:"amarat"`
}

func ValidateStages(stages []Stage) error {
	names := make(map[string]bool, len(stages))
	for i := range stages {
		stage := &stages[i]
		stage.Name = strings.TrimSpace(stage.Name)
		if stage.Name == "" || names[stage.Name] {
			return &InvalidPipelineError{fmt.Sprintf("stage %d needs a unique non-empty name", i+1)}
		}
		names[stage.Name] = true
		if err := stage.validate(); err != nil {
			return &InvalidPipelineError{fmt.Sprintf("stage %d: %s", i+1, err)}
		}
	}
	return nil
}

func (s Stage) validate() error {
	if s.Rule != "" && s.Rule != "completion" && s.Rule != "manual" && s.Rule != "threshold" {
		return errors.New("invalid advancement rule")
	}
	if s.Below != "" && s.Below != "hold" && s.Below != "reject" {
		return errors.New("invalid below-threshold action")
	}
	if s.ScoreThreshold != nil && (*s.ScoreThreshold < 0 || *s.ScoreThreshold > 100) {
		return errors.New("score threshold must be between 0 and 100")
	}
	if s.Rule == "threshold" && s.ScoreThreshold == nil && s.Interview == nil {
		return errors.New("a threshold rule requires a score threshold")
	}

	if s.OpensAt.IsZero() || !s.Deadline.After(s.OpensAt) {
		return errors.New("the deadline must be after the opening time")
	}
	switch s.Type {
	case AssessmentStage:
		if s.Interview != nil || len(s.ModuleIDs) == 0 || !uniqueIDs(s.ModuleIDs) {
			return errors.New("an assessment needs distinct modules and no live-interview settings")
		}
	case LiveInterviewStage:
		if len(s.ModuleIDs) != 0 || s.Interview == nil {
			return errors.New("a live interview needs interview settings and cannot contain test modules")
		}
		i := s.Interview
		if i.DurationMinutes <= 0 || float64(i.DurationMinutes) > s.Deadline.Sub(s.OpensAt).Minutes() {
			return errors.New("the interview duration must fit within the stage window")
		}
		if len(i.InterviewerIDs) == 0 || !uniqueIDs(i.InterviewerIDs) {
			return errors.New("a live interview needs distinct assigned interviewers")
		}
		if i.Behavioral && i.Threshold != nil {
			return errors.New("a behavioral interview cannot have a score threshold")
		}
		if !i.Behavioral && (i.Threshold == nil || *i.Threshold < 0 || *i.Threshold > 100) {
			return errors.New("a live interview requires a threshold from 0 to 100")
		}
		return validateRubric(i.Rubric, i.Amarat)
	default:
		return errors.New("invalid stage type")
	}
	return nil
}

func uniqueIDs(ids []string) bool {
	seen := make(map[string]bool, len(ids))
	for _, id := range ids {
		if id == "" || seen[id] {
			return false
		}
		seen[id] = true
	}
	return true
}

func validateRubric(rubric []Criterion, amarat []Amara) error {
	if len(rubric) == 0 {
		return errors.New("a live interview requires a rubric")
	}
	mappings, err := newAmaraMappings(amarat)
	if err != nil {
		return err
	}
	if err := mappings.markRubric(rubric); err != nil {
		return err
	}
	return mappings.validateUsed(amarat)
}
