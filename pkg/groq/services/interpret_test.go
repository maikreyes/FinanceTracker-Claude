package services

import (
	"encoding/json"
	"finance-tracker/pkg/config"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_Interpret_ParseaRespuestaCompleta(t *testing.T) {
	var gotBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)

		content, _ := json.Marshal(map[string]any{
			"monto":         45000,
			"concepto":      "Supermercado XYZ",
			"categoria":     "Comida",
			"medio_de_pago": "Credit Card",
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]any{"content": string(content)}},
			},
		})
	}))
	defer server.Close()

	original := apiBaseURL
	apiBaseURL = server.URL
	defer func() { apiBaseURL = original }()

	client := NewServices(&config.Config{GroqAPIKey: "fake-key"})
	result, err := client.Interpret([]byte("fake-jpeg-bytes"), []string{"Comida", "Ahorro"})
	if err != nil {
		t.Fatalf("Interpret() error inesperado: %v", err)
	}

	if result.Amount == nil || *result.Amount != 45000 {
		t.Errorf("Amount = %v, want 45000", result.Amount)
	}
	if result.Concept != "Supermercado XYZ" {
		t.Errorf("Concept = %q, want %q", result.Concept, "Supermercado XYZ")
	}
	if result.CategoryName != "Comida" {
		t.Errorf("CategoryName = %q, want %q", result.CategoryName, "Comida")
	}
	if result.PaymentMethod != "Credit Card" {
		t.Errorf("PaymentMethod = %q, want %q", result.PaymentMethod, "Credit Card")
	}

	messages, _ := gotBody["messages"].([]any)
	if len(messages) == 0 {
		t.Fatal("request sin messages")
	}
	firstMessage, _ := messages[0].(map[string]any)
	content, _ := firstMessage["content"].([]any)
	firstPart, _ := content[0].(map[string]any)
	promptText, _ := firstPart["text"].(string)
	if !strings.Contains(promptText, "Comida") || !strings.Contains(promptText, "Ahorro") {
		t.Errorf("prompt = %q, debería incluir las categorías reales", promptText)
	}
}

func TestClient_Interpret_SinMonto(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, _ := json.Marshal(map[string]any{
			"monto":         nil,
			"concepto":      "",
			"categoria":     nil,
			"medio_de_pago": nil,
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]any{"content": string(content)}},
			},
		})
	}))
	defer server.Close()

	original := apiBaseURL
	apiBaseURL = server.URL
	defer func() { apiBaseURL = original }()

	client := NewServices(&config.Config{GroqAPIKey: "fake-key"})
	result, err := client.Interpret([]byte("x"), nil)
	if err != nil {
		t.Fatalf("Interpret() error inesperado: %v", err)
	}
	if result.Amount != nil {
		t.Errorf("Amount = %v, want nil", result.Amount)
	}
}

func TestClient_Interpret_ErrorDeAPI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error": "rate limit exceeded"}`))
	}))
	defer server.Close()

	original := apiBaseURL
	apiBaseURL = server.URL
	defer func() { apiBaseURL = original }()

	client := NewServices(&config.Config{GroqAPIKey: "fake-key"})
	_, err := client.Interpret([]byte("x"), nil)
	if err == nil {
		t.Fatal("Interpret() esperaba error, no hubo")
	}
}
