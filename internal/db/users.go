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

	sql := "INSERT INTO users (id, first_name, last_name, email, hash) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := db.DB.Exec(context.Background(), sql, userID, user.FirstName, user.LastName, user.Email, hash)
	if err != nil {
		return fmt.Errorf("query failed: %v", err)
	}

	return nil
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

// EmailExistsInDatabase checks if an email exists in the database
func (db *DB) EmailExistsInDatabase(email string) bool {
	var foundEmail string
	sql := "SELECT email FROM users WHERE email=$1"

	err := db.DB.QueryRow(context.Background(), sql, email).Scan(&foundEmail)

	if err == pgx.ErrNoRows {
		return false
	} else if err != nil {
		log.Printf("Error querying database: %v", err)
		return false
	}
	return true
}

// GetExistingEmail fetches the existing email for a given user ID
func (db *DB) GetExistingEmail(userID string) (string, error) {
	var existingEmail string
	sql := "SELECT email FROM users WHERE id=$1"

	err := db.DB.QueryRow(context.Background(), sql, userID).Scan(&existingEmail)
	if err != nil {
		return "", err
	}
	return existingEmail, nil
}

// UniqueEmail is a custom validation function for unique email
func (db *DB) UniqueEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	userID := fl.Parent().FieldByName("ID").String()

	// Only check for uniqueness if the email has changed
	existingEmail, err := db.GetExistingEmail(userID)
	if err != nil {
		log.Printf("Error fetching existing email: %v", err)
		return false
	}

	if email != existingEmail {
		return !db.EmailExistsInDatabase(email)
	}
	return true
}
