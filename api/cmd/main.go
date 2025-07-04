package main

import (
	"fmt"
	"log"
	"net/http"

	"video-processor/internal/api"
	"video-processor/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.New()
	cfg.CreateDirectories()

	apiHandlers := api.NewAPIHandlers(cfg)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	r.Static("/uploads", "./"+cfg.UploadsDir)
	r.Static("/outputs", "./"+cfg.OutputsDir)

	r.GET("/health", apiHandlers.GetAPIHealth)

	apiV1 := r.Group("/api/v1")
	apiV1.GET("/health", apiHandlers.GetAPIHealth)
	apiV1.POST("/videos", apiHandlers.CreateVideo)
	apiV1.GET("/videos", apiHandlers.GetVideos)
	apiV1.GET("/videos/:filename/download", apiHandlers.GetVideoDownload)
	apiV1.DELETE("/videos/:filename", apiHandlers.DeleteVideo)

	fmt.Printf("🎬 API Service iniciado na porta %s\n", cfg.Port)
	fmt.Printf("🔧 Processor URL configurado: %s\n", cfg.ProcessorURL)

	log.Fatal(r.Run(":" + cfg.Port))
}
