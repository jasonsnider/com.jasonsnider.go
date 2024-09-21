package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jasonsnider/go.jasonsnider.com/admin"
	"github.com/jasonsnider/go.jasonsnider.com/api/v1"
	"github.com/jasonsnider/go.jasonsnider.com/pkg/passwords"
	"github.com/jasonsnider/go.jasonsnider.com/web"
	"github.com/joho/godotenv"
)

func main() {

	mode := flag.String("mode", "server", "Mode of operation: server, hash, check")
	password := flag.String("password", "", "The password to hash or check")
	hashValue := flag.String("hashvalue", "", "The hash to check the password against")
	flag.Parse()

	switch *mode {
	case "server":
		if err := runServer(); err != nil {
			log.Fatalf("Application error: %v", err)
		}
	case "hash":
		if *password == "" {
			log.Fatal("Please provide a password using the -password flag")
		}
		hashPassword(*password)
	case "check":
		if *password == "" || *hashValue == "" {
			log.Fatal("Please provide both -password and -hashvalue flags")
		}
		checkPassword(*password, *hashValue)
	default:
		log.Fatalf("Unknown mode: %s", *mode)
	}
}

func runServer() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbUrl := os.Getenv("DATABASE_URL")

	dbpool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	apiRouter := api.APIRouter(dbpool)
	webRouter := web.WebRouter(dbpool)
	adminRouter := admin.AdminRouter(dbpool)

	mainRouter := mux.NewRouter()
	mainRouter.PathPrefix("/api/v1/").Handler(http.StripPrefix("/api/v1", apiRouter))
	mainRouter.PathPrefix("/admin/").Handler(adminRouter)
	mainRouter.PathPrefix("/login").Handler(adminRouter)
	mainRouter.PathPrefix("/").Handler(webRouter)

	log.Fatal(http.ListenAndServe(":8080", mainRouter))

	return nil
}

func hashPassword(password string) {
	hash, err := passwords.HashPassword(password)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}
	fmt.Printf("Hashed password: %s\n", hash)
}

func checkPassword(password, hashValue string) {
	match := passwords.CheckPasswordHash(password, hashValue)
	if match {
		fmt.Println("Password matches the hash")
	} else {
		fmt.Println("Password does not match the hash")
	}
}
