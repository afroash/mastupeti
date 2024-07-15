package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/afroash/mastupeti/initializers"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

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

	return result.Location, nil
}
