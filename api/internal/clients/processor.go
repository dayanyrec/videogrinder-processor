package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"video-processor/api/internal/models"
)

type ProcessorClientInterface interface {
	ProcessVideo(filename string, videoFile io.Reader) (*models.ProcessingResult, error)
	HealthCheck() error
}

type ProcessorClient struct {
	baseURL string
	client  *http.Client
}

func NewProcessorClient(baseURL string) *ProcessorClient {
	return &ProcessorClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 300 * time.Second, // 5 minutes for video processing
		},
	}
}

func (pc *ProcessorClient) ProcessVideo(filename string, videoFile io.Reader) (*models.ProcessingResult, error) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("video", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, videoFile); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest("POST", pc.baseURL+"/process", &requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := pc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Warning: Failed to close response body: %v", err)
		}
	}()

	var result models.ProcessingResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (pc *ProcessorClient) HealthCheck() error {
	resp, err := pc.client.Get(pc.baseURL + "/health")
	if err != nil {
		return fmt.Errorf("failed to connect to processor service: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Warning: Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("processor service returned status %d", resp.StatusCode)
	}

	return nil
}
