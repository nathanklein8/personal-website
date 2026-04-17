package service

import (
	"context"
	"net/http"

	"nklein.xyz/backend/models"
	"nklein.xyz/backend/repository"
)

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) GetAll(ctx context.Context) ([]models.Project, error) {
	return s.repo.GetAll(ctx)
}

func (s *ProjectService) CreateOrUpdate(ctx context.Context, p *models.Project) error {
	return s.repo.CreateOrUpdate(ctx, p)
}

func (s *ProjectService) DeleteByID(ctx context.Context, id int) error {
	return s.repo.DeleteByID(ctx, id)
}

func (s *ProjectService) ValidateProject(p *models.Project) error {
	if p.Icon == "" {
		return &APIError{Status: http.StatusBadRequest, Message: "icon is required"}
	}
	if p.Title == "" {
		return &APIError{Status: http.StatusBadRequest, Message: "title is required"}
	}
	if p.Description == "" {
		return &APIError{Status: http.StatusBadRequest, Message: "description is required"}
	}
	return nil
}
