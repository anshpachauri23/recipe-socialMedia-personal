package services

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Service struct {
	client *s3.S3
	bucket string
	region string
}

func NewS3Service() *S3Service {
	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("S3_REGION")),
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create AWS session: %v", err))
	}

	return &S3Service{
		client: s3.New(sess),
		bucket: os.Getenv("S3_BUCKET"),
		region: os.Getenv("S3_REGION"),
	}
}

func (s *S3Service) UploadImage(imageData, filename string) (string, error) {
	var base64Data string

	// Check if it's a data URL format (e.g., "data:image/jpeg;base64,...")
	if strings.Contains(imageData, ",") {
		// Extract base64 data from data URL
		parts := strings.SplitN(imageData, ",", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid image data format")
		}
		base64Data = parts[1]
	} else {
		// It's already just the base64 data
		base64Data = imageData
	}

	// Decode base64 image data
	imageBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 image: %v", err)
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	key := fmt.Sprintf("recipe-images/%d-%s", timestamp, filename)

	// Upload to S3
	_, err = s.client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        strings.NewReader(string(imageBytes)),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %v", err)
	}

	// Return public URL
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		s.bucket,
		s.region,
		key), nil
}
