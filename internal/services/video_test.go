package services

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"video-processor/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVideoService_New(t *testing.T) {
	cfg := &config.Config{
		UploadsDir: "uploads",
		OutputsDir: "outputs",
		TempDir:    "temp",
	}

	service := NewVideoService(cfg)

	assert.NotNil(t, service)
	assert.Equal(t, cfg, service.config)
}

func TestVideoService_ProcessVideo_ValidationErrors(t *testing.T) {
	cfg := &config.Config{
		UploadsDir: "uploads",
		OutputsDir: "outputs",
		TempDir:    "temp",
	}
	service := NewVideoService(cfg)

	tests := []struct {
		name      string
		videoPath string
		timestamp string
		expectErr string
	}{
		{
			name:      "video path with directory traversal",
			videoPath: "uploads/../etc/passwd",
			timestamp: "20240101_120000",
			expectErr: "invalid path parameters",
		},
		{
			name:      "timestamp with directory traversal",
			videoPath: "uploads/test.mp4",
			timestamp: "../../../etc",
			expectErr: "invalid path parameters",
		},
		{
			name:      "both paths with directory traversal",
			videoPath: "../uploads/test.mp4",
			timestamp: "../timestamp",
			expectErr: "invalid path parameters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ProcessVideo(tt.videoPath, tt.timestamp)

			assert.False(t, result.Success)
			assert.Contains(t, result.Message, tt.expectErr)
			assert.Empty(t, result.ZipPath)
			assert.Zero(t, result.FrameCount)
			assert.Empty(t, result.Images)
		})
	}
}

func TestVideoService_ProcessVideo_TempDirCreationError(t *testing.T) {
	// Create a temporary file to use as the "directory" path
	// This should cause os.MkdirAll to fail since it can't create a directory with the same name as an existing file
	tempFile, err := os.CreateTemp("", "video_test_file")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	cfg := &config.Config{
		UploadsDir: "uploads",
		OutputsDir: "outputs",
		TempDir:    tempFile.Name(), // This should fail because it's a file, not a directory
	}
	service := NewVideoService(cfg)

	result := service.ProcessVideo("uploads/test.mp4", "20240101_120000")

	// Should fail due to temp directory creation error
	assert.False(t, result.Success)
	assert.NotEmpty(t, result.Message)
	// The error should contain information about directory creation failure
	assert.True(t,
		strings.Contains(result.Message, "erro ao criar diretório temporário") ||
			strings.Contains(result.Message, "not a directory") ||
			strings.Contains(result.Message, "file exists"),
		"Expected temp directory creation error, got: %s", result.Message)
}

func TestVideoService_extractFrames_PathSafetyValidation(t *testing.T) {
	cfg := &config.Config{
		TempDir: "temp",
	}
	service := NewVideoService(cfg)

	tests := []struct {
		name      string
		videoPath string
		tempDir   string
		expectErr string
	}{
		{
			name:      "video path with semicolon",
			videoPath: "uploads/test;rm -rf /.mp4",
			tempDir:   "temp/test",
			expectErr: "invalid characters in file path",
		},
		{
			name:      "video path with pipe",
			videoPath: "uploads/test|cat /etc/passwd.mp4",
			tempDir:   "temp/test",
			expectErr: "invalid characters in file path",
		},
		{
			name:      "video path with dollar sign",
			videoPath: "uploads/test$USER.mp4",
			tempDir:   "temp/test",
			expectErr: "invalid characters in file path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.extractFrames(tt.videoPath, tt.tempDir)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectErr)
		})
	}
}

func TestVideoService_createFramesZip_OutputPathValidation(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "video_service_test_zip")
	defer os.RemoveAll(tempDir)

	require.NoError(t, os.MkdirAll(tempDir, 0750))

	cfg := &config.Config{
		OutputsDir: tempDir,
	}
	service := NewVideoService(cfg)

	testFile := filepath.Join(tempDir, "test_frame.png")
	require.NoError(t, os.WriteFile(testFile, []byte("test"), 0644))

	frames := []string{testFile}

	tests := []struct {
		name        string
		timestamp   string
		expectError bool
	}{
		{
			name:        "valid timestamp",
			timestamp:   "20240101_120000",
			expectError: false,
		},
		{
			name:        "timestamp with directory traversal",
			timestamp:   "../../../etc",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zipPath, err := service.createFramesZip(frames, tt.timestamp)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, zipPath)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, zipPath)
				assert.FileExists(t, zipPath)
				defer os.Remove(zipPath)
			}
		})
	}
}

func TestVideoService_createZipFile(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "video_service_test_createzip")
	defer os.RemoveAll(tempDir)

	require.NoError(t, os.MkdirAll(tempDir, 0750))

	cfg := &config.Config{}
	service := NewVideoService(cfg)

	// Create test files
	testFiles := []string{
		filepath.Join(tempDir, "frame_001.png"),
		filepath.Join(tempDir, "frame_002.png"),
	}

	for _, file := range testFiles {
		require.NoError(t, os.WriteFile(file, []byte("test frame data"), 0644))
	}

	zipPath := filepath.Join(tempDir, "test.zip")

	err := service.createZipFile(testFiles, zipPath)

	assert.NoError(t, err)
	assert.FileExists(t, zipPath)

	// Verify zip file is not empty
	info, err := os.Stat(zipPath)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))
}

func TestVideoService_addFileToZip_InvalidPath(t *testing.T) {
	cfg := &config.Config{}
	service := NewVideoService(cfg)

	// This would be tested through createZipFile, but we can test the underlying logic
	// by testing with a non-existent file
	tempDir := filepath.Join(os.TempDir(), "video_service_test_addfile")
	defer os.RemoveAll(tempDir)

	require.NoError(t, os.MkdirAll(tempDir, 0750))

	zipPath := filepath.Join(tempDir, "test.zip")
	nonExistentFile := filepath.Join(tempDir, "non_existent.png")

	files := []string{nonExistentFile}
	err := service.createZipFile(files, zipPath)

	assert.Error(t, err)
}

func TestVideoService_ProcessVideo_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tempDir := filepath.Join(os.TempDir(), "video_service_integration_test")
	defer os.RemoveAll(tempDir)

	// Create test directories
	uploadsDir := filepath.Join(tempDir, "uploads")
	outputsDir := filepath.Join(tempDir, "outputs")
	tempVideoDir := filepath.Join(tempDir, "temp")

	require.NoError(t, os.MkdirAll(uploadsDir, 0750))
	require.NoError(t, os.MkdirAll(outputsDir, 0750))
	require.NoError(t, os.MkdirAll(tempVideoDir, 0750))

	cfg := &config.Config{
		UploadsDir: uploadsDir,
		OutputsDir: outputsDir,
		TempDir:    tempVideoDir,
	}

	service := NewVideoService(cfg)

	// Note: This test would require an actual video file and ffmpeg to be fully functional
	// For now, we test the structure and validation
	videoPath := filepath.Join(uploadsDir, "test.mp4")
	timestamp := "20240101_120000"

	// Test with non-existent video file (expected to fail at ffmpeg stage)
	result := service.ProcessVideo(videoPath, timestamp)

	// Should fail because file doesn't exist, but validation should pass
	assert.False(t, result.Success)
	assert.NotEmpty(t, result.Message)
}

// Benchmark for ProcessVideo performance.
func BenchmarkVideoService_ProcessVideo(b *testing.B) {
	cfg := &config.Config{
		UploadsDir: "uploads",
		OutputsDir: "outputs",
		TempDir:    "temp",
	}
	service := NewVideoService(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// This will fail quickly due to validation, but measures overhead
		_ = service.ProcessVideo("../invalid/path", "20240101_120000")
	}
}
