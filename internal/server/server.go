package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/qycxf/deploy-agent/config"
	"github.com/qycxf/deploy-agent/internal/api"
	"github.com/qycxf/deploy-agent/internal/db"
	"github.com/qycxf/deploy-agent/internal/repository"
)

func Start(port string) error {
	r := gin.Default()

	dbh, err := db.New()
	if err != nil {
		return err
	}
	userRepo := repository.NewUserRepository(dbh.DB)
	if err := userRepo.EnsureDemoUser(); err != nil {
		return err
	}

	handler := NewOpenAPIHandler(userRepo)
	api.RegisterHandlers(r, handler)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	if port == "" {
		port = config.GetPort()
	}
	log.Printf("Starting server on :%s", port)
	return r.Run(":" + port)
}
