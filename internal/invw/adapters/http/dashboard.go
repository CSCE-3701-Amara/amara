package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	authhttp "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/http"
	"github.com/CSCE-3701-Amara/amara/internal/invw/adapters/database"
)

func (h *AssessmentHandler) pipelineDashboard(c *gin.Context) {
	posting, ok := assessmentID(c, "id")
	if !ok {
		return
	}
	identity, _ := authhttp.Identity(c)
	dashboard, err := h.repo.PipelineDashboard(c.Request.Context(), posting, identity.CompanyID)
	writeAssessment(c, dashboard, err)
}

func (h *AssessmentHandler) stageThreshold(apply bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		stage, ok := assessmentID(c, "stageId")
		if !ok {
			return
		}
		var body struct {
			Threshold *int                        `json:"threshold" binding:"required,min=0,max=100"`
			Expected  *int                        `json:"expectedThreshold" binding:"required,min=0,max=100"`
			Preview   []database.ThresholdOutcome `json:"preview"`
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		if c.ShouldBindJSON(&body) != nil || (apply && body.Preview == nil) {
			writeError(c, http.StatusBadRequest, "invalid_threshold", "Supply thresholds between 0 and 100 and review the preview before applying.")
			return
		}
		identity, _ := authhttp.Identity(c)
		outcomes, err := h.repo.ChangeStageThreshold(c.Request.Context(), stage, identity.CompanyID, *body.Threshold, *body.Expected, apply, body.Preview)
		writeAssessment(c, outcomes, err)
	}
}
