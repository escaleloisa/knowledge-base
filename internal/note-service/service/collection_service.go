package service

import (
	"context"

	"github.com/escaleloisa/knowledge-base/pkg/models"
)

func (s *Service) CreateCollection(ctx context.Context, req models.CreateCollectionRequest) (*models.Collection, error) {
	return s.repo.CreateCollection(ctx, req)
}

func (s *Service) GetCollection(ctx context.Context, id string) (*models.Collection, error) {
	return s.repo.GetCollection(ctx, id)
}

func (s *Service) UpdateCollection(ctx context.Context, id string, req models.UpdateCollectionRequest) (*models.Collection, error) {
	return s.repo.UpdateCollection(ctx, id, req)
}

func (s *Service) DeleteCollection(ctx context.Context, id string) error {
	return s.repo.DeleteCollection(ctx, id)
}

func (s *Service) ListCollections(ctx context.Context) ([]models.Collection, error) {
	return s.repo.ListCollections(ctx)
}

func (s *Service) AddNotesToCollection(ctx context.Context, collectionID string, noteIDs []string) (int, error) {
	return s.repo.AddNotesToCollection(ctx, collectionID, noteIDs)
}

func (s *Service) RemoveNoteFromCollection(ctx context.Context, collectionID, noteID string) error {
	return s.repo.RemoveNoteFromCollection(ctx, collectionID, noteID)
}

func (s *Service) ListCollectionNotes(ctx context.Context, collectionID string, limit, offset int) ([]models.Note, error) {
	return s.repo.ListCollectionNotes(ctx, collectionID, limit, offset)
}
