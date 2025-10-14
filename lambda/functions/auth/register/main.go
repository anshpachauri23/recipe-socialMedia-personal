package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type User struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Email         string `json:"email"`
	FullName      string `json:"full_name"`
	Bio           *string `json:"bio"`
	ProfilePhotoURL *string `json:"profile_photo_url"`
	FollowerCount int    `json:"follower_count"`
	FollowingCount int   `json:"following_count"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	User    User   `json:"user"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Handle CORS preflight
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "POST, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
		}, nil
	}

	// Parse request body
	var req RegisterRequest
	if request.Body == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Request body is empty"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf(`{"error": "Invalid JSON: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Validate request
	if req.Username == "" || req.Email == "" || req.Password == "" || req.FullName == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "All fields are required"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to hash password"}`,
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

	// Create user
	user, err := createUser(db, req.Username, req.Email, string(hashedPassword), req.FullName)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to create user"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username":  user.Username,
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to generate token"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Prepare response
	response := RegisterResponse{
		Message: "User created successfully",
		Token:   tokenString,
		User: User{
			ID:             user.ID,
			Username:       user.Username,
			Email:          user.Email,
			FullName:       user.FullName,
			Bio:            user.Bio,
			ProfilePhotoURL: user.ProfilePhotoURL,
			FollowerCount:  user.FollowerCount,
			FollowingCount: user.FollowingCount,
		},
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
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

func connectToDatabase() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	return sql.Open("postgres", connStr)
}

func createUser(db *sql.DB, username, email, passwordHash, fullName string) (*User, error) {
	query := `INSERT INTO users (username, email, password_hash, full_name, is_public)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id, username, email, full_name, bio, profile_photo_url, 
						follower_count, following_count`
	
	user := &User{}
	err := db.QueryRow(query, username, email, passwordHash, fullName, true).Scan(
		&user.ID, &user.Username, &user.Email, &user.FullName,
		&user.Bio, &user.ProfilePhotoURL, &user.FollowerCount, &user.FollowingCount,
	)
	
	return user, err
}

func main() {
	// Recipe Social Media - User Registration Handler
	lambda.Start(handler)
}
