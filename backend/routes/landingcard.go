package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"nklein.xyz/backend/models"
	"nklein.xyz/backend/repository"
	"nklein.xyz/backend/server"
	"nklein.xyz/backend/service"
)

type LandingCardController struct {
	service *service.LandingCardService
}

func NewLandingCardController(svc *service.LandingCardService) *LandingCardController {
	return &LandingCardController{service: svc}
}

func LandingCardRoutes(s *server.Server) chi.Router {
	r := chi.NewRouter()
	controller := NewLandingCardController(service.NewLandingCardService(repository.NewLandingCardRepository(s.DB)))

	r.Get("/", controller.GetByID)
	r.Post("/", controller.CreateOrUpdate)
	r.Put("/", controller.CreateOrUpdate)

	return r
}

// GET /api/landingcard
func (c *LandingCardController) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	lc, err := c.service.GetByID(ctx, 1)
	if err != nil {
		http.Error(w, "failed to fetch from db", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(lc)
}

// POST / PUT /api/landingcard
func (c *LandingCardController) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var lc models.LandingCard
	if err := json.NewDecoder(r.Body).Decode(&lc); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := c.service.ValidateLandingCard(&lc); err != nil {
		if apiErr, ok := err.(*service.APIError); ok {
			http.Error(w, apiErr.Message, apiErr.Status)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := c.service.CreateOrUpdate(ctx, &lc); err != nil {
		http.Error(w, "failed to update landing card", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
