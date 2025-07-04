package main

import (
	"fmt"
	"log"
	"net/http"

	"video-processor/internal/config"
	"video-processor/internal/processor"
	"video-processor/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.New()
	cfg.CreateDirectories()

	// Override port for processor service
	if cfg.Port == "8080" {
		cfg.Port = "8082"
	}

	videoService := services.NewVideoService(cfg)
	processorHandlers := processor.NewProcessorHandlers(videoService, cfg)

	r := gin.Default()

	// CORS middleware
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

	// Processor endpoints
	r.POST("/process", processorHandlers.ProcessVideoUpload)
	r.GET("/health", processorHandlers.GetProcessorStatus)

	fmt.Println("🔧 Processor service iniciado na porta", cfg.Port)
	fmt.Printf("🔧 Health check: http://localhost:%s/health\n", cfg.Port)

	log.Fatal(r.Run(":" + cfg.Port))
}
