package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	authhttp "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/http"
	authdomain "github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	authports "github.com/CSCE-3701-Amara/amara/internal/auth/ports"
	"github.com/CSCE-3701-Amara/amara/internal/invw/adapters/database"
)

type AssessmentHandler struct {
	repo     *database.AssessmentRepository
	sessions authports.SessionValidator
}

func NewAssessmentHandler(repo *database.AssessmentRepository, sessions authports.SessionValidator) *AssessmentHandler {
	return &AssessmentHandler{repo: repo, sessions: sessions}
}

func (h *AssessmentHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	company := authhttp.CompanyAuth(h.sessions, writeError)
	v1.POST("/stage-executions/:stageExecutionId/decisions", company, h.stageDecision(false))
	v1.POST("/pipeline-stages/:stageId/decisions", company, h.stageDecision(true))
	v1.GET("/jobs/:id/pipeline/dashboard", company, h.pipelineDashboard)
	v1.POST("/pipeline-stages/:stageId/threshold/preview", company, h.stageThreshold(false))
	v1.PUT("/pipeline-stages/:stageId/threshold", company, h.stageThreshold(true))
	candidate := authhttp.CandidateAuthWithErrorWriter(h.sessions, writeError)
	v1.GET("/applications/:applicationId/assessment-instances", candidate, h.listInstances)
	v1.POST("/assessment-instances/:instanceId/start", candidate, h.attempt("instanceId", true, false))
	v1.GET("/assessment-attempts/:attemptId", candidate, h.attempt("attemptId", false, false))
	v1.POST("/assessment-attempts/:attemptId/submission", candidate, h.attempt("attemptId", false, true))
	v1.PUT("/assessment-attempts/:attemptId/submissions/:questionId", candidate, h.save)
	v1.POST("/assessment-attempts/:attemptId/evaluation", authhttp.CompanyAuth(h.sessions, writeError), h.evaluate)
	v1.GET("/assessment-attempts/:attemptId/evaluation", authhttp.CompanyAuth(h.sessions, writeError), h.review)
	v1.POST("/interviews/:interviewId/outcome", authhttp.AuthenticatedAuth(h.sessions), h.interviewOutcome)
}

func (h *AssessmentHandler) listInstances(c *gin.Context) {
	id, ok := assessmentID(c, "applicationId")
	if !ok {
		return
	}
	identity, _ := authhttp.Identity(c)
	items, err := h.repo.ListInstances(c.Request.Context(), identity.AccountID, id)
	writeAssessment(c, items, err)
}

func (h *AssessmentHandler) attempt(param string, start, submit bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := assessmentID(c, param)
		if !ok {
			return
		}
		identity, _ := authhttp.Identity(c)
		attempt, err := h.repo.CandidateAttempt(c.Request.Context(), identity.AccountID, id, start, submit)
		writeAssessment(c, attempt, err)
	}
}

func (h *AssessmentHandler) save(c *gin.Context) {
	attempt, ok := assessmentID(c, "attemptId")
	if !ok {
		return
	}
	question, ok := assessmentID(c, "questionId")
	if !ok {
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
	var body struct {
		Response json.RawMessage `json:"response" binding:"required"`
	}
	if c.ShouldBindJSON(&body) != nil {
		writeError(c, http.StatusBadRequest, "invalid_answer", "A response is required.")
		return
	}
	identity, _ := authhttp.Identity(c)
	at, err := h.repo.SaveAnswer(c.Request.Context(), identity.AccountID, attempt, question, body.Response)
	writeAssessment(c, gin.H{"id": question, "attemptId": attempt, "questionId": question, "response": body.Response, "savedAt": at}, err)
}

func (h *AssessmentHandler) evaluate(c *gin.Context) {
	id, ok := assessmentID(c, "attemptId")
	if !ok {
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
	var body struct {
		Evidence []database.QuestionEvidence `json:"evidence" binding:"required"`
	}
	if c.ShouldBindJSON(&body) != nil {
		writeError(c, http.StatusBadRequest, "invalid_evaluation", "Question evidence is required.")
		return
	}
	identity, _ := authhttp.Identity(c)
	if identity.CompanyID == "" {
		writeError(c, http.StatusForbidden, "forbidden", "A company account is required.")
		return
	}
	err := h.repo.FinalizeEvidence(c.Request.Context(), id, identity.CompanyID, body.Evidence)
	if err != nil {
		writeAssessment(c, nil, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AssessmentHandler) review(c *gin.Context) {
	id, ok := assessmentID(c, "attemptId")
	if !ok {
		return
	}
	identity, _ := authhttp.Identity(c)
	result, err := h.repo.ReviewEvaluation(c.Request.Context(), id, identity.CompanyID)
	writeAssessment(c, result, err)
}

func (h *AssessmentHandler) interviewOutcome(c *gin.Context) {
	id, ok := assessmentID(c, "interviewId")
	if !ok {
		return
	}
	identity, _ := authhttp.Identity(c)
	var company, interviewer string
	switch identity.AccountType {
	case authdomain.AccountTypeCompany:
		company = identity.CompanyID
	case authdomain.AccountTypeInterviewer:
		interviewer = identity.AccountID
	default:
		writeError(c, http.StatusForbidden, "forbidden", "Only interview staff can record the outcome.")
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 4096)
	var body struct {
		Outcome string `json:"outcome" binding:"required,oneof=held candidate_absent interviewer_absent"`
	}
	if c.ShouldBindJSON(&body) != nil {
		writeError(c, http.StatusBadRequest, "invalid_outcome", "An interview outcome is required.")
		return
	}
	result, err := h.repo.RecordInterviewOutcome(c.Request.Context(), id, company, interviewer, body.Outcome)
	writeAssessment(c, result, err)
}

func assessmentID(c *gin.Context, param string) (string, bool) {
	var id pgtype.UUID
	if id.Scan(c.Param(param)) != nil || !id.Valid {
		writeError(c, http.StatusBadRequest, "invalid_id", "The resource ID must be a UUID.")
		return "", false
	}
	return id.String(), true
}

func writeAssessment(c *gin.Context, result any, err error) {
	switch {
	case err == nil:
		c.JSON(http.StatusOK, result)
	case errors.Is(err, pgx.ErrNoRows):
		writeError(c, http.StatusNotFound, "assessment_not_found", "The requested resource was not found.")
	case errors.Is(err, database.ErrDecisionConflict):
		writeError(c, http.StatusConflict, "stage_decision_conflict", err.Error())
	case errors.Is(err, database.ErrThresholdConflict):
		writeError(c, http.StatusConflict, "stage_threshold_conflict", err.Error())
	case errors.Is(err, database.ErrAssessmentConflict):
		writeError(c, http.StatusConflict, "assessment_conflict", err.Error())
	default:
		slog.ErrorContext(c.Request.Context(), "assessment request failed", "error", err)
		writeError(c, http.StatusInternalServerError, "internal_error", "The assessment request failed.")
	}
}
