package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CSCE-3701-Amara/amara/internal/platform/apperrors"
	"github.com/gin-gonic/gin"
)

func TestProblemDetailsJSONAndContentType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/test-problem", func(c *gin.Context) {
		WriteProblem(c, ProblemDetails{
			Type:     "https://example.com/probs/custom",
			Title:    "Custom Problem",
			Status:   http.StatusForbidden,
			Detail:   "You are not allowed to access this resource.",
			Instance: "/test-problem",
			Code:     "custom_forbidden",
		})
	})

	router.GET("/test-defaults", func(c *gin.Context) {
		WriteProblem(c, ProblemDetails{
			Status: http.StatusNotFound,
		})
	})

	router.GET("/test-omit-empty", func(c *gin.Context) {
		WriteProblem(c, ProblemDetails{
			Status: http.StatusBadRequest,
			Title:  "Bad Request",
		})
	})

	t.Run("correct content type and fields", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test-problem", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
		}
		if ct := rec.Header().Get("Content-Type"); ct != ProblemContentType {
			t.Fatalf("Content-Type = %q, want %q", ct, ProblemContentType)
		}

		var prob ProblemDetails
		if err := json.Unmarshal(rec.Body.Bytes(), &prob); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}
		if prob.Type != "https://example.com/probs/custom" {
			t.Errorf("type = %q, want %q", prob.Type, "https://example.com/probs/custom")
		}
		if prob.Title != "Custom Problem" {
			t.Errorf("title = %q, want %q", prob.Title, "Custom Problem")
		}
		if prob.Status != http.StatusForbidden {
			t.Errorf("status = %d, want %d", prob.Status, http.StatusForbidden)
		}
		if prob.Detail != "You are not allowed to access this resource." {
			t.Errorf("detail = %q, want expected message", prob.Detail)
		}
		if prob.Instance != "/test-problem" {
			t.Errorf("instance = %q, want %q", prob.Instance, "/test-problem")
		}
		if prob.Code != "custom_forbidden" {
			t.Errorf("code = %q, want %q", prob.Code, "custom_forbidden")
		}
	})

	t.Run("defaults when fields omitted", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test-defaults", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}

		var prob ProblemDetails
		if err := json.Unmarshal(rec.Body.Bytes(), &prob); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}
		if prob.Type != "about:blank" {
			t.Errorf("type = %q, want %q", prob.Type, "about:blank")
		}
		if prob.Title != "Not Found" {
			t.Errorf("title = %q, want %q", prob.Title, "Not Found")
		}
		if prob.Status != http.StatusNotFound {
			t.Errorf("status = %d, want %d", prob.Status, http.StatusNotFound)
		}
		if prob.Instance != "/test-defaults" {
			t.Errorf("instance = %q, want %q", prob.Instance, "/test-defaults")
		}
	})

	t.Run("omits empty optional fields from JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test-omit-empty", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		body := rec.Body.String()
		if strings.Contains(body, `"code"`) {
			t.Errorf("expected code to be omitted when empty, got %s", body)
		}
		if strings.Contains(body, `"errors"`) {
			t.Errorf("expected errors to be omitted when empty, got %s", body)
		}
		if strings.Contains(body, `"detail"`) {
			t.Errorf("expected detail to be omitted when empty, got %s", body)
		}
	})
}

func TestApplicationErrorToHTTPStatusMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedTitle  string
		expectedCode   string
	}{
		{
			name:           "Invalid request syntax (sentinel)",
			err:            apperrors.ErrInvalidInput,
			expectedStatus: http.StatusBadRequest,
			expectedTitle:  "Bad Request",
			expectedCode:   "bad_request",
		},
		{
			name:           "Malformed JSON syntax error",
			err:            &json.SyntaxError{},
			expectedStatus: http.StatusBadRequest,
			expectedTitle:  "Bad Request",
			expectedCode:   "bad_request",
		},
		{
			name:           "JSON unmarshal type error",
			err:            &json.UnmarshalTypeError{},
			expectedStatus: http.StatusBadRequest,
			expectedTitle:  "Bad Request",
			expectedCode:   "bad_request",
		},
		{
			name:           "Missing or invalid authentication",
			err:            apperrors.ErrUnauthenticated,
			expectedStatus: http.StatusUnauthorized,
			expectedTitle:  "Unauthorized",
			expectedCode:   "unauthorized",
		},
		{
			name:           "Authenticated caller lacks permission",
			err:            apperrors.ErrForbidden,
			expectedStatus: http.StatusForbidden,
			expectedTitle:  "Forbidden",
			expectedCode:   "forbidden",
		},
		{
			name:           "Resource does not exist",
			err:            apperrors.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedTitle:  "Not Found",
			expectedCode:   "not_found",
		},
		{
			name:           "Resource conflict or duplicate operation",
			err:            apperrors.ErrConflict,
			expectedStatus: http.StatusConflict,
			expectedTitle:  "Conflict",
			expectedCode:   "conflict",
		},
		{
			name:           "Semantically invalid request (sentinel)",
			err:            apperrors.ErrValidation,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedTitle:  "Unprocessable Content",
			expectedCode:   "validation_failed",
		},
		{
			name:           "Request rate limit exceeded",
			err:            apperrors.ErrRateLimitExceeded,
			expectedStatus: http.StatusTooManyRequests,
			expectedTitle:  "Too Many Requests",
			expectedCode:   "rate_limit_exceeded",
		},
		{
			name:           "Unexpected internal failure (sentinel)",
			err:            apperrors.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
			expectedTitle:  "Internal Server Error",
			expectedCode:   "internal_error",
		},
		{
			name:           "Service temporarily unavailable",
			err:            apperrors.ErrServiceUnavailable,
			expectedStatus: http.StatusServiceUnavailable,
			expectedTitle:  "Service Unavailable",
			expectedCode:   "service_unavailable",
		},
		{
			name:           "Unknown generic error",
			err:            errors.New("something went wrong internally"),
			expectedStatus: http.StatusInternalServerError,
			expectedTitle:  "Internal Server Error",
			expectedCode:   "internal_error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/error", func(c *gin.Context) {
				HandleError(c, nil, tc.err)
			})

			req := httptest.NewRequest(http.MethodGet, "/error", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Fatalf("HTTP status = %d, want %d", rec.Code, tc.expectedStatus)
			}
			if ct := rec.Header().Get("Content-Type"); ct != ProblemContentType {
				t.Fatalf("Content-Type = %q, want %q", ct, ProblemContentType)
			}

			var prob ProblemDetails
			if err := json.Unmarshal(rec.Body.Bytes(), &prob); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if prob.Status != tc.expectedStatus {
				t.Errorf("JSON status = %d, want %d", prob.Status, tc.expectedStatus)
			}
			if prob.Title != tc.expectedTitle {
				t.Errorf("title = %q, want %q", prob.Title, tc.expectedTitle)
			}
			if prob.Code != tc.expectedCode {
				t.Errorf("code = %q, want %q", prob.Code, tc.expectedCode)
			}
			if prob.Type != "about:blank" {
				t.Errorf("type = %q, want %q", prob.Type, "about:blank")
			}
			if prob.Instance != "/error" {
				t.Errorf("instance = %q, want %q", prob.Instance, "/error")
			}
		})
	}
}

func TestValidationErrorsExtension(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	valErr := apperrors.Validation("One or more fields are invalid.",
		apperrors.FieldError{Field: "title", Message: "Title is required."},
		apperrors.FieldError{Field: "salary", Message: "Salary must be positive."},
	)

	router.POST("/jobs", func(c *gin.Context) {
		HandleError(c, nil, valErr)
	})

	req := httptest.NewRequest(http.MethodPost, "/jobs", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}

	var prob ProblemDetails
	if err := json.Unmarshal(rec.Body.Bytes(), &prob); err != nil {
		t.Fatal(err)
	}

	if prob.Status != 422 {
		t.Errorf("status = %d, want 422", prob.Status)
	}
	if prob.Title != "Unprocessable Content" {
		t.Errorf("title = %q, want 'Unprocessable Content'", prob.Title)
	}
	if prob.Detail != "One or more fields are invalid." {
		t.Errorf("detail = %q, want 'One or more fields are invalid.'", prob.Detail)
	}
	if prob.Code != "validation_failed" {
		t.Errorf("code = %q, want 'validation_failed'", prob.Code)
	}
	if len(prob.Errors) != 2 {
		t.Fatalf("expected 2 field errors, got %d", len(prob.Errors))
	}
	if prob.Errors[0].Field != "title" || prob.Errors[0].Message != "Title is required." {
		t.Errorf("unexpected first error: %+v", prob.Errors[0])
	}
	if prob.Errors[1].Field != "salary" || prob.Errors[1].Message != "Salary must be positive." {
		t.Errorf("unexpected second error: %+v", prob.Errors[1])
	}
}

func TestWrappedErrorsMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/wrapped-not-found", func(c *gin.Context) {
		err := fmt.Errorf("usecase failed: %w", apperrors.ErrNotFound)
		HandleError(c, nil, err)
	})

	router.GET("/deep-wrapped-conflict", func(c *gin.Context) {
		err := fmt.Errorf("api layer: %w",
			fmt.Errorf("service layer: %w",
				fmt.Errorf("database query: %w", apperrors.ErrConflict)))
		HandleError(c, nil, err)
	})

	router.GET("/wrapped-app-error", func(c *gin.Context) {
		appErr := apperrors.NotFound("The requested job does not exist.", "job_not_found")
		err := fmt.Errorf("handler: %w", appErr)
		HandleError(c, nil, err)
	})

	t.Run("single wrapped sentinel", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/wrapped-not-found", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
		var prob ProblemDetails
		if err := json.Unmarshal(rec.Body.Bytes(), &prob); err != nil {
			t.Fatal(err)
		}
		if prob.Status != 404 || prob.Title != "Not Found" || prob.Code != "not_found" {
			t.Fatalf("unexpected problem response: %+v", prob)
		}
	})

	t.Run("deep wrapped sentinel", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/deep-wrapped-conflict", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409", rec.Code)
		}
		var prob ProblemDetails
		if err := json.Unmarshal(rec.Body.Bytes(), &prob); err != nil {
			t.Fatal(err)
		}
		if prob.Status != 409 || prob.Title != "Conflict" || prob.Code != "conflict" {
			t.Fatalf("unexpected problem response: %+v", prob)
		}
	})

	t.Run("wrapped AppError preserves safe message and code", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/wrapped-app-error", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
		var prob ProblemDetails
		if err := json.Unmarshal(rec.Body.Bytes(), &prob); err != nil {
			t.Fatal(err)
		}
		if prob.Detail != "The requested job does not exist." {
			t.Errorf("detail = %q, want %q", prob.Detail, "The requested job does not exist.")
		}
		if prob.Code != "job_not_found" {
			t.Errorf("code = %q, want %q", prob.Code, "job_not_found")
		}
	})
}

func TestInternalDatabaseErrorsDoNotLeak(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))

	router := gin.New()
	router.GET("/db-fail", func(c *gin.Context) {
		dbErr := errors.New("pq: password authentication failed for user 'postgres' on host '10.0.0.5'")
		HandleError(c, logger, dbErr)
	})

	router.GET("/wrapped-db-not-found", func(c *gin.Context) {
		dbErr := errors.New("SELECT * FROM sensitive_table WHERE token = 'secret-auth-token': 0 rows returned")
		appErr := apperrors.Wrap(dbErr, apperrors.ErrNotFound, "The requested item was not found.", "item_not_found")
		HandleError(c, logger, appErr)
	})

	t.Run("unexpected 500 error hides details from client but logs them", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/db-fail", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500", rec.Code)
		}

		body := rec.Body.String()
		for _, sensitive := range []string{"password", "postgres", "10.0.0.5", "pq:"} {
			if strings.Contains(body, sensitive) {
				t.Fatalf("client response leaked sensitive info: %s", body)
			}
		}

		var prob ProblemDetails
		if err := json.Unmarshal(rec.Body.Bytes(), &prob); err != nil {
			t.Fatal(err)
		}
		if prob.Detail != "An internal server error occurred." {
			t.Errorf("detail = %q, want generic message", prob.Detail)
		}
		if prob.Code != "internal_error" {
			t.Errorf("code = %q, want 'internal_error'", prob.Code)
		}

		// Verify logger recorded the internal error
		logOutput := logBuf.String()
		if !strings.Contains(logOutput, "password authentication failed") {
			t.Fatalf("logger should contain underlying error, got: %s", logOutput)
		}
	})

	t.Run("wrapped error hides database cause while preserving safe detail", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/wrapped-db-not-found", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}

		body := rec.Body.String()
		for _, sensitive := range []string{"SELECT", "sensitive_table", "secret-auth-token"} {
			if strings.Contains(body, sensitive) {
				t.Fatalf("client response leaked sensitive query: %s", body)
			}
		}

		var prob ProblemDetails
		if err := json.Unmarshal(rec.Body.Bytes(), &prob); err != nil {
			t.Fatal(err)
		}
		if prob.Detail != "The requested item was not found." {
			t.Errorf("detail = %q, want %q", prob.Detail, "The requested item was not found.")
		}
		if prob.Code != "item_not_found" {
			t.Errorf("code = %q, want 'item_not_found'", prob.Code)
		}
	})
}

func TestClientErrorsDoNotLogAtErrorLevel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelError}))

	router := gin.New()
	router.GET("/client-err", func(c *gin.Context) {
		HandleError(c, logger, apperrors.ErrNotFound)
	})

	req := httptest.NewRequest(http.MethodGet, "/client-err", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}

	if logBuf.Len() > 0 {
		t.Fatalf("expected no ERROR logs for expected client errors, got: %s", logBuf.String())
	}
}

func TestDoubleWriteProtection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/double-write", func(c *gin.Context) {
		WriteProblem(c, ProblemDetails{
			Status: http.StatusNotFound,
			Detail: "First error",
		})
		// Attempt second write
		WriteProblem(c, ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail: "Second error should be ignored",
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/double-write", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}

	var prob ProblemDetails
	if err := json.Unmarshal(rec.Body.Bytes(), &prob); err != nil {
		t.Fatal(err)
	}
	if prob.Detail != "First error" {
		t.Fatalf("detail = %q, want 'First error'", prob.Detail)
	}
}

func TestPanicRecoveryProblemDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))

	router := gin.New()
	router.Use(Recovery(logger))
	router.GET("/panic", func(c *gin.Context) {
		panic("database connection string postgres://secret@db:5432")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != ProblemContentType {
		t.Fatalf("Content-Type = %q, want %q", ct, ProblemContentType)
	}

	body := rec.Body.String()
	if strings.Contains(body, "postgres://secret@db:5432") {
		t.Fatalf("panic value leaked to client: %s", body)
	}

	var prob ProblemDetails
	if err := json.Unmarshal(rec.Body.Bytes(), &prob); err != nil {
		t.Fatal(err)
	}
	if prob.Status != 500 || prob.Title != "Internal Server Error" || prob.Code != "internal_error" {
		t.Fatalf("unexpected problem response: %+v", prob)
	}
	if prob.Detail != "An internal server error occurred." {
		t.Fatalf("unexpected detail: %s", prob.Detail)
	}

	// Verify logger caught the panic and logged it
	logOutput := logBuf.String()
	if !strings.Contains(logOutput, "HTTP handler panic recovered") || !strings.Contains(logOutput, "stack") {
		t.Fatalf("logger should record recovered panic and stack, got: %s", logOutput)
	}
}

func TestPanicRecoveryAfterResponseCommitted(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))

	router := gin.New()
	router.Use(Recovery(logger))
	router.GET("/panic-after-write", func(c *gin.Context) {
		c.Status(http.StatusOK)
		_, _ = c.Writer.Write([]byte("committed body"))
		c.Writer.Flush()
		panic("panic after writing body")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic-after-write", nil)
	rec := httptest.NewRecorder()

	defer func() {
		r := recover()
		if r == nil && rec.Body.String() != "committed body" {
			t.Fatalf("unexpected body: %s", rec.Body.String())
		}
	}()

	router.ServeHTTP(rec, req)

	// Verify response body still starts with committed bytes and was not overwritten by a 500 JSON
	if !strings.HasPrefix(rec.Body.String(), "committed body") {
		t.Fatalf("expected committed response, got %s", rec.Body.String())
	}
	// Verify logger still recorded the panic
	if !strings.Contains(logBuf.String(), "HTTP handler panic recovered") {
		t.Fatalf("expected log of panic, got: %s", logBuf.String())
	}
}

func TestErrorHandlerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandler(nil))

	router.GET("/use-context-error", func(c *gin.Context) {
		_ = c.Error(apperrors.NotFound("resource not found", "resource_missing"))
	})

	req := httptest.NewRequest(http.MethodGet, "/use-context-error", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != ProblemContentType {
		t.Fatalf("Content-Type = %q, want %q", ct, ProblemContentType)
	}

	var prob ProblemDetails
	if err := json.Unmarshal(rec.Body.Bytes(), &prob); err != nil {
		t.Fatal(err)
	}
	if prob.Detail != "resource not found" || prob.Code != "resource_missing" {
		t.Fatalf("unexpected problem response: %+v", prob)
	}
}
