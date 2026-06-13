package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/escaleloisa/knowledge-base/pkg/models"
	"github.com/lib/pq"
)

func (r *Repository) CreateCollection(ctx context.Context, req models.CreateCollectionRequest) (*models.Collection, error) {
	c := &models.Collection{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO collections (name, description) VALUES ($1, $2)
		 RETURNING id, name, description, created_at, updated_at`,
		req.Name, req.Description,
	).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create collection: %w", err)
	}
	return c, nil
}

func (r *Repository) GetCollection(ctx context.Context, id string) (*models.Collection, error) {
	c := &models.Collection{}
	err := r.db.QueryRowContext(ctx,
		`SELECT c.id, c.name, c.description, COUNT(cn.note_id) as note_count, c.created_at, c.updated_at
	FROM collections c
	LEFT JOIN collection_notes cn ON cn.collection_id = c.id
	WHERE c.id = $1
	GROUP BY c.id`, id,
	).Scan(&c.ID, &c.Name, &c.Description, &c.NoteCount, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get collection: %w", err)
	}
	return c, nil
}

func (r *Repository) UpdateCollection(ctx context.Context, id string, req models.UpdateCollectionRequest) (*models.Collection, error) {
	c := &models.Collection{}
	err := r.db.QueryRowContext(ctx,
		`UPDATE collections SET name = $1, description = $2, updated_at = NOW()
	WHERE id = $3
	RETURNING id, name, description, created_at, updated_at`,
		req.Name, req.Description, id,
	).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update collection: %w", err)
	}
	return c, nil
}

func (r *Repository) DeleteCollection(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM collections WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) ListCollections(ctx context.Context) ([]models.Collection, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT c.id, c.name, c.description, COUNT(cn.note_id) as note_count, c.created_at, c.updated_at
	FROM collections c
	LEFT JOIN collection_notes cn ON cn.collection_id = c.id
	GROUP BY c.id
	ORDER BY c.name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list collections: %w", err)
	}
	defer rows.Close()

	var collections []models.Collection
	for rows.Next() {
		var c models.Collection
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.NoteCount, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan collection: %w", err)
		}
		collections = append(collections, c)
	}
	return collections, rows.Err()
}

func (r *Repository) AddNotesToCollection(ctx context.Context, collectionID string, noteIDs []string) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO collection_notes (collection_id, note_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
	)
	if err != nil {
		return 0, fmt.Errorf("prepare stmt: %w", err)
	}
	defer stmt.Close()

	added := 0
	for _, noteID := range noteIDs {
		result, err := stmt.ExecContext(ctx, collectionID, noteID)
		if err != nil {
			return 0, fmt.Errorf("add note %s: %w", noteID, err)
		}
		rows, _ := result.RowsAffected()
		added += int(rows)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit tx: %w", err)
	}
	return added, nil
}

func (r *Repository) RemoveNoteFromCollection(ctx context.Context, collectionID, noteID string) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM collection_notes WHERE collection_id = $1 AND note_id = $2`,
		collectionID, noteID,
	)
	if err != nil {
		return fmt.Errorf("remove note from collection: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) ListCollectionNotes(ctx context.Context, collectionID string, limit, offset int) ([]models.Note, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT n.id, n.title, n.content, n.tags, n.created_at, n.updated_at
		 FROM notes n
		 JOIN collection_notes cn ON cn.note_id = n.id
		 WHERE cn.collection_id = $1
		 ORDER BY cn.added_at DESC
		 LIMIT $2 OFFSET $3`,
		collectionID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list collection notes: %w", err)
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
