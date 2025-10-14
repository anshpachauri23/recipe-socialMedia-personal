package handlers

import (
	"net/http"
	"strconv"

	"recipe-social-backend/database"
	"recipe-social-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db *database.DB
}

func NewUserHandler(db *database.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	user, err := h.db.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	user, err := h.db.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	// Update fields if provided
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}
	if req.ProfilePhotoURL != nil {
		user.ProfilePhotoURL = req.ProfilePhotoURL
	}
	
	err = h.db.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}
	
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	err := h.db.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	user, err := h.db.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	// Check if current user is following this user
	currentUserID := c.GetInt("user_id")
	isFollowing := false
	if currentUserID != 0 {
		isFollowing, _ = h.db.IsFollowing(currentUserID, userID)
	}
	
	profile := models.UserProfile{
		User:         *user,
		IsFollowing:  isFollowing,
		IsOwnProfile: currentUserID == userID,
	}
	
	c.JSON(http.StatusOK, profile)
}

func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query required"})
		return
	}
	
	limit := 20
	offset := 0
	
	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil {
			limit = parsedLimit
		}
	}
	
	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil {
			offset = parsedOffset
		}
	}
	
	users, err := h.db.SearchUsers(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}
	
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) FollowUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	currentUserID := c.GetInt("user_id")
	if currentUserID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
		return
	}
	
	err = h.db.FollowUser(currentUserID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "User followed successfully"})
}

func (h *UserHandler) UnfollowUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	currentUserID := c.GetInt("user_id")
	err = h.db.UnfollowUser(currentUserID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "User unfollowed successfully"})
}

func (h *UserHandler) GetFollowing(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	limit := 20
	offset := 0
	
	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil {
			limit = parsedLimit
		}
	}
	
	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil {
			offset = parsedOffset
		}
	}
	
	users, err := h.db.GetFollowing(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get following"})
		return
	}
	
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetFollowers(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	limit := 20
	offset := 0
	
	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil {
			limit = parsedLimit
		}
	}
	
	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil {
			offset = parsedOffset
		}
	}
	
	users, err := h.db.GetFollowers(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get followers"})
		return
	}
	
	c.JSON(http.StatusOK, users)
}
