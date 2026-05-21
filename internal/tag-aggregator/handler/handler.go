package handler

import (
	"net/http"

	"github.com/escaleloisa/knowledge-base/pkg/response"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *Handler {
	return &Handler{rdb: rdb}
}

func (h *Handler) GetTags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	results, err := h.rdb.ZRevRangeWithScores(ctx, "tags:all", 0, -1).Result()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get tags")
		return
	}

	tags := make([]map[string]any, 0, len(results))
	for _, z := range results {
		tags = append(tags, map[string]any{
			"name":  z.Member,
			"count": int(z.Score),
		})
	}
	response.JSON(w, http.StatusOK, tags)
}
