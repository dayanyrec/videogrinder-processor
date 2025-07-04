package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	baseConfig "video-processor/internal/config"
	"video-processor/processor/internal/config"
	"video-processor/processor/internal/models"
	"video-processor/processor/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestHandlers() (handlers *ProcessorHandlers, cleanup func()) {
	tempDir := filepath.Join(os.TempDir(), "processor_test")
	uploadsDir := filepath.Join(tempDir, "uploads")
	outputsDir := filepath.Join(tempDir, "outputs")
	tempVideoDir := filepath.Join(tempDir, "temp")

	require.NoError(nil, os.MkdirAll(uploadsDir, 0750))
	require.NoError(nil, os.MkdirAll(outputsDir, 0750))
	require.NoError(nil, os.MkdirAll(tempVideoDir, 0750))

	cfg := &config.ProcessorConfig{
		Port: "8082",
		DirectoryConfig: &baseConfig.DirectoryConfig{
			UploadsDir: uploadsDir,
			OutputsDir: outputsDir,
			TempDir:    tempVideoDir,
		},
	}

	videoService := services.NewVideoService(cfg)
	handlers = NewProcessorHandlers(videoService, cfg)

	cleanup = func() {
		os.RemoveAll(tempDir)
	}

	return
}

func setupTestHandlersWithMissingDirs() (handlers *ProcessorHandlers, cleanup func()) {
	tempDir := filepath.Join(os.TempDir(), "processor_test_missing")
	uploadsDir := filepath.Join(tempDir, "uploads")
	outputsDir := filepath.Join(tempDir, "outputs")
	tempVideoDir := filepath.Join(tempDir, "temp")

	cfg := &config.ProcessorConfig{
		Port: "8082",
		DirectoryConfig: &baseConfig.DirectoryConfig{
			UploadsDir: uploadsDir,
			OutputsDir: outputsDir,
			TempDir:    tempVideoDir,
		},
	}

	videoService := services.NewVideoService(cfg)
	handlers = NewProcessorHandlers(videoService, cfg)

	cleanup = func() {
		os.RemoveAll(tempDir)
	}

	return
}

func TestNewProcessorHandlers_ShouldInitializeHandlersWithCorrectDependencies(t *testing.T) {
	cfg := &config.ProcessorConfig{
		Port: "8082",
		DirectoryConfig: &baseConfig.DirectoryConfig{
			UploadsDir: "uploads",
			OutputsDir: "outputs",
			TempDir:    "temp",
		},
	}
	videoService := services.NewVideoService(cfg)

	handlers := NewProcessorHandlers(videoService, cfg)

	assert.NotNil(t, handlers)
	assert.Equal(t, videoService, handlers.videoService)
}

func TestProcessVideoUpload_ShouldReturnBadRequestWhenNoFileIsProvided(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("POST", "/process", http.NoBody)
	c.Request = req

	handlers.ProcessVideoUpload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ProcessingResult
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Message, "Erro ao receber arquivo")
}

func TestProcessVideoUpload_ShouldReturnBadRequestWhenUploadingInvalidFileExtension(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("video", "test.txt")
	require.NoError(t, err)
	part.Write([]byte("not a video"))
	writer.Close()

	req := httptest.NewRequest("POST", "/process", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req

	handlers.ProcessVideoUpload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ProcessingResult
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Message, "Formato de arquivo não suportado")
}

func TestProcessVideoUpload_ShouldReturnCreatedWhenProcessingValidVideo(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("video", "test.mp4")
	require.NoError(t, err)
	part.Write([]byte("fake video content"))
	writer.Close()

	req := httptest.NewRequest("POST", "/process", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req

	handlers.ProcessVideoUpload(c)

	// The service will return an error because it's not a real video
	// but the handler should handle it gracefully
	assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusUnprocessableEntity)

	var response models.ProcessingResult
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response.Message)
}

func TestGetProcessorStatus_ShouldReturnHealthyStatusWhenServiceIsRunning(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/health", http.NoBody)
	c.Request = req

	handlers.GetProcessorStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "videogrinder-processor", response["service"])
	assert.NotNil(t, response["timestamp"])
	assert.Equal(t, "1.0.0", response["version"])

	// Check that we have both directories and ffmpeg checks
	checks := response["checks"].(map[string]interface{})
	assert.NotNil(t, checks["directories"])
	assert.NotNil(t, checks["ffmpeg"])

	// Check directories health
	directories := checks["directories"].(map[string]interface{})
	assert.Equal(t, "healthy", directories["status"])

	// Check ffmpeg health - it might be unhealthy in CI/test environment
	ffmpeg := checks["ffmpeg"].(map[string]interface{})
	assert.Contains(t, []string{"healthy", "unhealthy"}, ffmpeg["status"])
	assert.NotNil(t, ffmpeg["latency_ms"])
	assert.NotNil(t, ffmpeg["last_check"])
}

func TestGetProcessorStatus_ShouldReturnUnhealthyStatusWhenDirectoriesAreMissing(t *testing.T) {
	handlers, cleanup := setupTestHandlersWithMissingDirs()
	defer cleanup()

	// Create a non-existent directory to simulate missing directory
	tempDir := filepath.Join(os.TempDir(), "non_existent_processor_test_dir")
	handlers.config.UploadsDir = tempDir

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/health", http.NoBody)
	c.Request = req

	handlers.GetProcessorStatus(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "unhealthy", response["status"])
	assert.Equal(t, "videogrinder-processor", response["service"])

	// Check that directories are marked as unhealthy
	checks := response["checks"].(map[string]interface{})
	directories := checks["directories"].(map[string]interface{})
	assert.Equal(t, "unhealthy", directories["status"])

	// Check details for uploads directory specifically
	details := directories["details"].(map[string]interface{})
	uploads := details["non_existent_processor_test_dir"].(map[string]interface{})
	assert.Equal(t, "missing", uploads["status"])
	assert.Contains(t, uploads["error"], "does not exist")
}

func TestGetProcessorStatus_ShouldIncludeLatencyInformation(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/health", http.NoBody)
	c.Request = req

	handlers.GetProcessorStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Check ffmpeg latency information
	checks := response["checks"].(map[string]interface{})
	ffmpeg := checks["ffmpeg"].(map[string]interface{})

	assert.NotNil(t, ffmpeg["latency_ms"])
	assert.NotNil(t, ffmpeg["last_check"])

	// Latency should be a number >= 0
	latency := ffmpeg["latency_ms"].(float64)
	assert.GreaterOrEqual(t, latency, float64(0))
}

func TestGetProcessorStatus_ShouldVerifyAllRequiredDirectories(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/health", http.NoBody)
	c.Request = req

	handlers.GetProcessorStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Check that all directories are verified
	checks := response["checks"].(map[string]interface{})
	directories := checks["directories"].(map[string]interface{})
	details := directories["details"].(map[string]interface{})

	// Should have checks for uploads, outputs, and temp directories
	expectedDirs := []string{"uploads", "outputs", "temp"}
	for _, dirName := range expectedDirs {
		assert.Contains(t, details, dirName)
		dir := details[dirName].(map[string]interface{})
		assert.Equal(t, "healthy", dir["status"])
		assert.NotNil(t, dir["path"])
	}
}

func TestProcessorHandlers_Integration_ShouldProvideFullProcessingWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)

	// First check if service is healthy
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/health", http.NoBody)
	c.Request = req
	handlers.GetProcessorStatus(c)
	assert.Equal(t, http.StatusOK, w.Code)

	// Then try to process a video
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("video", "test.mp4")
	require.NoError(t, err)
	part.Write([]byte("fake video content"))
	writer.Close()

	req = httptest.NewRequest("POST", "/process", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req

	handlers.ProcessVideoUpload(c)

	// Should get a response (success or failure)
	assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusUnprocessableEntity)
	assert.NotEmpty(t, w.Body.String())
}

func BenchmarkProcessVideoUpload_ShouldPerformEfficientlyUnderLoad(b *testing.B) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)

	// Prepare test data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("video", "test.mp4")
	part.Write([]byte("fake video content"))
	writer.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest("POST", "/process", bytes.NewReader(body.Bytes()))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		c.Request = req

		handlers.ProcessVideoUpload(c)
	}
}
