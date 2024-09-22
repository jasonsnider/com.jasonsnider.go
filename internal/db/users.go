package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jasonsnider/com.jasonsnider.go/internal/types"
	"github.com/jasonsnider/com.jasonsnider.go/pkg/passwords"
)

func (db *DB) RegisterUser(user types.RegisterUser) error {

	userID := uuid.New()
	hash, _ := passwords.HashPassword(user.Password)
	admin := false

	sql := "INSERT INTO users (id, first_name, last_name, email, hash, admin) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := db.DB.Exec(context.Background(), sql, userID, user.FirstName, user.LastName, user.Email, hash, admin)
	if err != nil {
		return fmt.Errorf("query failed: %v", err)
	}

	return nil
}

func (db *DB) FetchUserById(id string) (types.User, error) {
	var user types.User
	err := db.DB.QueryRow(context.Background(), "SELECT id, first_name, last_name, email FROM users WHERE id=$1", id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		return user, fmt.Errorf("query failed: %v", err)
	}

	return user, nil
}
