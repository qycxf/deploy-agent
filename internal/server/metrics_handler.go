package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type metricsHandler struct{}

func (h *metricsHandler) GetMetrics(c *gin.Context) {
	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte("# HELP deploy_agent_up Service up metric\n# TYPE deploy_agent_up gauge\ndeploy_agent_up 1\n"))
}
