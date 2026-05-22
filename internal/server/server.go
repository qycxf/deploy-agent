package server

import (
    "log"

    "github.com/gin-gonic/gin"
    "github.com/qycxf/deploy-agent/config"
)

func Start(port string) error {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    if port == "" {
        port = config.GetPort()
    }
    log.Printf("Starting server on :%s", port)
    return r.Run(":" + port)
}
