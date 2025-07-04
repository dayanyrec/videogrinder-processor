package main

import (
	"fmt"
	"log"
	"net/http"

	"video-processor/internal/config"
	"video-processor/processor/internal/handlers"
	"video-processor/processor/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.New()
	cfg.CreateDirectories()

	videoService := services.NewVideoService(cfg)
	processorHandlers := handlers.NewProcessorHandlers(videoService, cfg)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	r.POST("/process", processorHandlers.ProcessVideoUpload)
	r.GET("/health", processorHandlers.GetProcessorStatus)

	fmt.Println("🔧 Processor service iniciado na porta", cfg.Port)
	fmt.Printf("🔧 Health check: http://localhost:%s/health\n", cfg.Port)

	log.Fatal(r.Run(":" + cfg.Port))
}
