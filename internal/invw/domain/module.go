package domain

import (
	"errors"
	"fmt"
	"strings"
)

const (
	Coding              = "coding"
	ProgrammingLanguage = "programming_language"
	Behavioral          = "behavioral"
	RecordedVideo       = "recorded_video"
)

var ErrModuleNotFound = errors.New("test module not found")
var ErrModuleLocked = errors.New("published and copied test modules are locked")

type InvalidModuleError struct {
	Reason string
}

func (e *InvalidModuleError) Error() string { return e.Reason }

type Module struct {
	ID         string  `json:"id"`
	CompanyID  string  `json:"company_id"`
	CopiedFrom *string `json:"copied_from_module_id,omitempty"`
	Published  bool    `json:"published"`
	Definition
}

type Definition struct {
	Name               string            `json:"name"`
	Description        string            `json:"description,omitempty"`
	Type               string            `json:"type"`
	ActiveMinutes      *int              `json:"active_minutes,omitempty"`
	Threshold          *int              `json:"threshold,omitempty"`
	MinimumTestsPassed *int              `json:"minimum_tests_passed,omitempty"`
	Questions          []Question        `json:"questions,omitempty"`
	Amarat             []Amara           `json:"amarat,omitempty"`
	Rubric             []Criterion       `json:"rubric,omitempty"`
	Library            *LibrarySelection `json:"library,omitempty"`
}

type Amara struct {
	Name        string `json:"name"`
	Crest       string `json:"crest"`
	Domain      string `json:"domain"`
	Description string `json:"description"`
}

type Question struct {
	Prompt             string     `json:"prompt"`
	Format             string     `json:"format"`
	Options            []string   `json:"options,omitempty"`
	Amara              string     `json:"amara,omitempty"`
	Difficulty         string     `json:"difficulty,omitempty"`
	TestCases          []TestCase `json:"test_cases,omitempty"`
	MinimumCasesPassed *int       `json:"minimum_cases_passed,omitempty"`
}

type TestCase struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	Weight         int    `json:"weight"`
}

type Criterion struct {
	Description string `json:"description"`
	Amara       string `json:"amara"`
}

type LibrarySelection struct {
	Pool          []ProblemRef   `json:"pool,omitempty"`
	Problems      []ProblemRef   `json:"problems,omitempty"`
	DifficultyMix map[string]int `json:"difficulty_mix,omitempty"`
	Amara         string         `json:"amara"`
}

type ProblemRef struct {
	ModuleID string `json:"module_id"`
	Position int    `json:"position"`
}

type LibraryProblem struct {
	ProblemRef
	Difficulty string `json:"difficulty"`
	Prompt     string `json:"prompt"`
}

func (d Definition) ValidateDraft() error {
	if strings.TrimSpace(d.Name) == "" {
		return errors.New("a module name is required")
	}
	switch d.Type {
	case Coding, ProgrammingLanguage, Behavioral, RecordedVideo:
	default:
		return errors.New("invalid module type")
	}
	if d.ActiveMinutes != nil && (*d.ActiveMinutes <= 0 || *d.ActiveMinutes > 2147483647 || d.Type == RecordedVideo) {
		return errors.New("active time must be positive and is available only for non-video modules")
	}
	if d.Threshold != nil && (*d.Threshold < 0 || *d.Threshold > 100 || d.Type == Behavioral) {
		return errors.New("threshold must be from 0 to 100 and is unavailable for behavioral modules")
	}
	if d.MinimumTestsPassed != nil && ((d.Type != Coding && d.Type != ProgrammingLanguage) || *d.MinimumTestsPassed < 1) {
		return errors.New("minimum tests passed must be positive and is available only for code modules")
	}
	return nil
}

func (d Definition) ValidatePublication(problems []LibraryProblem) error {
	if err := d.ValidateDraft(); err != nil {
		return err
	}
	if d.Type != Behavioral && d.Threshold == nil && d.MinimumTestsPassed == nil {
		return errors.New("scored modules require a threshold")
	}
	if len(d.Questions) == 0 && d.Library == nil {
		return errors.New("at least one question or library selection is required")
	}
	if d.Type == RecordedVideo && len(d.Rubric) == 0 {
		return errors.New("recorded video requires a rubric")
	}
	if d.Type != RecordedVideo && len(d.Rubric) != 0 {
		return errors.New("only recorded video supports a rubric")
	}
	mappings, err := newAmaraMappings(d.Amarat)
	if err != nil {
		return err
	}
	code := d.Type == Coding || d.Type == ProgrammingLanguage
	for i, question := range d.Questions {
		if err := question.validate(d.Type); err != nil {
			return fmt.Errorf("question %d: %w", i+1, err)
		}
		if d.MinimumTestsPassed != nil && question.MinimumCasesPassed == nil {
			return fmt.Errorf("question %d: minimum cases passed is required", i+1)
		}
		if code {
			if err := mappings.mark(question.Amara); err != nil {
				return fmt.Errorf("question %d: %w", i+1, err)
			}
		}
	}
	if err := mappings.markRubric(d.Rubric); err != nil {
		return err
	}
	if d.Library != nil {
		if !code {
			return errors.New("only code modules can select library problems")
		}
		if err := d.Library.validate(problems); err != nil {
			return err
		}
		if err := mappings.mark(d.Library.Amara); err != nil {
			return err
		}
	}
	if d.MinimumTestsPassed != nil {
		testCount := len(d.Questions)
		if d.Library != nil {
			testCount += len(d.Library.Problems)
			for _, count := range d.Library.DifficultyMix {
				testCount += count
			}
		}
		if *d.MinimumTestsPassed > testCount {
			return errors.New("minimum tests passed exceeds the module test count")
		}
	}
	return mappings.validateUsed(d.Amarat)
}

type amaraMappings map[string]bool

func newAmaraMappings(amarat []Amara) (amaraMappings, error) {
	mapped := make(amaraMappings, len(amarat))
	for _, a := range amarat {
		if strings.TrimSpace(a.Name) == "" || strings.TrimSpace(a.Crest) == "" ||
			strings.TrimSpace(a.Domain) == "" || strings.TrimSpace(a.Description) == "" {
			return nil, errors.New("each Amara requires a name, crest, domain, and description")
		}
		if _, exists := mapped[a.Name]; exists {
			return nil, errors.New("Amara names must be unique")
		}
		mapped[a.Name] = false
	}
	return mapped, nil
}

func (m amaraMappings) mark(name string) error {
	if _, exists := m[name]; !exists {
		return fmt.Errorf("unknown Amara %q", name)
	}
	m[name] = true
	return nil
}

func (m amaraMappings) markRubric(rubric []Criterion) error {
	for i, criterion := range rubric {
		if strings.TrimSpace(criterion.Description) == "" {
			return fmt.Errorf("rubric criterion %d requires a description", i+1)
		}
		if err := m.mark(criterion.Amara); err != nil {
			return fmt.Errorf("rubric criterion %d: %w", i+1, err)
		}
	}
	return nil
}

func (m amaraMappings) validateUsed(amarat []Amara) error {
	for _, a := range amarat {
		if !m[a.Name] {
			return fmt.Errorf("Amara %q requires a scored item", a.Name)
		}
	}
	return nil
}

func (q Question) validate(kind string) error {
	if strings.TrimSpace(q.Prompt) == "" {
		return errors.New("a prompt is required")
	}
	if q.Format != "multiple_choice" && len(q.Options) > 0 {
		return errors.New("options are only supported for multiple-choice questions")
	}
	if q.Format != "code" && len(q.TestCases) > 0 {
		return errors.New("test cases are only supported for code questions")
	}
	if q.MinimumCasesPassed != nil && (q.Format != "code" || *q.MinimumCasesPassed < 1 || *q.MinimumCasesPassed > len(q.TestCases)) {
		return errors.New("minimum cases passed must be within the code question's test case count")
	}
	switch kind {
	case Coding, ProgrammingLanguage:
		if q.Format != "code" || len(q.TestCases) == 0 || strings.TrimSpace(q.Difficulty) == "" {
			return errors.New("code questions require difficulty and weighted test cases")
		}
		for _, test := range q.TestCases {
			if test.Weight <= 0 {
				return errors.New("test case weights must be positive")
			}
		}
	case Behavioral:
		switch q.Format {
		case "multiple_choice":
			if len(q.Options) < 2 {
				return errors.New("multiple choice requires at least two options")
			}
			for _, option := range q.Options {
				if strings.TrimSpace(option) == "" {
					return errors.New("multiple-choice options must have text")
				}
			}
		case "short_text", "long_text":
		default:
			return errors.New("invalid behavioral response format")
		}
	case RecordedVideo:
		if q.Format != "video" {
			return errors.New("recorded video requires video prompts")
		}
	}
	return nil
}

func (s LibrarySelection) validate(problems []LibraryProblem) error {
	if (len(s.Problems) == 0) == (len(s.DifficultyMix) == 0) {
		return errors.New("select specific problems or a random difficulty mix")
	}
	available := make(map[ProblemRef]bool, len(problems))
	counts := make(map[string]int)
	for _, problem := range problems {
		available[problem.ProblemRef] = true
		counts[problem.Difficulty]++
	}
	for _, ref := range s.Problems {
		if !available[ref] {
			return errors.New("a selected problem is unavailable or was selected twice")
		}
		delete(available, ref)
	}
	for difficulty, count := range s.DifficultyMix {
		if count <= 0 || count > counts[difficulty] {
			return fmt.Errorf("invalid library problem count for difficulty %q", difficulty)
		}
	}
	return nil
}
