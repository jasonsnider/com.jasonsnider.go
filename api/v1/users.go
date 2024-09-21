package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (app *App) GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.Query(context.Background(), "SELECT id, username, email, first_name, last_name FROM users")
	if err != nil {
		http.Error(w, fmt.Sprintf("Query failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Row scan failed: %v", err), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		http.Error(w, fmt.Sprintf("Rows iteration failed: %v", rows.Err()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (app *App) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user User
	err := app.DB.QueryRow(context.Background(), "SELECT id, username, email, first_name, last_name FROM users WHERE id=$1", id).Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query failed: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Implementation for creating a user
}
