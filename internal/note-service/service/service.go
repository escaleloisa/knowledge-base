package service

import (
	"context"
	"database/sql"

	"github.com/escaleloisa/knowledge-base/internal/note-service/repository"
	"github.com/escaleloisa/knowledge-base/pkg/models"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req models.CreateNoteRequest) (*models.Note, error) {
	return s.repo.Create(ctx, req)
}

func (s *Service) Get(ctx context.Context, id string) (*models.Note, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Update(ctx context.Context, id string, req models.UpdateNoteRequest) (*models.Note, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]models.Note, error) {
	return s.repo.List(ctx, limit, offset)
}

// IsNotFound checks if the error is a not-found case.
func IsNotFound(err error) bool {
	return err == sql.ErrNoRows
}
