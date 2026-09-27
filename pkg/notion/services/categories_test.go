package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCategoryResolver_Create_CategoriaNueva(t *testing.T) {
	createCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodPost && r.URL.Path == "/pages" {
			createCalled = true
			json.NewEncoder(w).Encode(map[string]any{"id": "new-cat-id"})
			return
		}

		// Query de categorías (Resolve, antes de crear) — vacío, no existe todavía.
		json.NewEncoder(w).Encode(map[string]any{"results": []map[string]any{}, "has_more": false})
	}))
	defer server.Close()

	originalBaseURL := apiBaseURL
	apiBaseURL = server.URL
	defer func() { apiBaseURL = originalBaseURL }()

	resolver := newTestServices()

	id, err := resolver.Create("Mascotas")
	if err != nil {
		t.Fatalf("Create() error inesperado: %v", err)
	}
	if id != "new-cat-id" {
		t.Errorf("id = %q, want %q", id, "new-cat-id")
	}
	if !createCalled {
		t.Error("no se llamó a POST /pages para crear la categoría")
	}
}

func TestCategoryResolver_Create_CategoriaExistenteNoDuplica(t *testing.T) {
	createCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodPost && r.URL.Path == "/pages" {
			createCalled = true
			json.NewEncoder(w).Encode(map[string]any{"id": "should-not-be-created"})
			return
		}

		json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{
				{
					"id": "existing-cat-id",
					"properties": map[string]any{
						"Name": map[string]any{"title": []map[string]any{{"plain_text": "Comida"}}},
					},
				},
			},
			"has_more": false,
		})
	}))
	defer server.Close()

	originalBaseURL := apiBaseURL
	apiBaseURL = server.URL
	defer func() { apiBaseURL = originalBaseURL }()

	resolver := newTestServices()

	id, err := resolver.Create("comida") // distinto casing, misma categoría
	if err != nil {
		t.Fatalf("Create() error inesperado: %v", err)
	}
	if id != "existing-cat-id" {
		t.Errorf("id = %q, want %q (existente, sin crear)", id, "existing-cat-id")
	}
	if createCalled {
		t.Error("no debería haber llamado a POST /pages: la categoría ya existía")
	}
}
