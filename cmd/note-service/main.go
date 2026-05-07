package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/escaleloisa/knowledge-base/internal/note-service/handler"
	"github.com/escaleloisa/knowledge-base/internal/note-service/repository"
	"github.com/escaleloisa/knowledge-base/internal/note-service/routes"
	"github.com/escaleloisa/knowledge-base/internal/note-service/service"
	"github.com/escaleloisa/knowledge-base/pkg/config"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.PingContext(context.Background()); err != nil {
		log.Fatal("failed to ping database:", err)
	}

	repo := repository.New(db)
	svc := service.New(repo)
	h := handler.New(svc)

	mux := http.NewServeMux()
	routes.Register(mux, h)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Note service listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
