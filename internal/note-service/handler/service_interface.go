package handler

import (
	"context"

	"github.com/escaleloisa/knowledge-base/pkg/models"
)

// ServiceInterface defines the methods the handler needs from the service layer.
// This enables unit testing with mocks.
type ServiceInterface interface {
	// Notes
	Create(ctx context.Context, req models.CreateNoteRequest) (*models.Note, error)
	Get(ctx context.Context, id string) (*models.Note, error)
	Update(ctx context.Context, id string, req models.UpdateNoteRequest) (*models.Note, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]models.Note, error)

	// Collections
	CreateCollection(ctx context.Context, req models.CreateCollectionRequest) (*models.Collection, error)
	GetCollection(ctx context.Context, id string) (*models.Collection, error)
	UpdateCollection(ctx context.Context, id string, req models.UpdateCollectionRequest) (*models.Collection, error)
	DeleteCollection(ctx context.Context, id string) error
	ListCollections(ctx context.Context) ([]models.Collection, error)
	AddNotesToCollection(ctx context.Context, collectionID string, noteIDs []string) (int, error)
	RemoveNoteFromCollection(ctx context.Context, collectionID, noteID string) error
	ListCollectionNotes(ctx context.Context, collectionID string, limit, offset int) ([]models.Note, error)
}
