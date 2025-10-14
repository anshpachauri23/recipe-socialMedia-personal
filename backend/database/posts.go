package database

import (
	"database/sql"
	"recipe-social-backend/models"
)

func (db *DB) CreatePost(post *models.Post) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert post
	query := `
		INSERT INTO posts (user_id, title, description)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`
	
	err = tx.QueryRow(query, post.UserID, post.Title, post.Description).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}

	// Insert post images
	for _, image := range post.Images {
		imageQuery := `
			INSERT INTO post_images (post_id, image_url, step_description, step_number)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at`
		
		err = tx.QueryRow(imageQuery, post.ID, image.ImageURL, image.StepDescription, image.StepNumber).
			Scan(&image.ID, &image.CreatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (db *DB) GetPost(id int) (*models.Post, error) {
	post := &models.Post{}
	query := `
		SELECT p.id, p.user_id, p.title, p.description, p.total_likes, p.total_comments, 
			   p.created_at, p.updated_at
		FROM posts p
		WHERE p.id = $1`
	
	err := db.QueryRow(query, id).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Description,
		&post.TotalLikes, &post.TotalComments, &post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get post images
	images, err := db.GetPostImages(post.ID)
	if err != nil {
		return nil, err
	}
	post.Images = images

	// Get user info
	user, err := db.GetUserByID(post.UserID)
	if err != nil {
		return nil, err
	}
	post.User = user

	return post, nil
}

func (db *DB) GetPostImages(postID int) ([]models.PostImage, error) {
	query := `
		SELECT id, post_id, image_url, step_description, step_number, created_at
		FROM post_images
		WHERE post_id = $1
		ORDER BY step_number`
	
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var images []models.PostImage
	for rows.Next() {
		var image models.PostImage
		err := rows.Scan(&image.ID, &image.PostID, &image.ImageURL,
			&image.StepDescription, &image.StepNumber, &image.CreatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	
	return images, nil
}

func (db *DB) GetFeed(userID int, limit, offset int) ([]models.FeedPost, error) {
	query := `
		SELECT p.id, p.user_id, p.title, p.description, p.total_likes, p.total_comments,
			   p.created_at, p.updated_at,
			   u.id, u.username, u.full_name, u.profile_photo_url, u.follower_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id IN (
			SELECT following_id FROM follows WHERE follower_id = $1
		) OR u.follower_count > 1000
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3`
	
	rows, err := db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var posts []models.FeedPost
	for rows.Next() {
		var post models.FeedPost
		var user models.User
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Description,
			&post.TotalLikes, &post.TotalComments, &post.CreatedAt, &post.UpdatedAt,
			&user.ID, &user.Username, &user.FullName, &user.ProfilePhotoURL, &user.FollowerCount,
		)
		if err != nil {
			return nil, err
		}
		post.User = user
		posts = append(posts, post)
	}
	
	return posts, nil
}

func (db *DB) SearchPosts(query string, limit, offset int) ([]models.Post, error) {
	searchQuery := `
		SELECT p.id, p.user_id, p.title, p.description, p.total_likes, p.total_comments,
			   p.created_at, p.updated_at,
			   u.id, u.username, u.full_name, u.profile_photo_url
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.title ILIKE $1 OR p.description ILIKE $1
		ORDER BY p.total_likes DESC, p.created_at DESC
		LIMIT $2 OFFSET $3`
	
	rows, err := db.Query(searchQuery, "%"+query+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var user models.User
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Description,
			&post.TotalLikes, &post.TotalComments, &post.CreatedAt, &post.UpdatedAt,
			&user.ID, &user.Username, &user.FullName, &user.ProfilePhotoURL,
		)
		if err != nil {
			return nil, err
		}
		post.User = &user
		posts = append(posts, post)
	}
	
	return posts, nil
}

func (db *DB) UpdatePost(post *models.Post) error {
	query := `
		UPDATE posts 
		SET title = $1, description = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3`
	
	_, err := db.Exec(query, post.Title, post.Description, post.ID)
	return err
}

func (db *DB) DeletePost(id int) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func (db *DB) LikePost(userID, postID int) error {
	query := `INSERT INTO likes (user_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := db.Exec(query, userID, postID)
	return err
}

func (db *DB) UnlikePost(userID, postID int) error {
	query := `DELETE FROM likes WHERE user_id = $1 AND post_id = $2`
	_, err := db.Exec(query, userID, postID)
	return err
}

func (db *DB) IsLiked(userID, postID int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM likes WHERE user_id = $1 AND post_id = $2`
	err := db.QueryRow(query, userID, postID).Scan(&count)
	return count > 0, err
}

func (db *DB) CreateComment(comment *models.Comment) error {
	query := `
		INSERT INTO comments (user_id, post_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`
	
	return db.QueryRow(query, comment.UserID, comment.PostID, comment.Content).
		Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
}

func (db *DB) GetComments(postID int, limit, offset int) ([]models.Comment, error) {
	query := `
		SELECT c.id, c.user_id, c.post_id, c.content, c.created_at, c.updated_at,
			   u.id, u.username, u.full_name, u.profile_photo_url
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3`
	
	rows, err := db.Query(query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var user models.User
		err := rows.Scan(
			&comment.ID, &comment.UserID, &comment.PostID, &comment.Content,
			&comment.CreatedAt, &comment.UpdatedAt,
			&user.ID, &user.Username, &user.FullName, &user.ProfilePhotoURL,
		)
		if err != nil {
			return nil, err
		}
		comment.User = &user
		comments = append(comments, comment)
	}
	
	return comments, nil
}

func (db *DB) DeleteComment(id int) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}
