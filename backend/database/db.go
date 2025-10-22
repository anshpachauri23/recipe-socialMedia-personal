package database

import (
	"database/sql"
	"recipe-social-backend/models"
)

type DB struct {
	*sql.DB
}

func (db *DB) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, full_name, is_public)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	
	return db.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.FullName, user.IsPublic).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (db *DB) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password_hash, full_name, bio, profile_photo_url, 
			  is_public, follower_count, following_count, created_at, updated_at 
			  FROM users WHERE email = $1`
	
	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Bio, &user.ProfilePhotoURL, &user.IsPublic, &user.FollowerCount,
		&user.FollowingCount, &user.CreatedAt, &user.UpdatedAt,
	)
	
	return user, err
}

func (db *DB) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password_hash, full_name, bio, profile_photo_url, 
			  is_public, follower_count, following_count, created_at, updated_at 
			  FROM users WHERE id = $1`
	
	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Bio, &user.ProfilePhotoURL, &user.IsPublic, &user.FollowerCount,
		&user.FollowingCount, &user.CreatedAt, &user.UpdatedAt,
	)
	
	return user, err
}

func (db *DB) UpdateUser(user *models.User) error {
	query := `
		UPDATE users 
		SET full_name = $1, bio = $2, profile_photo_url = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4`
	
	_, err := db.Exec(query, user.FullName, user.Bio, user.ProfilePhotoURL, user.ID)
	return err
}

func (db *DB) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func (db *DB) SearchUsers(query string, limit, offset int) ([]models.User, error) {
	searchQuery := `
		SELECT id, username, full_name, bio, profile_photo_url, follower_count, following_count
		FROM users 
		WHERE username ILIKE $1 OR full_name ILIKE $1
		ORDER BY follower_count DESC
		LIMIT $2 OFFSET $3`
	
	rows, err := db.Query(searchQuery, "%"+query+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.Bio,
			&user.ProfilePhotoURL, &user.FollowerCount, &user.FollowingCount)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	
	return users, nil
}

func (db *DB) FollowUser(followerID, followingID int) error {
	query := `INSERT INTO follows (follower_id, following_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := db.Exec(query, followerID, followingID)
	return err
}

func (db *DB) UnfollowUser(followerID, followingID int) error {
	query := `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2`
	_, err := db.Exec(query, followerID, followingID)
	return err
}

func (db *DB) IsFollowing(followerID, followingID int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM follows WHERE follower_id = $1 AND following_id = $2`
	err := db.QueryRow(query, followerID, followingID).Scan(&count)
	return count > 0, err
}

func (db *DB) GetFollowing(userID int, limit, offset int) ([]models.User, error) {
	query := `
		SELECT u.id, u.username, u.full_name, u.bio, u.profile_photo_url, 
			   u.follower_count, u.following_count
		FROM users u
		JOIN follows f ON u.id = f.following_id
		WHERE f.follower_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3`
	
	rows, err := db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.Bio,
			&user.ProfilePhotoURL, &user.FollowerCount, &user.FollowingCount)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	
	return users, nil
}

func (db *DB) GetFollowers(userID int, limit, offset int) ([]models.User, error) {
	query := `
		SELECT u.id, u.username, u.full_name, u.bio, u.profile_photo_url, 
			   u.follower_count, u.following_count
		FROM users u
		JOIN follows f ON u.id = f.follower_id
		WHERE f.following_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3`
	
	rows, err := db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.Bio,
			&user.ProfilePhotoURL, &user.FollowerCount, &user.FollowingCount)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	
	return users, nil
}
