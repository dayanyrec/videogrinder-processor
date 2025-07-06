package config

import (
	"os"
	baseConfig "video-processor/internal/config"
)

type APIConfig struct {
	Port         string
	ProcessorURL string
	*baseConfig.DirectoryConfig
	*baseConfig.AWSConfig
	S3Service *baseConfig.S3Service
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func New() *APIConfig {
	awsConfig := baseConfig.NewAWSConfig()

	s3Service, err := baseConfig.NewS3Service(awsConfig)
	if err != nil {
		s3Service = nil
	}

	return &APIConfig{
		Port:            GetEnv("PORT", "8081"),
		ProcessorURL:    GetEnv("PROCESSOR_URL", "http://localhost:8082"),
		DirectoryConfig: baseConfig.NewDirectoryConfig(),
		AWSConfig:       awsConfig,
		S3Service:       s3Service,
	}
}

func (c *APIConfig) CreateDirectories() {
	c.DirectoryConfig.CreateDirectories()
}

func (c *APIConfig) IsS3Enabled() bool {
	return c.S3Service != nil
}
