package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"finance-tracker/pkg/transaction/model/receipt"
)

// Interpret le pasa la imagen a Groq y devuelve lo que pudo extraer.
// categoryNames son las categorías reales del usuario — se le piden en el
// prompt para que elija solo entre esas, nunca invente una nueva.
func (s *Services) Interpret(imageData []byte, categoryNames []string) (receipt.Interpreted, error) {
	payload := map[string]any{
		"model": s.model,
		"messages": []map[string]any{
			{
				"role": "user",
				"content": []map[string]any{
					{"type": "text", "text": buildPrompt(categoryNames)},
					{"type": "image_url", "image_url": map[string]any{
						"url": "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(imageData),
					}},
				},
			},
		},
		"response_format": map[string]any{"type": "json_object"},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return receipt.Interpreted{}, fmt.Errorf("groq: marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, apiBaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return receipt.Interpreted{}, fmt.Errorf("groq: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return receipt.Interpreted{}, fmt.Errorf("groq: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return receipt.Interpreted{}, fmt.Errorf("groq: read response: %w", err)
	}
	if resp.StatusCode >= 300 {
		return receipt.Interpreted{}, fmt.Errorf("groq: unexpected status %d: %s", resp.StatusCode, respBody)
	}

	var parsed chatCompletionResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return receipt.Interpreted{}, fmt.Errorf("groq: decode response: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return receipt.Interpreted{}, fmt.Errorf("groq: respuesta sin choices")
	}

	var data receiptData
	if err := json.Unmarshal([]byte(parsed.Choices[0].Message.Content), &data); err != nil {
		return receipt.Interpreted{}, fmt.Errorf("groq: decode structured data: %w", err)
	}

	return receipt.Interpreted{
		Amount:        data.Amount,
		Concept:       data.Concept,
		CategoryName:  data.CategoryName,
		PaymentMethod: data.PaymentMethod,
	}, nil
}

func buildPrompt(categoryNames []string) string {
	categories := "(el usuario no tiene categorías configuradas)"
	if len(categoryNames) > 0 {
		categories = strings.Join(categoryNames, ", ")
	}

	return fmt.Sprintf(`Analiza esta imagen de un recibo o factura de compra y extrae los siguientes datos. Responde ÚNICAMENTE con un objeto JSON con estas claves exactas, sin texto fuera del JSON:

- "monto": el valor total pagado, como número (sin símbolos de moneda ni separadores de miles). Si no puedes identificarlo con certeza, usa null.
- "concepto": una descripción breve del recibo (por ejemplo, el nombre del comercio). Si no es claro, usa una cadena vacía.
- "categoria": elige EXACTAMENTE una de estas categorías si aplica con claridad, o null si ninguna aplica con certeza: %s
- "medio_de_pago": si la imagen muestra claramente cómo se pagó, elige EXACTAMENTE uno de estos valores: "Cash", "Credit Card", "Debit Card", "Bank". Si no es claro, usa null.`, categories)
}

type chatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type receiptData struct {
	Amount        *float64 `json:"monto"`
	Concept       string   `json:"concepto"`
	CategoryName  string   `json:"categoria"`
	PaymentMethod string   `json:"medio_de_pago"`
}
