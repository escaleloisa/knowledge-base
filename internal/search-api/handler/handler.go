package handler

import (
	"net/http"

	"github.com/escaleloisa/knowledge-base/internal/search-api/service"
	"github.com/escaleloisa/knowledge-base/pkg/response"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		response.Error(w, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	params := service.SearchParams{
		Query: q,
		Tags:  r.URL.Query().Get("tags"),
		From:  r.URL.Query().Get("from"),
		To:    r.URL.Query().Get("to"),
	}

	result, err := h.svc.Search(params)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "search failed")
		return
	}

	response.JSON(w, http.StatusOK, result)
}
