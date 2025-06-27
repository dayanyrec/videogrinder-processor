package main

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
			err := validateProcessingInputs(tt.videoPath, tt.timestamp)
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
			err := validatePathSafety(tt.paths...)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid characters in file path")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsValidVideoFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{
			name:     "mp4 file",
			filename: "test.mp4",
			expected: true,
		},
		{
			name:     "avi file",
			filename: "test.avi",
			expected: true,
		},
		{
			name:     "mov file",
			filename: "test.mov",
			expected: true,
		},
		{
			name:     "mkv file",
			filename: "test.mkv",
			expected: true,
		},
		{
			name:     "wmv file",
			filename: "test.wmv",
			expected: true,
		},
		{
			name:     "flv file",
			filename: "test.flv",
			expected: true,
		},
		{
			name:     "webm file",
			filename: "test.webm",
			expected: true,
		},
		{
			name:     "uppercase mp4",
			filename: "test.MP4",
			expected: true,
		},
		{
			name:     "mixed case avi",
			filename: "test.AVI",
			expected: true,
		},
		{
			name:     "txt file",
			filename: "test.txt",
			expected: false,
		},
		{
			name:     "jpg file",
			filename: "test.jpg",
			expected: false,
		},
		{
			name:     "png file",
			filename: "test.png",
			expected: false,
		},
		{
			name:     "pdf file",
			filename: "test.pdf",
			expected: false,
		},
		{
			name:     "no extension",
			filename: "test",
			expected: false,
		},
		{
			name:     "empty filename",
			filename: "",
			expected: false,
		},
		{
			name:     "only extension",
			filename: ".mp4",
			expected: true,
		},
		{
			name:     "path with video file",
			filename: "/path/to/video.mp4",
			expected: true,
		},
		{
			name:     "windows path with video file",
			filename: "C:\\path\\to\\video.avi",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidVideoFile(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateOutputPath(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_outputs")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	oldDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(oldDir)

	testDir := filepath.Join(tempDir, "test_project")
	err = os.MkdirAll(filepath.Join(testDir, "outputs"), 0750)
	require.NoError(t, err)

	err = os.Chdir(testDir)
	require.NoError(t, err)

	tests := []struct {
		name        string
		zipPath     string
		expectError bool
	}{
		{
			name:        "valid output path",
			zipPath:     "outputs/test.zip",
			expectError: false,
		},
		{
			name:        "valid nested output path",
			zipPath:     "outputs/subfolder/test.zip",
			expectError: false,
		},
		{
			name:        "directory traversal attempt",
			zipPath:     "outputs/../etc/passwd",
			expectError: true,
		},
		{
			name:        "directory traversal with multiple levels",
			zipPath:     "outputs/../../root/.ssh/id_rsa",
			expectError: true,
		},
		{
			name:        "absolute path outside outputs",
			zipPath:     "/etc/passwd",
			expectError: true,
		},
		{
			name:        "relative path outside outputs",
			zipPath:     "../secrets.txt",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOutputPath(tt.zipPath)
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
	tempBase, err := os.MkdirTemp("", "test_setup")
	require.NoError(t, err)
	defer os.RemoveAll(tempBase)

	tests := []struct {
		name        string
		tempDir     string
		expectError bool
	}{
		{
			name:        "create new directory",
			tempDir:     filepath.Join(tempBase, "new_temp"),
			expectError: false,
		},
		{
			name:        "create nested directory",
			tempDir:     filepath.Join(tempBase, "nested", "temp", "dir"),
			expectError: false,
		},
		{
			name:        "create directory with special chars in name",
			tempDir:     filepath.Join(tempBase, "temp_dir_123"),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setupTempDirectory(tt.tempDir)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.DirExists(t, tt.tempDir)

				info, err := os.Stat(tt.tempDir)
				assert.NoError(t, err)
				assert.Equal(t, os.FileMode(0750), info.Mode().Perm())
			}
		})
	}
}

func TestCleanupTempDirectory(t *testing.T) {
	tempBase, err := os.MkdirTemp("", "test_cleanup")
	require.NoError(t, err)
	defer os.RemoveAll(tempBase)

	testDir := filepath.Join(tempBase, "to_cleanup")
	err = os.MkdirAll(testDir, 0750)
	require.NoError(t, err)

	testFile := filepath.Join(testDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	assert.DirExists(t, testDir)
	assert.FileExists(t, testFile)

	cleanupTempDirectory(testDir)

	assert.NoDirExists(t, testDir)
	assert.NoFileExists(t, testFile)
}

func TestCleanupTempDirectoryNonExistent(t *testing.T) {
	nonExistentDir := "/tmp/definitely_does_not_exist_12345"

	assert.NotPanics(t, func() {
		cleanupTempDirectory(nonExistentDir)
	})
}
