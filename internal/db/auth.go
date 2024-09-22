package db

import (
	"context"
	"fmt"

	"github.com/jasonsnider/com.jasonsnider.go/internal/types"
)

func (db *DB) FetchAuth(email string) (types.AuthUser, error) {

	var user types.AuthUser
	sql := "SELECT id, first_name, last_name, email, hash FROM users WHERE email=$1"

	err := db.DB.QueryRow(context.Background(), sql, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Hash)
	if err != nil {
		return user, fmt.Errorf("query failed: %v", err)
	}

	return user, nil
}
