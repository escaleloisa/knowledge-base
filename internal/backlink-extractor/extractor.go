package backlinkextractor

import (
	"context"
	"database/sql"
	"log"
	"regexp"

	kafkapkg "github.com/escaleloisa/knowledge-base/pkg/kafka"
	"github.com/segmentio/kafka-go"
)

var linkPattern = regexp.MustCompile(`\[\[([^\]]+)\]\]`)

type Extractor struct {
	reader *kafka.Reader
	db     *sql.DB
}

func New(reader *kafka.Reader, db *sql.DB) *Extractor {
	return &Extractor{reader: reader, db: db}
}

func (e *Extractor) Run(ctx context.Context) error {
	log.Println("Backlink extractor started, consuming events...")
	for {
		msg, err := e.reader.ReadMessage(ctx)
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
			if err := e.processNote(ctx, event); err != nil {
				log.Printf("process error: %v", err)
			}
		case kafkapkg.EventNoteDeleted:
			if err := e.deleteBacklinks(ctx, event.NoteID); err != nil {
				log.Printf("delete error: %v", err)
			}
		}
	}
}

func (e *Extractor) processNote(ctx context.Context, event kafkapkg.Event) error {
	// Delete existing backlinks for this note
	if err := e.deleteBacklinks(ctx, event.NoteID); err != nil {
		return err
	}

	// Extract [[links]] from content
	matches := linkPattern.FindAllStringSubmatch(event.Note.Content, -1)
	if len(matches) == 0 {
		return nil
	}

	for _, match := range matches {
		title := match[1]
		// Find target note by title
		var targetID string
		err := e.db.QueryRowContext(ctx,
			`SELECT id FROM notes WHERE title = $1 AND id != $2`, title, event.NoteID,
		).Scan(&targetID)
		if err != nil {
			continue // Target note doesn't exist, skip
		}

		// Insert backlink
		_, err = e.db.ExecContext(ctx,
			`INSERT INTO backlinks (source_note_id, target_note_id) VALUES ($1, $2)
			 ON CONFLICT DO NOTHING`, event.NoteID, targetID)
		if err != nil {
			log.Printf("insert backlink error: %v", err)
		}
	}

	log.Printf("processed backlinks for note %s", event.NoteID)
	return nil
}

func (e *Extractor) deleteBacklinks(ctx context.Context, noteID string) error {
	_, err := e.db.ExecContext(ctx, `DELETE FROM backlinks WHERE source_note_id = $1`, noteID)
	return err
}
