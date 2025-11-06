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
	"github.com/go-chi/cors"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Server struct {
	db *sql.DB
}

type TestItem struct {
	ID    int    `json:"id"`
	Count int    `json:"count"`
	Note  string `json:"note"`
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
	// Basic CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173", // dev server
			"http://atlas:2989",     // test server
			"https://nklein.xyz",    // prod server
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	}))

	s := &Server{db: db}

	// Health check endpoint
	r.Get("/health", s.handleHealth)

	// register connectivity testing endpoints
	r.Post("/api/test", s.createTestItem)
	r.Post("/api/test/{id}/increment", s.incrementCount)

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

// POST /api/test
// Body: { "note": "optional string" }
func (s *Server) createTestItem(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Note string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	var id, count int
	err := s.db.QueryRow(`INSERT INTO test_items (note) VALUES ($1) RETURNING id, count`, input.Note).Scan(&id, &count)
	if err != nil {
		http.Error(w, "insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]any{"id": id, "count": count, "note": input.Note}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// POST /api/test/{id}/increment
func (s *Server) incrementCount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := s.db.Exec(`UPDATE test_items SET count = count + 1 WHERE id = $1`, id)
	if err != nil {
		http.Error(w, "update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var updated struct {
		ID    int    `json:"id"`
		Count int    `json:"count"`
		Note  string `json:"note"`
	}
	err = s.db.QueryRow(`SELECT id, count, note FROM test_items WHERE id = $1`, id).Scan(&updated.ID, &updated.Count, &updated.Note)
	if err != nil {
		http.Error(w, "fetch failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}
