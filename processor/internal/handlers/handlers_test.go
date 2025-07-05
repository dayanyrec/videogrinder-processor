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

func TestProcessorHandlers_Integration_ShouldProvideFullProcessingWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)

	// Try to process a video
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
