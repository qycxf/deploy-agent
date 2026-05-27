package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qycxf/deploy-agent/internal/api"
)

type envHandler struct {
	store *memoryStore
}

func (h *envHandler) GetHealth(c *gin.Context) {
	uptime := time.Since(h.store.startedAt).Seconds()
	status := "ok"
	c.JSON(http.StatusOK, api.HealthResponse{
		Status:        &status,
		UptimeSeconds: &uptime,
	})
}
