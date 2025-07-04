package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ValidateProcessingInputs(videoPath, timestamp string) error {
	if strings.Contains(videoPath, "..") || strings.Contains(timestamp, "..") {
		return fmt.Errorf("invalid path parameters")
	}
	return nil
}

func ValidatePathSafety(paths ...string) error {
	dangerousChars := []string{";", "&", "|", "$", "`", "(", ")", "{", "}", "[", "]", "*", "?", "<", ">", "~"}
	for _, path := range paths {
		for _, char := range dangerousChars {
			if strings.Contains(path, char) {
				return fmt.Errorf("invalid characters in file path")
			}
		}
	}
	return nil
}

func ValidateOutputPath(zipPath, outputsDir string) error {
	cleanZipPath := filepath.Clean(zipPath)
	outputsDirAbs, _ := filepath.Abs(outputsDir)
	absZipPath, _ := filepath.Abs(cleanZipPath)
	if !strings.HasPrefix(absZipPath, outputsDirAbs+string(filepath.Separator)) {
		return fmt.Errorf("invalid zip path")
	}
	return nil
}

func SetupTempDirectory(tempDir string) error {
	if err := os.MkdirAll(tempDir, 0750); err != nil {
		return fmt.Errorf("erro ao criar diretório temporário: %w", err)
	}
	return nil
}

func CleanupTempDirectory(tempDir string) {
	if err := os.RemoveAll(tempDir); err != nil {
		log.Printf("Warning: Failed to remove temp directory %s: %v", tempDir, err)
	}
}
