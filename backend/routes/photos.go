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
	"nklein.xyz/backend/server"
)

type Photo struct {
	ID           int     `json:"id,omitempty"`
	Title        string  `json:"title"`
	FilePath     string  `json:"filePath"`
	AltText      *string `json:"altText,omitempty"`
	DateTaken    *string `json:"dateTaken,omitempty"`
	Location     *string `json:"location,omitempty"`
	Camera       *string `json:"camera,omitempty"`
	Lens         *string `json:"lens,omitempty"`
	Aperture     *string `json:"aperture,omitempty"`
	ShutterSpeed *string `json:"shutterSpeed,omitempty"`
	ISO          *string `json:"iso,omitempty"`
	Visible      bool    `json:"visible"`
	Featured     bool    `json:"featured,omitempty"`
	SortOrder    int     `json:"sortOrder"`
}


type PhotoCreateRequest struct {
	Title        *string   `json:"title"`
	FilePath     string    `json:"file_path"`
	AltText      *string   `json:"altText,omitempty"`
	DateTaken    *string   `json:"dateTaken,omitempty"`
	Location     *string   `json:"location,omitempty"`
	Camera       *string   `json:"camera,omitempty"`
	Lens         *string   `json:"lens,omitempty"`
	Aperture     *string   `json:"aperture,omitempty"`
	ShutterSpeed *string   `json:"shutterSpeed,omitempty"`
	ISO          *string   `json:"iso,omitempty"`
	Visible      *bool     `json:"visible,omitempty"`
	Featured     *bool     `json:"featured,omitempty"`
	SortOrder    *int      `json:"sortOrder,omitempty"`
}


type PhotoUpdateRequest struct {
	Title        *string   `json:"title"`
	FilePath     string    `json:"file_path,omitempty"`
	AltText      *string   `json:"altText,omitempty"`
	DateTaken    *string   `json:"dateTaken,omitempty"`
	Location     *string   `json:"location,omitempty"`
	Camera       *string   `json:"camera,omitempty"`
	Lens         *string   `json:"lens,omitempty"`
	Aperture     *string   `json:"aperture,omitempty"`
	ShutterSpeed *string   `json:"shutterSpeed,omitempty"`
	ISO          *string   `json:"iso,omitempty"`
	Visible      bool      `json:"visible,omitempty"`
	Featured     bool      `json:"featured,omitempty"`
	SortOrder    int       `json:"sortOrder,omitempty"`
}
func PhotoRoutes(s *server.Server) chi.Router {
	r := chi.NewRouter()

	// get list of all photo objs
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handleGetPhotos(s, w, r)
	})

	// manage existing photos (delete / edit)
	r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handleEditPhoto(s, w, r)
	})
	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handleDeletePhoto(s, w, r)
	})

	// create new photo record
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handleCreatePhoto(s, w, r)
	})

	// serve actual image file from volume
	r.Get("/{id}/image", func(w http.ResponseWriter, r *http.Request) {
		handleServeImage(s, w, r)
	})

	return r
}

// GET /api/photos
func handleGetPhotos(s *server.Server, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, `
        SELECT id, title, file_path, alt_text, date_taken, location, 
               camera, lens, aperture, shutter_speed, iso, visible, sort_order 
        FROM photos 
        ORDER BY sort_order ASC, id ASC
    `)
	if err != nil {
		http.Error(w, "failed to fetch photos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var photos []Photo

	for rows.Next() {
		var p Photo
		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.FilePath,
			&p.AltText,
			&p.DateTaken,
			&p.Location,
			&p.Camera,
			&p.Lens,
			&p.Aperture,
			&p.ShutterSpeed,
			&p.ISO,
			&p.Visible,
			&p.SortOrder,
		); err != nil {
			http.Error(w, "failed to scan photo row", http.StatusInternalServerError)
			return
		}
		photos = append(photos, p)
	}

	json.NewEncoder(w).Encode(photos)
}

// POST /api/photos
func handleCreatePhoto(s *server.Server, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var p Photo
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	err := s.DB.QueryRowContext(ctx, `
        INSERT INTO photos (title, file_path, alt_text, date_taken, location, camera, lens, aperture, shutter_speed, iso, visible, sort_order)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING id
    `,
		p.Title,
		p.FilePath,
		p.AltText,
		p.DateTaken,
		p.Location,
		p.Camera,
		p.Lens,
		p.Aperture,
		p.ShutterSpeed,
		p.ISO,
		p.Visible,
		p.SortOrder,
	).Scan(&p.ID)

	if err != nil {
		http.Error(w, "failed to insert photo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// PUT /api/photos/{id}
func handleEditPhoto(s *server.Server, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid photo ID", http.StatusBadRequest)
		return
	}

	var p Photo
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	res, err := s.DB.ExecContext(ctx, `
        UPDATE photos
        SET title = $1,
            file_path = $2,
            alt_text = $3,
            date_taken = $4,
            location = $5,
            camera = $6,
            lens = $7,
            aperture = $8,
            shutter_speed = $9,
            iso = $10,
            visible = $11,
            sort_order = $12
        WHERE id = $13
    `,
		p.Title,
		p.FilePath,
		p.AltText,
		p.DateTaken,
		p.Location,
		p.Camera,
		p.Lens,
		p.Aperture,
		p.ShutterSpeed,
		p.ISO,
		p.Visible,
		p.SortOrder,
		id,
	)

	if err != nil {
		http.Error(w, "failed to update photo", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "photo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/photos/{id}
func handleDeletePhoto(s *server.Server, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid photo ID", http.StatusBadRequest)
		return
	}

	res, err := s.DB.ExecContext(ctx, `DELETE FROM photos WHERE id = $1`, id)
	if err != nil {
		http.Error(w, "failed to delete photo", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "photo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /api/photos/{id}/image
func handleServeImage(s *server.Server, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid photo ID", http.StatusBadRequest)
		return
	}

	var filePath string
	err = s.DB.QueryRowContext(ctx, `SELECT file_path FROM photos WHERE id = $1`, id).Scan(&filePath)
	if err != nil {
		http.Error(w, "photo not found", http.StatusNotFound)
		return
	}

	// The volume is mounted at /photos
	absPath := filepath.Join("/photos", filePath)
	
	// Log the path for debugging purposes
	fmt.Printf("Serving image from path: %s\n", absPath)

	http.ServeFile(w, r, absPath)
}
