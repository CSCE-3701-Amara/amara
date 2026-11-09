package http

import (
	authhttp "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/http"
	"github.com/CSCE-3701-Amara/amara/internal/invw/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *PipelineHandler) graph(save bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := pipelineID(c, "versionId")
		if !ok {
			return
		}
		identity, _ := authhttp.Identity(c)
		var g domain.Graph
		var err error
		if save {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 2<<20)
			if c.ShouldBindJSON(&g) != nil {
				writeError(c, 400, "invalid_graph", "A graph is required.")
				return
			}
			for _, s := range g.Stages {
				if !normalizeIDs(s.ModuleIDs) || (s.Interview != nil && !normalizeIDs(s.Interview.InterviewerIDs)) {
					writeError(c, 422, "invalid_graph", "Module and interviewer IDs must be UUIDs.")
					return
				}
			}
			g, err = h.repo.SaveGraph(c.Request.Context(), identity.CompanyID, id, g)
		} else {
			g, err = h.repo.Graph(c.Request.Context(), identity.CompanyID, id)
		}
		writePipelineResult(c, 200, g, err)
	}
}
