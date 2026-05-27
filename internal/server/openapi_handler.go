package server

import (
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qycxf/deploy-agent/internal/api"
)

type OpenAPIHandler struct {
	mu          sync.RWMutex
	startedAt   time.Time
	deployments map[string]api.Deployment
	nextID      int64
}

func NewOpenAPIHandler() *OpenAPIHandler {
	return &OpenAPIHandler{
		startedAt:   time.Now(),
		deployments: make(map[string]api.Deployment),
	}
}

func (h *OpenAPIHandler) PostAuthLogin(c *gin.Context) {
	var req api.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendAPIError(c, http.StatusBadRequest, 400, "invalid request body")
		return
	}

	if req.Username != "demo" || req.Password != "secret" {
		sendAPIError(c, http.StatusUnauthorized, 401, "invalid credentials")
		return
	}

	token := "demo-token"
	expiresIn := 3600
	c.JSON(http.StatusOK, api.AuthResponse{
		Token:     &token,
		ExpiresIn: &expiresIn,
	})
}

func (h *OpenAPIHandler) GetDeployments(c *gin.Context, params api.GetDeploymentsParams) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	items := make([]api.Deployment, 0, len(h.deployments))
	ids := make([]string, 0, len(h.deployments))
	for id := range h.deployments {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		items = append(items, h.deployments[id])
	}

	page := 1
	if params.Page != nil && *params.Page > 0 {
		page = *params.Page
	}
	pageSize := 20
	if params.PageSize != nil && *params.PageSize > 0 {
		pageSize = *params.PageSize
	}

	start := (page - 1) * pageSize
	if start > len(items) {
		start = len(items)
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}

	paged := items[start:end]
	total := len(items)

	c.JSON(http.StatusOK, api.DeploymentList{
		Total: &total,
		Items: &paged,
	})
}

func (h *OpenAPIHandler) PostDeployments(c *gin.Context) {
	var req api.CreateDeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendAPIError(c, http.StatusBadRequest, 400, "invalid request body")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.nextID++
	id := "dep-" + strconv.FormatInt(h.nextID, 10)
	now := time.Now()
	status := api.Running

	deployment := api.Deployment{
		Id:        &id,
		Name:      &req.Name,
		Version:   &req.Version,
		Status:    &status,
		CreatedAt: &now,
	}

	h.deployments[id] = deployment
	c.JSON(http.StatusCreated, deployment)
}

func (h *OpenAPIHandler) DeleteDeploymentsId(c *gin.Context, id string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.deployments[id]; !ok {
		sendAPIError(c, http.StatusNotFound, 404, "deployment not found")
		return
	}

	delete(h.deployments, id)
	c.Status(http.StatusNoContent)
}

func (h *OpenAPIHandler) GetDeploymentsId(c *gin.Context, id string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	deployment, ok := h.deployments[id]
	if !ok {
		sendAPIError(c, http.StatusNotFound, 404, "deployment not found")
		return
	}

	c.JSON(http.StatusOK, deployment)
}

func (h *OpenAPIHandler) GetHealth(c *gin.Context) {
	uptime := time.Since(h.startedAt).Seconds()
	status := "ok"
	c.JSON(http.StatusOK, api.HealthResponse{
		Status:        &status,
		UptimeSeconds: &uptime,
	})
}

func (h *OpenAPIHandler) GetMetrics(c *gin.Context) {
	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte("# HELP deploy_agent_up Service up metric\n# TYPE deploy_agent_up gauge\ndeploy_agent_up 1\n"))
}

func sendAPIError(c *gin.Context, statusCode int, code int, message string) {
	c.JSON(statusCode, api.Error{
		Code:    &code,
		Message: &message,
	})
}
