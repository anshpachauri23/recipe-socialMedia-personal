package models

import (
	"time"
)

type User struct {
	ID             int       `json:"id" db:"id"`
	Username       string    `json:"username" db:"username"`
	Email          string    `json:"email" db:"email"`
	PasswordHash   string    `json:"-" db:"password_hash"`
	FullName       string    `json:"full_name" db:"full_name"`
	Bio            *string   `json:"bio" db:"bio"`
	ProfilePhotoURL *string  `json:"profile_photo_url" db:"profile_photo_url"`
	IsPublic       bool      `json:"is_public" db:"is_public"`
	FollowerCount  int       `json:"follower_count" db:"follower_count"`
	FollowingCount int       `json:"following_count" db:"following_count"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type UserProfile struct {
	User
	IsFollowing bool `json:"is_following"`
	IsOwnProfile bool `json:"is_own_profile"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required,min=1,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	FullName       *string `json:"full_name"`
	Bio            *string `json:"bio"`
	ProfilePhotoURL *string `json:"profile_photo_url"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}
