package config

import (
	"os"
	baseConfig "video-processor/internal/config"
)

type ProcessorConfig struct {
	Port string
	*baseConfig.DirectoryConfig
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func New() *ProcessorConfig {
	return &ProcessorConfig{
		Port:            GetEnv("PORT", "8082"),
		DirectoryConfig: baseConfig.NewDirectoryConfig(),
	}
}

func (c *ProcessorConfig) CreateDirectories() {
	c.DirectoryConfig.CreateDirectories()
}
