package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jasonsnider/com.jasonsnider.go/internal/db"
	"github.com/jasonsnider/com.jasonsnider.go/internal/types"
)

func (app *App) GetUsers(w http.ResponseWriter, r *http.Request) {

	db := db.DB{DB: app.DB}
	users, err := db.FetchUsers()

	if err != nil {
		http.Error(w, fmt.Sprintf("FetchArticlesByType failed: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (app *App) GetUser(w http.ResponseWriter, r *http.Request) {
	db := db.DB{DB: app.DB}
	vars := mux.Vars(r)
	id := vars["id"]

	var user types.User
	user, err := db.FetchUserById(id)

	if err != nil {
		http.Error(w, fmt.Sprintf("FetchUserById failed: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Implementation for creating a user
}

func (app *App) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Implementation for updating a user
}

func (app *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Implementation for deleting a user
}
