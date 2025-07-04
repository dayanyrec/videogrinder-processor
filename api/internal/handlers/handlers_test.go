package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"video-processor/api/internal/config"
	baseConfig "video-processor/internal/config"
	"video-processor/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockProcessorClient struct {
	healthCheckFunc  func() error
	processVideoFunc func(string, io.Reader) (*models.ProcessingResult, error)
}

func (m *MockProcessorClient) HealthCheck() error {
	if m.healthCheckFunc != nil {
		return m.healthCheckFunc()
	}
	return nil
}

func (m *MockProcessorClient) ProcessVideo(filename string, fileReader io.Reader) (*models.ProcessingResult, error) {
	if m.processVideoFunc != nil {
		return m.processVideoFunc(filename, fileReader)
	}
	return &models.ProcessingResult{
		Success:    true,
		Message:    "Processamento concluído! 5 frames extraídos.",
		ZipPath:    "frames_test.zip",
		FrameCount: 5,
		Images:     []string{"frame_001.png", "frame_002.png"},
	}, nil
}

func setupTestHandlers() (handlers *APIHandlers, cleanup func()) {
	tempDir := filepath.Join(os.TempDir(), "api_test")
	uploadsDir := filepath.Join(tempDir, "uploads")
	outputsDir := filepath.Join(tempDir, "outputs")
	tempAPIDir := filepath.Join(tempDir, "temp")

	require.NoError(nil, os.MkdirAll(uploadsDir, 0750))
	require.NoError(nil, os.MkdirAll(outputsDir, 0750))
	require.NoError(nil, os.MkdirAll(tempAPIDir, 0750))

	cfg := &config.APIConfig{
		Port:         "8081",
		ProcessorURL: "http://localhost:8082",
		DirectoryConfig: &baseConfig.DirectoryConfig{
			UploadsDir: uploadsDir,
			OutputsDir: outputsDir,
			TempDir:    tempAPIDir,
		},
	}

	handlers = NewAPIHandlers(cfg)

	cleanup = func() {
		os.RemoveAll(tempDir)
	}

	return
}

func TestNewAPIHandlers_ShouldInitializeHandlersWithCorrectDependencies(t *testing.T) {
	cfg := &config.APIConfig{
		Port:         "8081",
		ProcessorURL: "http://localhost:8082",
		DirectoryConfig: &baseConfig.DirectoryConfig{
			UploadsDir: "uploads",
			OutputsDir: "outputs",
			TempDir:    "temp",
		},
	}

	handlers := NewAPIHandlers(cfg)

	assert.NotNil(t, handlers)
	assert.NotNil(t, handlers.processorClient)
	assert.Equal(t, cfg, handlers.config)
}

func TestCreateVideo_ShouldReturnBadRequestWhenUploadingInvalidFileExtension(t *testing.T) {
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

	req := httptest.NewRequest("POST", "/api/v1/videos", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req

	handlers.CreateVideo(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ProcessingResult
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Message, "Formato de arquivo não suportado")
}

func TestCreateVideo_ShouldReturnBadRequestWhenNoFileIsProvided(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("POST", "/api/v1/videos", http.NoBody)
	c.Request = req

	handlers.CreateVideo(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ProcessingResult
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Message, "Erro ao receber arquivo")
}

func TestCreateVideo_ShouldReturnServiceUnavailableWhenProcessorIsDown(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	mockClient := &MockProcessorClient{
		healthCheckFunc: func() error {
			return assert.AnError
		},
	}
	handlers.processorClient = mockClient

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("video", "test.mp4")
	require.NoError(t, err)
	part.Write([]byte("fake video content"))
	writer.Close()

	req := httptest.NewRequest("POST", "/api/v1/videos", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req

	handlers.CreateVideo(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response models.ProcessingResult
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Message, "Serviço de processamento indisponível")
}

func TestCreateVideo_ShouldReturnCreatedWhenVideoProcessingSucceeds(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	mockClient := &MockProcessorClient{
		healthCheckFunc: func() error {
			return nil
		},
		processVideoFunc: func(filename string, fileReader io.Reader) (*models.ProcessingResult, error) {
			return &models.ProcessingResult{
				Success:    true,
				Message:    "Processamento concluído! 5 frames extraídos.",
				ZipPath:    "frames_test.zip",
				FrameCount: 5,
				Images:     []string{"frame_001.png", "frame_002.png"},
			}, nil
		},
	}
	handlers.processorClient = mockClient

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("video", "test.mp4")
	require.NoError(t, err)
	part.Write([]byte("fake video content"))
	writer.Close()

	req := httptest.NewRequest("POST", "/api/v1/videos", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req

	handlers.CreateVideo(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.ProcessingResult
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "frames_test.zip", response.ZipPath)
	assert.Equal(t, 5, response.FrameCount)
}

func TestCreateVideo_ShouldReturnUnprocessableEntityWhenProcessingFails(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	mockClient := &MockProcessorClient{
		healthCheckFunc: func() error {
			return nil
		},
		processVideoFunc: func(filename string, fileReader io.Reader) (*models.ProcessingResult, error) {
			return &models.ProcessingResult{
				Success: false,
				Message: "Erro ao processar vídeo: formato inválido",
			}, nil
		},
	}
	handlers.processorClient = mockClient

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("video", "test.mp4")
	require.NoError(t, err)
	part.Write([]byte("fake video content"))
	writer.Close()

	req := httptest.NewRequest("POST", "/api/v1/videos", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req

	handlers.CreateVideo(c)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var response models.ProcessingResult
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Message, "formato inválido")
}

func TestGetVideoDownload_ShouldReturnNotFoundWhenRequestedFileDoesNotExist(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{gin.Param{Key: "filename", Value: "nonexistent.zip"}}

	handlers.GetVideoDownload(c)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Arquivo não encontrado", response["error"])
}

func TestGetVideoDownload_ShouldReturnFileWithCorrectHeadersWhenFileExists(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	testFile := filepath.Join(handlers.config.OutputsDir, "test.zip")
	err := os.WriteFile(testFile, []byte("test zip content"), 0644)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/api/v1/videos/test.zip/download", http.NoBody)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "filename", Value: "test.zip"}}

	handlers.GetVideoDownload(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/zip", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment; filename=test.zip")
	assert.Equal(t, "test zip content", w.Body.String())
}

func TestGetVideos_ShouldReturnEmptyListWhenNoProcessedVideosExist(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handlers.GetVideos(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, float64(0), response["total"])
	assert.Empty(t, response["videos"])
}

func TestGetVideos_ShouldReturnListWithCorrectCountWhenMultipleVideosExist(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	testFile1 := filepath.Join(handlers.config.OutputsDir, "test1.zip")
	testFile2 := filepath.Join(handlers.config.OutputsDir, "test2.zip")
	err := os.WriteFile(testFile1, []byte("test content 1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(testFile2, []byte("test content 2"), 0644)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handlers.GetVideos(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, float64(2), response["total"])

	videos := response["videos"].([]interface{})
	assert.Len(t, videos, 2)
}

func TestDeleteVideo_ShouldReturnNotFoundWhenAttemptingToDeleteNonExistentFile(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{gin.Param{Key: "filename", Value: "nonexistent.zip"}}

	handlers.DeleteVideo(c)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Arquivo não encontrado", response["error"])
}

func TestDeleteVideo_ShouldReturnNoContentAndRemoveFileWhenDeletingExistingFile(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	testFile := filepath.Join(handlers.config.OutputsDir, "test.zip")
	err := os.WriteFile(testFile, []byte("test zip content"), 0644)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{gin.Param{Key: "filename", Value: "test.zip"}}

	handlers.DeleteVideo(c)

	assert.Equal(t, http.StatusNoContent, w.Code)

	_, err = os.Stat(testFile)
	assert.True(t, os.IsNotExist(err))
}

func TestIsValidVideoFile_ShouldValidateVideoFileExtensionsCorrectly(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{
			name:     "should accept mp4 file",
			filename: "test.mp4",
			expected: true,
		},
		{
			name:     "should accept avi file",
			filename: "test.avi",
			expected: true,
		},
		{
			name:     "should accept mov file",
			filename: "test.mov",
			expected: true,
		},
		{
			name:     "should accept mkv file",
			filename: "test.mkv",
			expected: true,
		},
		{
			name:     "should accept wmv file",
			filename: "test.wmv",
			expected: true,
		},
		{
			name:     "should accept flv file",
			filename: "test.flv",
			expected: true,
		},
		{
			name:     "should accept webm file",
			filename: "test.webm",
			expected: true,
		},
		{
			name:     "should accept uppercase mp4 extension",
			filename: "test.MP4",
			expected: true,
		},
		{
			name:     "should accept mixed case avi extension",
			filename: "test.AVI",
			expected: true,
		},
		{
			name:     "should reject txt file",
			filename: "test.txt",
			expected: false,
		},
		{
			name:     "should reject jpg file",
			filename: "test.jpg",
			expected: false,
		},
		{
			name:     "should reject png file",
			filename: "test.png",
			expected: false,
		},
		{
			name:     "should reject pdf file",
			filename: "test.pdf",
			expected: false,
		},
		{
			name:     "should reject file without extension",
			filename: "test",
			expected: false,
		},
		{
			name:     "should reject empty filename",
			filename: "",
			expected: false,
		},
		{
			name:     "should accept extension-only mp4 filename",
			filename: ".mp4",
			expected: true,
		},
		{
			name:     "should accept video file with path",
			filename: "path/to/test.mp4",
			expected: true,
		},
		{
			name:     "should accept video file with windows path",
			filename: "C:\\path\\to\\test.mp4",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidVideoFile(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAPIHandlers_Integration_ShouldProvideFullWorkflowBehavior(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handlers.GetVideos(c)
	assert.Equal(t, http.StatusOK, w.Code)
	var statusResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &statusResponse)
	assert.Equal(t, float64(0), statusResponse["total"])

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "filename", Value: "test.zip"}}
	handlers.GetVideoDownload(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func BenchmarkGetVideos_ShouldPerformEfficientlyUnderLoad(b *testing.B) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handlers.GetVideos(c)
	}
}

func TestGetAPIHealth_ShouldReturnHealthyStatusWhenAllServicesAreOperational(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	mockClient := &MockProcessorClient{
		healthCheckFunc: func() error {
			return nil
		},
	}
	handlers.processorClient = mockClient

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/health", http.NoBody)
	c.Request = req

	handlers.GetAPIHealth(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "videogrinder-api", response["service"])
	assert.NotNil(t, response["timestamp"])
	assert.Equal(t, "1.0.0", response["version"])

	checks := response["checks"].(map[string]interface{})
	assert.NotNil(t, checks["directories"])
	assert.NotNil(t, checks["processor"])

	directories := checks["directories"].(map[string]interface{})
	assert.Equal(t, "healthy", directories["status"])

	processor := checks["processor"].(map[string]interface{})
	assert.Equal(t, "healthy", processor["status"])
}

func TestGetAPIHealth_ShouldReturnUnhealthyStatusWhenProcessorIsDown(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	mockClient := &MockProcessorClient{
		healthCheckFunc: func() error {
			return assert.AnError
		},
	}
	handlers.processorClient = mockClient

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/health", http.NoBody)
	c.Request = req

	handlers.GetAPIHealth(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "unhealthy", response["status"])
	assert.Equal(t, "videogrinder-api", response["service"])

	checks := response["checks"].(map[string]interface{})
	processor := checks["processor"].(map[string]interface{})
	assert.Equal(t, "unhealthy", processor["status"])
	assert.NotNil(t, processor["error"])
	assert.NotNil(t, processor["latency_ms"])
}

func TestGetAPIHealth_ShouldReturnUnhealthyStatusWhenDirectoriesAreMissing(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	tempDir := filepath.Join(os.TempDir(), "non_existent_test_dir")
	handlers.config.UploadsDir = tempDir

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/health", http.NoBody)
	c.Request = req

	handlers.GetAPIHealth(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "unhealthy", response["status"])

	checks := response["checks"].(map[string]interface{})
	directories := checks["directories"].(map[string]interface{})
	assert.Equal(t, "unhealthy", directories["status"])

	details := directories["details"].(map[string]interface{})
	uploads := details["non_existent_test_dir"].(map[string]interface{})
	assert.Equal(t, "missing", uploads["status"])
	assert.Contains(t, uploads["error"], "does not exist")
}

func TestGetAPIHealth_ShouldIncludeLatencyInformation(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	mockClient := &MockProcessorClient{
		healthCheckFunc: func() error {
			return nil
		},
	}
	handlers.processorClient = mockClient

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/health", http.NoBody)
	c.Request = req

	handlers.GetAPIHealth(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	checks := response["checks"].(map[string]interface{})
	processor := checks["processor"].(map[string]interface{})

	assert.NotNil(t, processor["latency_ms"])
	assert.NotNil(t, processor["last_check"])
	assert.NotNil(t, processor["url"])

	// Latency should be a number >= 0
	latency := processor["latency_ms"].(float64)
	assert.GreaterOrEqual(t, latency, float64(0))
}

func TestGetAPIHealth_ShouldVerifyAllRequiredDirectories(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	mockClient := &MockProcessorClient{
		healthCheckFunc: func() error {
			return nil
		},
	}
	handlers.processorClient = mockClient

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/health", http.NoBody)
	c.Request = req

	handlers.GetAPIHealth(c)

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
