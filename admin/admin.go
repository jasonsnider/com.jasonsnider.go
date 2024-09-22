package admin

import (
	"log"
	"net/http"

	"github.com/boj/redistore"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jasonsnider/com.jasonsnider.go/pkg/auth"
	"github.com/jasonsnider/com.jasonsnider.go/pkg/cache"
)

type App struct {
	DB           *pgxpool.Pool
	BustCssCache string
	BustJsCache  string
	SessionStore *redistore.RediStore
	//SessionStore *sessions.CookieStore
}

func AdminRouter(dbpool *pgxpool.Pool) *mux.Router {

	store, err := redistore.NewRediStore(10, "tcp", "redis:6379", "", []byte("your-secret-key"))

	if err != nil {
		log.Fatalf("Failed to initialize Redis store: %v", err)
	}

	// Initialize middleware
	auth := &auth.AuthMiddleware{SessionStore: store}

	app := &App{
		DB:           dbpool,
		BustCssCache: cache.BustCssCache(),
		BustJsCache:  cache.BustJsCache(),
		SessionStore: store,
	}

	router := mux.NewRouter()

	router.HandleFunc("/login", app.Authenticate).Methods("GET")
	router.HandleFunc("/login", app.Authenticate).Methods("POST")

	protected := router.PathPrefix("/admin").Subrouter()
	protected.Use(auth.AuthRequired)

	protected.HandleFunc("/dashboard", app.Dashboard).Methods("GET")

	protected.HandleFunc("/register", app.RegisterUser).Methods("GET")
	protected.HandleFunc("/register", app.RegisterUser).Methods("POST")

	protected.HandleFunc("/users", app.ListUsers).Methods("GET")
	protected.HandleFunc("/users/{id}", app.ViewUser).Methods("GET")
	protected.HandleFunc("/users/{id}/edit", app.UpdateUser).Methods("GET")
	protected.HandleFunc("/users/{id}/edit", app.UpdateUser).Methods("POST")
	protected.HandleFunc("/users/{id}/delete", app.DeleteUser).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("404 Not Found: %s", r.URL.Path)
		http.NotFound(w, r)
	})

	return router
}
