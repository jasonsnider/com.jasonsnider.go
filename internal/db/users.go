package db

import (
	"context"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jasonsnider/com.jasonsnider.go/internal/types"
	"github.com/jasonsnider/com.jasonsnider.go/pkg/passwords"
)

func (db *DB) RegisterUser(user types.RegisterUser) error {

	userID := uuid.New()
	hash, _ := passwords.HashPassword(user.Password)

	sql := "INSERT INTO users (id, first_name, last_name, email, hash) VALUES ($1, $2, $3, $4, $5)"
	_, err := db.DB.Exec(context.Background(), sql, userID, user.FirstName, user.LastName, user.Email, hash)
	if err != nil {
		return fmt.Errorf("query failed: %v", err)
	}

	return nil
}

func (db *DB) CreateUser(user types.User) (string, error) {

	userID := uuid.New().String()

	sql := "INSERT INTO users (id, first_name, last_name, email, role) VALUES ($1, $2, $3, $4, $5)"
	_, err := db.DB.Exec(context.Background(), sql, userID, user.FirstName, user.LastName, user.Email, user.Role)
	if err != nil {
		return "", fmt.Errorf("query failed: %v", err)
	}

	return userID, nil
}

func (db *DB) FetchUsers() ([]types.User, error) {
	sql := "SELECT id, first_name, last_name, email, role FROM users"
	rows, err := db.DB.Query(context.Background(), sql)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var users []types.User
	for rows.Next() {
		var user types.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role)
		if err != nil {
			return nil, fmt.Errorf("row scan failed: %v", err)
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration failed: %v", rows.Err())
	}

	return users, nil
}

func (db *DB) FetchUserById(id string) (types.User, error) {
	var user types.User
	sql := "SELECT id, first_name, last_name, email, role FROM users WHERE id=$1"
	err := db.DB.QueryRow(context.Background(), sql, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role)
	if err != nil {
		return user, fmt.Errorf("query failed: %v", err)
	}

	return user, nil
}

func (db *DB) DeleteUser(id string) error {
	sql := "DELETE FROM users WHERE id=$1"
	_, err := db.DB.Exec(context.Background(), sql, id)
	if err != nil {
		return fmt.Errorf("query failed: %v", err)
	}

	return nil
}

// fetchEmail is a helper function that queries the database for an email by a given column and value.
func (db *DB) fetchEmail(column, value string) (string, error) {
	var email string
	query := fmt.Sprintf("SELECT email FROM users WHERE %s=$1", column)
	err := db.DB.QueryRow(context.Background(), query, value).Scan(&email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil // No rows means no email, return empty string and no error
		}
		return "", err // Return other errors to be handled by the caller
	}
	return email, nil
}

// EmailExistsInDatabase checks if an email exists in the database
func (db *DB) EmailExistsInDatabase(email string) bool {
	foundEmail, err := db.fetchEmail("email", email)
	if err != nil {
		log.Printf("Error querying database for email: %v", err)
		return false
	}
	return foundEmail != ""
}

// GetExistingEmail fetches the existing email for a given user ID
func (db *DB) GetExistingEmail(userID string) (string, error) {
	return db.fetchEmail("id", userID)
}

// UniqueEmail is a custom validation function for unique email
func (db *DB) UniqueEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	userID := fl.Parent().FieldByName("ID").String()

	// If userID is provided, check if the email has changed
	if userID != "" {
		existingEmail, err := db.GetExistingEmail(userID)
		if err != nil {
			log.Printf("Error fetching existing email: %v", err)
			return false
		}
		if email == existingEmail {
			return true // Email is unchanged
		}
	}

	// Check if the email exists in the database
	return !db.EmailExistsInDatabase(email)
}
