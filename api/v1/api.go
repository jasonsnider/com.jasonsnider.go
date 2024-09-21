package api

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	DB *pgxpool.Pool
}

func APIRouter(dbpool *pgxpool.Pool) *mux.Router {
	app := &App{DB: dbpool}

	router := mux.NewRouter()
	router.HandleFunc("/users", app.GetUsers).Methods("GET")
	router.HandleFunc("/user/{id}", app.GetUser).Methods("GET")
	router.HandleFunc("/users", app.CreateUser).Methods("POST")

	return router
}
