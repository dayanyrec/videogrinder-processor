package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateProcessingInputs(t *testing.T) {
	tests := []struct {
		name        string
		videoPath   string
		timestamp   string
		expectError bool
	}{
		{
			name:        "valid inputs",
			videoPath:   "uploads/test.mp4",
			timestamp:   "20240101_120000",
			expectError: false,
		},
		{
			name:        "valid nested path",
			videoPath:   "uploads/folder/test.mp4",
			timestamp:   "20240101_120000",
			expectError: false,
		},
		{
			name:        "video path with directory traversal",
			videoPath:   "uploads/../etc/passwd",
			timestamp:   "20240101_120000",
			expectError: true,
		},
		{
			name:        "timestamp with directory traversal",
			videoPath:   "uploads/test.mp4",
			timestamp:   "../../../etc",
			expectError: true,
		},
		{
			name:        "both paths with directory traversal",
			videoPath:   "../uploads/test.mp4",
			timestamp:   "../timestamp",
			expectError: true,
		},
		{
			name:        "empty inputs",
			videoPath:   "",
			timestamp:   "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProcessingInputs(tt.videoPath, tt.timestamp)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid path parameters")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePathSafety(t *testing.T) {
	tests := []struct {
		name        string
		paths       []string
		expectError bool
	}{
		{
			name:        "safe paths",
			paths:       []string{"uploads/test.mp4", "temp/frame_001.png", "outputs/result.zip"},
			expectError: false,
		},
		{
			name:        "path with semicolon",
			paths:       []string{"uploads/test;rm -rf /.mp4"},
			expectError: true,
		},
		{
			name:        "path with ampersand",
			paths:       []string{"uploads/test&whoami.mp4"},
			expectError: true,
		},
		{
			name:        "path with pipe",
			paths:       []string{"uploads/test|cat /etc/passwd.mp4"},
			expectError: true,
		},
		{
			name:        "path with dollar sign",
			paths:       []string{"uploads/test$USER.mp4"},
			expectError: true,
		},
		{
			name:        "path with backticks",
			paths:       []string{"uploads/test`whoami`.mp4"},
			expectError: true,
		},
		{
			name:        "path with parentheses",
			paths:       []string{"uploads/test(command).mp4"},
			expectError: true,
		},
		{
			name:        "path with brackets",
			paths:       []string{"uploads/test[file].mp4"},
			expectError: true,
		},
		{
			name:        "path with braces",
			paths:       []string{"uploads/test{var}.mp4"},
			expectError: true,
		},
		{
			name:        "path with asterisk",
			paths:       []string{"uploads/test*.mp4"},
			expectError: true,
		},
		{
			name:        "path with question mark",
			paths:       []string{"uploads/test?.mp4"},
			expectError: true,
		},
		{
			name:        "path with less than",
			paths:       []string{"uploads/test<file.mp4"},
			expectError: true,
		},
		{
			name:        "path with greater than",
			paths:       []string{"uploads/test>file.mp4"},
			expectError: true,
		},
		{
			name:        "path with tilde",
			paths:       []string{"uploads/test~.mp4"},
			expectError: true,
		},
		{
			name:        "multiple paths, one dangerous",
			paths:       []string{"uploads/safe.mp4", "uploads/dangerous;rm.mp4"},
			expectError: true,
		},
		{
			name:        "empty path",
			paths:       []string{""},
			expectError: false,
		},
		{
			name:        "no paths",
			paths:       []string{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePathSafety(tt.paths...)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid characters in file path")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateOutputPath(t *testing.T) {
	tests := []struct {
		name        string
		zipPath     string
		outputsDir  string
		expectError bool
	}{
		{
			name:        "valid output path",
			zipPath:     "outputs/test.zip",
			outputsDir:  "outputs",
			expectError: false,
		},
		{
			name:        "valid nested output path",
			zipPath:     "outputs/subfolder/test.zip",
			outputsDir:  "outputs",
			expectError: false,
		},
		{
			name:        "directory traversal attempt",
			zipPath:     "outputs/../etc/passwd",
			outputsDir:  "outputs",
			expectError: true,
		},
		{
			name:        "directory traversal with multiple levels",
			zipPath:     "outputs/../../etc/passwd",
			outputsDir:  "outputs",
			expectError: true,
		},
		{
			name:        "absolute path outside outputs",
			zipPath:     "/etc/passwd",
			outputsDir:  "outputs",
			expectError: true,
		},
		{
			name:        "relative path outside outputs",
			zipPath:     "../outputs/test.zip",
			outputsDir:  "outputs",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOutputPath(tt.zipPath, tt.outputsDir)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid zip path")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSetupTempDirectory(t *testing.T) {
	tests := []struct {
		name        string
		tempDir     string
		expectError bool
	}{
		{
			name:        "create new directory",
			tempDir:     filepath.Join(os.TempDir(), "setup_test_new"),
			expectError: false,
		},
		{
			name:        "create nested directory",
			tempDir:     filepath.Join(os.TempDir(), "setup_test_parent", "child"),
			expectError: false,
		},
		{
			name:        "create directory with special chars in name",
			tempDir:     filepath.Join(os.TempDir(), "setup_test_special-chars_123"),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before test
			os.RemoveAll(tt.tempDir)
			defer os.RemoveAll(tt.tempDir)

			err := SetupTempDirectory(tt.tempDir)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.DirExists(t, tt.tempDir)
			}
		})
	}
}

func TestCleanupTempDirectory(t *testing.T) {
	testDir := filepath.Join(os.TempDir(), "cleanup_test_directory")
	defer os.RemoveAll(testDir)

	require.NoError(t, os.MkdirAll(testDir, 0750))

	testFile := filepath.Join(testDir, "test_file.txt")
	require.NoError(t, os.WriteFile(testFile, []byte("test"), 0644))

	assert.DirExists(t, testDir)
	assert.FileExists(t, testFile)

	CleanupTempDirectory(testDir)

	assert.NoDirExists(t, testDir)
}

func TestCleanupTempDirectoryNonExistent(t *testing.T) {
	nonExistentDir := filepath.Join(os.TempDir(), "non_existent_directory_cleanup_test")

	assert.NotPanics(t, func() {
		CleanupTempDirectory(nonExistentDir)
	})
}
