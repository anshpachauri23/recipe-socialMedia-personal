package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
)

type CreatePostRequest struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Images      []ImageRequest `json:"images"`
}

type ImageRequest struct {
	ImageData       string `json:"image_data"`       // Base64 encoded image
	ImageURL        string `json:"image_url"`        // Will be generated after S3 upload
	StepDescription string `json:"step_description"`
	StepNumber      int    `json:"step_number"`
}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extract user ID from JWT token
	userID, err := extractUserIDFromToken(request.Headers["Authorization"])
	if err != nil {
		log.Printf("Error extracting user ID: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       `{"error": "Invalid token"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Connect to database
	db, err := connectToDatabase()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Database connection failed"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}
	defer db.Close()

	// Parse request body
	var req CreatePostRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Invalid request body"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Validate request
	if req.Title == "" || len(req.Images) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Title and images are required"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to begin transaction"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Insert post
	insertPostSQL := `
		INSERT INTO posts (user_id, title, description)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	var postID int64
	var createdAt, updatedAt string
	err = tx.QueryRow(insertPostSQL, userID, req.Title, req.Description).Scan(&postID, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("Error inserting post: %v", err)
		tx.Rollback()
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to create post"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Upload images to S3 and insert into database
	for _, image := range req.Images {
		// Upload image to S3 if image data is provided
		var imageURL string
		if image.ImageData != "" {
			// Generate filename
			filename := fmt.Sprintf("step-%d.jpg", image.StepNumber)
			
			// Upload to S3
			s3URL, err := uploadImageToS3(image.ImageData, filename)
			if err != nil {
				tx.Rollback()
				return events.APIGatewayProxyResponse{
					StatusCode: 500,
					Body:       fmt.Sprintf(`{"error": "Failed to upload image: %s"}`, err.Error()),
					Headers: map[string]string{
						"Content-Type": "application/json",
						"Access-Control-Allow-Origin": "*",
					},
				}, nil
			}
			imageURL = s3URL
		} else {
			// Use provided URL if no image data
			imageURL = image.ImageURL
		}

		// Insert image record
		insertImageSQL := `
			INSERT INTO post_images (post_id, image_url, step_description, step_number)
			VALUES ($1, $2, $3, $4)`

		_, err := tx.Exec(insertImageSQL, postID, imageURL, image.StepDescription, image.StepNumber)
		if err != nil {
			log.Printf("Error inserting image: %v", err)
			tx.Rollback()
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"error": "Failed to create post images"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
					"Access-Control-Allow-Origin": "*",
				},
			}, nil
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to commit transaction"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Return success response
	response := map[string]interface{}{
		"message": "Post created successfully",
		"post_id": postID,
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to process response"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func extractUserIDFromToken(authHeader string) (int, error) {
	if authHeader == "" {
		return 0, fmt.Errorf("no authorization header")
	}

	// Remove "Bearer " prefix
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return 0, fmt.Errorf("invalid authorization header format")
	}

	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(float64); ok {
			return int(userID), nil
		}
		return 0, fmt.Errorf("user_id not found in token")
	}

	return 0, fmt.Errorf("invalid token")
}

func connectToDatabase() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	return sql.Open("postgres", connStr)
}

func uploadImageToS3(imageData, filename string) (string, error) {
	// Decode base64 image data
	imageBytes, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 image: %v", err)
	}

	// Create S3 session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create S3 session: %v", err)
	}

	// Create S3 service
	s3Client := s3.New(sess)

	// Generate unique filename
	timestamp := time.Now().Unix()
	key := fmt.Sprintf("recipe-images/%d-%s", timestamp, filename)

	// Upload to S3
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("S3_BUCKET")),
		Key:         aws.String(key),
		Body:        strings.NewReader(string(imageBytes)),
		ContentType: aws.String("image/jpeg"),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %v", err)
	}

	// Return public URL
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", 
		os.Getenv("S3_BUCKET"), 
		os.Getenv("AWS_REGION"), 
		key), nil
}

func main() {
	// Recipe Social Media - Create Post Handler
	lambda.Start(handler)
}
