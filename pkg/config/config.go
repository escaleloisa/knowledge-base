package config

import (
	"os"
	"strings"
)

type Config struct {
	Port             string
	DatabaseURL      string
	KafkaBrokers     []string
	ElasticsearchURL string
	RedisURL         string
	NoteServiceURL   string
	ImportPath       string
}

func Load() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/knowledgebase?sslmode=disable"),
		KafkaBrokers:     strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		ElasticsearchURL: getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
		RedisURL:         getEnv("REDIS_URL", "redis://localhost:6379"),
		NoteServiceURL:   getEnv("NOTE_SERVICE_URL", "http://localhost:8080"),
		ImportPath:       getEnv("IMPORT_PATH", "./import"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
