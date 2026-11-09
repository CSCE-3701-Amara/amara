package database

import (
	"encoding/json"
	"testing"

	"github.com/CSCE-3701-Amara/amara/internal/invw/domain"
)

func TestCodeModuleAdvancementRequiresFinalizedCaseEvidence(t *testing.T) {
	for _, kind := range []string{domain.Coding, domain.ProgrammingLanguage} {
		for _, tc := range []struct {
			name      string
			evaluated bool
			passed    int
			want      string
		}{
			{"ungraded", false, 0, "pending_grading"},
			{"below minimum", true, 0, "held"},
			{"exactly minimum", true, 1, "passed"},
		} {
			t.Run(kind+"/"+tc.name, func(t *testing.T) {
				s := advancementStage{rule: "completion", below: "hold", submitted: true, evaluated: tc.evaluated, outcome: "scored", questionIDs: map[string]string{"1:1": "question"}, evidence: []QuestionEvidence{{QuestionID: "question", PassedCases: &tc.passed, TotalCases: new(1)}}}
				if err := json.Unmarshal([]byte(`{"modules":[{"id":"module","configuration":{"definition":{"type":"`+kind+`","minimum_tests_passed":1}}}]}`), &s.configuration); err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal([]byte(`[{"module_id":"module","questions":[{"format":"code","minimum_cases_passed":1,"test_cases":[{}]}]}]`), &s.content); err != nil {
					t.Fatal(err)
				}
				if got := s.automaticDecision(); got.status != tc.want {
					t.Fatalf("decision=%+v, want %s", got, tc.want)
				}
			})
		}
	}
}

func TestLegacyCodeThreshold(t *testing.T) {
	for _, tc := range []struct {
		name          string
		passed, total int
		score         *int
		weights       []int
		rule          string
		stageScore    *int
		want          string
	}{
		{"equal weight historical", 1, 1, nil, []int{1}, "completion", nil, "passed"},
		{"failed historical", 0, 1, nil, []int{1}, "completion", nil, "held"},
		{"weighted exact minimum", 1, 2, new(70), []int{9, 1}, "completion", nil, "passed"},
		{"weighted below minimum", 1, 2, new(69), []int{9, 1}, "completion", nil, "held"},
		{"unequal historical needs regrading", 1, 2, nil, []int{9, 1}, "completion", nil, "pending_grading"},
		{"stage threshold enforced", 1, 1, new(80), []int{1}, "threshold", new(80), "held"},
		{"stage threshold met", 1, 1, new(90), []int{1}, "threshold", new(90), "passed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := advancementStage{rule: tc.rule, threshold: new(90), score: tc.stageScore, below: "hold", submitted: true, evaluated: true, outcome: "scored", questionIDs: map[string]string{"1:1": "q"}, evidence: []QuestionEvidence{{QuestionID: "q", PassedCases: &tc.passed, TotalCases: &tc.total, Score: tc.score}}}
			if err := json.Unmarshal([]byte(`{"modules":[{"id":"m","configuration":{"definition":{"type":"coding","threshold":70}}}]}`), &s.configuration); err != nil {
				t.Fatal(err)
			}
			var cases []domain.TestCase
			for _, weight := range tc.weights {
				cases = append(cases, domain.TestCase{Weight: weight})
			}
			raw, _ := json.Marshal([]map[string]any{{"module_id": "m", "questions": []domain.Question{{Format: "code", TestCases: cases}}}})
			if err := json.Unmarshal(raw, &s.content); err != nil {
				t.Fatal(err)
			}
			if got := s.automaticDecision(); got.status != tc.want {
				t.Fatalf("decision %+v, want %s", got, tc.want)
			}
		})
	}
}
