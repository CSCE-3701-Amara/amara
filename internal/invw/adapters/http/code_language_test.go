package http

import (
	"context"
	"encoding/json"
	"testing"
)

type languageRunner struct{}

func (languageRunner) Run(context.Context, string, string, string) (string, bool, error) {
	return "", true, nil
}
func (languageRunner) SupportedLanguages() []string { return []string{"go", "python"} }

// An unsupported language must be refused at save time; once submitted, the
// grader could never run it and the attempt would stay pending forever.
func TestSaveRejectsUnsupportedCodeLanguage(t *testing.T) {
	h := &AssessmentHandler{runner: languageRunner{}}
	for response, want := range map[string]bool{
		`{"code":"print(1)","language":"python"}`: true,
		`{"code":"x","language":" Go "}`:          true,
		`{"code":"x","language":"cobol"}`:         false,
		`"a short text answer"`:                   true,
		`[0, 2]`:                                  true,
		`{"videoUrl":"https://example.test/v"}`:   true,
	} {
		if got := h.languageSupported(json.RawMessage(response)); got != want {
			t.Errorf("languageSupported(%s) = %v, want %v", response, got, want)
		}
	}
}
