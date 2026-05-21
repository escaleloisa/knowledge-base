package searchindexer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	kafkapkg "github.com/escaleloisa/knowledge-base/pkg/kafka"
	"github.com/segmentio/kafka-go"
)

type Indexer struct {
	reader *kafka.Reader
	esURL  string
}

func New(reader *kafka.Reader, esURL string) *Indexer {
	return &Indexer{reader: reader, esURL: esURL}
}

func (idx *Indexer) Run(ctx context.Context) error {
	log.Println("Search indexer started, consuming events...")
	for {
		msg, err := idx.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Printf("read error: %v", err)
			continue
		}

		event, err := kafkapkg.UnmarshalEvent(msg.Value)
		if err != nil {
			log.Printf("unmarshal error: %v", err)
			continue
		}

		switch event.Type {
		case kafkapkg.EventNoteCreated, kafkapkg.EventNoteUpdated:
			if err := idx.indexNote(event); err != nil {
				log.Printf("index error: %v", err)
			}
		case kafkapkg.EventNoteDeleted:
			if err := idx.deleteNote(event.NoteID); err != nil {
				log.Printf("delete error: %v", err)
			}
		}
	}
}

func (idx *Indexer) indexNote(event kafkapkg.Event) error {
	doc, _ := json.Marshal(event.Note)
	url := fmt.Sprintf("%s/notes/_doc/%s", idx.esURL, event.NoteID)
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(doc))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("elasticsearch request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("elasticsearch returned %d", resp.StatusCode)
	}
	log.Printf("indexed note %s", event.NoteID)
	return nil
}

func (idx *Indexer) deleteNote(id string) error {
	url := fmt.Sprintf("%s/notes/_doc/%s", idx.esURL, id)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("elasticsearch request: %w", err)
	}
	defer resp.Body.Close()
	log.Printf("deleted note %s from index", id)
	return nil
}
