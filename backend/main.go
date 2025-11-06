package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Server struct {
	db *sql.DB
}

func main() {
	// Load env vars
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println("DATABASE_URL not set")
		os.Exit(1)
	}

	// Connect to Postgres
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		panic(fmt.Errorf("failed to connect to DB: %w", err))
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		panic(fmt.Errorf("database ping failed: %w", err))
	}

	fmt.Println("Connected to PostgreSQL")

	// Initialize router
	r := chi.NewRouter()
	s := &Server{db: db}

	// Health check endpoint
	r.Get("/health", s.handleHealth)

	// Start server
	addr := ":8080"
	fmt.Printf("🚀 Starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(err)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		http.Error(w, "database connection failed", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"status": "ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
