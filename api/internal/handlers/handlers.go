package handlers

import (
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
		"status":    "healthy",
		"service":   "videogrinder-api",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
		"checks": gin.H{
			"directories": ah.checkDirectories(),
			"processor":   ah.checkProcessorConnectivity(),
		},
	}

	dirCheck := health["checks"].(gin.H)["directories"].(gin.H)
	procCheck := health["checks"].(gin.H)["processor"].(gin.H)

	if dirCheck["status"] != "healthy" || procCheck["status"] != "healthy" {
		health["status"] = "unhealthy"
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
				"status": "healthy",
				"path":   dir,
			}
		}
	}

	return gin.H{
		"status":  map[bool]string{true: "healthy", false: "unhealthy"}[allHealthy],
		"details": details,
	}
}

func (ah *APIHandlers) checkProcessorConnectivity() gin.H {
	start := time.Now()
	err := ah.processorClient.HealthCheck()
	latency := time.Since(start)

	if err != nil {
		return gin.H{
			"status":     "unhealthy",
			"url":        ah.config.ProcessorURL,
			"error":      err.Error(),
			"latency_ms": latency.Milliseconds(),
			"last_check": time.Now().Unix(),
		}
	}

	return gin.H{
		"status":     "healthy",
		"url":        ah.config.ProcessorURL,
		"latency_ms": latency.Milliseconds(),
		"last_check": time.Now().Unix(),
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

	result, err := ah.processorClient.ProcessVideo(header.Filename, file)
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

		results = append(results, map[string]interface{}{
			"filename":     file,
			"size":         *info.ContentLength,
			"created_at":   info.LastModified.Format("2006-01-02 15:04:05"),
			"download_url": "/api/v1/videos/" + file + "/download",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total":  len(results),
		"videos": results,
	})
}

func (ah *APIHandlers) GetVideoDownload(c *gin.Context) {
	filename := c.Param("filename")

	exists, err := ah.config.S3Service.FileExists(ah.config.S3Buckets.OutputsBucket, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar arquivo no S3: " + err.Error()})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		return
	}

	reader, err := ah.config.S3Service.DownloadFile(ah.config.S3Buckets.OutputsBucket, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao baixar arquivo do S3: " + err.Error()})
		return
	}
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("Warning: Failed to close S3 reader: %v", err)
		}
	}()

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/zip")

	if _, err := io.Copy(c.Writer, reader); err != nil {
		log.Printf("Error streaming file from S3: %v", err)
		return
	}
}

func (ah *APIHandlers) DeleteVideo(c *gin.Context) {
	filename := c.Param("filename")

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

	c.JSON(http.StatusOK, gin.H{"message": "Arquivo deletado com sucesso"})
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
