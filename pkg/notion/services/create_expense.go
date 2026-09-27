package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"finance-tracker/pkg/transaction/model/transaction"
)

// CreateExpense crea una página nueva en la data source "Expenses" a partir
// de una transaction.Transaction. Devuelve el ID de la página creada.
func (s *Services) CreateExpense(tx transaction.Transaction) (string, error) {
	payload := map[string]any{
		"parent":     map[string]any{"type": "data_source_id", "data_source_id": s.expensesDataSourceID},
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
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Notion-Version", apiVersion)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
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

// mapTransactionToProperties mapea transaction.Transaction al formato de
// properties que espera la API de páginas de Notion, según el esquema real
// de la data source "Expenses" documentado en
// .claude/rules/notion-rules.md. "Add to Month" (formula) queda fuera: es
// de solo lectura, Notion la calcula sola.
func mapTransactionToProperties(tx transaction.Transaction) map[string]any {
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

	if len(tx.ReceiptFileIDs) > 0 {
		files := make([]map[string]any, 0, len(tx.ReceiptFileIDs))
		for _, id := range tx.ReceiptFileIDs {
			files = append(files, map[string]any{
				"type":        "file_upload",
				"file_upload": map[string]any{"id": id},
			})
		}
		properties["Receipt"] = map[string]any{"files": files}
	}

	return properties
}
