package config

import (
	"os"
	"testing"
	"time"
)

func TestNewAWSConfig(t *testing.T) {
	config := NewAWSConfig()

	if config.Region != "us-east-1" {
		t.Errorf("Expected default region us-east-1, got %s", config.Region)
	}

	if config.S3Buckets.UploadsBucket != "videogrinder-uploads" {
		t.Errorf("Expected default uploads bucket videogrinder-uploads, got %s", config.S3Buckets.UploadsBucket)
	}

	if config.DynamoDB.VideoJobsTable != "video-jobs" {
		t.Errorf("Expected default table video-jobs, got %s", config.DynamoDB.VideoJobsTable)
	}

	if config.PresignedTimeout != time.Hour {
		t.Errorf("Expected default presigned timeout 1h, got %v", config.PresignedTimeout)
	}
}

func TestAWSConfigWithEnvironmentVariables(t *testing.T) {
	os.Setenv("AWS_REGION", "us-west-2")
	os.Setenv("AWS_ENDPOINT_URL", DefaultLocalStackEndpoint)
	os.Setenv("AWS_EXTERNAL_URL", DefaultLocalStackEndpoint)
	os.Setenv("AWS_PRESIGNED_TIMEOUT", "30m")
	os.Setenv("S3_BUCKET_UPLOADS", "test-uploads")
	defer func() {
		os.Unsetenv("AWS_REGION")
		os.Unsetenv("AWS_ENDPOINT_URL")
		os.Unsetenv("AWS_EXTERNAL_URL")
		os.Unsetenv("AWS_PRESIGNED_TIMEOUT")
		os.Unsetenv("S3_BUCKET_UPLOADS")
	}()

	config := NewAWSConfig()

	if config.Region != "us-west-2" {
		t.Errorf("Expected region us-west-2, got %s", config.Region)
	}

	if config.EndpointURL != DefaultLocalStackEndpoint {
		t.Errorf("Expected endpoint %s, got %s", DefaultLocalStackEndpoint, config.EndpointURL)
	}

	if config.ExternalURL != DefaultLocalStackEndpoint {
		t.Errorf("Expected external URL %s, got %s", DefaultLocalStackEndpoint, config.ExternalURL)
	}

	if config.PresignedTimeout != 30*time.Minute {
		t.Errorf("Expected presigned timeout 30m, got %v", config.PresignedTimeout)
	}

	if config.S3Buckets.UploadsBucket != "test-uploads" {
		t.Errorf("Expected uploads bucket test-uploads, got %s", config.S3Buckets.UploadsBucket)
	}
}

func TestIsLocalStack(t *testing.T) {
	config := &AWSConfig{EndpointURL: ""}
	if config.IsLocalStack() {
		t.Error("Expected IsLocalStack to return false when EndpointURL is empty")
	}

	config.EndpointURL = DefaultLocalStackEndpoint
	if !config.IsLocalStack() {
		t.Error("Expected IsLocalStack to return true when EndpointURL is set")
	}
}

func TestGetExternalEndpoint(t *testing.T) {
	tests := []struct {
		name        string
		config      *AWSConfig
		expected    string
		description string
	}{
		{
			name: "external URL set",
			config: &AWSConfig{
				ExternalURL: "http://custom-endpoint:9000",
				EndpointURL: "http://localstack:4566",
				Region:      "us-east-1",
			},
			expected:    "http://custom-endpoint:9000",
			description: "should return custom external URL when set",
		},
		{
			name: "localstack default",
			config: &AWSConfig{
				ExternalURL: "",
				EndpointURL: "http://localstack:4566",
				Region:      "us-east-1",
			},
			expected:    DefaultLocalStackEndpoint,
			description: "should return default LocalStack URL when external URL not set",
		},
		{
			name: "production AWS",
			config: &AWSConfig{
				ExternalURL: "",
				EndpointURL: "",
				Region:      "us-east-1",
			},
			expected:    "https://s3.us-east-1.amazonaws.com",
			description: "should return AWS S3 endpoint for production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetExternalEndpoint()
			if result != tt.expected {
				t.Errorf("GetExternalEndpoint() = %v, want %v (%s)", result, tt.expected, tt.description)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name      string
		config    *AWSConfig
		url       string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid HTTPS URL for production",
			config:    &AWSConfig{EndpointURL: ""},
			url:       "https://s3.amazonaws.com/bucket/key",
			wantError: false,
		},
		{
			name:      "valid HTTP URL for LocalStack",
			config:    &AWSConfig{EndpointURL: DefaultLocalStackEndpoint},
			url:       DefaultLocalStackEndpoint + "/bucket/key",
			wantError: false,
		},
		{
			name:      "invalid URL format",
			config:    &AWSConfig{EndpointURL: ""},
			url:       "://invalid-url",
			wantError: true,
			errorMsg:  "invalid URL format",
		},
		{
			name:      "HTTP in production should fail",
			config:    &AWSConfig{EndpointURL: ""},
			url:       "http://s3.amazonaws.com/bucket/key",
			wantError: true,
			errorMsg:  "HTTPS required in production",
		},
		{
			name:      "URL without host",
			config:    &AWSConfig{EndpointURL: ""},
			url:       "https:///bucket/key",
			wantError: true,
			errorMsg:  "URL must have a valid host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.ValidateURL(tt.url)
			if tt.wantError && err == nil {
				t.Errorf("ValidateURL() expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("ValidateURL() unexpected error: %v", err)
			}
			if tt.wantError && err != nil && tt.errorMsg != "" {
				if !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("ValidateURL() error = %v, want error containing %v", err, tt.errorMsg)
				}
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{"1h", time.Hour},
		{"30m", 30 * time.Minute},
		{"invalid", time.Hour}, // fallback to default
		{"", time.Hour},        // fallback to default
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseDuration(tt.input)
			if result != tt.expected {
				t.Errorf("parseDuration(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				indexString(s, substr) >= 0)))
}

func indexString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func TestGetEndpoints(t *testing.T) {
	config := &AWSConfig{
		Region:      "us-east-1",
		EndpointURL: "",
	}

	expectedS3 := "https://s3.us-east-1.amazonaws.com"
	if config.GetS3Endpoint() != expectedS3 {
		t.Errorf("Expected S3 endpoint %s, got %s", expectedS3, config.GetS3Endpoint())
	}

	config.EndpointURL = DefaultLocalStackEndpoint
	if config.GetS3Endpoint() != DefaultLocalStackEndpoint {
		t.Errorf("Expected LocalStack endpoint %s, got %s", DefaultLocalStackEndpoint, config.GetS3Endpoint())
	}
}

func TestGetDynamoDBEndpoint(t *testing.T) {
	config := &AWSConfig{
		Region:      "us-east-1",
		EndpointURL: "",
	}

	expectedDynamoDB := "https://dynamodb.us-east-1.amazonaws.com"
	if config.GetDynamoDBEndpoint() != expectedDynamoDB {
		t.Errorf("Expected DynamoDB endpoint %s, got %s", expectedDynamoDB, config.GetDynamoDBEndpoint())
	}

	config.EndpointURL = DefaultLocalStackEndpoint
	if config.GetDynamoDBEndpoint() != DefaultLocalStackEndpoint {
		t.Errorf("Expected LocalStack endpoint %s, got %s", DefaultLocalStackEndpoint, config.GetDynamoDBEndpoint())
	}
}

func TestGetSQSEndpoint(t *testing.T) {
	config := &AWSConfig{
		Region:      "us-east-1",
		EndpointURL: "",
	}

	expectedSQS := "https://sqs.us-east-1.amazonaws.com"
	if config.GetSQSEndpoint() != expectedSQS {
		t.Errorf("Expected SQS endpoint %s, got %s", expectedSQS, config.GetSQSEndpoint())
	}

	config.EndpointURL = DefaultLocalStackEndpoint
	if config.GetSQSEndpoint() != DefaultLocalStackEndpoint {
		t.Errorf("Expected LocalStack endpoint %s, got %s", DefaultLocalStackEndpoint, config.GetSQSEndpoint())
	}
}
