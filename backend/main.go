package main

import (
	"log"
	"os"

	"recipe-social-backend/database"
	"recipe-social-backend/handlers"
	"recipe-social-backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	sqlDB, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer sqlDB.Close()
	
	// Wrap in custom DB type
	db := &database.DB{sqlDB}

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	userHandler := handlers.NewUserHandler(db)
	postHandler := handlers.NewPostHandler(db)

	// Public routes
	api := r.Group("/api")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.GET("/users/search", userHandler.SearchUsers)
		api.GET("/posts/search", postHandler.SearchPosts)
		api.GET("/posts/:id", postHandler.GetPost)
		api.GET("/users/:id", userHandler.GetUserProfile)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/users/me", userHandler.GetCurrentUser)
		protected.PUT("/users/me", userHandler.UpdateProfile)
		protected.DELETE("/users/me", userHandler.DeleteAccount)
		protected.POST("/users/:id/follow", userHandler.FollowUser)
		protected.DELETE("/users/:id/follow", userHandler.UnfollowUser)
		protected.GET("/users/me/following", userHandler.GetFollowing)
		protected.GET("/users/me/followers", userHandler.GetFollowers)

		// Post routes
		protected.GET("/posts/feed", postHandler.GetFeed)
		protected.POST("/posts", postHandler.CreatePost)
		protected.PUT("/posts/:id", postHandler.UpdatePost)
		protected.DELETE("/posts/:id", postHandler.DeletePost)
		protected.POST("/posts/:id/like", postHandler.LikePost)
		protected.DELETE("/posts/:id/like", postHandler.UnlikePost)
		protected.POST("/posts/:id/comments", postHandler.CreateComment)
		protected.GET("/posts/:id/comments", postHandler.GetComments)
		protected.DELETE("/comments/:id", postHandler.DeleteComment)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
