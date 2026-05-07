package routes

import (
	"net/http"

	"github.com/escaleloisa/knowledge-base/internal/note-service/handler"
)

func Register(mux *http.ServeMux, h *handler.Handler) {
	mux.HandleFunc("POST /api/notes", h.Create)
	mux.HandleFunc("GET /api/notes/{id}", h.Get)
	mux.HandleFunc("PUT /api/notes/{id}", h.Update)
	mux.HandleFunc("DELETE /api/notes/{id}", h.Delete)
	mux.HandleFunc("GET /api/notes", h.List)
}
