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

type LandingCard struct {
	Bio      string     `json:"bio"`
	Email    string     `json:"email"`
	Linkedin string     `json:"linkedin"`
	Github   string     `json:"github"`
	Skills   [][]string `json:"skills"`
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
	r.Get("/api/health", s.handleHealth)

	// register landing card endpoints
	r.HandleFunc("/api/landingcard", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handleGetLandingCard(w, r)
		case http.MethodPost, http.MethodPut:
			s.handleSetLandingCard(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

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

func (s *Server) handleGetLandingCard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, `
        SELECT bio, email, linkedin, github, skills
        FROM landing_card
        WHERE id = 1
    `)

	var landingCard LandingCard
	var skillsJSON []byte

	if err := row.Scan(&landingCard.Bio, &landingCard.Email, &landingCard.Linkedin, &landingCard.Github, &skillsJSON); err != nil {
		http.Error(w, "failed to fetch landing card content", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(skillsJSON, &landingCard.Skills); err != nil {
		http.Error(w, "invalid skills format in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(landingCard)
}

func (s *Server) handleSetLandingCard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var landingCard LandingCard
	if err := json.NewDecoder(r.Body).Decode(&landingCard); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	skillsJSON, err := json.Marshal(landingCard.Skills)
	if err != nil {
		http.Error(w, "invalid skills format", http.StatusBadRequest)
		return
	}

	_, err = s.db.ExecContext(ctx, `
        INSERT INTO landing_card (id, bio, email, linkedin, github, skills)
        VALUES (1, $1, $2, $3, $4, $5)
        ON CONFLICT (id)
        DO UPDATE SET
          bio = EXCLUDED.bio,
          email = EXCLUDED.email,
          linkedin = EXCLUDED.linkedin,
		  github = EXCLUDED.github,
          skills = EXCLUDED.skills
    `,
		landingCard.Bio,
		landingCard.Email,
		landingCard.Linkedin,
		landingCard.Github,
		skillsJSON,
	)

	if err != nil {
		http.Error(w, "failed to update landing card info", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
