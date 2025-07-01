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

	"video-processor/internal/config"
	"video-processor/internal/models"
	"video-processor/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestHandlers() (handlers *WebHandlers, cleanup func()) {
	tempDir := filepath.Join(os.TempDir(), "webhandlers_test")
	uploadsDir := filepath.Join(tempDir, "uploads")
	outputsDir := filepath.Join(tempDir, "outputs")
	tempVideoDir := filepath.Join(tempDir, "temp")

	os.MkdirAll(uploadsDir, 0750)
	os.MkdirAll(outputsDir, 0750)
	os.MkdirAll(tempVideoDir, 0750)

	cfg := &config.Config{
		UploadsDir: uploadsDir,
		OutputsDir: outputsDir,
		TempDir:    tempVideoDir,
	}

	videoService := services.NewVideoService(cfg)
	handlers = NewWebHandlers(videoService, cfg)

	cleanup = func() {
		os.RemoveAll(tempDir)
	}

	return
}

func TestNewWebHandlers(t *testing.T) {
	cfg := &config.Config{
		UploadsDir: "uploads",
		OutputsDir: "outputs",
		TempDir:    "temp",
	}
	videoService := services.NewVideoService(cfg)

	handlers := NewWebHandlers(videoService, cfg)

	assert.NotNil(t, handlers)
	assert.Equal(t, videoService, handlers.videoService)
	assert.Equal(t, cfg, handlers.config)
}

func TestWebHandlers_HandleVideoUpload_InvalidFile(t *testing.T) {
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

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req

	handlers.HandleVideoUpload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ProcessingResult
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Message, "Formato de arquivo não suportado")
}

func TestWebHandlers_HandleVideoUpload_NoFile(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("POST", "/upload", http.NoBody)
	c.Request = req

	handlers.HandleVideoUpload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ProcessingResult
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Message, "Erro ao receber arquivo")
}

func TestWebHandlers_HandleVideoUpload_ValidFile(t *testing.T) {
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

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req

	handlers.HandleVideoUpload(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.ProcessingResult
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotEmpty(t, response.Message)
}

func TestWebHandlers_HandleDownload_FileNotFound(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{gin.Param{Key: "filename", Value: "nonexistent.zip"}}

	handlers.HandleDownload(c)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Arquivo não encontrado", response["error"])
}

func TestWebHandlers_HandleDownload_FileExists(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	testFile := filepath.Join(handlers.config.OutputsDir, "test.zip")
	err := os.WriteFile(testFile, []byte("test zip content"), 0644)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/download/test.zip", http.NoBody)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "filename", Value: "test.zip"}}

	handlers.HandleDownload(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/zip", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment; filename=test.zip")
	assert.Equal(t, "test zip content", w.Body.String())
}

func TestWebHandlers_HandleStatus_NoFiles(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handlers.HandleStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, float64(0), response["total"])
	assert.Empty(t, response["files"])
}

func TestWebHandlers_HandleStatus_WithFiles(t *testing.T) {
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

	handlers.HandleStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, float64(2), response["total"])

	files := response["files"].([]interface{})
	assert.Len(t, files, 2)
}

func TestWebHandlers_HandleHome(t *testing.T) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handlers.HandleHome(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/html", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "FIAP X - Processador de Vídeos")
	assert.Contains(t, w.Body.String(), "<!DOCTYPE html>")
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
			filename: "path/to/test.mp4",
			expected: true,
		},
		{
			name:     "windows path with video file",
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

func TestGetHTMLForm(t *testing.T) {
	html := GetHTMLForm()

	assert.NotEmpty(t, html)
	assert.Contains(t, html, "<!DOCTYPE html>")
	assert.Contains(t, html, "FIAP X - Processador de Vídeos")
	assert.Contains(t, html, "Upload do Vídeo")
	assert.Contains(t, html, "Arquivos Disponíveis")
	assert.Contains(t, html, "fetch('/upload'")
	assert.Contains(t, html, "fetch('/api/status'")
}

func TestWebHandlers_Integration_FullWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handlers.HandleHome(c)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	handlers.HandleStatus(c)
	assert.Equal(t, http.StatusOK, w.Code)
	var statusResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &statusResponse)
	assert.Equal(t, float64(0), statusResponse["total"])

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "filename", Value: "test.zip"}}
	handlers.HandleDownload(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func BenchmarkWebHandlers_HandleStatus(b *testing.B) {
	handlers, cleanup := setupTestHandlers()
	defer cleanup()

	gin.SetMode(gin.TestMode)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handlers.HandleStatus(c)
	}
}
