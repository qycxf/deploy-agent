package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qycxf/deploy-agent/internal/api"
	"github.com/qycxf/deploy-agent/internal/repository"
)

type authHandler struct {
	userRepo *repository.UserRepository
}

func (h *authHandler) PostAuthLogin(c *gin.Context) {
	var req api.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendAPIError(c, http.StatusBadRequest, 400, "invalid request body")
		return
	}

	u, err := h.userRepo.FindByUsername(req.Username)
	if err != nil {
		sendAPIError(c, http.StatusInternalServerError, 500, "failed to query user")
		return
	}
	if u == nil || repository.VerifyPassword(req.Password, u.PasswordHash) != nil {
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
