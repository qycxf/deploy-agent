package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qycxf/deploy-agent/internal/api"
)

type authHandler struct {
	store *memoryStore
}

func (h *authHandler) PostAuthLogin(c *gin.Context) {
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
