package service

import (
	"context"
	"net/http"

	"nklein.xyz/backend/models"
	"nklein.xyz/backend/repository"
)

type PhotoService struct {
	repo *repository.PhotoRepository
}

func NewPhotoService(repo *repository.PhotoRepository) *PhotoService {
	return &PhotoService{repo: repo}
}

func (s *PhotoService) GetAll(ctx context.Context) ([]models.Photo, error) {
	return s.repo.GetAll(ctx)
}

func (s *PhotoService) Create(ctx context.Context, p *models.Photo) error {
	return s.repo.Create(ctx, p)
}

func (s *PhotoService) Update(ctx context.Context, p *models.Photo) error {
	return s.repo.Update(ctx, p)
}

func (s *PhotoService) DeleteByID(ctx context.Context, id int) error {
	return s.repo.DeleteByID(ctx, id)
}

func (s *PhotoService) GetByID(ctx context.Context, id int) (*models.Photo, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PhotoService) ValidatePhotoCreate(p *models.PhotoCreateRequest) error {
	if p.Title == nil || *p.Title == "" {
		return &APIError{Status: http.StatusBadRequest, Message: "title is required"}
	}
	if p.FilePath == "" {
		return &APIError{Status: http.StatusBadRequest, Message: "file_path is required"}
	}
	return nil
}

func (s *PhotoService) ValidatePhotoUpdate(p *models.PhotoUpdateRequest) error {
	if p.Title == nil || *p.Title == "" {
		return &APIError{Status: http.StatusBadRequest, Message: "title is required"}
	}
	if p.FilePath == "" {
		return &APIError{Status: http.StatusBadRequest, Message: "file_path is required"}
	}
	return nil
}
