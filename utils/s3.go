package utils

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/afroash/mastupeti/initializers"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ProgressReader struct {
	io.Reader
	Total      int64
	ReadSoFar  int64
	OnProgress func(percentage int)
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.ReadSoFar += int64(n)
	if pr.Total > 0 {
		percentage := int(float64(pr.ReadSoFar) / float64(pr.Total) * 100)
		pr.OnProgress(percentage)
	}
	return n, err
}

func UploadToS3(file io.Reader, fileHeader *multipart.FileHeader) (string, error) {
	initializers.LoadEnvVariables()
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return "", fmt.Errorf("failed to load configuration, %v", err)
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String("mastupeti-videos-store"),
		Key:         aws.String(fileHeader.Filename),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
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

// Below section details functions for handling our CloudFront URLs.

// LoadPrivateKey reads the PEM encoded private Key.
func LoadPrivateKey(fileName string) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile("private_key.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file, %v", err)
	}

	//fmt to be deleted
	//fmt.Println("Private Key file read successfully")

	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" && block.Type != "PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing Private Key")
	}
	//fmt.Println("Private Key decoded successfully")
	var privateKey *rsa.PrivateKey
	if block.Type == "PRIVATE KEY" {
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKCS8 private key, %v", err)
		}
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("not an RSA private key")
		}
	} else {
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKCS1 private key, %v", err)
		}
		privateKey = key
	}
	//fmt.Println("Private Key parsed successfully")
	return privateKey, nil

}

// generatecloudfronturl generates a signed URL for our clouddistribution.
func GenerateCloudFrontSignedURL(key, distributionDomain, keyPairID string, privateKey *rsa.PrivateKey, expiration time.Time) (string, error) {
	url := fmt.Sprintf("https://%s/%s", distributionDomain, key)
	signer := sign.NewURLSigner(keyPairID, privateKey)
	signedURL, err := signer.Sign(url, expiration)
	if err != nil {
		return "", fmt.Errorf("failed to sign URL, %v", err)
	}
	return signedURL, nil
}

// GenerateCloudFrontSignedURLs generates signed URLs for multiple CloudFront resources.
func GenerateCloudFrontSignedURLs(keys []string, distributionDomain, keyPairID string, privateKey *rsa.PrivateKey, expiration time.Time) (map[string]string, error) {
	urls := make(map[string]string)
	for _, key := range keys {
		signedURL, err := GenerateCloudFrontSignedURL(key, distributionDomain, keyPairID, privateKey, expiration)
		if err != nil {
			return nil, fmt.Errorf("failed to generate signed URL for key %s: %v", key, err)
		}
		urls[key] = signedURL
	}
	return urls, nil
}
