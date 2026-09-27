package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"finance-tracker/pkg/transaction/model/summary"
	"finance-tracker/pkg/transaction/model/transaction"
)

// FindLatest devuelve el movimiento registrado más recientemente en
// "Expenses" (por fecha de creación de la página, no por la propiedad
// Date — dos movimientos pueden compartir fecha). Si Category tiene una
// relación, un segundo request resuelve su nombre para mostrar. Un data
// source vacío no es un error: devuelve found=false.
func (s *Services) FindLatest() (summary.Latest, bool, error) {
	payload := map[string]any{
		"sorts":     []map[string]any{{"timestamp": "created_time", "direction": "descending"}},
		"page_size": 1,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return summary.Latest{}, false, fmt.Errorf("notion: marshal latest query: %w", err)
	}

	url := apiBaseURL + "/data_sources/" + s.expensesDataSourceID + "/query"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return summary.Latest{}, false, fmt.Errorf("notion: build latest query: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Notion-Version", apiVersion)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return summary.Latest{}, false, fmt.Errorf("notion: latest query failed: %w", err)
	}
	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return summary.Latest{}, false, fmt.Errorf("notion: read latest query response: %w", err)
	}
	if resp.StatusCode >= 300 {
		return summary.Latest{}, false, fmt.Errorf("notion: unexpected status %d querying latest: %s", resp.StatusCode, respBody)
	}

	var parsed latestQueryResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return summary.Latest{}, false, fmt.Errorf("notion: decode latest query response: %w", err)
	}
	if len(parsed.Results) == 0 {
		return summary.Latest{}, false, nil
	}

	page := parsed.Results[0]
	name := ""
	if len(page.Properties.Name.Title) > 0 {
		name = page.Properties.Name.Title[0].PlainText
	}

	lt := summary.Latest{
		PageID: page.ID,
		Name:   name,
		Amount: page.Properties.Amount.Number,
	}
	if page.Properties.Type.Select != nil {
		lt.Type = transaction.Type(page.Properties.Type.Select.Name)
	}
	if page.Properties.PaymentMethod.Select != nil {
		lt.PaymentMethod = transaction.PaymentMethod(page.Properties.PaymentMethod.Select.Name)
	}
	if page.Properties.Date.Date != nil {
		lt.Date = page.Properties.Date.Date.Start
	}

	if len(page.Properties.Category.Relation) > 0 {
		categoryName, err := s.fetchPageTitle(page.Properties.Category.Relation[0].ID)
		if err != nil {
			return summary.Latest{}, false, fmt.Errorf("notion: resolve category name of latest: %w", err)
		}
		lt.CategoryName = categoryName
	}

	return lt, true, nil
}

type latestQueryResponse struct {
	Results []struct {
		ID         string `json:"id"`
		Properties struct {
			Name struct {
				Title []struct {
					PlainText string `json:"plain_text"`
				} `json:"title"`
			} `json:"Name"`
			Amount struct {
				Number float64 `json:"number"`
			} `json:"Amount"`
			Type struct {
				Select *struct {
					Name string `json:"name"`
				} `json:"select"`
			} `json:"Type"`
			PaymentMethod struct {
				Select *struct {
					Name string `json:"name"`
				} `json:"select"`
			} `json:"Payment Method"`
			Date struct {
				Date *struct {
					Start string `json:"start"`
				} `json:"date"`
			} `json:"Date"`
			Category struct {
				Relation []struct {
					ID string `json:"id"`
				} `json:"relation"`
			} `json:"Category"`
		} `json:"properties"`
	} `json:"results"`
}

// fetchPageTitle devuelve el título (propiedad Name) de una página
// cualquiera de Notion — se usa para resolver el nombre de la categoría
// relacionada al movimiento más reciente.
func (s *Services) fetchPageTitle(pageID string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, apiBaseURL+"/pages/"+pageID, nil)
	if err != nil {
		return "", fmt.Errorf("notion: build get page request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Notion-Version", apiVersion)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("notion: get page failed: %w", err)
	}
	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("notion: read get page response: %w", err)
	}
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("notion: unexpected status %d getting page: %s", resp.StatusCode, respBody)
	}

	var parsed struct {
		Properties struct {
			Name struct {
				Title []struct {
					PlainText string `json:"plain_text"`
				} `json:"title"`
			} `json:"Name"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", fmt.Errorf("notion: decode get page response: %w", err)
	}
	if len(parsed.Properties.Name.Title) == 0 {
		return "", nil
	}
	return parsed.Properties.Name.Title[0].PlainText, nil
}
