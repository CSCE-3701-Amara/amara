package domain

import (
	"errors"
	"time"
)

var ErrPipelineProgressNotFound = errors.New("pipeline progress not found")

type PipelineStageStatus string

const (
	PipelineStageCompleted PipelineStageStatus = "completed"
	PipelineStageActive    PipelineStageStatus = "active"
	PipelineStageLocked    PipelineStageStatus = "upcoming_locked"
)

type CandidatePipelineProgress struct {
	ApplicationID string
	Pipeline      PipelineVersion
	Stages        []PipelineStage
	Links         []PipelineLink
}

type PipelineVersion struct {
	ID          string
	Version     int32
	PublishedAt *time.Time
}

type PipelineStage struct {
	ID       string
	Name     string
	Type     string
	Progress *StageProgress
	Awards   []AmaraAward
}

type StageProgress struct {
	Completed   bool
	CompletedAt *time.Time
	Passed      bool
	Active      bool
}

type PipelineLink struct {
	FromStageID string `json:"from_stage_id"`
	ToStageID   string `json:"to_stage_id"`
}

type AmaraAward struct {
	ID               string    `json:"id"`
	CompanyID        string    `json:"company_id"`
	CompanyName      string    `json:"company_name"`
	Name             string    `json:"name"`
	Crest            string    `json:"crest"`
	SkillDomain      string    `json:"skill_domain"`
	Level            *string   `json:"level"`
	LevelDescription *string   `json:"level_description"`
	Strengths        []string  `json:"strengths"`
	DevelopmentAreas []string  `json:"development_areas"`
	Evidence         []string  `json:"evidence"`
	AwardedAt        time.Time `json:"awarded_at"`
}
