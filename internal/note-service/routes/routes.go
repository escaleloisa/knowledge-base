package routes

import (
	"net/http"

	"github.com/escaleloisa/knowledge-base/internal/note-service/handler"
)

func Register(mux *http.ServeMux, h *handler.Handler) {
	// Notes
	mux.HandleFunc("POST /api/notes", h.Create)
	mux.HandleFunc("GET /api/notes/{id}", h.Get)
	mux.HandleFunc("PUT /api/notes/{id}", h.Update)
	mux.HandleFunc("DELETE /api/notes/{id}", h.Delete)
	mux.HandleFunc("GET /api/notes", h.List)

	// Collections
	mux.HandleFunc("POST /api/collections", h.CreateCollection)
	mux.HandleFunc("GET /api/collections", h.ListCollections)
	mux.HandleFunc("GET /api/collections/{id}", h.GetCollection)
	mux.HandleFunc("PUT /api/collections/{id}", h.UpdateCollection)
	mux.HandleFunc("DELETE /api/collections/{id}", h.DeleteCollection)

	// Collection-Note relationships
	mux.HandleFunc("GET /api/collections/{id}/notes", h.ListCollectionNotes)
	mux.HandleFunc("POST /api/collections/{id}/notes", h.AddNotesToCollection)
	mux.HandleFunc("DELETE /api/collections/{id}/notes/{noteId}", h.RemoveNoteFromCollection)

}
