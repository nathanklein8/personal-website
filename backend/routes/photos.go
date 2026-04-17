package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"database/sql"
	"nklein.xyz/backend/models"
	"nklein.xyz/backend/repository"
	"nklein.xyz/backend/server"
	"nklein.xyz/backend/service"
)

type PhotoController struct {
	service *service.PhotoService
}

func NewPhotoController(svc *service.PhotoService) *PhotoController {
	return &PhotoController{service: svc}
}

func PhotoRoutes(s *server.Server) chi.Router {
	r := chi.NewRouter()
	controller := NewPhotoController(service.NewPhotoService(repository.NewPhotoRepository(s.DB)))

	r.Get("/", controller.GetAll)
	r.Post("/", controller.Create)
	r.Put("/{id}", controller.Update)
	r.Delete("/{id}", controller.DeleteByID)
	r.Get("/{id}/image", controller.ServeImage)

	return r
}

// GET /api/photos
func (c *PhotoController) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	photos, err := c.service.GetAll(ctx)
	if err != nil {
		http.Error(w, "failed to fetch photos", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(photos)
}

// POST /api/photos
func (c *PhotoController) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var p models.Photo
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := c.service.ValidatePhotoCreate(&models.PhotoCreateRequest{
		Title:    &p.Title,
		FilePath: p.FilePath,
	}); err != nil {
		if apiErr, ok := err.(*service.APIError); ok {
			http.Error(w, apiErr.Message, apiErr.Status)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := c.service.Create(ctx, &p); err != nil {
		http.Error(w, "failed to insert photo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// PUT /api/photos/{id}
func (c *PhotoController) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid photo ID", http.StatusBadRequest)
		return
	}

	var p models.PhotoUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := c.service.ValidatePhotoUpdate(&p); err != nil {
		if apiErr, ok := err.(*service.APIError); ok {
			http.Error(w, apiErr.Message, apiErr.Status)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Build full photo object with ID for update
	p.ID = id
	fullPhoto := &models.Photo{
		ID:           p.ID,
		Title:        *p.Title,
		FilePath:     p.FilePath,
		AltText:      p.AltText,
		DateTaken:    p.DateTaken,
		Location:     p.Location,
		Camera:       p.Camera,
		Lens:         p.Lens,
		Aperture:     p.Aperture,
		ShutterSpeed: p.ShutterSpeed,
		ISO:          p.ISO,
		Visible:      p.Visible,
		SortOrder:    p.SortOrder,
	}

	if err := c.service.Update(ctx, fullPhoto); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "photo not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed to update photo", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/photos/{id}
func (c *PhotoController) DeleteByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid photo ID", http.StatusBadRequest)
		return
	}

	if err := c.service.DeleteByID(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "photo not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed to delete photo", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /api/photos/{id}/image
func (c *PhotoController) ServeImage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid photo ID", http.StatusBadRequest)
		return
	}

	var p *models.Photo
	p, err = c.service.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "photo not found", http.StatusNotFound)
		return
	}

	absPath := filepath.Join("/photos", p.FilePath)
	fmt.Printf("Serving image from path: %s\n", absPath)

	http.ServeFile(w, r, absPath)
}
