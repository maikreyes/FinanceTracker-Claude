package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_FindLatest_ConCategoria(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if strings.Contains(r.URL.Path, "/pages/cat-page-1") {
			json.NewEncoder(w).Encode(map[string]any{
				"properties": map[string]any{
					"Name": map[string]any{"title": []map[string]any{{"plain_text": "Comida"}}},
				},
			})
			return
		}

		json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{
				{
					"id": "page-latest",
					"properties": map[string]any{
						"Name":           map[string]any{"title": []map[string]any{{"plain_text": "mecato"}}},
						"Amount":         map[string]any{"number": 13000},
						"Type":           map[string]any{"select": map[string]any{"name": "Egreso"}},
						"Payment Method": map[string]any{"select": map[string]any{"name": "Cash"}},
						"Date":           map[string]any{"date": map[string]any{"start": "2026-09-20"}},
						"Category":       map[string]any{"relation": []map[string]any{{"id": "cat-page-1"}}},
					},
				},
			},
		})
	}))
	defer server.Close()

	originalBaseURL := apiBaseURL
	apiBaseURL = server.URL
	defer func() { apiBaseURL = originalBaseURL }()

	client := newTestServices()

	lt, found, err := client.FindLatest()
	if err != nil {
		t.Fatalf("FindLatest() error inesperado: %v", err)
	}
	if !found {
		t.Fatal("found = false, want true")
	}
	if lt.PageID != "page-latest" || lt.Name != "mecato" || lt.Amount != 13000 {
		t.Errorf("lt = %+v, datos básicos no coinciden", lt)
	}
	if lt.CategoryName != "Comida" {
		t.Errorf("CategoryName = %q, want %q", lt.CategoryName, "Comida")
	}
	if lt.Date != "2026-09-20" {
		t.Errorf("Date = %q, want %q", lt.Date, "2026-09-20")
	}
}

func TestClient_FindLatest_SinCategoria(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{
				{
					"id": "page-latest-2",
					"properties": map[string]any{
						"Name":   map[string]any{"title": []map[string]any{{"plain_text": "Ingreso"}}},
						"Amount": map[string]any{"number": 3000000},
						"Type":   map[string]any{"select": map[string]any{"name": "Ingreso"}},
					},
				},
			},
		})
	}))
	defer server.Close()

	originalBaseURL := apiBaseURL
	apiBaseURL = server.URL
	defer func() { apiBaseURL = originalBaseURL }()

	client := newTestServices()

	lt, found, err := client.FindLatest()
	if err != nil {
		t.Fatalf("FindLatest() error inesperado: %v", err)
	}
	if !found {
		t.Fatal("found = false, want true")
	}
	if lt.CategoryName != "" {
		t.Errorf("CategoryName = %q, want vacío", lt.CategoryName)
	}
}

func TestClient_FindLatest_SinMovimientos(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"results": []map[string]any{}})
	}))
	defer server.Close()

	originalBaseURL := apiBaseURL
	apiBaseURL = server.URL
	defer func() { apiBaseURL = originalBaseURL }()

	client := newTestServices()

	_, found, err := client.FindLatest()
	if err != nil {
		t.Fatalf("FindLatest() error inesperado: %v", err)
	}
	if found {
		t.Error("found = true, want false (sin movimientos)")
	}
}
