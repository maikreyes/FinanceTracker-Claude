package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const categoryCacheTTL = 5 * time.Minute

// CategoryResolver resuelve nombres de categoría contra la data source
// "Category" real de Notion, cacheando el resultado en memoria (ver
// .claude/rules/notion-rules.md y
// specs/features/001-registrar-movimiento-por-mensaje/plan.md). Implementa
// usecase.CategoryResolver.
type CategoryResolver struct {
	apiKey       string
	dataSourceID string // data source "Category"
	httpClient   *http.Client

	mu       sync.Mutex
	cache    map[string]string // nombre en minúsculas -> page ID
	cachedAt time.Time
}

func NewCategoryResolver(apiKey, dataSourceID string) *CategoryResolver {
	return &CategoryResolver{
		apiKey:       apiKey,
		dataSourceID: dataSourceID,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

// Resolve busca name (case-insensitive) entre las categorías reales de
// Notion y devuelve el ID de su página.
func (r *CategoryResolver) Resolve(name string) (string, bool, error) {
	categories, err := r.categories()
	if err != nil {
		return "", false, err
	}
	id, ok := categories[strings.ToLower(strings.TrimSpace(name))]
	return id, ok, nil
}

func (r *CategoryResolver) categories() (map[string]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.cache != nil && time.Since(r.cachedAt) < categoryCacheTTL {
		return r.cache, nil
	}

	fresh, err := r.fetchCategories()
	if err != nil {
		return nil, err
	}
	r.cache = fresh
	r.cachedAt = time.Now()
	return fresh, nil
}

type categoryQueryResponse struct {
	Results []struct {
		ID         string `json:"id"`
		Properties struct {
			Name struct {
				Title []struct {
					PlainText string `json:"plain_text"`
				} `json:"title"`
			} `json:"Name"`
		} `json:"properties"`
	} `json:"results"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor"`
}

func (r *CategoryResolver) fetchCategories() (map[string]string, error) {
	result := make(map[string]string)
	cursor := ""

	for {
		payload := map[string]any{}
		if cursor != "" {
			payload["start_cursor"] = cursor
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("notion: marshal category query: %w", err)
		}

		url := apiBaseURL + "/data_sources/" + r.dataSourceID + "/query"
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("notion: build category query: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+r.apiKey)
		req.Header.Set("Notion-Version", apiVersion)
		req.Header.Set("Content-Type", "application/json")

		resp, err := r.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("notion: category query failed: %w", err)
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("notion: read category query response: %w", err)
		}
		if resp.StatusCode >= 300 {
			return nil, fmt.Errorf("notion: unexpected status %d querying categories: %s", resp.StatusCode, respBody)
		}

		var parsed categoryQueryResponse
		if err := json.Unmarshal(respBody, &parsed); err != nil {
			return nil, fmt.Errorf("notion: decode category query response: %w", err)
		}

		for _, page := range parsed.Results {
			if len(page.Properties.Name.Title) == 0 {
				continue
			}
			name := page.Properties.Name.Title[0].PlainText
			result[strings.ToLower(name)] = page.ID
		}

		if !parsed.HasMore {
			break
		}
		cursor = parsed.NextCursor
	}

	return result, nil
}
