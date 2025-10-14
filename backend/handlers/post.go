package handlers

import (
	"net/http"
	"strconv"

	"recipe-social-backend/database"
	"recipe-social-backend/models"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	db *database.DB
}

func NewPostHandler(db *database.DB) *PostHandler {
	return &PostHandler{db: db}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Create post
	post := &models.Post{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
	}
	
	// Convert request images to post images
	for _, imgReq := range req.Images {
		post.Images = append(post.Images, models.PostImage{
			ImageURL:        imgReq.ImageURL,
			StepDescription: imgReq.StepDescription,
			StepNumber:      imgReq.StepNumber,
		})
	}
	
	err := h.db.CreatePost(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}
	
	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetPost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	
	post, err := h.db.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	
	// Check if current user liked this post
	currentUserID := c.GetInt("user_id")
	if currentUserID != 0 {
		post.IsLiked, _ = h.db.IsLiked(currentUserID, postID)
	}
	
	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) GetFeed(c *gin.Context) {
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
	
	posts, err := h.db.GetFeed(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feed"})
		return
	}
	
	// Check which posts are liked by current user
	for i := range posts {
		posts[i].IsLiked, _ = h.db.IsLiked(userID, posts[i].ID)
	}
	
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	
	userID := c.GetInt("user_id")
	
	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get existing post
	post, err := h.db.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	
	// Check if user owns the post
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this post"})
		return
	}
	
	// Update fields if provided
	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.Description != nil {
		post.Description = req.Description
	}
	
	err = h.db.UpdatePost(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}
	
	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	
	userID := c.GetInt("user_id")
	
	// Get existing post
	post, err := h.db.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	
	// Check if user owns the post
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this post"})
		return
	}
	
	err = h.db.DeletePost(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func (h *PostHandler) LikePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	
	userID := c.GetInt("user_id")
	
	err = h.db.LikePost(userID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Post liked successfully"})
}

func (h *PostHandler) UnlikePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	
	userID := c.GetInt("user_id")
	
	err = h.db.UnlikePost(userID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike post"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Post unliked successfully"})
}

func (h *PostHandler) CreateComment(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	
	userID := c.GetInt("user_id")
	
	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	comment := &models.Comment{
		UserID: userID,
		PostID: postID,
		Content: req.Content,
	}
	
	err = h.db.CreateComment(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}
	
	// Get user info for the comment
	user, err := h.db.GetUserByID(userID)
	if err == nil {
		comment.User = user
	}
	
	c.JSON(http.StatusCreated, comment)
}

func (h *PostHandler) GetComments(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
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
	
	comments, err := h.db.GetComments(postID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}
	
	c.JSON(http.StatusOK, comments)
}

func (h *PostHandler) DeleteComment(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}
	
	userID := c.GetInt("user_id")
	
	// TODO: Add authorization check to ensure user owns the comment
	
	err = h.db.DeleteComment(commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

func (h *PostHandler) SearchPosts(c *gin.Context) {
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
	
	posts, err := h.db.SearchPosts(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search posts"})
		return
	}
	
	c.JSON(http.StatusOK, posts)
}
