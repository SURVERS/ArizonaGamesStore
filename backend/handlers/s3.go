package handlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func UploadFileToS3(filePath string, fileName string) (string, error) {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("AWS_S3_BUCKET")
	endpoint := os.Getenv("AWS_ENDPOINT_URL")

	if accessKey == "" || secretKey == "" || region == "" || bucket == "" {
		return "", fmt.Errorf("AWS credentials or bucket name not set in environment variables")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
		o.UsePathStyle = true
	})

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %q, %v", filePath, err)
	}
	defer file.Close()

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   file,
		ACL:    "public-read",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}

	publicURL := fmt.Sprintf("https://%s.hb.vkcloud-storage.ru/%s", bucket, fileName)
	return publicURL, nil
}

func DeleteFileFromS3(fileKey string) error {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("AWS_S3_BUCKET")
	endpoint := os.Getenv("AWS_ENDPOINT_URL")

	if accessKey == "" || secretKey == "" || region == "" || bucket == "" {
		return fmt.Errorf("AWS credentials or bucket name not set")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
		o.UsePathStyle = true
	})

	_, err = s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	})

	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %v", err)
	}

	fmt.Printf("Successfully deleted %q from S3\n", fileKey)
	return nil
}

func extractS3KeyFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 3 {
		return strings.Join(parts[3:], "/")
	}
	return ""
}

func generateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	uniqueID := uuid.New().String()
	return fmt.Sprintf("%s%s", uniqueID, ext)
}
