package config

import (
	"os"
	"testing"
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
}

func TestAWSConfigWithEnvironmentVariables(t *testing.T) {
	os.Setenv("AWS_REGION", "us-west-2")
	os.Setenv("AWS_ENDPOINT_URL", "http://localhost:4566")
	os.Setenv("S3_BUCKET_UPLOADS", "test-uploads")
	defer func() {
		os.Unsetenv("AWS_REGION")
		os.Unsetenv("AWS_ENDPOINT_URL")
		os.Unsetenv("S3_BUCKET_UPLOADS")
	}()

	config := NewAWSConfig()

	if config.Region != "us-west-2" {
		t.Errorf("Expected region us-west-2, got %s", config.Region)
	}

	if config.EndpointURL != "http://localhost:4566" {
		t.Errorf("Expected endpoint http://localhost:4566, got %s", config.EndpointURL)
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

	config.EndpointURL = "http://localhost:4566"
	if !config.IsLocalStack() {
		t.Error("Expected IsLocalStack to return true when EndpointURL is set")
	}
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

	config.EndpointURL = "http://localhost:4566"
	if config.GetS3Endpoint() != "http://localhost:4566" {
		t.Errorf("Expected LocalStack endpoint http://localhost:4566, got %s", config.GetS3Endpoint())
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

	config.EndpointURL = "http://localhost:4566"
	if config.GetDynamoDBEndpoint() != "http://localhost:4566" {
		t.Errorf("Expected LocalStack endpoint http://localhost:4566, got %s", config.GetDynamoDBEndpoint())
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

	config.EndpointURL = "http://localhost:4566"
	if config.GetSQSEndpoint() != "http://localhost:4566" {
		t.Errorf("Expected LocalStack endpoint http://localhost:4566, got %s", config.GetSQSEndpoint())
	}
}
