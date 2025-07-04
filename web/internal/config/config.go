package config

import (
	"os"
)

type WebConfig struct {
	Port      string
	APIURL    string
	StaticDir string
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func New() *WebConfig {
	return &WebConfig{
		Port:      GetEnv("PORT", "8080"),
		APIURL:    GetEnv("API_URL", "http://localhost:8081"),
		StaticDir: GetEnv("STATIC_DIR", "./web/static"),
	}
}
