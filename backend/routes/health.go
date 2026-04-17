package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"nklein.xyz/backend/server"
)

type HealthController struct {
	s *server.Server
}

func NewHealthController(s *server.Server) *HealthController {
	return &HealthController{s: s}
}

func HealthRoutes(s *server.Server) chi.Router {
	r := chi.NewRouter()
	controller := NewHealthController(s)

	r.Get("/", controller.HealthCheck)

	return r
}

// GET /api/health
func (c *HealthController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	if err := c.checkDB(ctx); err != nil {
		http.Error(w, "database connection failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

func (c *HealthController) checkDB(ctx context.Context) error {
	_, err := c.s.DB.ExecContext(ctx, "SELECT 1")
	return err
}
