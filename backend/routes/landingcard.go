package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"nklein.xyz/backend/server"
)

type LandingCard struct {
	Bio      string     `json:"bio"`
	Email    string     `json:"email"`
	Linkedin string     `json:"linkedin"`
	Github   string     `json:"github"`
	Skills   [][]string `json:"skills"`
}

func LandingCardRoutes(s *server.Server) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		handleGetLandingCard(s, w, req)
	})

	r.Post("/", func(w http.ResponseWriter, req *http.Request) {
		handleSetLandingCard(s, w, req)
	})

	r.Put("/", func(w http.ResponseWriter, req *http.Request) {
		handleSetLandingCard(s, w, req)
	})

	return r
}

func handleGetLandingCard(s *server.Server, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	row := s.DB.QueryRowContext(ctx, `
        SELECT bio, email, linkedin, github, skills
        FROM landing_card
        WHERE id = 1
    `)

	var lc LandingCard
	var skillsJSON []byte

	if err := row.Scan(&lc.Bio, &lc.Email, &lc.Linkedin, &lc.Github, &skillsJSON); err != nil {
		http.Error(w, "failed to fetch from db", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(skillsJSON, &lc.Skills); err != nil {
		http.Error(w, "invalid skills format", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(lc)
}

func handleSetLandingCard(s *server.Server, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var lc LandingCard
	if err := json.NewDecoder(r.Body).Decode(&lc); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	skillsJSON, _ := json.Marshal(lc.Skills)

	_, err := s.DB.ExecContext(ctx, `
        INSERT INTO landing_card (id, bio, email, linkedin, github, skills)
        VALUES (1, $1, $2, $3, $4, $5)
        ON CONFLICT (id) DO UPDATE SET
            bio = EXCLUDED.bio,
            email = EXCLUDED.email,
            linkedin = EXCLUDED.linkedin,
            github = EXCLUDED.github,
            skills = EXCLUDED.skills
    `,
		lc.Bio,
		lc.Email,
		lc.Linkedin,
		lc.Github,
		skillsJSON,
	)

	if err != nil {
		http.Error(w, "failed to update landing card", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
