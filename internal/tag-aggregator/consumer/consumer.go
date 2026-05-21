package consumer

import (
	"context"
	"fmt"
	"log"

	kafkapkg "github.com/escaleloisa/knowledge-base/pkg/kafka"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	rdb    *redis.Client
}

func New(reader *kafka.Reader, rdb *redis.Client) *Consumer {
	return &Consumer{reader: reader, rdb: rdb}
}

func (c *Consumer) Run(ctx context.Context) error {
	log.Println("Tag aggregator started, consuming events...")
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Printf("read error: %v", err)
			continue
		}

		event, err := kafkapkg.UnmarshalEvent(msg.Value)
		if err != nil {
			log.Printf("unmarshal error: %v", err)
			continue
		}

		switch event.Type {
		case kafkapkg.EventNoteCreated:
			c.incrementTags(ctx, event.Note.Tags)
		case kafkapkg.EventNoteUpdated:
			c.incrementTags(ctx, event.Note.Tags)
		case kafkapkg.EventNoteDeleted:
			// For simplicity, we don't decrement on delete in this version
		}
	}
}

func (c *Consumer) incrementTags(ctx context.Context, tags []string) {
	for _, tag := range tags {
		c.rdb.Incr(ctx, fmt.Sprintf("tag:%s", tag))
		c.rdb.ZIncrBy(ctx, "tags:all", 1, tag)
	}
	log.Printf("updated tags: %v", tags)
}
