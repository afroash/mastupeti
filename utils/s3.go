package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/afroash/mastupeti/initializers"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// func to upload video to s3 should return metadata of the video. Update1.0
func UploadToS3(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	initializers.LoadEnvVariables()
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return "", fmt.Errorf("failed to load configuration, %v", err)
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("mastupeti-videos-store"),
		Key:    aws.String(fileHeader.Filename),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}

	return *result.Key, nil
}

// Function to get a pre-signed URL for a single video
func GetSignedURL(key string) (string, error) {
	initializers.LoadEnvVariables()
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return "", fmt.Errorf("failed to load configuration, %v", err)
	}

	client := s3.NewFromConfig(cfg)

	presigner := s3.NewPresignClient(client)
	req, err := presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("mastupeti-videos-store"),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(15*time.Minute)) // URL valid for 15 minutes
	if err != nil {
		return "", fmt.Errorf("failed to sign url, %v", err)
	}

	return req.URL, nil
}

// Function to get pre-signed URLs for multiple videos
func GetSignedURLs(keys []string) (map[string]string, error) {
	urls := make(map[string]string)
	for _, key := range keys {
		url, err := GetSignedURL(key)
		if err != nil {
			return nil, fmt.Errorf("failed to get signed URL for key %s: %v", key, err)
		}
		urls[key] = url
	}
	return urls, nil
}
