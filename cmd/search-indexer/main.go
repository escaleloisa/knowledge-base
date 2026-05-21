package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	searchindexer "github.com/escaleloisa/knowledge-base/internal/search-indexer"
	"github.com/escaleloisa/knowledge-base/pkg/config"
	kafkapkg "github.com/escaleloisa/knowledge-base/pkg/kafka"
)

func main() {
	cfg := config.Load()

	reader := kafkapkg.NewReader(cfg.KafkaBrokers, "search-indexer")
	defer reader.Close()

	indexer := searchindexer.New(reader, cfg.ElasticsearchURL)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Println("Search indexer starting...")
	if err := indexer.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
