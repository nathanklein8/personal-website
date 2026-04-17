package routes

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"nklein.xyz/backend/models"
	"nklein.xyz/backend/repository"
	"nklein.xyz/backend/server"
	"nklein.xyz/backend/service"
)

type ProjectController struct {
	service *service.ProjectService
}

func NewProjectController(svc *service.ProjectService) *ProjectController {
	return &ProjectController{service: svc}
}

func ProjectRoutes(s *server.Server) chi.Router {
	r := chi.NewRouter()
	controller := NewProjectController(service.NewProjectService(repository.NewProjectRepository(s.DB)))

	r.Get("/", controller.GetAll)
	r.Post("/", controller.CreateOrUpdate)
	r.Put("/", controller.CreateOrUpdate)
	r.Delete("/{id}", controller.DeleteByID)
	r.Put("/{id}", controller.UpdateByID)

	return r
}

// GET /api/projects
func (c *ProjectController) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	projects, err := c.service.GetAll(ctx)
	if err != nil {
		http.Error(w, "failed to fetch projects", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(projects)
}

// POST / PUT /api/projects
func (c *ProjectController) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var p models.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := c.service.ValidateProject(&p); err != nil {
		if apiErr, ok := err.(*service.APIError); ok {
			http.Error(w, apiErr.Message, apiErr.Status)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := c.service.CreateOrUpdate(ctx, &p); err != nil {
		http.Error(w, "failed to save project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/projects/{id}
func (c *ProjectController) DeleteByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid project ID", http.StatusBadRequest)
		return
	}

	if err := c.service.DeleteByID(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "project not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed to delete project", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PUT /api/projects/{id}
func (c *ProjectController) UpdateByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid project ID", http.StatusBadRequest)
		return
	}

	var p models.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	p.ID = id

	if err := c.service.ValidateProject(&p); err != nil {
		if apiErr, ok := err.(*service.APIError); ok {
			http.Error(w, apiErr.Message, apiErr.Status)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := c.service.CreateOrUpdate(ctx, &p); err != nil {
		http.Error(w, "failed to update project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
