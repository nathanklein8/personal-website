package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"nklein.xyz/backend/server"
)

func HealthRoutes(s *server.Server) chi.Router {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
		defer cancel()

		if err := s.DB.PingContext(ctx); err != nil {
			http.Error(w, "database connection failed", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	return r
}
