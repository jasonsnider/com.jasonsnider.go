package web

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jasonsnider/go.jasonsnider.com/pkg/cache"
)

type App struct {
	DB           *pgxpool.Pool
	BustCssCache string
	BustJsCache  string
}

func WebRouter(dbpool *pgxpool.Pool) *mux.Router {
	app := &App{
		DB:           dbpool,
		BustCssCache: cache.BustCssCache(),
		BustJsCache:  cache.BustJsCache(),
	}

	router := mux.NewRouter()
	router.HandleFunc("/", app.Home).Methods("GET")
	router.HandleFunc("/articles", app.ListArticles).Methods("GET")
	router.HandleFunc("/articles/{slug}", app.ViewArticle).Methods("GET")

	router.HandleFunc("/games", app.ListGames).Methods("GET")
	router.HandleFunc("/games/{slug}", app.ViewGame).Methods("GET")

	router.HandleFunc("/tools", app.ListTools).Methods("GET")
	router.HandleFunc("/tools/{slug}", app.ViewTool).Methods("GET")

	router.HandleFunc("/contact", app.Contact).Methods("GET")
	router.HandleFunc("/contact", app.Contact).Methods("POST")

	return router
}
