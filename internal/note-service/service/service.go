package service

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/escaleloisa/knowledge-base/internal/note-service/repository"
	"github.com/escaleloisa/knowledge-base/pkg/models"

	kafkapkg "github.com/escaleloisa/knowledge-base/pkg/kafka"
	kafka "github.com/segmentio/kafka-go"
)

type Service struct {
	repo   *repository.Repository
	writer *kafka.Writer
}

func New(repo *repository.Repository, writer *kafka.Writer) *Service {
	return &Service{repo: repo, writer: writer}
}

func (s *Service) Create(ctx context.Context, req models.CreateNoteRequest) (*models.Note, error) {
	note, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	s.publishEvent(ctx, kafkapkg.EventNoteCreated, note)
	return note, nil
}

func (s *Service) Get(ctx context.Context, id string) (*models.Note, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Update(ctx context.Context, id string, req models.UpdateNoteRequest) (*models.Note, error) {
	note, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}
	if note != nil {
		s.publishEvent(ctx, kafkapkg.EventNoteUpdated, note)
	}
	return note, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	s.publishEvent(ctx, kafkapkg.EventNoteDeleted, &models.Note{ID: id})
	return nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]models.Note, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *Service) publishEvent(ctx context.Context, eventType kafkapkg.EventType, note *models.Note) {
	if s.writer == nil {
		return
	}
	event := kafkapkg.Event{
		Type:      eventType,
		Note:      note,
		NoteID:    note.ID,
		Timestamp: time.Now(),
	}
	data, err := kafkapkg.MarshalEvent(event)
	if err != nil {
		log.Printf("failed to marshal event: %v", err)
		return
	}
	err = s.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(note.ID),
		Value: data,
	})
	if err != nil {
		log.Printf("failed to publish event: %v", err)
	}
}

func IsNotFound(err error) bool {
	return err == sql.ErrNoRows
}
