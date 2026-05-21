package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Service struct {
	esURL string
}

func New(esURL string) *Service {
	return &Service{esURL: esURL}
}

type SearchParams struct {
	Query string
	Tags  string
	From  string
	To    string
}

type SearchResult struct {
	Total   int              `json:"total"`
	Results []map[string]any `json:"results"`
}

func (s *Service) Search(params SearchParams) (*SearchResult, error) {
	query := buildQuery(params)
	body, _ := json.Marshal(query)

	url := fmt.Sprintf("%s/notes/_search", s.esURL)
	resp, err := http.Post(url, "application/json", strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("elasticsearch request: %w", err)
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var esResp esSearchResponse
	json.Unmarshal(data, &esResp)

	results := make([]map[string]any, 0, len(esResp.Hits.Hits))
	for _, hit := range esResp.Hits.Hits {
		result := hit.Source
		result["_score"] = hit.Score
		if hit.Highlight != nil {
			result["_highlight"] = hit.Highlight
		}
		results = append(results, result)
	}

	return &SearchResult{Total: esResp.Hits.Total.Value, Results: results}, nil
}

func buildQuery(p SearchParams) map[string]any {
	must := []map[string]any{
		{"multi_match": map[string]any{
			"query":  p.Query,
			"fields": []string{"title^2", "content"},
		}},
	}

	var filters []map[string]any
	if p.Tags != "" {
		tagList := strings.Split(p.Tags, ",")
		filters = append(filters, map[string]any{"terms": map[string]any{"tags": tagList}})
	}
	if p.From != "" || p.To != "" {
		rangeFilter := map[string]any{}
		if p.From != "" {
			rangeFilter["gte"] = p.From
		}
		if p.To != "" {
			rangeFilter["lte"] = p.To
		}
		filters = append(filters, map[string]any{"range": map[string]any{"created_at": rangeFilter}})
	}

	boolQuery := map[string]any{"must": must}
	if len(filters) > 0 {
		boolQuery["filter"] = filters
	}

	return map[string]any{
		"query":     map[string]any{"bool": boolQuery},
		"highlight": map[string]any{"fields": map[string]any{"content": map[string]any{}}},
	}
}

type esSearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source    map[string]any      `json:"_source"`
			Score     float64             `json:"_score"`
			Highlight map[string][]string `json:"highlight"`
		} `json:"hits"`
	} `json:"hits"`
}
