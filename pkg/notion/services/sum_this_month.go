package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "time/tzdata" // embebe la base de zonas horarias — no depender de que el SO la tenga instalada

	"finance-tracker/pkg/transaction/model/summary"
	"finance-tracker/pkg/transaction/model/transaction"
)

const bogotaTimeZone = "America/Bogota"

// SumThisMonth suma Amount y cuenta filas por Type (Egreso/Ingreso) sobre
// las páginas de "Expenses" cuya Date cae en el mes calendario actual, en
// hora de Bogotá (ver riesgo de zona horaria en
// specs/features/005-consultar-resumen-gastos/plan.md). No usa la fórmula
// "Add to Month" — es por fila, no hay forma de pedirle a la API un total
// agregado, así que se suma del lado del cliente.
func (s *Services) SumThisMonth() (summary.Monthly, error) {
	loc, err := time.LoadLocation(bogotaTimeZone)
	if err != nil {
		return summary.Monthly{}, fmt.Errorf("notion: cargar zona horaria %s: %w", bogotaTimeZone, err)
	}

	now := time.Now().In(loc)
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)

	var total summary.Monthly
	cursor := ""
	for {
		payload := map[string]any{
			"filter": map[string]any{
				"property": "Date",
				"date": map[string]any{
					"on_or_after": firstOfMonth.Format("2006-01-02"),
				},
			},
		}
		if cursor != "" {
			payload["start_cursor"] = cursor
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return summary.Monthly{}, fmt.Errorf("notion: marshal summary query: %w", err)
		}

		url := apiBaseURL + "/data_sources/" + s.expensesDataSourceID + "/query"
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return summary.Monthly{}, fmt.Errorf("notion: build summary query: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
		req.Header.Set("Notion-Version", apiVersion)
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.httpClient.Do(req)
		if err != nil {
			return summary.Monthly{}, fmt.Errorf("notion: summary query failed: %w", err)
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return summary.Monthly{}, fmt.Errorf("notion: read summary query response: %w", err)
		}
		if resp.StatusCode >= 300 {
			return summary.Monthly{}, fmt.Errorf("notion: unexpected status %d querying summary: %s", resp.StatusCode, respBody)
		}

		var parsed summaryQueryResponse
		if err := json.Unmarshal(respBody, &parsed); err != nil {
			return summary.Monthly{}, fmt.Errorf("notion: decode summary query response: %w", err)
		}

		for _, page := range parsed.Results {
			amount := page.Properties.Amount.Number
			if page.Properties.Type.Select == nil {
				// Filas históricas anteriores a que existiera la propiedad
				// Type quedan sin clasificar — no se pueden sumar a ningún
				// lado (ver .claude/rules/notion-rules.md).
				continue
			}
			switch page.Properties.Type.Select.Name {
			case string(transaction.TypeExpense):
				total.Egreso += amount
				total.CountEgreso++
			case string(transaction.TypeIncome):
				total.Ingreso += amount
				total.CountIngreso++
			}
		}

		if !parsed.HasMore {
			break
		}
		cursor = parsed.NextCursor
	}

	return total, nil
}

type summaryQueryResponse struct {
	Results []struct {
		Properties struct {
			Amount struct {
				Number float64 `json:"number"`
			} `json:"Amount"`
			Type struct {
				Select *struct {
					Name string `json:"name"`
				} `json:"select"`
			} `json:"Type"`
		} `json:"properties"`
	} `json:"results"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor"`
}
