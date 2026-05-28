package server

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qycxf/deploy-agent/internal/api"
	"github.com/qycxf/deploy-agent/internal/repository"
)

type memoryStore struct {
	mu          sync.RWMutex
	startedAt   time.Time
	deployments map[string]api.Deployment
	nextID      int64
}

func newMemoryStore() *memoryStore {
	return &memoryStore{
		startedAt:   time.Now(),
		deployments: make(map[string]api.Deployment),
	}
}

// OpenAPIHandler aggregates individual handlers to satisfy api.ServerInterface.
type OpenAPIHandler struct {
	store   *memoryStore
	auth    *authHandler
	service *serviceHandler
	env     *envHandler
	metrics *metricsHandler
}

func NewOpenAPIHandler(userRepo *repository.UserRepository) *OpenAPIHandler {
	store := newMemoryStore()
	return &OpenAPIHandler{
		store:   store,
		auth:    &authHandler{userRepo: userRepo},
		service: &serviceHandler{store: store},
		env:     &envHandler{store: store},
		metrics: &metricsHandler{},
	}
}

func (h *OpenAPIHandler) PostAuthLogin(c *gin.Context) {
	h.auth.PostAuthLogin(c)
}

func (h *OpenAPIHandler) GetDeployments(c *gin.Context, params api.GetDeploymentsParams) {
	h.service.GetDeployments(c, params)
}

func (h *OpenAPIHandler) PostDeployments(c *gin.Context) {
	h.service.PostDeployments(c)
}

func (h *OpenAPIHandler) DeleteDeploymentsId(c *gin.Context, id string) {
	h.service.DeleteDeploymentsId(c, id)
}

func (h *OpenAPIHandler) GetDeploymentsId(c *gin.Context, id string) {
	h.service.GetDeploymentsId(c, id)
}

func (h *OpenAPIHandler) GetHealth(c *gin.Context) {
	h.env.GetHealth(c)
}

func (h *OpenAPIHandler) GetMetrics(c *gin.Context) {
	h.metrics.GetMetrics(c)
}

func sendAPIError(c *gin.Context, statusCode int, code int, message string) {
	c.JSON(statusCode, api.Error{
		Code:    &code,
		Message: &message,
	})
}
