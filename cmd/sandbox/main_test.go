package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubmissionLimits(t *testing.T) {
	encode := base64.StdEncoding.EncodeToString
	for _, tc := range []struct {
		name   string
		input  submission
		origin string
		busy   bool
		want   int
	}{
		{name: "browser", origin: "http://example.test", want: http.StatusForbidden},
		{name: "network", input: submission{Network: true}, want: http.StatusBadRequest},
		{name: "bad encoding", input: submission{Source: "!"}, want: http.StatusBadRequest},
		{name: "large source", input: submission{Source: encode(make([]byte, 65537))}, want: http.StatusBadRequest},
		{name: "large stdin", input: submission{Input: encode(make([]byte, 16385))}, want: http.StatusBadRequest},
		{name: "unknown language", input: submission{Language: 999}, want: http.StatusUnprocessableEntity},
		{name: "capacity", input: submission{Language: 71, Source: encode([]byte("print(1)"))}, busy: true, want: http.StatusServiceUnavailable},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := &service{jobs: map[string]result{}, slots: make(chan struct{}, 1)}
			if tc.busy {
				s.slots <- struct{}{}
			}
			body, _ := json.Marshal(tc.input)
			req := httptest.NewRequest(http.MethodPost, "/submissions", bytes.NewReader(body))
			req.Header.Set("Origin", tc.origin)
			w := httptest.NewRecorder()
			s.submit(w, req)
			if w.Code != tc.want {
				t.Fatalf("status %d, want %d: %s", w.Code, tc.want, w.Body.String())
			}
			if len(s.jobs) != 0 {
				t.Fatal("rejected request queued a job")
			}
		})
	}
}
func TestOutputBound(t *testing.T) {
	var out output
	if n, err := out.Write(make([]byte, maxOutputBytes)); err != nil || n != maxOutputBytes {
		t.Fatalf("boundary: %d %v", n, err)
	}
	if _, err := out.Write([]byte("x")); err == nil {
		t.Fatal("unbounded output")
	}
	if len(out.data) != maxOutputBytes {
		t.Fatal("retained excess output")
	}
}

func TestDockerArgsUseLanguageConfiguration(t *testing.T) {
	config := languageConfigs[languageTypeScript]
	args := dockerArgs("token", config, "/tmp/source")

	for _, want := range []string{
		"--network", "none",
		"--read-only",
		"--cap-drop", "ALL",
		"--volume", "/tmp/source:/code.ts:ro,Z",
		"node:24-alpine",
		"timeout", "-k", "1s", "8s",
		"--experimental-strip-types",
	} {
		if !contains(args, want) {
			t.Errorf("docker args missing %q: %v", want, args)
		}
	}
}

func TestLanguages(t *testing.T) {
	recorder := httptest.NewRecorder()
	languages(recorder, httptest.NewRequest(http.MethodGet, "/languages", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status %d, want %d", recorder.Code, http.StatusOK)
	}
	if got, want := recorder.Body.String(), `[{"id":71,"name":"Python"},{"id":63,"name":"JavaScript"},{"id":74,"name":"TypeScript (native Node type stripping)"}]`+"\n"; got != want {
		t.Fatalf("body %q, want %q", got, want)
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
