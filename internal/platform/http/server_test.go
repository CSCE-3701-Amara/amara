package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CSCE-3701-Amara/amara/internal/platform/config"
	"github.com/CSCE-3701-Amara/amara/internal/platform/database"
	"github.com/CSCE-3701-Amara/amara/internal/platform/logging"
)

func testConfig() config.HTTPConfig {
	return config.HTTPConfig{
		Address: "127.0.0.1:0", ReadHeaderTimeout: time.Second,
		ReadTimeout: time.Second, WriteTimeout: time.Second,
		IdleTimeout: time.Second, ShutdownTimeout: time.Second,
		AllowedOrigins: []string{"http://localhost:5173"},
	}
}

func testServer(t *testing.T, cfg config.HTTPConfig, output io.Writer, registrars ...RouteRegistrar) *Server {
	t.Helper()
	return testServerWithDB(t, cfg, output, nil, registrars...)
}

func testServerWithDB(t *testing.T, cfg config.HTTPConfig, output io.Writer, dbSvc database.Service, registrars ...RouteRegistrar) *Server {
	t.Helper()
	gin.SetMode(gin.TestMode)
	server, err := New(cfg, logging.New(output, slog.LevelInfo, "json"), dbSvc, registrars...)
	if err != nil {
		t.Fatal(err)
	}
	return server
}

func TestRoutesAndRequestIDs(t *testing.T) {
	server := testServer(t, testConfig(), io.Discard, func(v1 *gin.RouterGroup) {
		v1.GET("/example", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	})
	seenIDs := make(map[string]bool)
	for _, tc := range []struct {
		method, path, code string
		status             int
	}{
		{http.MethodGet, "/health", "", http.StatusOK},
		{http.MethodGet, "/health/live", "", http.StatusOK},
		{http.MethodGet, "/health/ready", "", http.StatusOK},
		{http.MethodGet, "/v1/example", "", http.StatusNoContent},
		{http.MethodGet, "/missing", "not_found", http.StatusNotFound},
		{http.MethodPost, "/health", "method_not_allowed", http.StatusMethodNotAllowed},
	} {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, tc.path, nil)
			request.Header.Set("X-Request-ID", "caller-controlled")
			response := httptest.NewRecorder()
			server.server.Handler.ServeHTTP(response, request)
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d: %s", response.Code, tc.status, response.Body)
			}
			id := response.Header().Get("X-Request-ID")
			if id == "" || id == "caller-controlled" || seenIDs[id] {
				t.Fatalf("request must have a unique server-generated ID, got %q", id)
			}
			seenIDs[id] = true
			if (tc.path == "/health" || tc.path == "/health/live") && tc.method == http.MethodGet && response.Body.String() != `{"status":"ok"}` {
				t.Fatalf("unexpected liveness response: %s", response.Body)
			}
			if tc.code != "" {
				var body struct {
					Error     struct{ Code string } `json:"error"`
					RequestID string                `json:"request_id"`
				}
				if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}
				if body.Error.Code != tc.code || body.RequestID != id {
					t.Fatalf("unexpected error response: %+v", body)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Health endpoint behaviour
// ---------------------------------------------------------------------------

// stubDBService is an in-process test double for database.Service.
type stubDBService struct {
	status database.HealthStatus
}

func (s *stubDBService) Health(_ context.Context) database.HealthStatus { return s.status }
func (s *stubDBService) Pool() *pgxpool.Pool                             { return nil }
func (s *stubDBService) Close()                                          {}

// stubDB returns a database.Service fake with the given pre-canned status.
func stubDB(status string, errMsg string) database.Service {
	return &stubDBService{status: database.HealthStatus{
		Status: status,
		Error:  errMsg,
		Stats:  database.PoolStats{MaxConns: 4},
	}}
}

func TestLivenessIsIndependentOfDatabase(t *testing.T) {
	// Even with an "unhealthy" database stub the liveness endpoint must return 200.
	server := testServerWithDB(t, testConfig(), io.Discard, stubDB("down", "ping failed"))
	for _, path := range []string{"/health", "/health/live"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		server.server.Handler.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("%s: status = %d, want 200 (liveness must never depend on DB)", path, response.Code)
		}
		if response.Body.String() != `{"status":"ok"}` {
			t.Fatalf("%s: unexpected body: %s", path, response.Body)
		}
	}
}

func TestReadinessReturns200WhenDatabaseIsUp(t *testing.T) {
	server := testServerWithDB(t, testConfig(), io.Discard, stubDB("up", ""))
	request := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	response := httptest.NewRecorder()
	server.server.Handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", response.Code, response.Body)
	}
	var body struct {
		Status   string `json:"status"`
		Database struct {
			Status string `json:"status"`
		} `json:"database"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Status != "up" || body.Database.Status != "up" {
		t.Fatalf("unexpected body: %+v", body)
	}
}

func TestReadinessReturns503WhenDatabaseIsDown(t *testing.T) {
	server := testServerWithDB(t, testConfig(), io.Discard, stubDB("down", "database ping timed out"))
	request := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	response := httptest.NewRecorder()
	server.server.Handler.ServeHTTP(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503: %s", response.Code, response.Body)
	}
	var body struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Status != "down" {
		t.Fatalf("body status = %q, want %q", body.Status, "down")
	}
	if body.Error == "" {
		t.Fatal("error field must be present when database is down")
	}
}

func TestReadinessWithNoDBServiceReturnsOK(t *testing.T) {
	// nil dbSvc — no database wired. Should be treated as live-only.
	server := testServerWithDB(t, testConfig(), io.Discard, nil)
	request := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	response := httptest.NewRecorder()
	server.server.Handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 when no DB service configured", response.Code)
	}
}

func TestReadinessDoesNotExposeInternalDetails(t *testing.T) {
	server := testServerWithDB(t, testConfig(), io.Discard, stubDB("down", "database ping timed out"))
	request := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	response := httptest.NewRecorder()
	server.server.Handler.ServeHTTP(response, request)
	body := response.Body.String()
	for _, sensitive := range []string{"127.0.0.1", "password", "secret", "postgres://"} {
		if strings.Contains(body, sensitive) {
			t.Fatalf("response exposed sensitive string %q: %s", sensitive, body)
		}
	}
}

// ---------------------------------------------------------------------------
// CORS
// ---------------------------------------------------------------------------

func TestCORS(t *testing.T) {
	server := testServer(t, testConfig(), io.Discard)
	for _, tc := range []struct {
		name, method, origin string
		status               int
		allowed              bool
	}{
		{"no origin", http.MethodGet, "", http.StatusOK, false},
		{"allowed", http.MethodGet, "http://localhost:5173", http.StatusOK, true},
		{"preflight", http.MethodOptions, "http://localhost:5173", http.StatusNoContent, true},
		{"disallowed", http.MethodGet, "https://untrusted.example", http.StatusForbidden, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, "/health", nil)
			request.Header.Set("Origin", tc.origin)
			if tc.method == http.MethodOptions {
				request.Header.Set("Access-Control-Request-Method", "GET")
				request.Header.Set("Access-Control-Request-Headers", "Authorization")
			}
			response := httptest.NewRecorder()
			server.server.Handler.ServeHTTP(response, request)
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d", response.Code, tc.status)
			}
			origin := response.Header().Get("Access-Control-Allow-Origin")
			if (tc.allowed && origin != tc.origin) || (!tc.allowed && origin != "") {
				t.Fatalf("unexpected allowed origin: %q", origin)
			}
			if tc.method == http.MethodOptions && !strings.Contains(response.Header().Get("Access-Control-Allow-Headers"), "Authorization") {
				t.Fatal("preflight must allow the Authorization header")
			}
		})
	}
	cfg := testConfig()
	cfg.AllowedOrigins = nil
	server = testServer(t, cfg, io.Discard)
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.Header.Set("Origin", "http://localhost:5173")
	response := httptest.NewRecorder()
	server.server.Handler.ServeHTTP(response, request)
	if response.Code != http.StatusForbidden {
		t.Fatal("empty origin list must deny cross-origin requests")
	}
}

// ---------------------------------------------------------------------------
// Recovery and logging
// ---------------------------------------------------------------------------

func TestRecoveryAndLogsDoNotExposeRequestData(t *testing.T) {
	var output bytes.Buffer
	server := testServer(t, testConfig(), &output, func(v1 *gin.RouterGroup) {
		v1.GET("/items/:id", func(c *gin.Context) {
			if RequestID(c.Request.Context()) == "" {
				t.Error("request context must carry a correlation ID")
			}
			c.Status(http.StatusNoContent)
		})
		v1.GET("/panic/:id", func(c *gin.Context) { panic("sensitive-panic-value") })
	})
	for _, tc := range []struct {
		path   string
		status int
	}{
		{"/v1/items/private-identifier?token=sensitive-query", http.StatusNoContent},
		{"/v1/panic/private-identifier", http.StatusInternalServerError},
		{"/unknown/private-identifier", http.StatusNotFound},
		{"/health", http.StatusOK},
	} {
		request := httptest.NewRequest(http.MethodGet, tc.path, strings.NewReader("sensitive-body"))
		request.Header.Set("Authorization", "Bearer sensitive-header")
		response := httptest.NewRecorder()
		server.server.Handler.ServeHTTP(response, request)
		if response.Code != tc.status {
			t.Fatalf("%s: status = %d, want %d", tc.path, response.Code, tc.status)
		}
		if strings.Contains(response.Body.String(), "sensitive-panic-value") {
			t.Fatal("panic values must not be returned to clients")
		}
	}
	for _, sensitive := range []string{"private-identifier", "sensitive-query", "sensitive-body", "sensitive-header", "sensitive-panic-value"} {
		if strings.Contains(output.String(), sensitive) {
			t.Fatalf("logs exposed %q", sensitive)
		}
	}
	if !strings.Contains(output.String(), `"route":"/v1/items/:id"`) ||
		!strings.Contains(output.String(), `"route":"unmatched"`) ||
		!strings.Contains(output.String(), `"status":500`) {
		t.Fatalf("logs must include route templates and final statuses: %s", &output)
	}
	for line := range strings.SplitSeq(strings.TrimSpace(output.String()), "\n") {
		if !json.Valid([]byte(line)) {
			t.Fatalf("log line is not JSON: %s", line)
		}
	}
}

// ---------------------------------------------------------------------------
// Trusted proxies
// ---------------------------------------------------------------------------

func TestTrustedProxies(t *testing.T) {
	for _, trusted := range []bool{false, true} {
		cfg := testConfig()
		expected := "203.0.113.10"
		if trusted {
			cfg.TrustedProxies = []string{"203.0.113.10"}
			expected = "198.51.100.9"
		}
		server := testServer(t, cfg, io.Discard, func(v1 *gin.RouterGroup) {
			v1.GET("/client-ip", func(c *gin.Context) { c.String(http.StatusOK, c.ClientIP()) })
		})
		request := httptest.NewRequest(http.MethodGet, "/v1/client-ip", nil)
		request.RemoteAddr = "203.0.113.10:12345"
		request.Header.Set("X-Forwarded-For", "198.51.100.9")
		response := httptest.NewRecorder()
		server.server.Handler.ServeHTTP(response, request)
		if response.Body.String() != expected {
			t.Fatalf("trusted=%t: client IP = %q, want %q", trusted, response.Body.String(), expected)
		}
	}
}

// ---------------------------------------------------------------------------
// Shutdown
// ---------------------------------------------------------------------------

func TestServerShutdown(t *testing.T) {
	for _, force := range []bool{false, true} {
		name := "drain"
		if force {
			name = "deadline"
		}
		t.Run(name, func(t *testing.T) {
			started, release, stopped := make(chan struct{}), make(chan struct{}), make(chan struct{})
			var releaseOnce sync.Once
			defer releaseOnce.Do(func() { close(release) })
			cfg := testConfig()
			if force {
				cfg.ShutdownTimeout = 20 * time.Millisecond
			}
			server := testServer(t, cfg, io.Discard, func(v1 *gin.RouterGroup) {
				v1.GET("/slow", func(c *gin.Context) {
					defer close(stopped)
					close(started)
					select {
					case <-release:
						c.String(http.StatusOK, "completed")
					case <-c.Request.Context().Done():
					}
				})
			})
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer listener.Close()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			result := make(chan error, 1)
			go func() { result <- server.serve(ctx, listener) }()

			client := &http.Client{Timeout: 3 * time.Second, Transport: &http.Transport{}}
			defer client.CloseIdleConnections()
			clientResult := make(chan error, 1)
			go func() {
				response, err := client.Get("http://" + listener.Addr().String() + "/v1/slow")
				if err == nil {
					defer response.Body.Close()
					body, readErr := io.ReadAll(response.Body)
					if readErr != nil {
						err = readErr
					} else if response.StatusCode != http.StatusOK || string(body) != "completed" {
						err = errors.New("active request did not complete")
					}
				}
				clientResult <- err
			}()
			select {
			case <-started:
			case <-time.After(3 * time.Second):
				t.Fatal("request did not start")
			}
			shutdownStarted := make(chan struct{})
			server.server.RegisterOnShutdown(func() { close(shutdownStarted) })
			cancel()
			select {
			case <-shutdownStarted:
			case <-time.After(3 * time.Second):
				t.Fatal("shutdown did not start")
			}
			if !force {
				select {
				case err := <-result:
					t.Fatalf("server exited before draining its active request: %v", err)
				default:
				}
				releaseOnce.Do(func() { close(release) })
			}
			select {
			case err := <-result:
				if (force && !errors.Is(err, context.DeadlineExceeded)) || (!force && err != nil) {
					t.Fatalf("unexpected shutdown result: %v", err)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("server did not stop")
			}
			select {
			case err := <-clientResult:
				if (force && err == nil) || (!force && err != nil) {
					t.Fatalf("unexpected client result: %v", err)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("client did not finish")
			}
			select {
			case <-stopped:
			case <-time.After(3 * time.Second):
				t.Fatal("handler did not stop")
			}
		})
	}
}

func TestRunReportsListenErrors(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	cfg := testConfig()
	cfg.Address = listener.Addr().String()
	server := testServer(t, cfg, io.Discard)
	if err := server.Run(context.Background()); err == nil {
		t.Fatal("an occupied port must return a startup error")
	}
}

