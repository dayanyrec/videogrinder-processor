package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"video-processor/internal/models"
)

// ProcessorClientInterface defines the interface for processor client operations
type ProcessorClientInterface interface {
	ProcessVideo(filename string, fileReader io.Reader) (models.ProcessingResult, error)
	HealthCheck() error
}

type ProcessorClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewProcessorClient(baseURL string) *ProcessorClient {
	return &ProcessorClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute, // Video processing can take time
		},
	}
}

// ProcessVideo sends a video file to the processor service for processing
func (pc *ProcessorClient) ProcessVideo(filename string, fileReader io.Reader) (models.ProcessingResult, error) {
	// Create multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Create form file
	part, err := writer.CreateFormFile("video", filename)
	if err != nil {
		return models.ProcessingResult{}, err
	}

	// Copy file content
	_, err = io.Copy(part, fileReader)
	if err != nil {
		return models.ProcessingResult{}, err
	}

	// Close writer to finalize the form
	err = writer.Close()
	if err != nil {
		return models.ProcessingResult{}, err
	}

	// Create request
	url := fmt.Sprintf("%s/process", pc.baseURL)
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return models.ProcessingResult{}, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return models.ProcessingResult{}, fmt.Errorf("failed to send request to processor: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Warning: Failed to close response body: %v\n", err)
		}
	}()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ProcessingResult{}, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var result models.ProcessingResult
	if err := json.Unmarshal(body, &result); err != nil {
		return models.ProcessingResult{}, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}

// HealthCheck checks if the processor service is healthy
func (pc *ProcessorClient) HealthCheck() error {
	url := fmt.Sprintf("%s/health", pc.baseURL)
	resp, err := pc.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("processor service unavailable: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Warning: Failed to close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("processor service unhealthy: status %d", resp.StatusCode)
	}

	return nil
}
