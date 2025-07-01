package api

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"video-processor/internal/config"
	"video-processor/internal/models"
	"video-processor/internal/services"

	"github.com/gin-gonic/gin"
)

type APIHandlers struct {
	videoService *services.VideoService
	config       *config.Config
}

func NewAPIHandlers(videoService *services.VideoService, cfg *config.Config) *APIHandlers {
	return &APIHandlers{
		videoService: videoService,
		config:       cfg,
	}
}

func (ah *APIHandlers) HandleVideoUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(400, models.ProcessingResult{
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
		c.JSON(400, models.ProcessingResult{
			Success: false,
			Message: "Formato de arquivo não suportado. Use: mp4, avi, mov, mkv",
		})
		return
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s", timestamp, filepath.Base(header.Filename))
	videoPath := filepath.Join(ah.config.UploadsDir, filename)

	cleanVideoPath := filepath.Clean(videoPath)
	uploadsDir, _ := filepath.Abs(ah.config.UploadsDir)
	absVideoPath, _ := filepath.Abs(cleanVideoPath)
	if !strings.HasPrefix(absVideoPath, uploadsDir+string(filepath.Separator)) {
		c.JSON(400, models.ProcessingResult{
			Success: false,
			Message: "Invalid file path",
		})
		return
	}

	out, err := os.Create(filepath.Clean(videoPath))
	if err != nil {
		c.JSON(500, models.ProcessingResult{
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
		c.JSON(500, models.ProcessingResult{
			Success: false,
			Message: "Erro ao salvar arquivo: " + err.Error(),
		})
		return
	}

	result := ah.videoService.ProcessVideo(videoPath, timestamp)

	if result.Success {
		if err := os.Remove(videoPath); err != nil {
			log.Printf("Warning: Failed to remove video file %s: %v", videoPath, err)
		}
	}

	c.JSON(200, result)
}

func (ah *APIHandlers) HandleDownload(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(ah.config.OutputsDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": "Arquivo não encontrado"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/zip")

	c.File(filePath)
}

func (ah *APIHandlers) HandleStatus(c *gin.Context) {
	files, err := filepath.Glob(filepath.Join(ah.config.OutputsDir, "*.zip"))
	if err != nil {
		c.JSON(500, gin.H{"error": "Erro ao listar arquivos"})
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
			"download_url": "/download/" + filepath.Base(file),
		})
	}

	c.JSON(200, gin.H{
		"files": results,
		"total": len(results),
	})
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
