package models

import (
	"time"
)

type Post struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	Title        string    `json:"title" db:"title"`
	Description  *string   `json:"description" db:"description"`
	TotalLikes   int       `json:"total_likes" db:"total_likes"`
	TotalComments int      `json:"total_comments" db:"total_comments"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	User         *User     `json:"user,omitempty"`
	Images       []PostImage `json:"images,omitempty"`
	IsLiked      bool      `json:"is_liked"`
}

type PostImage struct {
	ID              int    `json:"id" db:"id"`
	PostID          int    `json:"post_id" db:"post_id"`
	ImageURL        string `json:"image_url" db:"image_url"`
	StepDescription *string `json:"step_description" db:"step_description"`
	StepNumber      int    `json:"step_number" db:"step_number"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type CreatePostRequest struct {
	Title       string                `json:"title" binding:"required,max=200"`
	Description *string               `json:"description"`
	Images      []CreateImageRequest  `json:"images" binding:"required,min=1"`
}

type CreateImageRequest struct {
	ImageURL        string  `json:"image_url" binding:"required"`
	StepDescription *string `json:"step_description"`
	StepNumber      int     `json:"step_number" binding:"required"`
}

type UpdatePostRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type Comment struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	PostID    int       `json:"post_id" db:"post_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	User      *User     `json:"user,omitempty"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,max=500"`
}

type FeedPost struct {
	Post
	User User `json:"user"`
}
