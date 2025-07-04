package main

import (
	"fmt"
	"log"
	"net/http"

	"video-processor/web/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.New()

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

	r.Static("/static", cfg.StaticDir)

	r.GET("/", func(c *gin.Context) {
		c.File(cfg.StaticDir + "/index.html")
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "web",
		})
	})

	fmt.Printf("🎬 Web Service iniciado na porta %s\n", cfg.Port)
	fmt.Printf("🌐 Serving static files from %s\n", cfg.StaticDir)
	fmt.Printf("🔧 API URL configurado: %s\n", cfg.APIURL)

	log.Fatal(r.Run(":" + cfg.Port))
}
