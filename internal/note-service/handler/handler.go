package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/escaleloisa/knowledge-base/internal/note-service/service"
	"github.com/escaleloisa/knowledge-base/pkg/models"
	"github.com/escaleloisa/knowledge-base/pkg/response"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" || req.Content == "" {
		response.Error(w, http.StatusBadRequest, "title and content are required")
		return
	}

	note, err := h.svc.Create(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create note")
		return
	}
	response.JSON(w, http.StatusCreated, note)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	note, err := h.svc.Get(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get note")
		return
	}
	if note == nil {
		response.Error(w, http.StatusNotFound, "note not found")
		return
	}
	response.JSON(w, http.StatusOK, note)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req models.UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	note, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update note")
		return
	}
	if note == nil {
		response.Error(w, http.StatusNotFound, "note not found")
		return
	}
	response.JSON(w, http.StatusOK, note)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.svc.Delete(r.Context(), id)
	if service.IsNotFound(err) {
		response.Error(w, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to delete note")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	notes, err := h.svc.List(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list notes")
		return
	}
	response.JSON(w, http.StatusOK, notes)
}
