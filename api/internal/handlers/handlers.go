package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"video-processor/api/internal/clients"
	"video-processor/api/internal/config"
	"video-processor/api/internal/models"

	"github.com/gin-gonic/gin"
)

const (
	StatusHealthy   = "healthy"
	StatusUnhealthy = "unhealthy"
)

type APIHandlers struct {
	processorClient clients.ProcessorClientInterface
	config          *config.APIConfig
}

func NewAPIHandlers(cfg *config.APIConfig) *APIHandlers {
	return &APIHandlers{
		processorClient: clients.NewProcessorClient(cfg.ProcessorURL),
		config:          cfg,
	}
}

func (ah *APIHandlers) GetAPIHealth(c *gin.Context) {
	health := gin.H{
		"status":    StatusHealthy,
		"service":   "videogrinder-api",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
		"checks": gin.H{
			"directories": ah.checkDirectories(),
			"processor":   ah.checkProcessorConnectivity(),
			"s3":          ah.checkS3Connectivity(),
		},
	}

	dirCheck := health["checks"].(gin.H)["directories"].(gin.H)
	procCheck := health["checks"].(gin.H)["processor"].(gin.H)
	s3Check := health["checks"].(gin.H)["s3"].(gin.H)

	s3Status := s3Check["status"].(string)
	s3Healthy := s3Status == StatusHealthy || s3Status == "disabled"

	if dirCheck["status"] != StatusHealthy || procCheck["status"] != StatusHealthy || !s3Healthy {
		health["status"] = StatusUnhealthy
		c.JSON(http.StatusServiceUnavailable, health)
		return
	}

	c.JSON(http.StatusOK, health)
}

func (ah *APIHandlers) checkDirectories() gin.H {
	directories := []string{
		ah.config.UploadsDir,
		ah.config.OutputsDir,
		ah.config.TempDir,
	}

	allHealthy := true
	details := make(gin.H)

	for _, dir := range directories {
		dirName := filepath.Base(dir)

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			details[dirName] = gin.H{
				"status": "missing",
				"path":   dir,
				"error":  "Directory does not exist",
			}
			allHealthy = false
			continue
		}

		testFile := filepath.Join(dir, ".health_check_test")
		if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
			details[dirName] = gin.H{
				"status": "not_writable",
				"path":   dir,
				"error":  "Directory is not writable: " + err.Error(),
			}
			allHealthy = false
		} else {
			if err := os.Remove(testFile); err != nil {
				log.Printf("Warning: Failed to remove test file %s: %v", testFile, err)
			}
			details[dirName] = gin.H{
				"status": StatusHealthy,
				"path":   dir,
			}
		}
	}

	return gin.H{
		"status":  map[bool]string{true: StatusHealthy, false: StatusUnhealthy}[allHealthy],
		"details": details,
	}
}

func (ah *APIHandlers) checkProcessorConnectivity() gin.H {
	start := time.Now()
	err := ah.processorClient.HealthCheck()
	latency := time.Since(start)

	if err != nil {
		return gin.H{
			"status":     StatusUnhealthy,
			"url":        ah.config.ProcessorURL,
			"error":      err.Error(),
			"latency_ms": latency.Milliseconds(),
			"last_check": time.Now().Unix(),
		}
	}

	return gin.H{
		"status":     StatusHealthy,
		"url":        ah.config.ProcessorURL,
		"latency_ms": latency.Milliseconds(),
		"last_check": time.Now().Unix(),
	}
}

func (ah *APIHandlers) checkS3Connectivity() gin.H {
	if !ah.config.IsS3Enabled() {
		return gin.H{
			"status":  "disabled",
			"message": "S3 integration is disabled",
		}
	}

	start := time.Now()
	err := ah.config.AWSConfig.CheckHealth()
	latency := time.Since(start)

	if err != nil {
		return gin.H{
			"status":     StatusUnhealthy,
			"error":      err.Error(),
			"latency_ms": latency.Milliseconds(),
			"last_check": time.Now().Unix(),
			"endpoint":   ah.config.AWSConfig.GetS3Endpoint(),
		}
	}

	return gin.H{
		"status":     StatusHealthy,
		"latency_ms": latency.Milliseconds(),
		"last_check": time.Now().Unix(),
		"endpoint":   ah.config.AWSConfig.GetS3Endpoint(),
	}
}

func (ah *APIHandlers) CreateVideo(c *gin.Context) {
	file, header, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ProcessingResult{
			Success: false,
			Message: "Erro ao receber arquivo: " + err.Error(),
		})
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Warning: Failed to close uploaded file: %v", err)
		}
	}()

	if !IsValidVideoFile(header.Filename) {
		c.JSON(http.StatusBadRequest, models.ProcessingResult{
			Success: false,
			Message: "Formato de arquivo não suportado. Use: mp4, avi, mov, mkv",
		})
		return
	}

	if err := ah.processorClient.HealthCheck(); err != nil {
		c.JSON(http.StatusServiceUnavailable, models.ProcessingResult{
			Success: false,
			Message: "Serviço de processamento indisponível: " + err.Error(),
		})
		return
	}

	if ah.config.IsS3Enabled() {
		ah.processVideoWithS3(c, file, header.Filename)
	} else {
		ah.processVideoDirectly(c, file, header.Filename)
	}
}

func (ah *APIHandlers) processVideoWithS3(c *gin.Context, file io.Reader, filename string) {
	timestamp := time.Now().Format("20060102_150405")
	s3Key := fmt.Sprintf("%s_%s", timestamp, filepath.Base(filename))

	if err := ah.config.S3Service.UploadFile(ah.config.S3Buckets.UploadsBucket, s3Key, file); err != nil {
		c.JSON(http.StatusInternalServerError, models.ProcessingResult{
			Success: false,
			Message: "Erro ao fazer upload do vídeo para S3: " + err.Error(),
		})
		return
	}

	log.Printf("Video uploaded to S3: s3://%s/%s", ah.config.S3Buckets.UploadsBucket, s3Key)

	result, err := ah.processorClient.ProcessVideoFromS3(s3Key)
	if err != nil {
		if cleanupErr := ah.config.S3Service.DeleteFile(ah.config.S3Buckets.UploadsBucket, s3Key); cleanupErr != nil {
			log.Printf("Warning: Failed to cleanup uploaded video from S3: %v", cleanupErr)
		}
		c.JSON(http.StatusInternalServerError, models.ProcessingResult{
			Success: false,
			Message: "Erro ao processar vídeo: " + err.Error(),
		})
		return
	}

	if result.Success {
		// Generate presigned URL for download
		if result.ZipPath != "" {
			downloadURL, err := ah.config.S3Service.GeneratePresignedURL(ah.config.S3Buckets.OutputsBucket, result.ZipPath, time.Hour)
			if err != nil {
				log.Printf("Warning: Failed to generate presigned URL for %s: %v", result.ZipPath, err)
				// Fallback to API endpoint
				downloadURL = "/api/v1/videos/" + result.ZipPath + "/download"
			}
			result.DownloadURL = downloadURL
		}
		c.JSON(http.StatusCreated, result)
	} else {
		c.JSON(http.StatusUnprocessableEntity, result)
	}
}

func (ah *APIHandlers) processVideoDirectly(c *gin.Context, file io.Reader, filename string) {
	result, err := ah.processorClient.ProcessVideo(filename, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ProcessingResult{
			Success: false,
			Message: "Erro ao processar vídeo: " + err.Error(),
		})
		return
	}

	if result.Success {
		c.JSON(http.StatusCreated, result)
	} else {
		c.JSON(http.StatusUnprocessableEntity, result)
	}
}

func (ah *APIHandlers) GetVideos(c *gin.Context) {
	if ah.config.IsS3Enabled() {
		ah.getVideosFromS3(c)
	} else {
		ah.getVideosFromFilesystem(c)
	}
}

func (ah *APIHandlers) getVideosFromS3(c *gin.Context) {
	files, err := ah.config.S3Service.ListFiles(ah.config.S3Buckets.OutputsBucket, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar arquivos do S3: " + err.Error()})
		return
	}

	results := make([]map[string]interface{}, 0, len(files))
	for _, file := range files {
		if !strings.HasSuffix(file, ".zip") {
			continue
		}

		info, err := ah.config.S3Service.GetFileInfo(ah.config.S3Buckets.OutputsBucket, file)
		if err != nil {
			log.Printf("Warning: Failed to get file info for %s: %v", file, err)
			continue
		}

		downloadURL, err := ah.config.S3Service.GeneratePresignedURL(ah.config.S3Buckets.OutputsBucket, file, time.Hour)
		if err != nil {
			log.Printf("Warning: Failed to generate presigned URL for %s: %v", file, err)
			downloadURL = "/api/v1/videos/" + filepath.Base(file) + "/download"
		}

		results = append(results, map[string]interface{}{
			"filename":     filepath.Base(file),
			"size":         *info.ContentLength,
			"created_at":   info.LastModified.Format("2006-01-02 15:04:05"),
			"download_url": downloadURL,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"videos": results,
		"total":  len(results),
	})
}

func (ah *APIHandlers) getVideosFromFilesystem(c *gin.Context) {
	files, err := filepath.Glob(filepath.Join(ah.config.OutputsDir, "*.zip"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar arquivos"})
		return
	}

	results := make([]map[string]interface{}, 0, len(files))
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"filename":     filepath.Base(file),
			"size":         info.Size(),
			"created_at":   info.ModTime().Format("2006-01-02 15:04:05"),
			"download_url": "/api/v1/videos/" + filepath.Base(file) + "/download",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"videos": results,
		"total":  len(results),
	})
}

// GetVideoDownload handles video download requests with flexible response modes.
// Supports both direct redirect to presigned URL and JSON response with URL.
// Query parameter 'redirect=true' enables direct redirect mode for better UX.
func (ah *APIHandlers) GetVideoDownload(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome do arquivo é obrigatório"})
		return
	}

	if ah.config.IsS3Enabled() {
		ah.downloadVideoFromS3(c, filename)
	} else {
		ah.downloadVideoFromFilesystem(c, filename)
	}
}

func (ah *APIHandlers) downloadVideoFromS3(c *gin.Context, filename string) {
	exists, err := ah.config.S3Service.FileExists(ah.config.S3Buckets.OutputsBucket, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar arquivo no S3: " + err.Error()})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		return
	}

	presignedURL, err := ah.config.S3Service.GeneratePresignedURL(ah.config.S3Buckets.OutputsBucket, filename, time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar URL de download: " + err.Error()})
		return
	}

	// Support both redirect and JSON response modes
	if c.Query("redirect") == "true" {
		// Direct redirect to S3 - better UX, faster download
		c.Redirect(http.StatusFound, presignedURL)
		return
	}

	// JSON response - allows frontend to handle URL as needed
	c.JSON(http.StatusOK, gin.H{
		"download_url": presignedURL,
		"filename":     filename,
		"expires_in":   3600,
	})
}

func (ah *APIHandlers) downloadVideoFromFilesystem(c *gin.Context, filename string) {
	filePath := filepath.Join(ah.config.OutputsDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/zip")
	c.File(filePath)
}

func (ah *APIHandlers) DeleteVideo(c *gin.Context) {
	filename := c.Param("filename")

	if ah.config.IsS3Enabled() {
		ah.deleteVideoFromS3(c, filename)
	} else {
		ah.deleteVideoFromFilesystem(c, filename)
	}
}

func (ah *APIHandlers) deleteVideoFromS3(c *gin.Context, filename string) {
	exists, err := ah.config.S3Service.FileExists(ah.config.S3Buckets.OutputsBucket, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar arquivo no S3: " + err.Error()})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		return
	}

	if err := ah.config.S3Service.DeleteFile(ah.config.S3Buckets.OutputsBucket, filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar arquivo do S3: " + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (ah *APIHandlers) deleteVideoFromFilesystem(c *gin.Context, filename string) {
	filePath := filepath.Join(ah.config.OutputsDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		return
	}

	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar arquivo"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func IsValidVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".mp4", ".avi", ".mov", ".mkv", ".wmv", ".flv", ".webm"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}
