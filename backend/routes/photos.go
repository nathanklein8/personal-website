package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
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

	// File explorer browsing routes
	r.Group(func(r chi.Router) {
		r.Get("/available", controller.ListYears)
		r.Get("/available/{year}", controller.ListEvents)
		r.Get("/available/{year}/{event}", controller.ListPhotos)
		r.Get("/available/{year}/{event}/{filename}/preview", controller.ServePreview)
	})

	// CRUD routes
	r.Get("/", controller.GetAll)
	r.Get("/{id}/thumbnail", controller.ServeThumbnail)
	r.Post("/", controller.AddPhoto)
	r.Put("/{id}", controller.Update)
	r.Delete("/{id}", controller.DeleteByID)
	r.Post("/regenerate-thumbnails", controller.RegenerateThumbnails)

	return r
}

// --- File Explorer Browsing ---

// GET /api/photos/available - List year directories
func (c *PhotoController) ListYears(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	years, err := c.service.ListYears(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(years)
}

// GET /api/photos/available/{year} - List event directories within a year
func (c *PhotoController) ListEvents(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	year := chi.URLParam(r, "year")
	events, err := c.service.ListEvents(ctx, year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(events)
}

// GET /api/photos/available/{year}/{event} - List .jpg filenames
func (c *PhotoController) ListPhotos(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	year := chi.URLParam(r, "year")
	event := chi.URLParam(r, "event")
	photos, err := c.service.ListPhotos(ctx, year, event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(photos)
}

// GET /api/photos/available/{year}/{event}/{filename}/preview - Serve thumbnail, medium, or full preview
// Query param: size=thumb|med|full (default: med)
func (c *PhotoController) ServePreview(w http.ResponseWriter, r *http.Request) {
	year := chi.URLParam(r, "year")
	event := chi.URLParam(r, "event")
	filename := chi.URLParam(r, "filename")

	size := r.URL.Query().Get("size")
	if size == "" {
		size = "med"
	}

	// Build the relative source path
	sourceRelPath := filepath.Join(year, event, filename)
	thumbRelPath := filepath.Join(year, event, filename+"_thumb.jpg")
	mediumRelPath := filepath.Join(year, event, filename+"_med.jpg")

	var servePath string
	switch size {
	case "thumb":
		servePath = filepath.Join("/thumbnails", thumbRelPath)
	case "med":
		servePath = filepath.Join("/thumbnails", mediumRelPath)
	case "full":
		servePath = filepath.Join("/photos", sourceRelPath)
	default:
		http.Error(w, "invalid size parameter", http.StatusBadRequest)
		return
	}

	if _, err := os.Stat(servePath); os.IsNotExist(err) {
		http.Error(w, "preview not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, servePath)
}

// --- CRUD Operations ---

// GET /api/photos/{id}/thumbnail - Serve medium preview for a photo
func (c *PhotoController) ServeThumbnail(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid photo ID", http.StatusBadRequest)
		return
	}

	p, err := c.service.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "photo not found", http.StatusNotFound)
		return
	}

	servePath := filepath.Join("/thumbnails", p.MediumPath())
	if _, err := os.Stat(servePath); os.IsNotExist(err) {
		http.Error(w, "thumbnail not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, servePath)
}

// GET /api/photos - List photos based on query parameter
// ?type=all -> all photos (admin)
// ?type=visible -> only visible photos (public gallery)
// ?type=featured -> only visible and featured photos (landing page)
// no parameter -> all photos (default for backwards compatibility)
func (c *PhotoController) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	photoType := r.URL.Query().Get("type")
	var photos []models.Photo
	var err error

	switch photoType {
	case "visible":
		photos, err = c.service.GetVisible(ctx)
	case "featured":
		photos, err = c.service.GetFeatured(ctx)
	default:
		photos, err = c.service.GetAll(ctx)
	}

	if err != nil {
		http.Error(w, "failed to fetch photos", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(photos)
}

// POST /api/photos - Add a photo from the file explorer
func (c *PhotoController) AddPhoto(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var req models.AddPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Filename == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	photo, err := c.service.AddPhoto(ctx, req)
	if err != nil {
		if apiErr, ok := err.(*service.APIError); ok {
			http.Error(w, apiErr.Message, apiErr.Status)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(photo)
}

// PUT /api/photos/{id} - Update caption and ordering
func (c *PhotoController) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid photo ID", http.StatusBadRequest)
		return
	}

	var req models.PhotoUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := c.service.Update(ctx, id, req); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "photo not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed to update photo", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/photos/{id} - Delete photo DB row and generated thumbnails
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

// POST /api/photos/regenerate-thumbnails - Regenerate all thumbnails
func (c *PhotoController) RegenerateThumbnails(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	if err := c.service.RegenerateAllThumbnails(ctx); err != nil {
		http.Error(w, "failed to regenerate thumbnails", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
