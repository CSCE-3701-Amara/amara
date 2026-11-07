package judge0

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRun(t *testing.T) {
	for _, tc := range []struct {
		name      string
		status    int
		succeeded bool
		wantError bool
	}{
		{"accepted", 3, true, false},
		{"compile error", 6, false, false},
		{"provider error", 13, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("X-Auth-Token") != "test-token" || r.URL.Query().Get("base64_encoded") != "true" {
					t.Error("missing authentication or encoding option")
				}
				switch {
				case r.Method == http.MethodPost && r.URL.Path == "/submissions":
					var body map[string]any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Error(err)
					}
					if body["source_code"] != "cHJpbnQoKQ==" || body["stdin"] != "aW5wdXQ=" ||
						body["language_id"] != float64(71) || body["enable_network"] != false ||
						body["expected_output"] != nil || r.URL.Query().Get("wait") != "false" {
						t.Errorf("unexpected submission: %v", body)
					}
					fmt.Fprint(w, "{\"token\":\"run-1\"}")
				case r.Method == http.MethodGet && r.URL.Path == "/submissions/run-1":
					fmt.Fprintf(w, "{\"status\":{\"id\":%d},\"stdout\":\"b3V0cHV0Cg==\"}", tc.status)
				default:
					http.NotFound(w, r)
				}
			}))
			defer server.Close()
			client, err := New(server.URL, "test-token", nil)
			if err != nil {
				t.Fatal(err)
			}
			output, succeeded, err := client.Run(t.Context(), "python", "print()", "input")
			if succeeded != tc.succeeded || (err != nil) != tc.wantError || (succeeded && output != "output\n") {
				t.Fatalf("Run = %q, %v, %v", output, succeeded, err)
			}
		})
	}
}
