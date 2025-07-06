package config

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type AWSConfig struct {
	Region           string
	AccessKeyID      string
	SecretAccessKey  string
	EndpointURL      string
	ExternalURL      string // Browser-accessible URL (different from internal EndpointURL)
	S3Buckets        S3Config
	DynamoDB         DynamoDBConfig
	SQS              SQSConfig
	PresignedTimeout time.Duration // Configurable timeout for presigned URLs
}

type S3Config struct {
	UploadsBucket string
	OutputsBucket string
}

type DynamoDBConfig struct {
	VideoJobsTable string
}

type SQSConfig struct {
	VideoProcessingQueue    string
	VideoProcessingDLQQueue string
}

func NewAWSConfig() *AWSConfig {
	return &AWSConfig{
		Region:          GetEnv("AWS_REGION", "us-east-1"),
		AccessKeyID:     GetEnv("AWS_ACCESS_KEY_ID", ""),
		SecretAccessKey: GetEnv("AWS_SECRET_ACCESS_KEY", ""),
		EndpointURL:     GetEnv("AWS_ENDPOINT_URL", ""),
		ExternalURL:     GetEnv("AWS_EXTERNAL_URL", ""), // New: browser-accessible URL
		S3Buckets: S3Config{
			UploadsBucket: GetEnv("S3_BUCKET_UPLOADS", "videogrinder-uploads"),
			OutputsBucket: GetEnv("S3_BUCKET_OUTPUTS", "videogrinder-outputs"),
		},
		DynamoDB: DynamoDBConfig{
			VideoJobsTable: GetEnv("DYNAMODB_TABLE_VIDEO_JOBS", "video-jobs"),
		},
		SQS: SQSConfig{
			VideoProcessingQueue:    GetEnv("SQS_QUEUE_VIDEO_PROCESSING", "video-processing-queue"),
			VideoProcessingDLQQueue: GetEnv("SQS_QUEUE_VIDEO_PROCESSING_DLQ", "video-processing-dlq"),
		},
		PresignedTimeout: parseDuration(GetEnv("AWS_PRESIGNED_TIMEOUT", "1h")),
	}
}

// parseDuration safely parses duration with fallback
func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Printf("Warning: Invalid duration %s, using default 1h", s)
		return time.Hour
	}
	return d
}

func (c *AWSConfig) IsLocalStack() bool {
	return c.EndpointURL != ""
}

func (c *AWSConfig) GetS3Endpoint() string {
	if c.IsLocalStack() {
		return c.EndpointURL
	}
	return fmt.Sprintf("https://s3.%s.amazonaws.com", c.Region)
}

func (c *AWSConfig) GetDynamoDBEndpoint() string {
	if c.IsLocalStack() {
		return c.EndpointURL
	}
	return fmt.Sprintf("https://dynamodb.%s.amazonaws.com", c.Region)
}

func (c *AWSConfig) GetSQSEndpoint() string {
	if c.IsLocalStack() {
		return c.EndpointURL
	}
	return fmt.Sprintf("https://sqs.%s.amazonaws.com", c.Region)
}

func (c *AWSConfig) CheckHealth() error {
	if !c.IsLocalStack() {
		return nil
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	healthURL := fmt.Sprintf("%s/health", c.EndpointURL)
	resp, err := client.Get(healthURL)
	if err != nil {
		return fmt.Errorf("failed to connect to LocalStack: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Warning: Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LocalStack health check failed: status %d", resp.StatusCode)
	}

	return nil
}

// GetExternalEndpoint returns browser-accessible endpoint
func (c *AWSConfig) GetExternalEndpoint() string {
	if c.ExternalURL != "" {
		return c.ExternalURL
	}
	if c.IsLocalStack() {
		// Default LocalStack external URL for browser access
		return "http://localhost:4566"
	}
	return c.GetS3Endpoint()
}

// ValidateURL validates URLs following Security First mandate
func (c *AWSConfig) ValidateURL(rawURL string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Security First: Enforce HTTPS in production
	if !c.IsLocalStack() && parsedURL.Scheme != "https" {
		return fmt.Errorf("HTTPS required in production, got: %s", parsedURL.Scheme)
	}

	// Security First: Validate hostname
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a valid host")
	}

	return nil
}
