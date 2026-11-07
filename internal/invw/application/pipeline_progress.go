package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/invw/domain"
	"github.com/CSCE-3701-Amara/amara/internal/invw/ports"
)

var (
	ErrInvalidPipelineProgressQuery = errors.New("invalid pipeline progress query")
	ErrPipelineProgressNotFound     = errors.New("pipeline progress not found")
)

type CandidatePipelineProgress struct {
	ApplicationID string                `json:"application_id"`
	Pipeline      PipelineVersion       `json:"pipeline"`
	Stages        []CandidateStage      `json:"stages"`
	Links         []domain.PipelineLink `json:"links"`
}

type PipelineVersion struct {
	ID          string     `json:"id"`
	Version     int32      `json:"version"`
	PublishedAt *time.Time `json:"published_at"`
}

type CandidateStage struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Type        string                     `json:"type"`
	Status      domain.PipelineStageStatus `json:"status"`
	IsCurrent   bool                       `json:"is_current"`
	CompletedAt *time.Time                 `json:"completed_at,omitempty"`
	Passed      *bool                      `json:"passed,omitempty"`
	Awards      []domain.AmaraAward        `json:"amarat"`
}

type CandidatePipelineProgressService struct {
	repository ports.CandidatePipelineProgressRepository
}

func NewCandidatePipelineProgress(repository ports.CandidatePipelineProgressRepository) *CandidatePipelineProgressService {
	return &CandidatePipelineProgressService{repository: repository}
}

func (s *CandidatePipelineProgressService) Get(ctx context.Context, candidateID, applicationID string) (CandidatePipelineProgress, error) {
	if s == nil || s.repository == nil {
		return CandidatePipelineProgress{}, fmt.Errorf("candidate pipeline progress dependencies are not configured")
	}
	if !validUUID(candidateID) || !validUUID(applicationID) {
		return CandidatePipelineProgress{}, fmt.Errorf("%w: candidate and application IDs must be UUIDs", ErrInvalidPipelineProgressQuery)
	}
	pipeline, err := s.repository.GetByCandidateApplication(ctx, candidateID, applicationID)
	if err != nil {
		return CandidatePipelineProgress{}, fmt.Errorf("get candidate pipeline progress: %w", err)
	}
	if len(pipeline.Stages) == 0 {
		return CandidatePipelineProgress{}, fmt.Errorf("%w: the published pipeline has no stages", ErrPipelineProgressNotFound)
	}

	incoming := make(map[string]int, len(pipeline.Stages))
	passed := make(map[string]bool, len(pipeline.Stages))
	for _, stage := range pipeline.Stages {
		if stage.Progress != nil && stage.Progress.Completed {
			passed[stage.ID] = stage.Progress.Passed
		}
	}
	for _, link := range pipeline.Links {
		incoming[link.ToStageID]++
	}

	result := CandidatePipelineProgress{
		ApplicationID: pipeline.ApplicationID,
		Pipeline:      PipelineVersion{ID: pipeline.Pipeline.ID, Version: pipeline.Pipeline.Version, PublishedAt: pipeline.Pipeline.PublishedAt},
		Stages:        make([]CandidateStage, 0, len(pipeline.Stages)),
		Links:         nonNilLinks(pipeline.Links),
	}
	for _, stage := range pipeline.Stages {
		status := domain.PipelineStageLocked
		var completedAt *time.Time
		var stagePassed *bool
		if stage.Progress != nil && stage.Progress.Completed {
			status = domain.PipelineStageCompleted
			completedAt = stage.Progress.CompletedAt
			value := stage.Progress.Passed
			stagePassed = &value
		} else if stage.Progress != nil && stage.Progress.Active {
			status = domain.PipelineStageActive
		} else if incoming[stage.ID] == 0 || allPredecessorsPassed(stage.ID, pipeline.Links, passed) {
			status = domain.PipelineStageActive
		}
		awards := []domain.AmaraAward{}
		if status == domain.PipelineStageCompleted && stagePassed != nil && *stagePassed {
			awards = nonNilAwards(stage.Awards)
		}
		result.Stages = append(result.Stages, CandidateStage{
			ID: stage.ID, Name: stage.Name, Type: stage.Type,
			Status: status, IsCurrent: status == domain.PipelineStageActive,
			CompletedAt: completedAt, Passed: stagePassed,
			Awards: awards,
		})
	}
	return result, nil
}

func allPredecessorsPassed(stageID string, links []domain.PipelineLink, passed map[string]bool) bool {
	foundPredecessor := false
	for _, link := range links {
		if link.ToStageID != stageID {
			continue
		}
		foundPredecessor = true
		if !passed[link.FromStageID] {
			return false
		}
	}
	return foundPredecessor
}

func nonNilAwards(awards []domain.AmaraAward) []domain.AmaraAward {
	if awards == nil {
		return []domain.AmaraAward{}
	}
	return awards
}

func nonNilLinks(links []domain.PipelineLink) []domain.PipelineLink {
	if links == nil {
		return []domain.PipelineLink{}
	}
	return links
}

func validUUID(value string) bool {
	if len(value) != 36 || value[8] != '-' || value[13] != '-' || value[18] != '-' || value[23] != '-' {
		return false
	}
	for i, r := range value {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue
		}
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return false
		}
	}
	return true
}
