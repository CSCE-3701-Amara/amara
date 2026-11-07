package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	authhttp "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/http"
)

func (h *AssessmentHandler) stageDecision(bulk bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := "stageExecutionId"
		if bulk {
			param = "stageId"
		}
		id, ok := assessmentID(c, param)
		if !ok {
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 16384)
		var body struct {
			Decision     string   `json:"decision" binding:"required,oneof=pass reject"`
			DecisionType string   `json:"decisionType" binding:"required,eq=manual"`
			Reason       string   `json:"reason" binding:"required"`
			Expected     string   `json:"expectedStatus" binding:"omitempty,oneof=locked available in_progress pending_grading held passed rejected cancelled missed"`
			IDs          []string `json:"stageExecutionIds" binding:"max=100"`
		}
		if c.ShouldBindJSON(&body) != nil || strings.TrimSpace(body.Reason) == "" ||
			(bulk && (len(body.IDs) == 0 || body.Expected != "")) ||
			(!bulk && (len(body.IDs) != 0 || body.Expected == "")) {
			writeError(c, http.StatusBadRequest, "invalid_stage_decision", "Supply a manual pass or reject decision and a reason.")
			return
		}
		ids, stage := []string{id}, ""
		if bulk {
			ids, stage = body.IDs, id
			seen := make(map[string]bool, len(ids))
			for index, raw := range ids {
				var parsed pgtype.UUID
				if parsed.Scan(raw) != nil || !parsed.Valid || seen[parsed.String()] {
					writeError(c, http.StatusBadRequest, "invalid_stage_execution_ids", "Supply distinct stage execution UUIDs.")
					return
				}
				ids[index] = parsed.String()
				seen[ids[index]] = true
			}
		}
		identity, _ := authhttp.Identity(c)
		if identity.CompanyID == "" {
			writeError(c, http.StatusForbidden, "forbidden", "A Hiring Manager account is required.")
			return
		}
		results, err := h.repo.DecideStages(c.Request.Context(), ids, identity.CompanyID, stage, body.Decision, strings.TrimSpace(body.Reason), body.Expected)
		if err != nil || bulk {
			writeAssessment(c, results, err)
			return
		}
		writeAssessment(c, results[0], nil)
	}
}
