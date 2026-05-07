package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/escaleloisa/knowledge-base/pkg/models"
	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, req models.CreateNoteRequest) (*models.Note, error) {
	note := &models.Note{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO notes (title, content, tags) VALUES ($1, $2, $3)
		 RETURNING id, title, content, tags, created_at, updated_at`,
		req.Title, req.Content, pq.Array(req.Tags),
	).Scan(&note.ID, &note.Title, &note.Content, pq.Array(&note.Tags), &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create note: %w", err)
	}
	return note, nil
}

func (r *Repository) Get(ctx context.Context, id string) (*models.Note, error) {
	note := &models.Note{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, content, tags, created_at, updated_at FROM notes WHERE id = $1`, id,
	).Scan(&note.ID, &note.Title, &note.Content, pq.Array(&note.Tags), &note.CreatedAt, &note.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get note: %w", err)
	}
	return note, nil
}

func (r *Repository) Update(ctx context.Context, id string, req models.UpdateNoteRequest) (*models.Note, error) {
	note := &models.Note{}
	err := r.db.QueryRowContext(ctx,
		`UPDATE notes SET title = $1, content = $2, tags = $3, updated_at = NOW()
		 WHERE id = $4
		 RETURNING id, title, content, tags, created_at, updated_at`,
		req.Title, req.Content, pq.Array(req.Tags), id,
	).Scan(&note.ID, &note.Title, &note.Content, pq.Array(&note.Tags), &note.CreatedAt, &note.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update note: %w", err)
	}
	return note, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM notes WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]models.Note, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, content, tags, created_at, updated_at FROM notes
		 ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, pq.Array(&n.Tags), &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan note: %w", err)
		}
		notes = append(notes, n)
	}
	return notes, rows.Err()
}
