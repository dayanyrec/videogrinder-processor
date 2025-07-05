package config

import (
	"fmt"
	"net/http"
	"time"
)

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	EndpointURL     string
	S3Buckets       S3Config
	DynamoDB        DynamoDBConfig
	SQS             SQSConfig
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
	}
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LocalStack health check failed: status %d", resp.StatusCode)
	}

	return nil
}
