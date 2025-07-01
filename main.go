package main

import (
	"fmt"
	"log"

	"video-processor/internal/api"
	"video-processor/internal/config"
	"video-processor/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.New()
	cfg.CreateDirectories()

	videoService := services.NewVideoService(cfg)
	apiHandlers := api.NewAPIHandlers(videoService, cfg)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.Static("/uploads", "./"+cfg.UploadsDir)
	r.Static("/outputs", "./"+cfg.OutputsDir)
	r.Static("/static", "./static")

	r.GET("/", apiHandlers.HandleHome)
	r.POST("/upload", apiHandlers.HandleVideoUpload)
	r.GET("/download/:filename", apiHandlers.HandleDownload)
	r.GET("/api/status", apiHandlers.HandleStatus)

	fmt.Println("🎬 Servidor iniciado na porta", cfg.Port)
	fmt.Printf("📂 Acesse: http://localhost:%s\n", cfg.Port)

	log.Fatal(r.Run(":" + cfg.Port))
}
