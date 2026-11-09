package http

import (
	"context"
	"encoding/json"
	authhttp "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/http"
	"github.com/CSCE-3701-Amara/amara/internal/invw/adapters/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (h *AssessmentHandler) ConfigureRunner(r database.CodeRunner) {
	h.runner = r
	h.runs = make(chan struct{}, 2)
}
func (h *AssessmentHandler) runCode(c *gin.Context) {
	id, ok := assessmentID(c, "attemptId")
	if !ok {
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 96<<10)
	var body struct {
		QuestionID string `json:"questionId"`
		Language   string `json:"language"`
		Source     string `json:"source"`
		Stdin      string `json:"stdin"`
	}
	if c.ShouldBindJSON(&body) != nil || len(body.Source) > 65536 || len(body.Stdin) > 16384 || body.Language == "" {
		writeError(c, 400, "invalid_run", "Supply a question, language, source and input within the execution limits.")
		return
	}
	identity, _ := authhttp.Identity(c)
	attempt, err := h.repo.CandidateAttempt(c.Request.Context(), identity.AccountID, id, false, false)
	if err != nil {
		writeAssessment(c, nil, err)
		return
	}
	if attempt.Status != "in_progress" || !attempt.Deadline.After(time.Now()) {
		writeError(c, 409, "assessment_conflict", "The attempt is closed.")
		return
	}
	var questions []struct {
		ID     string `json:"id"`
		Format string `json:"responseFormat"`
	}
	if err := json.Unmarshal(attempt.Questions, &questions); err != nil {
		writeAssessment(c, nil, err)
		return
	}
	found := false
	for _, q := range questions {
		if q.ID == body.QuestionID && q.Format == "code" {
			found = true
		}
	}
	if !found {
		writeError(c, 404, "question_not_found", "Code question not found.")
		return
	}
	allowed, err := h.repo.CanRunCode(c.Request.Context(), identity.AccountID, id)
	if err != nil {
		writeAssessment(c, nil, err)
		return
	}
	if !allowed {
		writeError(c, 409, "assessment_conflict", "The stage is closed.")
		return
	}
	if h.runner == nil {
		writeError(c, 503, "execution_unavailable", "Code execution is not configured.")
		return
	}
	select {
	case h.runs <- struct{}{}:
		defer func() { <-h.runs }()
	default:
		writeError(c, 429, "execution_busy", "Execution capacity is busy. Retry shortly.")
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 25*time.Second)
	defer cancel()
	out, success, err := h.runner.Run(ctx, body.Language, body.Source, body.Stdin)
	if err != nil {
		writeError(c, 503, "execution_unavailable", "Execution could not complete. Retry shortly.")
		return
	}
	c.JSON(http.StatusOK, gin.H{"stdout": out, "succeeded": success})
}
