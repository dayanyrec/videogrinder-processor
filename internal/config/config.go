package config

import (
	"log"
	"os"
)

type Config struct {
	Port         string
	UploadsDir   string
	OutputsDir   string
	TempDir      string
	ProcessorURL string
	APIURL       string
}

func New() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		UploadsDir:   getEnv("UPLOADS_DIR", "uploads"),
		OutputsDir:   getEnv("OUTPUTS_DIR", "outputs"),
		TempDir:      getEnv("TEMP_DIR", "temp"),
		ProcessorURL: getEnv("PROCESSOR_URL", "http://localhost:8082"),
		APIURL:       getEnv("API_URL", "http://localhost:8081"),
	}
}

func (c *Config) CreateDirectories() {
	dirs := []string{c.UploadsDir, c.OutputsDir, c.TempDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0750); err != nil {
			log.Printf("Warning: Failed to create directory %s: %v", dir, err)
		}
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
