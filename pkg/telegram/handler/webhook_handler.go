package handler

import (
	"crypto/subtle"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"finance-tracker/pkg/telegram/model/update"
)

const (
	secretTokenHeader = "X-Telegram-Bot-Api-Secret-Token"
	maxWebhookBody    = 1 << 20
)

// WebhookHandler recibe el POST que manda Telegram por cada update. Con el
// secreto incorrecto responde 403 sin tocar el store ni ninguna API.
//
// Responde 200 aunque el procesamiento falle: Telegram reintenta cualquier
// respuesta que no sea 2xx, y reprocesar un update que ya registró un
// movimiento lo registraría dos veces. Los fallos se informan al usuario
// por el propio chat.
func (h *Handler) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !h.validSecret(r.Header.Get(secretTokenHeader)) {
		log.Printf("webhook rechazado: secret token inválido")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, maxWebhookBody))
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var u update.Update
	if err := json.Unmarshal(body, &u); err != nil {
		http.Error(w, "Invalid update", http.StatusBadRequest)
		return
	}

	firstTime, err := h.Store.MarkUpdateSeen(u.UpdateID)
	if err != nil {
		// Sin poder deduplicar es mejor procesar (posible duplicado) que
		// perder el update.
		log.Printf("error marcando update %d: %v", u.UpdateID, err)
	} else if !firstTime {
		log.Printf("update %d repetido, se ignora", u.UpdateID)
		w.WriteHeader(http.StatusOK)
		return
	}

	h.Dispatch(r.Context(), u)
	w.WriteHeader(http.StatusOK)
}

// validSecret compara en tiempo constante. Un secreto no configurado nunca
// es válido: el webhook falla cerrado.
func (h *Handler) validSecret(got string) bool {
	if h.WebhookSecret == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(h.WebhookSecret)) == 1
}
