package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"nklein.xyz/backend/server"
)

type Project struct {
	ID             int      `json:"id,omitempty"`
	Icon           string   `json:"icon"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Technologies   []string `json:"technologies"`
	DeploymentLink *string  `json:"deploymentLink,omitempty"`
	Image          *string  `json:"image,omitempty"`
	AltText        *string  `json:"altText,omitempty"`
}

func ProjectRoutes(s *server.Server) chi.Router {
	r := chi.NewRouter()

	// get list of all project objs
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handleGetProjects(s, w, r)
	})

	// insert a new project obj
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handleSetProject(s, w, r)
	})
	r.Put("/", func(w http.ResponseWriter, r *http.Request) {
		handleSetProject(s, w, r)
	})

	// manage existing projects (delete / edit)
	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handleDeleteProject(s, w, r)
	})
	r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handleEditProject(s, w, r)
	})

	return r
}

// GET /api/projects
func handleGetProjects(s *server.Server, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, `
        SELECT id, icon, title, description, technologies, deployment_link, image, alt_text
        FROM projects
        ORDER BY id ASC
    `)
	if err != nil {
		http.Error(w, "failed to fetch projects", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var projects []Project

	for rows.Next() {
		var p Project
		var techJSON []byte

		if err := rows.Scan(
			&p.ID,
			&p.Icon,
			&p.Title,
			&p.Description,
			&techJSON,
			&p.DeploymentLink,
			&p.Image,
			&p.AltText,
		); err != nil {
			http.Error(w, "failed to scan project row", http.StatusInternalServerError)
			return
		}

		if err := json.Unmarshal(techJSON, &p.Technologies); err != nil {
			http.Error(w, "invalid technologies format", http.StatusInternalServerError)
			return
		}

		projects = append(projects, p)
	}

	json.NewEncoder(w).Encode(projects)
}

// POST / PUT /api/projects
func handleSetProject(s *server.Server, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var p Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	techJSON, _ := json.Marshal(p.Technologies)

	// If ID is 0 → insert new project; otherwise, update existing
	if p.ID == 0 {
		_, err := s.DB.ExecContext(ctx, `
            INSERT INTO projects (icon, title, description, technologies, deployment_link, image, alt_text)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
        `,
			p.Icon,
			p.Title,
			p.Description,
			techJSON,
			p.DeploymentLink,
			p.Image,
			p.AltText,
		)

		if err != nil {
			http.Error(w, "failed to insert project", http.StatusInternalServerError)
			return
		}
	} else {
		_, err := s.DB.ExecContext(ctx, `
            UPDATE projects
            SET icon = $1,
                title = $2,
                description = $3,
                technologies = $4,
                deployment_link = $5,
                image = $6,
                alt_text = $7
            WHERE id = $8
        `,
			p.Icon,
			p.Title,
			p.Description,
			techJSON,
			p.DeploymentLink,
			p.Image,
			p.AltText,
			p.ID,
		)

		if err != nil {
			http.Error(w, "failed to update project", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleDeleteProject(s *server.Server, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		http.Error(w, "missing project ID", http.StatusBadRequest)
		return
	}

	// Convert ID to int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid project ID", http.StatusBadRequest)
		return
	}

	// Execute DELETE query
	res, err := s.DB.ExecContext(ctx, `
        DELETE FROM projects WHERE id = $1
    `, id)
	if err != nil {
		http.Error(w, "failed to delete project", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "failed to confirm deletion", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleEditProject(s *server.Server, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		http.Error(w, "missing project ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid project ID", http.StatusBadRequest)
		return
	}

	var p Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	techJSON, _ := json.Marshal(p.Technologies)

	res, err := s.DB.ExecContext(ctx, `
        UPDATE projects
        SET icon = $1,
            title = $2,
            description = $3,
            technologies = $4,
            deployment_link = $5,
            image = $6,
            alt_text = $7
        WHERE id = $8
    `,
		p.Icon,
		p.Title,
		p.Description,
		techJSON,
		p.DeploymentLink,
		p.Image,
		p.AltText,
		id,
	)

	if err != nil {
		http.Error(w, "failed to update project", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "failed to confirm update", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
