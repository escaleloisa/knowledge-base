package kafka

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/escaleloisa/knowledge-base/pkg/models"
	kafka "github.com/segmentio/kafka-go"
)

const Topic = "notes"

type EventType string

const (
	EventNoteCreated EventType = "note.created"
	EventNoteUpdated EventType = "note.updated"
	EventNoteDeleted EventType = "note.deleted"
)

type Event struct {
	Type      EventType    `json:"type"`
	Note      *models.Note `json:"note,omitempty"`
	NoteID    string       `json:"note_id"`
	Timestamp time.Time    `json:"timestamp"`
}

func NewWriter(brokers []string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        Topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	}
}

func NewReader(brokers []string, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    Topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
}

func MarshalEvent(e Event) ([]byte, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("marshal event: %w", err)
	}
	return data, nil
}

func UnmarshalEvent(data []byte) (Event, error) {
	var e Event
	if err := json.Unmarshal(data, &e); err != nil {
		return e, fmt.Errorf("unmarshal event: %w", err)
	}
	return e, nil
}
