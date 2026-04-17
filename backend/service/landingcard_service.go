package service

import (
	"context"
	"net/http"

	"nklein.xyz/backend/models"
	"nklein.xyz/backend/repository"
)

type LandingCardService struct {
	repo *repository.LandingCardRepository
}

func NewLandingCardService(repo *repository.LandingCardRepository) *LandingCardService {
	return &LandingCardService{repo: repo}
}

func (s *LandingCardService) GetByID(ctx context.Context, id int) (*models.LandingCard, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *LandingCardService) CreateOrUpdate(ctx context.Context, lc *models.LandingCard) error {
	return s.repo.CreateOrUpdate(ctx, lc)
}

func (s *LandingCardService) ValidateLandingCard(lc *models.LandingCard) error {
	if lc.Bio == "" {
		return &APIError{Status: http.StatusBadRequest, Message: "bio is required"}
	}
	if lc.Email == "" {
		return &APIError{Status: http.StatusBadRequest, Message: "email is required"}
	}
	return nil
}
