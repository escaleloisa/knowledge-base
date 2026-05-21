package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/escaleloisa/knowledge-base/internal/search-api/handler"
	"github.com/escaleloisa/knowledge-base/internal/search-api/service"
	"github.com/escaleloisa/knowledge-base/pkg/config"
)

func main() {
	cfg := config.Load()

	svc := service.New(cfg.ElasticsearchURL)
	h := handler.New(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/search", h.Search)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Search API listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
