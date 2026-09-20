package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"finance-tracker/internal/domain"
)

const (
	apiBaseURL = "https://api.notion.com/v1"
	// "Expenses" es una multi-source database — parent debe ser
	// data_source_id, no database_id. Validado end-to-end con la REST API
	// real (token de integración, no solo MCP) el 2026-09-20 (ver
	// .claude/history/2026-09-20-notion-client-real-end-to-end.md).
	apiVersion = "2025-09-03"
)

// Client es el cliente HTTP de la API de Notion, sin dependencias externas.
type Client struct {
	apiKey       string // NOTION_API_KEY
	dataSourceID string // NOTION_DATABASE_ID — ID de la data source "Expenses"
	httpClient   *http.Client
}

func NewClient(apiKey, dataSourceID string) *Client {
	return &Client{
		apiKey:       apiKey,
		dataSourceID: dataSourceID,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

// CreateExpense crea una página nueva en la data source "Expenses" a partir
// de una domain.Transaction. Devuelve el ID de la página creada.
func (c *Client) CreateExpense(tx domain.Transaction) (string, error) {
	payload := map[string]any{
		"parent":     map[string]any{"type": "data_source_id", "data_source_id": c.dataSourceID},
		"properties": mapTransactionToProperties(tx),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("notion: marshal transaction: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, apiBaseURL+"/pages", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("notion: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Notion-Version", apiVersion)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("notion: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("notion: read response: %w", err)
	}

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("notion: unexpected status %d: %s", resp.StatusCode, respBody)
	}

	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(respBody, &created); err != nil {
		return "", fmt.Errorf("notion: decode response: %w", err)
	}

	return created.ID, nil
}

// mapTransactionToProperties mapea domain.Transaction al formato de
// properties que espera la API de páginas de Notion, según el esquema real
// de la data source "Expenses" documentado en
// .claude/rules/notion-rules.md. "Add to Month" (formula) queda fuera: es
// de solo lectura, Notion la calcula sola.
func mapTransactionToProperties(tx domain.Transaction) map[string]any {
	properties := map[string]any{
		"Name": map[string]any{
			"title": []map[string]any{
				{"text": map[string]any{"content": tx.Name}},
			},
		},
		"Amount": map[string]any{
			"number": tx.Amount,
		},
		"Date": map[string]any{
			"date": map[string]any{"start": tx.Date.Format("2006-01-02")},
		},
	}

	if tx.Type != "" {
		properties["Type"] = map[string]any{
			"select": map[string]any{"name": string(tx.Type)},
		}
	}

	if tx.PaymentMethod != "" {
		properties["Payment Method"] = map[string]any{
			"select": map[string]any{"name": string(tx.PaymentMethod)},
		}
	}

	if tx.Notes != "" {
		properties["Notes"] = map[string]any{
			"rich_text": []map[string]any{
				{"text": map[string]any{"content": tx.Notes}},
			},
		}
	}

	if tx.CategoryID != "" {
		properties["Category"] = map[string]any{
			"relation": []map[string]any{{"id": tx.CategoryID}},
		}
	}

	return properties
}
