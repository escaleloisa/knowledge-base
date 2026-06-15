# Elasticsearch Full-Text Search

Elasticsearch is a distributed search and analytics engine built on Apache Lucene.

## Why Elasticsearch?

- Full-text search with relevance scoring
- Near real-time indexing
- Scales horizontally
- Rich query DSL

## Index Mapping

```json
{
  "mappings": {
    "properties": {
      "title": { "type": "text", "boost": 2 },
      "content": { "type": "text" },
      "tags": { "type": "keyword" },
      "created_at": { "type": "date" }
    }
  }
}
```

## Search Query Example

```json
{
  "query": {
    "bool": {
      "must": [
        { "multi_match": {
            "query": "docker kubernetes",
            "fields": ["title^2", "content"]
        }}
      ],
      "filter": [
        { "terms": { "tags": ["devops"] }},
        { "range": { "created_at": { "gte": "2026-01-01" }}}
      ]
    }
  },
  "highlight": {
    "fields": { "content": {} }
  }
}
```

## Scoring

- **TF-IDF**: Term frequency × inverse document frequency
- **BM25**: Improved scoring algorithm (default in ES 5+)
- **Boost**: Increase weight of specific fields (title^2)

## Integration with Go

```go
resp, err := http.Post(esURL+"/notes/_search", "application/json", body)
```

See [[kafka-events]] for indexing events into Elasticsearch.
