package http

import (
	"net/http"
)

const RegistrationPath = "/api/v1/auth/registrations"

func NewRouter(registrationHandler http.Handler) (*http.ServeMux, error) {
	if registrationHandler == nil {
		return nil, ErrMissingRegistrationHandler
	}
	mux := http.NewServeMux()
	mux.Handle(RegistrationPath, registrationHandler)
	mux.HandleFunc("/health", func(writer http.ResponseWriter, _ *http.Request) {
		if writer.Header().Get("Content-Type") == "" {
			writer.Header().Set("Content-Type", "application/json")
		}
		_, _ = writer.Write([]byte(`{"status":"ok"}`))
	})
	return mux, nil
}
