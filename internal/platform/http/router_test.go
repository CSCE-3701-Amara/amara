package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouterRequiresRegistrationHandler(t *testing.T) {
	if _, err := NewRouter(nil); err != ErrMissingRegistrationHandler {
		t.Fatalf("expected missing handler error, got %v", err)
	}
}

func TestNewRouterRegistersHealthAndRegistrationRoutes(t *testing.T) {
	handler := http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusCreated)
	})
	router, err := NewRouter(handler, handler)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		path string
		want int
	}{
		{path: "/health", want: http.StatusOK},
		{path: RegistrationPath, want: http.StatusCreated},
		{path: CompanyRegistrationPath, want: http.StatusCreated},
	} {
		request := httptest.NewRequest(http.MethodPost, test.path, nil)
		if test.path == "/health" {
			request = httptest.NewRequest(http.MethodGet, test.path, nil)
		}
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		if recorder.Code != test.want {
			t.Fatalf("%s: expected %d, got %d", test.path, test.want, recorder.Code)
		}
	}
}
