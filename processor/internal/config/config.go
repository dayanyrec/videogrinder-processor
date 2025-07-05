package config

import (
	"os"
	baseConfig "video-processor/internal/config"
)

type ProcessorConfig struct {
	Port string
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

func New() *ProcessorConfig {
	awsConfig := baseConfig.NewAWSConfig()

	s3Service, err := baseConfig.NewS3Service(awsConfig)
	if err != nil {
		// Log error but don't fail - fallback to filesystem
		// This allows development without S3 if needed
		s3Service = nil
	}

	return &ProcessorConfig{
		Port:            GetEnv("PORT", "8082"),
		DirectoryConfig: baseConfig.NewDirectoryConfig(),
		AWSConfig:       awsConfig,
		S3Service:       s3Service,
	}
}

func (c *ProcessorConfig) CreateDirectories() {
	c.DirectoryConfig.CreateDirectories()
}

func (c *ProcessorConfig) IsS3Enabled() bool {
	return c.S3Service != nil
}
