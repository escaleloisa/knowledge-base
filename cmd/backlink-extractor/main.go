package main

import (
	"context"
	"database/sql"
	"log"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	backlinkextractor "github.com/escaleloisa/knowledge-base/internal/backlink-extractor"
	"github.com/escaleloisa/knowledge-base/pkg/config"
	kafkapkg "github.com/escaleloisa/knowledge-base/pkg/kafka"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	reader := kafkapkg.NewReader(cfg.KafkaBrokers, "backlink-extractor")
	defer reader.Close()

	extractor := backlinkextractor.New(reader, db)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Println("Backlink extractor starting...")
	if err := extractor.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
