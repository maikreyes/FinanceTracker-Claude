package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

const categoryCacheTTL = 5 * time.Minute

// Las categorías se cachean en memoria (ver .claude/rules/notion-rules.md
// y specs/features/001-registrar-movimiento-por-mensaje/plan.md).

type categoryEntry struct {
	name string // nombre real, con su casing original
	id   string // page ID
}

// Resolve busca name (case-insensitive) entre las categorías reales de
// Notion y devuelve el ID de su página.
func (s *Services) Resolve(name string) (string, bool, error) {
	categories, err := s.categories()
	if err != nil {
		return "", false, err
	}
	entry, ok := categories[strings.ToLower(strings.TrimSpace(name))]
	return entry.id, ok, nil
}

// ListNames devuelve los nombres reales de todas las categorías (con su
// casing original), ordenados alfabéticamente — se usa para armar el
// mensaje de error cuando una categoría no se reconoce.
func (s *Services) ListNames() ([]string, error) {
	categories, err := s.categories()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(categories))
	for _, entry := range categories {
		names = append(names, entry.name)
	}
	sort.Strings(names)
	return names, nil
}

// Create crea una categoría nueva en la data source "Category". Si ya
// existe una con ese nombre (case-insensitive, vía el caché de Resolve),
// no crea una página duplicada — devuelve el ID de la existente. Tras
// crear una nueva, invalida el caché para que quede visible de inmediato
// sin esperar categoryCacheTTL.
func (s *Services) Create(name string) (string, error) {
	name = strings.TrimSpace(name)

	if id, ok, err := s.Resolve(name); err != nil {
		return "", err
	} else if ok {
		return id, nil
	}

	payload := map[string]any{
		"parent": map[string]any{"type": "data_source_id", "data_source_id": s.categoryDataSourceID},
		"properties": map[string]any{
			"Name": map[string]any{
				"title": []map[string]any{
					{"text": map[string]any{"content": name}},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("notion: marshal category: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, apiBaseURL+"/pages", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("notion: build create category request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Notion-Version", apiVersion)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("notion: create category failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("notion: read create category response: %w", err)
	}
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("notion: unexpected status %d creating category: %s", resp.StatusCode, respBody)
	}

	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(respBody, &created); err != nil {
		return "", fmt.Errorf("notion: decode create category response: %w", err)
	}

	s.categoriesMu.Lock()
	s.categoriesCache = nil
	s.categoriesMu.Unlock()

	return created.ID, nil
}

func (s *Services) categories() (map[string]categoryEntry, error) {
	s.categoriesMu.Lock()
	defer s.categoriesMu.Unlock()

	if s.categoriesCache != nil && time.Since(s.categoriesCachedAt) < categoryCacheTTL {
		return s.categoriesCache, nil
	}

	fresh, err := s.fetchCategories()
	if err != nil {
		return nil, err
	}
	s.categoriesCache = fresh
	s.categoriesCachedAt = time.Now()
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

func (s *Services) fetchCategories() (map[string]categoryEntry, error) {
	result := make(map[string]categoryEntry)
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

		url := apiBaseURL + "/data_sources/" + s.categoryDataSourceID + "/query"
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("notion: build category query: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
		req.Header.Set("Notion-Version", apiVersion)
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.httpClient.Do(req)
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
			result[strings.ToLower(name)] = categoryEntry{name: name, id: page.ID}
		}

		if !parsed.HasMore {
			break
		}
		cursor = parsed.NextCursor
	}

	return result, nil
}
