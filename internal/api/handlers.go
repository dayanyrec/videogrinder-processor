package api

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"video-processor/internal/clients"
	"video-processor/internal/config"
	"video-processor/internal/models"

	"github.com/gin-gonic/gin"
)

type APIHandlers struct {
	processorClient clients.ProcessorClientInterface
	config          *config.Config
}

func NewAPIHandlers(cfg *config.Config) *APIHandlers {
	return &APIHandlers{
		processorClient: clients.NewProcessorClient(cfg.ProcessorURL),
		config:          cfg,
	}
}

// GetAPIHealth provides comprehensive health check for the API service
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

	// Determine overall status based on checks
	dirCheck := health["checks"].(gin.H)["directories"].(gin.H)
	procCheck := health["checks"].(gin.H)["processor"].(gin.H)

	if dirCheck["status"] != "healthy" || procCheck["status"] != "healthy" {
		health["status"] = "unhealthy"
		c.JSON(http.StatusServiceUnavailable, health)
		return
	}

	c.JSON(http.StatusOK, health)
}

// checkDirectories verifies that all required directories exist and are writable
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

		// Check if directory exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			details[dirName] = gin.H{
				"status": "missing",
				"path":   dir,
				"error":  "Directory does not exist",
			}
			allHealthy = false
			continue
		}

		// Check if directory is writable
		testFile := filepath.Join(dir, ".health_check_test")
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			details[dirName] = gin.H{
				"status": "not_writable",
				"path":   dir,
				"error":  "Directory is not writable: " + err.Error(),
			}
			allHealthy = false
		} else {
			// Clean up test file
			os.Remove(testFile)
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

// checkProcessorConnectivity verifies connectivity to the Processor service
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

	// Check processor health before processing
	if err := ah.processorClient.HealthCheck(); err != nil {
		c.JSON(http.StatusServiceUnavailable, models.ProcessingResult{
			Success: false,
			Message: "Serviço de processamento indisponível: " + err.Error(),
		})
		return
	}

	// Send file to processor service
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

func (ah *APIHandlers) GetVideoDownload(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(ah.config.OutputsDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/zip")

	c.File(filePath)
}

func (ah *APIHandlers) DeleteVideo(c *gin.Context) {
	filename := c.Param("filename")
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
