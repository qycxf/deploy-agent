package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/qycxf/deploy-agent/config"
	"github.com/qycxf/deploy-agent/internal/api"
)

func Start(port string) error {
	r := gin.Default()

	handler := NewOpenAPIHandler()
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
