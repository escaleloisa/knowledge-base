package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/escaleloisa/knowledge-base/pkg/models"
)

// mockService implements ServiceInterface for testing
type mockService struct {
	// Collection methods
	createCollectionFn     func(ctx context.Context, req models.CreateCollectionRequest) (*models.Collection, error)
	getCollectionFn        func(ctx context.Context, id string) (*models.Collection, error)
	updateCollectionFn     func(ctx context.Context, id string, req models.UpdateCollectionRequest) (*models.Collection, error)
	deleteCollectionFn     func(ctx context.Context, id string) error
	listCollectionsFn      func(ctx context.Context) ([]models.Collection, error)
	addNotesToCollectionFn func(ctx context.Context, collectionID string, noteIDs []string) (int, error)
	removeNoteFromCollFn   func(ctx context.Context, collectionID, noteID string) error
	listCollectionNotesFn  func(ctx context.Context, collectionID string, limit, offset int) ([]models.Note, error)

	// Note methods (required by interface)
	createFn func(ctx context.Context, req models.CreateNoteRequest) (*models.Note, error)
	getFn    func(ctx context.Context, id string) (*models.Note, error)
	updateFn func(ctx context.Context, id string, req models.UpdateNoteRequest) (*models.Note, error)
	deleteFn func(ctx context.Context, id string) error
	listFn   func(ctx context.Context, limit, offset int) ([]models.Note, error)
}

func (m *mockService) Create(ctx context.Context, req models.CreateNoteRequest) (*models.Note, error) {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}
	return nil, nil
}

func (m *mockService) Get(ctx context.Context, id string) (*models.Note, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return nil, nil
}

func (m *mockService) Update(ctx context.Context, id string, req models.UpdateNoteRequest) (*models.Note, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, req)
	}
	return nil, nil
}

func (m *mockService) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockService) List(ctx context.Context, limit, offset int) ([]models.Note, error) {
	if m.listFn != nil {
		return m.listFn(ctx, limit, offset)
	}
	return nil, nil
}

func (m *mockService) CreateCollection(ctx context.Context, req models.CreateCollectionRequest) (*models.Collection, error) {
	if m.createCollectionFn != nil {
		return m.createCollectionFn(ctx, req)
	}
	return nil, nil
}

func (m *mockService) GetCollection(ctx context.Context, id string) (*models.Collection, error) {
	if m.getCollectionFn != nil {
		return m.getCollectionFn(ctx, id)
	}
	return nil, nil
}

func (m *mockService) UpdateCollection(ctx context.Context, id string, req models.UpdateCollectionRequest) (*models.Collection, error) {
	if m.updateCollectionFn != nil {
		return m.updateCollectionFn(ctx, id, req)
	}
	return nil, nil
}

func (m *mockService) DeleteCollection(ctx context.Context, id string) error {
	if m.deleteCollectionFn != nil {
		return m.deleteCollectionFn(ctx, id)
	}
	return nil
}

func (m *mockService) ListCollections(ctx context.Context) ([]models.Collection, error) {
	if m.listCollectionsFn != nil {
		return m.listCollectionsFn(ctx)
	}
	return nil, nil
}

func (m *mockService) AddNotesToCollection(ctx context.Context, collectionID string, noteIDs []string) (int, error) {
	if m.addNotesToCollectionFn != nil {
		return m.addNotesToCollectionFn(ctx, collectionID, noteIDs)
	}
	return 0, nil
}

func (m *mockService) RemoveNoteFromCollection(ctx context.Context, collectionID, noteID string) error {
	if m.removeNoteFromCollFn != nil {
		return m.removeNoteFromCollFn(ctx, collectionID, noteID)
	}
	return nil
}

func (m *mockService) ListCollectionNotes(ctx context.Context, collectionID string, limit, offset int) ([]models.Note, error) {
	if m.listCollectionNotesFn != nil {
		return m.listCollectionNotesFn(ctx, collectionID, limit, offset)
	}
	return nil, nil
}

// Helper to create a request with path values (Go 1.22 ServeMux)
func newRequest(method, path string, body any) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func sampleCollection() *models.Collection {
	return &models.Collection{
		ID:          "col-123",
		Name:        "Engineering Docs",
		Description: "Technical documentation",
		NoteCount:   3,
		CreatedAt:   time.Date(2026, 6, 13, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2026, 6, 13, 0, 0, 0, 0, time.UTC),
	}
}

// --- Tests ---

func TestCreateCollection_Success(t *testing.T) {
	mock := &mockService{
		createCollectionFn: func(ctx context.Context, req models.CreateCollectionRequest) (*models.Collection, error) {
			return sampleCollection(), nil
		},
	}
	h := New(mock)

	body := models.CreateCollectionRequest{Name: "Engineering Docs", Description: "Technical documentation"}
	req := newRequest(http.MethodPost, "/api/collections", body)
	w := httptest.NewRecorder()

	h.CreateCollection(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp models.Collection
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Name != "Engineering Docs" {
		t.Errorf("expected name 'Engineering Docs', got '%s'", resp.Name)
	}
}

func TestCreateCollection_MissingName(t *testing.T) {
	mock := &mockService{}
	h := New(mock)

	body := models.CreateCollectionRequest{Name: "", Description: "no name"}
	req := newRequest(http.MethodPost, "/api/collections", body)
	w := httptest.NewRecorder()

	h.CreateCollection(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateCollection_InvalidBody(t *testing.T) {
	mock := &mockService{}
	h := New(mock)

	req := httptest.NewRequest(http.MethodPost, "/api/collections", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateCollection(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetCollection_Success(t *testing.T) {
	mock := &mockService{
		getCollectionFn: func(ctx context.Context, id string) (*models.Collection, error) {
			return sampleCollection(), nil
		},
	}
	h := New(mock)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/collections/{id}", h.GetCollection)

	req := httptest.NewRequest(http.MethodGet, "/api/collections/col-123", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp models.Collection
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.ID != "col-123" {
		t.Errorf("expected id 'col-123', got '%s'", resp.ID)
	}
}

func TestGetCollection_NotFound(t *testing.T) {
	mock := &mockService{
		getCollectionFn: func(ctx context.Context, id string) (*models.Collection, error) {
			return nil, nil
		},
	}
	h := New(mock)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/collections/{id}", h.GetCollection)

	req := httptest.NewRequest(http.MethodGet, "/api/collections/nonexistent", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteCollection_Success(t *testing.T) {
	mock := &mockService{
		deleteCollectionFn: func(ctx context.Context, id string) error {
			return nil
		},
	}
	h := New(mock)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/collections/{id}", h.DeleteCollection)

	req := httptest.NewRequest(http.MethodDelete, "/api/collections/col-123", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestDeleteCollection_NotFound(t *testing.T) {
	mock := &mockService{
		deleteCollectionFn: func(ctx context.Context, id string) error {
			return sql.ErrNoRows
		},
	}
	h := New(mock)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/collections/{id}", h.DeleteCollection)

	req := httptest.NewRequest(http.MethodDelete, "/api/collections/nonexistent", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestListCollections_Success(t *testing.T) {
	mock := &mockService{
		listCollectionsFn: func(ctx context.Context) ([]models.Collection, error) {
			return []models.Collection{*sampleCollection()}, nil
		},
	}
	h := New(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/collections", nil)
	w := httptest.NewRecorder()

	h.ListCollections(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp []models.Collection
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp) != 1 {
		t.Errorf("expected 1 collection, got %d", len(resp))
	}
}

func TestAddNotesToCollection_Success(t *testing.T) {
	mock := &mockService{
		addNotesToCollectionFn: func(ctx context.Context, collectionID string, noteIDs []string) (int, error) {
			return 2, nil
		},
	}
	h := New(mock)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/collections/{id}/notes", h.AddNotesToCollection)

	body := models.AddNotesRequest{NoteIDs: []string{"note-1", "note-2"}}
	req := newRequest(http.MethodPost, "/api/collections/col-123/notes", body)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]int
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["added"] != 2 {
		t.Errorf("expected added=2, got %d", resp["added"])
	}
}

func TestAddNotesToCollection_EmptyNoteIDs(t *testing.T) {
	mock := &mockService{}
	h := New(mock)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/collections/{id}/notes", h.AddNotesToCollection)

	body := models.AddNotesRequest{NoteIDs: []string{}}
	req := newRequest(http.MethodPost, "/api/collections/col-123/notes", body)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRemoveNoteFromCollection_Success(t *testing.T) {
	mock := &mockService{
		removeNoteFromCollFn: func(ctx context.Context, collectionID, noteID string) error {
			return nil
		},
	}
	h := New(mock)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/collections/{id}/notes/{noteId}", h.RemoveNoteFromCollection)

	req := httptest.NewRequest(http.MethodDelete, "/api/collections/col-123/notes/note-1", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestRemoveNoteFromCollection_NotFound(t *testing.T) {
	mock := &mockService{
		removeNoteFromCollFn: func(ctx context.Context, collectionID, noteID string) error {
			return sql.ErrNoRows
		},
	}
	h := New(mock)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/collections/{id}/notes/{noteId}", h.RemoveNoteFromCollection)

	req := httptest.NewRequest(http.MethodDelete, "/api/collections/col-123/notes/note-1", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestListCollectionNotes_Success(t *testing.T) {
	mock := &mockService{
		listCollectionNotesFn: func(ctx context.Context, collectionID string, limit, offset int) ([]models.Note, error) {
			return []models.Note{
				{ID: "note-1", Title: "Test Note", Content: "content", Tags: []string{"go"}},
			}, nil
		},
	}
	h := New(mock)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/collections/{id}/notes", h.ListCollectionNotes)

	req := httptest.NewRequest(http.MethodGet, "/api/collections/col-123/notes", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp []models.Note
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp) != 1 {
		t.Errorf("expected 1 note, got %d", len(resp))
	}
}

func TestListCollectionNotes_DefaultLimit(t *testing.T) {
	var capturedLimit int
	mock := &mockService{
		listCollectionNotesFn: func(ctx context.Context, collectionID string, limit, offset int) ([]models.Note, error) {
			capturedLimit = limit
			return []models.Note{}, nil
		},
	}
	h := New(mock)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/collections/{id}/notes", h.ListCollectionNotes)

	req := httptest.NewRequest(http.MethodGet, "/api/collections/col-123/notes", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if capturedLimit != 20 {
		t.Errorf("expected default limit 20, got %d", capturedLimit)
	}
}
