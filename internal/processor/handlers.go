package processor

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
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

// GetProcessorStatus returns comprehensive processor service status
func (ph *ProcessorHandlers) GetProcessorStatus(c *gin.Context) {
	health := gin.H{
		"status":    "healthy",
		"service":   "videogrinder-processor",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
		"checks": gin.H{
			"directories": ph.checkDirectories(),
			"ffmpeg":      ph.checkFFmpegAvailability(),
		},
	}

	// Determine overall status based on checks
	dirCheck := health["checks"].(gin.H)["directories"].(gin.H)
	ffmpegCheck := health["checks"].(gin.H)["ffmpeg"].(gin.H)

	if dirCheck["status"] != "healthy" || ffmpegCheck["status"] != "healthy" {
		health["status"] = "unhealthy"
		c.JSON(http.StatusServiceUnavailable, health)
		return
	}

	c.JSON(http.StatusOK, health)
}

// checkDirectories verifies that all required directories exist and are writable
func (ph *ProcessorHandlers) checkDirectories() gin.H {
	directories := []string{
		ph.config.UploadsDir,
		ph.config.OutputsDir,
		ph.config.TempDir,
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

// checkFFmpegAvailability verifies FFmpeg is available and functional
func (ph *ProcessorHandlers) checkFFmpegAvailability() gin.H {
	start := time.Now()

	// Try to run ffmpeg -version to check availability
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()
	latency := time.Since(start)

	if err != nil {
		return gin.H{
			"status":     "unhealthy",
			"error":      "FFmpeg not available: " + err.Error(),
			"latency_ms": latency.Milliseconds(),
			"last_check": time.Now().Unix(),
		}
	}

	return gin.H{
		"status":     "healthy",
		"latency_ms": latency.Milliseconds(),
		"last_check": time.Now().Unix(),
	}
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
