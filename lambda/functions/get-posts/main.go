package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
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

	// Get query parameters
	userID := request.QueryStringParameters["user_id"]
	limit := request.QueryStringParameters["limit"]
	offset := request.QueryStringParameters["offset"]

	// Set defaults
	if limit == "" {
		limit = "20"
	}
	if offset == "" {
		offset = "0"
	}

	// Parse limit and offset
	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	// Build SQL query
	var query string
	var args []interface{}

	if userID != "" {
		// Get feed for specific user
		query = `
			SELECT p.id, p.user_id, p.title, p.description, p.total_likes, p.total_comments,
				   p.created_at, u.id, u.username, u.full_name, u.profile_photo_url, u.follower_count
			FROM posts p
			JOIN users u ON p.user_id = u.id
			WHERE p.user_id IN (
				SELECT following_id FROM follows WHERE follower_id = $1
			) OR u.follower_count > 1000
			ORDER BY p.created_at DESC
			LIMIT $2 OFFSET $3`
		args = []interface{}{userID, limitInt, offsetInt}
	} else {
		// Get all posts
		query = `
			SELECT p.id, p.user_id, p.title, p.description, p.total_likes, p.total_comments,
				   p.created_at, u.id, u.username, u.full_name, u.profile_photo_url, u.follower_count
			FROM posts p
			JOIN users u ON p.user_id = u.id
			ORDER BY p.created_at DESC
			LIMIT $1 OFFSET $2`
		args = []interface{}{limitInt, offsetInt}
	}

	// Execute query
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to fetch posts"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}
	defer rows.Close()

	// Process results
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
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"error": "Failed to scan post data"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
					"Access-Control-Allow-Origin": "*",
				},
			}, nil
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

func connectToDatabase() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	return sql.Open("postgres", connStr)
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
	// Recipe Social Media - Get Posts Handler
	lambda.Start(handler)
}
