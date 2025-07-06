package config

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Service struct {
	client   *s3.S3
	uploader *s3manager.Uploader
	config   *AWSConfig
}

func NewS3Service(awsConfig *AWSConfig) (*S3Service, error) {
	sess, err := createAWSSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	client := s3.New(sess)
	uploader := s3manager.NewUploader(sess)

	return &S3Service{
		client:   client,
		uploader: uploader,
		config:   awsConfig,
	}, nil
}

func createAWSSession(awsConfig *AWSConfig) (*session.Session, error) {
	config := &aws.Config{
		Region: aws.String(awsConfig.Region),
	}

	if awsConfig.IsLocalStack() {
		config.Endpoint = aws.String(awsConfig.EndpointURL)
		config.S3ForcePathStyle = aws.Bool(true)
		config.Credentials = credentials.NewStaticCredentials(
			awsConfig.AccessKeyID,
			awsConfig.SecretAccessKey,
			"",
		)
	}

	return session.NewSession(config)
}

func (s *S3Service) UploadFile(bucket, key string, body io.Reader) error {
	return s.UploadFileWithContentType(bucket, key, body, "")
}

func (s *S3Service) UploadFileWithContentType(bucket, key string, body io.Reader, contentType string) error {
	// Set appropriate content type based on file extension if not provided
	if contentType == "" {
		if strings.HasSuffix(strings.ToLower(key), ".zip") {
			contentType = "application/zip"
		} else if strings.HasSuffix(strings.ToLower(key), ".mp4") {
			contentType = "video/mp4"
		} else {
			contentType = "binary/octet-stream"
		}
	}

	_, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	log.Printf("Successfully uploaded file to s3://%s/%s with content-type: %s", bucket, key, contentType)
	return nil
}

func (s *S3Service) DownloadFile(bucket, key string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file from S3: %w", err)
	}

	return result.Body, nil
}

func (s *S3Service) DeleteFile(bucket, key string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	log.Printf("Successfully deleted file from s3://%s/%s", bucket, key)
	return nil
}

func (s *S3Service) ListFiles(bucket, prefix string) ([]string, error) {
	result, err := s.client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list files from S3: %w", err)
	}

	var files []string
	for _, obj := range result.Contents {
		if obj.Key != nil {
			files = append(files, *obj.Key)
		}
	}

	return files, nil
}

func (s *S3Service) FileExists(bucket, key string) (bool, error) {
	_, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if file exists in S3: %w", err)
	}

	return true, nil
}

func (s *S3Service) GetFileInfo(bucket, key string) (*s3.HeadObjectOutput, error) {
	result, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file info from S3: %w", err)
	}

	return result, nil
}

func (s *S3Service) GeneratePresignedURL(bucket, key string, expiration time.Duration) (string, error) {
	// Use configurable timeout instead of hardcoded expiration
	if expiration == 0 {
		expiration = s.config.PresignedTimeout
	}

	if s.config.IsLocalStack() {
		// For LocalStack, generate proper presigned URL using AWS SDK
		req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

		// Override endpoint for browser-accessible URL
		if s.config.ExternalURL != "" {
			req.HTTPRequest.URL.Host = strings.Replace(s.config.ExternalURL, "http://", "", 1)
			req.HTTPRequest.URL.Host = strings.Replace(req.HTTPRequest.URL.Host, "https://", "", 1)
		} else {
			// Default LocalStack external URL for browser access
			req.HTTPRequest.URL.Host = "localhost:4566"
		}
		req.HTTPRequest.URL.Scheme = "http" // LocalStack uses HTTP

		urlStr, err := req.Presign(expiration)
		if err != nil {
			return "", fmt.Errorf("failed to generate LocalStack presigned URL: %w", err)
		}

		// Security First: Validate generated URL
		if err := s.config.ValidateURL(urlStr); err != nil {
			return "", fmt.Errorf("generated URL failed security validation: %w", err)
		}

		return urlStr, nil
	}

	// For real AWS, generate standard presigned URL
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	urlStr, err := req.Presign(expiration)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	// Security First: Validate generated URL
	if err := s.config.ValidateURL(urlStr); err != nil {
		return "", fmt.Errorf("generated URL failed security validation: %w", err)
	}

	log.Printf("Generated presigned URL for s3://%s/%s (expires in %v)", bucket, key, expiration)
	return urlStr, nil
}
