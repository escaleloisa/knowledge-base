package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/escaleloisa/knowledge-base/internal/tag-aggregator/consumer"
	"github.com/escaleloisa/knowledge-base/internal/tag-aggregator/handler"
	"github.com/escaleloisa/knowledge-base/pkg/config"
	kafkapkg "github.com/escaleloisa/knowledge-base/pkg/kafka"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatal("invalid redis url:", err)
	}
	rdb := redis.NewClient(opt)
	defer rdb.Close()

	reader := kafkapkg.NewReader(cfg.KafkaBrokers, "tag-aggregator")
	defer reader.Close()

	c := consumer.New(reader, rdb)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start HTTP server for tag queries
	h := handler.New(rdb)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/tags", h.GetTags)

	go func() {
		addr := fmt.Sprintf(":%s", cfg.Port)
		log.Printf("Tag aggregator HTTP on %s", addr)
		http.ListenAndServe(addr, mux)
	}()

	log.Println("Tag aggregator starting...")
	if err := c.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
