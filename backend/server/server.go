package server

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Server struct {
	DB *sql.DB
}

// New initializes the DB and returns a Server instance.
func New(dbURL string) (*Server, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	fmt.Println("Connected to PostgreSQL")
	return &Server{DB: db}, nil
}
