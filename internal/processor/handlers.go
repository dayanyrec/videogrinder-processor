package processor

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"video-processor/internal/config"
	"video-processor/internal/models"
	"video-processor/internal/services"

	"github.com/gin-gonic/gin"
)

type ProcessorHandlers struct {
	videoService *services.VideoService
	config       *config.Config
}

func NewProcessorHandlers(videoService *services.VideoService, cfg *config.Config) *ProcessorHandlers {
	return &ProcessorHandlers{
		videoService: videoService,
		config:       cfg,
	}
}

// ProcessVideoUpload handles video upload and processing
func (ph *ProcessorHandlers) ProcessVideoUpload(c *gin.Context) {
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

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s", timestamp, filepath.Base(header.Filename))
	videoPath := filepath.Join(ph.config.UploadsDir, filename)

	// Save uploaded file
	out, err := os.Create(filepath.Clean(videoPath))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ProcessingResult{
			Success: false,
			Message: "Erro ao salvar arquivo: " + err.Error(),
		})
		return
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Printf("Warning: Failed to close output file: %v", err)
		}
	}()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ProcessingResult{
			Success: false,
			Message: "Erro ao salvar arquivo: " + err.Error(),
		})
		return
	}

	// Process video
	result := ph.videoService.ProcessVideo(videoPath, timestamp)

	// Clean up uploaded file
	if err := os.Remove(videoPath); err != nil {
		log.Printf("Warning: Failed to remove video file %s: %v", videoPath, err)
	}

	if result.Success {
		c.JSON(http.StatusCreated, result)
	} else {
		c.JSON(http.StatusUnprocessableEntity, result)
	}
}

// GetProcessorStatus returns processor service status
func (ph *ProcessorHandlers) GetProcessorStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "video-processor",
		"timestamp": time.Now().Unix(),
	})
}

// IsValidVideoFile validates video file extensions
func IsValidVideoFile(filename string) bool {
	ext := filepath.Ext(filename)
	if len(ext) > 0 {
		ext = ext[1:] // Remove the dot
	}
	ext = strings.ToLower(ext)
	validExts := []string{"mp4", "avi", "mov", "mkv", "wmv", "flv", "webm"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}
