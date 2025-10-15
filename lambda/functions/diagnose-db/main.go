package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
)

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

	// Debug JWT token parsing
	debugInfo := map[string]interface{}{
		"authorization_header": request.Headers["Authorization"],
		"jwt_secret_set":       os.Getenv("JWT_SECRET") != "",
		"jwt_secret_length":    len(os.Getenv("JWT_SECRET")),
		"headers_count":        len(request.Headers),
		"all_headers":          request.Headers,
	}

	// Extract user ID from JWT token
	userID, err := extractUserIDFromToken(request.Headers["Authorization"])
	if err != nil {
		debugInfo["token_error"] = err.Error()
		debugInfo["user_id"] = 0
	} else {
		debugInfo["user_id"] = userID
		debugInfo["token_error"] = "none"
	}

	responseBody, _ := json.Marshal(debugInfo)

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

	// Extract token from "Bearer <token>" (handle leading/trailing spaces)
	tokenString := strings.TrimSpace(authHeader)
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	if tokenString == "" || tokenString == strings.TrimSpace(authHeader) {
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

func main() {
	lambda.Start(handler)
}
