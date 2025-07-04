package config

import (
	"os"
	baseConfig "video-processor/internal/config"
)

type APIConfig struct {
	Port         string
	ProcessorURL string
	*baseConfig.DirectoryConfig
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func New() *APIConfig {
	return &APIConfig{
		Port:            GetEnv("PORT", "8081"),
		ProcessorURL:    GetEnv("PROCESSOR_URL", "http://localhost:8082"),
		DirectoryConfig: baseConfig.NewDirectoryConfig(),
	}
}

func (c *APIConfig) CreateDirectories() {
	c.DirectoryConfig.CreateDirectories()
}
