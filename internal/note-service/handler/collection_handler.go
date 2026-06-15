package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/escaleloisa/knowledge-base/internal/note-service/service"
	"github.com/escaleloisa/knowledge-base/pkg/models"
	"github.com/escaleloisa/knowledge-base/pkg/response"
)

func (h *Handler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	collection, err := h.svc.CreateCollection(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create collection")
		return
	}
	response.JSON(w, http.StatusCreated, collection)
}

func (h *Handler) GetCollection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	collection, err := h.svc.GetCollection(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get collection")
		return
	}
	if collection == nil {
		response.Error(w, http.StatusNotFound, "collection not found")
		return
	}
	response.JSON(w, http.StatusOK, collection)
}

func (h *Handler) UpdateCollection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req models.UpdateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	collection, err := h.svc.UpdateCollection(r.Context(), id, req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update collection")
		return
	}
	if collection == nil {
		response.Error(w, http.StatusNotFound, "collection not found")
		return
	}
	response.JSON(w, http.StatusOK, collection)
}

func (h *Handler) DeleteCollection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.svc.DeleteCollection(r.Context(), id)
	if service.IsNotFound(err) {
		response.Error(w, http.StatusNotFound, "collection not found")
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to delete collection")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListCollections(w http.ResponseWriter, r *http.Request) {
	collections, err := h.svc.ListCollections(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list collections")
		return
	}
	if collections == nil {
		collections = []models.Collection{}
	}
	response.JSON(w, http.StatusOK, collections)
}

func (h *Handler) AddNotesToCollection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req models.AddNotesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.NoteIDs) == 0 {
		response.Error(w, http.StatusBadRequest, "note_ids is required")
		return
	}

	added, err := h.svc.AddNotesToCollection(r.Context(), id, req.NoteIDs)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to add notes to collection")
		return
	}
	response.JSON(w, http.StatusOK, map[string]int{"added": added})
}

func (h *Handler) RemoveNoteFromCollection(w http.ResponseWriter, r *http.Request) {
	collectionID := r.PathValue("id")
	noteID := r.PathValue("noteId")

	err := h.svc.RemoveNoteFromCollection(r.Context(), collectionID, noteID)
	if service.IsNotFound(err) {
		response.Error(w, http.StatusNotFound, "note not in collection")
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to remove note from collection")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListCollectionNotes(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	notes, err := h.svc.ListCollectionNotes(r.Context(), id, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list collection notes")
		return
	}
	response.JSON(w, http.StatusOK, notes)
}
