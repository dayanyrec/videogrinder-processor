package services

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"video-processor/processor/internal/config"
	"video-processor/processor/internal/models"
	"video-processor/processor/internal/utils"
)

type VideoService struct {
	config *config.ProcessorConfig
}

func NewVideoService(cfg *config.ProcessorConfig) *VideoService {
	return &VideoService{
		config: cfg,
	}
}

func (vs *VideoService) ProcessVideo(videoPath, timestamp string) models.ProcessingResult {
	fmt.Printf("Iniciando processamento: %s\n", videoPath)

	if err := utils.ValidateProcessingInputs(videoPath, timestamp); err != nil {
		return models.ProcessingResult{Success: false, Message: err.Error()}
	}

	tempDir := filepath.Join(vs.config.TempDir, timestamp)
	if err := utils.SetupTempDirectory(tempDir); err != nil {
		return models.ProcessingResult{Success: false, Message: err.Error()}
	}
	defer utils.CleanupTempDirectory(tempDir)

	frames, err := vs.extractFrames(videoPath, tempDir)
	if err != nil {
		return models.ProcessingResult{Success: false, Message: err.Error()}
	}

	fmt.Printf("📸 Extraídos %d frames\n", len(frames))

	zipPath, err := vs.createFramesZip(frames, timestamp)
	if err != nil {
		return models.ProcessingResult{Success: false, Message: err.Error()}
	}

	fmt.Printf("✅ ZIP criado: %s\n", zipPath)

	imageNames := make([]string, len(frames))
	for i, frame := range frames {
		imageNames[i] = filepath.Base(frame)
	}

	return models.ProcessingResult{
		Success:    true,
		Message:    fmt.Sprintf("Processamento concluído! %d frames extraídos.", len(frames)),
		ZipPath:    filepath.Base(zipPath),
		FrameCount: len(frames),
		Images:     imageNames,
	}
}

func (vs *VideoService) extractFrames(videoPath, tempDir string) ([]string, error) {
	framePattern := filepath.Join(tempDir, "frame_%04d.png")

	videoPath = filepath.Clean(videoPath)
	framePattern = filepath.Clean(framePattern)

	if err := utils.ValidatePathSafety(videoPath, framePattern); err != nil {
		return nil, err
	}

	absVideoPath, err := filepath.Abs(videoPath)
	if err != nil {
		return nil, fmt.Errorf("error resolving video path: %w", err)
	}
	absFramePattern, err := filepath.Abs(framePattern)
	if err != nil {
		return nil, fmt.Errorf("error resolving frame pattern path: %w", err)
	}

	cmd := exec.Command("ffmpeg", // #nosec G204
		"-i", absVideoPath,
		"-vf", "fps=1",
		"-y",
		absFramePattern,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("erro no ffmpeg: %s\nOutput: %s", err.Error(), string(output))
	}

	frames, err := filepath.Glob(filepath.Join(tempDir, "*.png"))
	if err != nil || len(frames) == 0 {
		return nil, fmt.Errorf("nenhum frame foi extraído do vídeo")
	}

	return frames, nil
}

func (vs *VideoService) createFramesZip(frames []string, timestamp string) (string, error) {
	zipFilename := fmt.Sprintf("frames_%s.zip", timestamp)

	if vs.config.IsS3Enabled() {
		return vs.createFramesZipToS3(frames, zipFilename)
	} else {
		return vs.createFramesZipToFilesystem(frames, zipFilename)
	}
}

func (vs *VideoService) createFramesZipToS3(frames []string, zipFilename string) (string, error) {
	var zipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuffer)

	for _, file := range frames {
		if err := vs.addFileToZip(zipWriter, file); err != nil {
			if closeErr := zipWriter.Close(); closeErr != nil {
				log.Printf("Warning: Failed to close ZIP writer: %v", closeErr)
			}
			return "", fmt.Errorf("erro ao adicionar arquivo ao ZIP: %w", err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		return "", fmt.Errorf("erro ao finalizar ZIP: %w", err)
	}

	reader := bytes.NewReader(zipBuffer.Bytes())
	if err := vs.config.S3Service.UploadFile(vs.config.S3Buckets.OutputsBucket, zipFilename, reader); err != nil {
		return "", fmt.Errorf("erro ao fazer upload do ZIP para S3: %w", err)
	}

	fmt.Printf("✅ ZIP salvo no S3: s3://%s/%s\n", vs.config.S3Buckets.OutputsBucket, zipFilename)
	return zipFilename, nil
}

func (vs *VideoService) createFramesZipToFilesystem(frames []string, zipFilename string) (string, error) {
	zipPath := filepath.Join(vs.config.OutputsDir, zipFilename)

	if err := utils.ValidateOutputPath(zipPath, vs.config.OutputsDir); err != nil {
		return "", err
	}

	if err := vs.createZipFile(frames, zipPath); err != nil {
		return "", fmt.Errorf("erro ao criar arquivo ZIP: %w", err)
	}

	return zipPath, nil
}

func (vs *VideoService) createZipFile(files []string, zipPath string) error {
	zipFile, err := os.Create(filepath.Clean(zipPath))
	if err != nil {
		return err
	}
	defer func() {
		if err := zipFile.Close(); err != nil {
			log.Printf("Warning: Failed to close ZIP file: %v", err)
		}
	}()

	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		if err := zipWriter.Close(); err != nil {
			log.Printf("Warning: Failed to close ZIP writer: %v", err)
		}
	}()

	for _, file := range files {
		if err := vs.addFileToZip(zipWriter, file); err != nil {
			return err
		}
	}

	return nil
}

func (vs *VideoService) addFileToZip(zipWriter *zip.Writer, filename string) error {
	file, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Warning: Failed to close file %s: %v", filename, err)
		}
	}()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filepath.Base(filename)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}
