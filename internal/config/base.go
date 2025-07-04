package config

import (
	"log"
	"os"
)

type DirectoryConfig struct {
	UploadsDir string
	OutputsDir string
	TempDir    string
}

func NewDirectoryConfig() *DirectoryConfig {
	return &DirectoryConfig{
		UploadsDir: GetEnv("UPLOADS_DIR", "uploads"),
		OutputsDir: GetEnv("OUTPUTS_DIR", "outputs"),
		TempDir:    GetEnv("TEMP_DIR", "temp"),
	}
}

func (dc *DirectoryConfig) CreateDirectories() {
	dirs := []string{dc.UploadsDir, dc.OutputsDir, dc.TempDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0750); err != nil {
			log.Printf("Warning: Failed to create directory %s: %v", dir, err)
		}
	}
}

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
