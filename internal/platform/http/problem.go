package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/CSCE-3701-Amara/amara/internal/platform/apperrors"
	"github.com/gin-gonic/gin"
)

// ProblemContentType is the MIME type defined by RFC 9457 for Problem Details.
const ProblemContentType = "application/problem+json"

// FieldError describes a specific field validation error.
type FieldError = apperrors.FieldError

// ProblemDetails represents an RFC 9457 Problem Details HTTP error response payload.
type ProblemDetails struct {
	Type     string       `json:"type"`
	Title    string       `json:"title"`
	Status   int          `json:"status"`
	Detail   string       `json:"detail,omitempty"`
	Instance string       `json:"instance,omitempty"`
	Code     string       `json:"code,omitempty"`
	Errors   []FieldError `json:"errors,omitempty"`
}

// Error implements the error interface for ProblemDetails.
func (p ProblemDetails) Error() string {
	if p.Detail != "" {
		return fmt.Sprintf("%s (%d): %s", p.Title, p.Status, p.Detail)
	}
	return fmt.Sprintf("%s (%d)", p.Title, p.Status)
}

// WriteProblem renders an RFC 9457 Problem Details response to the client.
// It sets Content-Type to application/problem+json, sets the matching HTTP status code,
// and ensures the handler chain is aborted.
// If the response has already been committed (written), it avoids writing again.
func WriteProblem(c *gin.Context, prob ProblemDetails) {
	if c == nil || c.Writer == nil || c.Writer.Written() {
		return
	}

	if prob.Type == "" {
		prob.Type = "about:blank"
	}
	if prob.Status == 0 {
		prob.Status = http.StatusInternalServerError
	}
	if prob.Title == "" {
		prob.Title = defaultTitle(prob.Status)
	}
	if prob.Instance == "" && c.Request != nil && c.Request.URL != nil {
		prob.Instance = c.Request.URL.Path
	}

	c.Abort()
	c.Header("Content-Type", ProblemContentType)
	c.Status(prob.Status)
	if err := json.NewEncoder(c.Writer).Encode(prob); err != nil {
		_ = err
	}
}

// MapError translates any Go error into a standardized RFC 9457 ProblemDetails struct.
// Known application errors are mapped to appropriate HTTP status codes and safe details.
// Unexpected errors are mapped to a generic 500 response to prevent internal details from leaking.
func MapError(err error) ProblemDetails {
	if err == nil {
		return ProblemDetails{}
	}

	// 1. Explicit ProblemDetails
	var pd ProblemDetails
	if errors.As(err, &pd) {
		if pd.Type == "" {
			pd.Type = "about:blank"
		}
		if pd.Status == 0 {
			pd.Status = http.StatusInternalServerError
		}
		if pd.Title == "" {
			pd.Title = defaultTitle(pd.Status)
		}
		return pd
	}

	// 2. Semantically invalid request (ValidationError with field errors)
	var valErr *apperrors.ValidationError
	if errors.As(err, &valErr) {
		detail := valErr.Message
		if detail == "" {
			detail = defaultDetail(http.StatusUnprocessableEntity)
		}
		return ProblemDetails{
			Type:   "about:blank",
			Title:  defaultTitle(http.StatusUnprocessableEntity),
			Status: http.StatusUnprocessableEntity,
			Detail: detail,
			Code:   "validation_failed",
			Errors: valErr.Errors,
		}
	}

	// 3. Structured application error (AppError with safe Message and Code)
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		status := mapCategoryToStatus(appErr.Category)
		detail := appErr.Message
		if detail == "" {
			detail = defaultDetail(status)
		}
		code := appErr.Code
		if code == "" {
			code = defaultCode(status)
		}
		return ProblemDetails{
			Type:   "about:blank",
			Title:  defaultTitle(status),
			Status: status,
			Detail: detail,
			Code:   code,
		}
	}

	// 4. Malformed JSON / Request syntax error
	var syntaxErr *json.SyntaxError
	var unmarshalTypeErr *json.UnmarshalTypeError
	if errors.As(err, &syntaxErr) || errors.As(err, &unmarshalTypeErr) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
		return ProblemDetails{
			Type:   "about:blank",
			Title:  defaultTitle(http.StatusBadRequest),
			Status: http.StatusBadRequest,
			Detail: "The request body contains malformed JSON or invalid syntax.",
			Code:   defaultCode(http.StatusBadRequest),
		}
	}

	// 5. Sentinel errors / wrapped errors mapped via errors.Is
	status := mapCategoryToStatus(err)
	return ProblemDetails{
		Type:   "about:blank",
		Title:  defaultTitle(status),
		Status: status,
		Detail: defaultDetail(status),
		Code:   defaultCode(status),
	}
}

// HandleError maps err to RFC 9457 Problem Details, logs unexpected internal errors,
// and writes the response to c. If headers have already been committed, it does nothing.
func HandleError(c *gin.Context, logger *slog.Logger, err error) {
	if err == nil || c == nil || c.Writer == nil || c.Writer.Written() {
		return
	}

	prob := MapError(err)

	if prob.Status >= http.StatusInternalServerError {
		if logger == nil {
			logger = slog.Default()
		}
		logger.ErrorContext(c.Request.Context(), "HTTP handler unexpected error",
			"request_id", RequestID(c.Request.Context()),
			"status", prob.Status,
			"error", err,
		)
	}

	WriteProblem(c, prob)
}

func mapCategoryToStatus(err error) int {
	switch {
	case errors.Is(err, apperrors.ErrInvalidInput):
		return http.StatusBadRequest // 400
	case errors.Is(err, apperrors.ErrUnauthenticated):
		return http.StatusUnauthorized // 401
	case errors.Is(err, apperrors.ErrForbidden):
		return http.StatusForbidden // 403
	case errors.Is(err, apperrors.ErrNotFound):
		return http.StatusNotFound // 404
	case errors.Is(err, apperrors.ErrConflict):
		return http.StatusConflict // 409
	case errors.Is(err, apperrors.ErrValidation):
		return http.StatusUnprocessableEntity // 422
	case errors.Is(err, apperrors.ErrRateLimitExceeded):
		return http.StatusTooManyRequests // 429
	case errors.Is(err, apperrors.ErrServiceUnavailable):
		return http.StatusServiceUnavailable // 503
	case errors.Is(err, apperrors.ErrInternal):
		return http.StatusInternalServerError // 500
	default:
		return http.StatusInternalServerError // 500
	}
}

func defaultTitle(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "Bad Request"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusForbidden:
		return "Forbidden"
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusMethodNotAllowed:
		return "Method Not Allowed"
	case http.StatusConflict:
		return "Conflict"
	case http.StatusUnprocessableEntity:
		return "Unprocessable Content"
	case http.StatusTooManyRequests:
		return "Too Many Requests"
	case http.StatusInternalServerError:
		return "Internal Server Error"
	case http.StatusServiceUnavailable:
		return "Service Unavailable"
	default:
		text := http.StatusText(status)
		if text == "" {
			return "Error"
		}
		return text
	}
}

func defaultDetail(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "The request payload or parameters are invalid."
	case http.StatusUnauthorized:
		return "Authentication is required to access this resource."
	case http.StatusForbidden:
		return "You do not have permission to perform this action."
	case http.StatusNotFound:
		return "The requested resource was not found."
	case http.StatusMethodNotAllowed:
		return "The request method is not allowed."
	case http.StatusConflict:
		return "The request could not be completed due to a conflict with the current state of the target resource."
	case http.StatusUnprocessableEntity:
		return "The request was well-formed but could not be processed due to semantic errors."
	case http.StatusTooManyRequests:
		return "Too many requests. Please try again later."
	case http.StatusInternalServerError:
		return "An internal server error occurred."
	case http.StatusServiceUnavailable:
		return "The service is temporarily unavailable. Please try again later."
	default:
		return "An error occurred while processing the request."
	}
}

func defaultCode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "bad_request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusForbidden:
		return "forbidden"
	case http.StatusNotFound:
		return "not_found"
	case http.StatusMethodNotAllowed:
		return "method_not_allowed"
	case http.StatusConflict:
		return "conflict"
	case http.StatusUnprocessableEntity:
		return "validation_failed"
	case http.StatusTooManyRequests:
		return "rate_limit_exceeded"
	case http.StatusInternalServerError:
		return "internal_error"
	case http.StatusServiceUnavailable:
		return "service_unavailable"
	default:
		return "error"
	}
}
