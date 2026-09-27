package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_SumThisMonth_PaginaYSumaPorTipo(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")

		if calls == 1 {
			json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{
					{"properties": map[string]any{
						"Amount": map[string]any{"number": 1000},
						"Type":   map[string]any{"select": map[string]any{"name": "Egreso"}},
					}},
					{"properties": map[string]any{
						"Amount": map[string]any{"number": 500},
						"Type":   map[string]any{"select": nil}, // fila histórica sin Type
					}},
				},
				"has_more":    true,
				"next_cursor": "cursor-2",
			})
			return
		}

		json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{
				{"properties": map[string]any{
					"Amount": map[string]any{"number": 2000},
					"Type":   map[string]any{"select": map[string]any{"name": "Ingreso"}},
				}},
			},
			"has_more": false,
		})
	}))
	defer server.Close()

	originalBaseURL := apiBaseURL
	apiBaseURL = server.URL
	defer func() { apiBaseURL = originalBaseURL }()

	client := newTestServices()

	summary, err := client.SumThisMonth()
	if err != nil {
		t.Fatalf("SumThisMonth() error inesperado: %v", err)
	}
	if summary.Egreso != 1000 {
		t.Errorf("Egreso = %v, want 1000", summary.Egreso)
	}
	if summary.Ingreso != 2000 {
		t.Errorf("Ingreso = %v, want 2000", summary.Ingreso)
	}
	if summary.CountEgreso != 1 {
		t.Errorf("CountEgreso = %d, want 1", summary.CountEgreso)
	}
	if summary.CountIngreso != 1 {
		t.Errorf("CountIngreso = %d, want 1", summary.CountIngreso)
	}
	if summary.Balance() != 1000 {
		t.Errorf("Balance() = %v, want 1000 (2000 ingreso - 1000 egreso)", summary.Balance())
	}
	if calls != 2 {
		t.Errorf("calls = %d, want 2 (paginación)", calls)
	}
}

func TestClient_SumThisMonth_SinMovimientos(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"results":  []map[string]any{},
			"has_more": false,
		})
	}))
	defer server.Close()

	originalBaseURL := apiBaseURL
	apiBaseURL = server.URL
	defer func() { apiBaseURL = originalBaseURL }()

	client := newTestServices()

	summary, err := client.SumThisMonth()
	if err != nil {
		t.Fatalf("SumThisMonth() error inesperado: %v", err)
	}
	if summary.Egreso != 0 || summary.Ingreso != 0 {
		t.Errorf("Egreso=%v Ingreso=%v, want 0 y 0", summary.Egreso, summary.Ingreso)
	}
	if summary.CountEgreso != 0 || summary.CountIngreso != 0 {
		t.Errorf("CountEgreso=%d CountIngreso=%d, want 0 y 0", summary.CountEgreso, summary.CountIngreso)
	}
	if summary.Balance() != 0 {
		t.Errorf("Balance() = %v, want 0", summary.Balance())
	}
}
