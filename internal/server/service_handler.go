package server

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qycxf/deploy-agent/internal/api"
)

type serviceHandler struct {
	store *memoryStore
}

func (h *serviceHandler) GetDeployments(c *gin.Context, params api.GetDeploymentsParams) {
	h.store.mu.RLock()
	defer h.store.mu.RUnlock()

	items := make([]api.Deployment, 0, len(h.store.deployments))
	ids := make([]string, 0, len(h.store.deployments))
	for id := range h.store.deployments {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		items = append(items, h.store.deployments[id])
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

func (h *serviceHandler) PostDeployments(c *gin.Context) {
	var req api.CreateDeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendAPIError(c, http.StatusBadRequest, 400, "invalid request body")
		return
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	h.store.nextID++
	id := "dep-" + strconv.FormatInt(h.store.nextID, 10)
	now := nowTime()
	status := api.Running

	deployment := api.Deployment{
		Id:        &id,
		Name:      &req.Name,
		Version:   &req.Version,
		Status:    &status,
		CreatedAt: &now,
	}

	h.store.deployments[id] = deployment
	c.JSON(http.StatusCreated, deployment)
}

func (h *serviceHandler) DeleteDeploymentsId(c *gin.Context, id string) {
	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	if _, ok := h.store.deployments[id]; !ok {
		sendAPIError(c, http.StatusNotFound, 404, "deployment not found")
		return
	}

	delete(h.store.deployments, id)
	c.Status(http.StatusNoContent)
}

func (h *serviceHandler) GetDeploymentsId(c *gin.Context, id string) {
	h.store.mu.RLock()
	defer h.store.mu.RUnlock()

	deployment, ok := h.store.deployments[id]
	if !ok {
		sendAPIError(c, http.StatusNotFound, 404, "deployment not found")
		return
	}

	c.JSON(http.StatusOK, deployment)
}
