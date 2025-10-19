package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
)

type Post struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	TotalLikes   int    `json:"total_likes"`
	TotalComments int   `json:"total_comments"`
	CreatedAt    string `json:"created_at"`
	User         User   `json:"user"`
	Images       []PostImage `json:"images"`
	IsLiked      bool   `json:"is_liked"`
}

type User struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	FullName       string `json:"full_name"`
	ProfilePhotoURL string `json:"profile_photo_url"`
	FollowerCount  int    `json:"follower_count"`
}

type PostImage struct {
	ID              int    `json:"id"`
	ImageURL        string `json:"image_url"`
	StepDescription string `json:"step_description"`
	StepNumber      int    `json:"step_number"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Handle CORS preflight
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
		}, nil
	}

	// Extract user ID from JWT token
	userID, err := extractUserIDFromToken(request.Headers["Authorization"])
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       `{"error": "Invalid token"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Get query parameters
	limit := request.QueryStringParameters["limit"]
	offset := request.QueryStringParameters["offset"]

	// Set defaults
	if limit == "" {
		limit = "20"
	}
	if offset == "" {
		offset = "0"
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

	// Get feed posts
	posts, err := getFeedPosts(db, userID, limit, offset)
	if err != nil {
		log.Printf("Error getting feed posts: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to fetch feed"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Convert to JSON
	responseBody, err := json.Marshal(posts)
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
		StatusCode: 200,
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

	// Extract token from "Bearer <token>"
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

func getFeedPosts(db *sql.DB, userID int, limit, offset string) ([]Post, error) {
	query := `
		SELECT p.id, p.user_id, p.title, p.description, p.total_likes, p.total_comments,
			   p.created_at, u.id, u.username, u.full_name, u.profile_photo_url, u.follower_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id IN (
			SELECT following_id FROM follows WHERE follower_id = $1
		) OR u.follower_count > 1000
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3`

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	rows, err := db.Query(query, userID, limitInt, offsetInt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		post := Post{}
		user := User{}
		
		// Use sql.NullString for nullable fields
		var profilePhotoURL sql.NullString
		
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Description,
			&post.TotalLikes, &post.TotalComments, &post.CreatedAt,
			&user.ID, &user.Username, &user.FullName, &profilePhotoURL, &user.FollowerCount,
		)
		if err != nil {
			return nil, err
		}
		
		// Handle NULL profile photo URL
		if profilePhotoURL.Valid {
			user.ProfilePhotoURL = profilePhotoURL.String
		} else {
			user.ProfilePhotoURL = ""
		}
		
		post.User = user
		posts = append(posts, post)
	}

	// Get images for each post
	for i := range posts {
		images, err := getPostImages(db, posts[i].ID)
		if err == nil {
			posts[i].Images = images
		}
	}

	return posts, nil
}

func getPostImages(db *sql.DB, postID int) ([]PostImage, error) {
	query := `
		SELECT id, image_url, step_description, step_number
		FROM post_images
		WHERE post_id = $1
		ORDER BY step_number`

	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []PostImage
	for rows.Next() {
		image := PostImage{}
		err := rows.Scan(&image.ID, &image.ImageURL, &image.StepDescription, &image.StepNumber)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}

func main() {
	// Recipe Social Media - User Feed Handler
	lambda.Start(handler)
}
