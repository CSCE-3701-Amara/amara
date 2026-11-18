package domain

import (
	"testing"
	"time"
)

func TestValidateOrderedStages(t *testing.T) {
	opens := time.Date(2026, time.November, 10, 8, 0, 0, 0, time.UTC)
	threshold := 70
	assessment := Stage{
		Name: "Assessment", Type: AssessmentStage, OpensAt: opens, Deadline: opens.Add(time.Hour),
		ModuleIDs: []string{"coding", "behavioral"},
	}
	interview := Stage{
		Name: "Interview", Type: LiveInterviewStage, OpensAt: opens, Deadline: opens.Add(time.Hour),
		Interview: &LiveInterview{
			DurationMinutes: 30, InterviewerIDs: []string{"interviewer"}, Threshold: &threshold,
			Rubric: []Criterion{{Description: "Explains choices", Amara: "Communication"}},
			Amarat: []Amara{{Name: "Communication", Crest: "crest", Domain: "communication", Description: "Explanation"}},
		},
	}
	secondInterview := interview
	secondInterview.Name = "Final interview"
	if err := ValidateStages([]Stage{assessment, interview, secondInterview}); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name   string
		stages []Stage
	}{
		{"duplicate stage names", []Stage{assessment, assessment}},
		{"duplicate modules", []Stage{{Name: "Assessment", Type: AssessmentStage, OpensAt: opens, Deadline: opens.Add(time.Hour), ModuleIDs: []string{"coding", "coding"}}}},
		{"no assessment modules", []Stage{{Name: "Assessment", Type: AssessmentStage, OpensAt: opens, Deadline: opens.Add(time.Hour)}}},
		{"deadline at opening", []Stage{{Name: "Assessment", Type: AssessmentStage, OpensAt: opens, Deadline: opens, ModuleIDs: []string{"coding"}}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateStages(tc.stages); err == nil {
				t.Fatal("expected invalid stage configuration")
			}
		})
	}

	interview.ModuleIDs = []string{"coding"}
	if err := ValidateStages([]Stage{interview}); err == nil {
		t.Fatal("a live interview cannot include test modules")
	}
	interview.ModuleIDs = nil
	interview.Interview.Rubric[0].Amara = "Unknown"
	if err := ValidateStages([]Stage{interview}); err == nil {
		t.Fatal("the rubric must map to a configured Amara")
	}
}

// A behavioral interview has no score, so a threshold rule on it must name one
// explicitly; otherwise the stage could never be decided by threshold.
func TestThresholdRuleOnBehavioralInterviewRequiresThreshold(t *testing.T) {
	opens := time.Date(2026, time.November, 10, 8, 0, 0, 0, time.UTC)
	stage := Stage{
		Name: "Behavioral", Type: LiveInterviewStage, OpensAt: opens, Deadline: opens.Add(time.Hour),
		Rule: "threshold",
		Interview: &LiveInterview{
			Behavioral: true, DurationMinutes: 30, InterviewerIDs: []string{"interviewer"},
			Rubric: []Criterion{{Description: "Explains choices", Amara: "Communication"}},
			Amarat: []Amara{{Name: "Communication", Crest: "crest", Domain: "communication", Description: "Explanation"}},
		},
	}
	if err := ValidateStages([]Stage{stage}); err == nil {
		t.Fatal("threshold rule without a threshold was accepted")
	}
	threshold := 60
	stage.ScoreThreshold = &threshold
	if err := ValidateStages([]Stage{stage}); err != nil {
		t.Fatalf("explicit threshold rejected: %v", err)
	}
}
