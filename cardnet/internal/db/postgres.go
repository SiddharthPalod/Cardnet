package db

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgres() (*sql.DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://cardnet:cardnet@localhost:5432/cardnet?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Set connection pool limits (prevent exhaustion)
	// With async writes, we need fewer connections
	db.SetMaxOpenConns(100)  // Max 100 connections
	db.SetMaxIdleConns(20)   // Keep 20 idle (increased for better connection reuse)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute) // Close idle connections after 1 minute

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
