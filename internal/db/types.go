package db

import "github.com/jackc/pgx/v5/pgxpool"

type DB struct {
	DB *pgxpool.Pool
}
