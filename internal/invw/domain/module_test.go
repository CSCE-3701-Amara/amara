package domain

import "testing"

func TestBehavioralPublicationRequiresNeitherScoreNorThreshold(t *testing.T) {
	module := Definition{
		Name: "Behavioral questions",
		Type: Behavioral,
		Questions: []Question{
			{Prompt: "How do you share feedback?", Format: "long_text"},
			{Prompt: "Which example applies?", Format: "multiple_choice", Options: []string{"First", "Second"}},
		},
	}
	if err := module.ValidatePublication(nil); err != nil {
		t.Fatal(err)
	}
	threshold := 70
	module.Threshold = &threshold
	if err := module.ValidatePublication(nil); err == nil {
		t.Fatal("behavioral modules must reject a score threshold")
	}
}

func TestRecordedVideoPublicationRequiresMappedRubric(t *testing.T) {
	threshold := 70
	module := Definition{
		Name:      "Recorded interview",
		Type:      RecordedVideo,
		Threshold: &threshold,
		Questions: []Question{{Prompt: "Describe a project.", Format: "video"}},
		Rubric:    []Criterion{{Description: "Explain trade-offs", Amara: "Communication"}},
	}
	if err := module.ValidatePublication(nil); err == nil {
		t.Fatal("rubric criteria must reference a configured Amara")
	}
	module.Amarat = []Amara{{
		Name: "Communication", Crest: "speech", Domain: "Communication", Description: "Explaining trade-offs",
	}}
	if err := module.ValidatePublication(nil); err != nil {
		t.Fatal(err)
	}
	module.ActiveMinutes = &threshold
	if err := module.ValidatePublication(nil); err == nil {
		t.Fatal("recorded video must reject an active timer")
	}
}

func TestLibrarySelectionRejectsInsufficientDifficultyPool(t *testing.T) {
	ref := ProblemRef{ModuleID: "module", Position: 0}
	problems := []LibraryProblem{{ProblemRef: ref, Difficulty: "easy"}}
	selection := LibrarySelection{DifficultyMix: map[string]int{"easy": 2}}
	if err := selection.validate(problems); err == nil {
		t.Fatal("random selection cannot ask for more distinct problems than the library contains")
	}
	selection = LibrarySelection{Problems: []ProblemRef{ref, ref}}
	if err := selection.validate(problems); err == nil {
		t.Fatal("specific selection must reject duplicate problems")
	}
	selection = LibrarySelection{Problems: []ProblemRef{ref}}
	if err := selection.validate(problems); err != nil {
		t.Fatal(err)
	}
}
