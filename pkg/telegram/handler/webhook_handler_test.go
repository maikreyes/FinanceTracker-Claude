package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const textUpdate = `{"update_id": 7, "message": {"message_id": 1, "chat": {"id": 1}, "text": "egreso mecato 13000 efectivo"}}`

func postWebhook(h *Handler, secret, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/webhook", strings.NewReader(body))
	if secret != "" {
		req.Header.Set(secretTokenHeader, secret)
	}
	rec := httptest.NewRecorder()
	h.WebhookHandler(rec, req)
	return rec
}

func TestWebhook_SecretIncorrectoResponde403SinTocarNada(t *testing.T) {
	h := newHarness(t)

	for _, secret := range []string{"", "otro"} {
		rec := postWebhook(h.handler, secret, textUpdate)
		if rec.Code != http.StatusForbidden {
			t.Errorf("secret %q: status = %d, want 403", secret, rec.Code)
		}
	}
	if h.writer.created != 0 || len(h.messenger.texts) != 0 || len(h.store.seen) != 0 {
		t.Error("un webhook rechazado no debería registrar, responder ni marcar el update")
	}
}

func TestWebhook_SinSecretConfiguradoFallaCerrado(t *testing.T) {
	h := newHarness(t)
	h.handler.WebhookSecret = ""

	if rec := postWebhook(h.handler, "", textUpdate); rec.Code != http.StatusForbidden {
		t.Errorf("status = %d, want 403 con el secreto sin configurar", rec.Code)
	}
}

func TestWebhook_SoloAceptaPost(t *testing.T) {
	h := newHarness(t)
	req := httptest.NewRequest(http.MethodGet, "/api/webhook", nil)
	req.Header.Set(secretTokenHeader, testSecret)
	rec := httptest.NewRecorder()

	h.handler.WebhookHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", rec.Code)
	}
}

func TestWebhook_UpdateValidoProduceLoMismoQueDispatch(t *testing.T) {
	h := newHarness(t)

	rec := postWebhook(h.handler, testSecret, textUpdate)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if h.writer.created != 1 || h.writer.got.Name != "mecato" || h.writer.got.Amount != 13000 {
		t.Errorf("registro = %+v (creados: %d), want mecato 13000", h.writer.got, h.writer.created)
	}
	if !strings.Contains(h.messenger.lastText(), "Registrado") {
		t.Errorf("respuesta = %q, debería confirmar el registro", h.messenger.lastText())
	}
}

func TestWebhook_UpdateRepetidoNoSeProcesaDosVeces(t *testing.T) {
	h := newHarness(t)

	postWebhook(h.handler, testSecret, textUpdate)
	rec := postWebhook(h.handler, testSecret, textUpdate)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200 (Telegram reintenta cualquier otro código)", rec.Code)
	}
	if h.writer.created != 1 {
		t.Errorf("registros = %d, want 1", h.writer.created)
	}
}

func TestWebhook_JSONInvalidoResponde400(t *testing.T) {
	h := newHarness(t)

	if rec := postWebhook(h.handler, testSecret, "no es json"); rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestWebhook_ChatNoAutorizadoResponde200SinRegistrar(t *testing.T) {
	h := newHarness(t)
	body := `{"update_id": 8, "message": {"message_id": 1, "chat": {"id": 999}, "text": "egreso mecato 13000 efectivo"}}`

	rec := postWebhook(h.handler, testSecret, body)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	if h.writer.created != 0 || len(h.messenger.texts) != 0 {
		t.Error("un chat no autorizado no debería registrar ni recibir respuesta")
	}
}

func TestWebhook_CallbackDeBotonResuelveLaFotoPendiente(t *testing.T) {
	h := newHarness(t)
	h.interpreter.result.Amount = amountOf(45000)
	h.sendPhoto("file-1", "")

	body := `{"update_id": 9, "callback_query": {"id": "cq", "data": "tipo:Egreso", "message": {"message_id": 1, "chat": {"id": 1}}}}`
	rec := postWebhook(h.handler, testSecret, body)

	if rec.Code != http.StatusOK || h.writer.created != 1 {
		t.Errorf("status = %d, registros = %d, want 200 y 1", rec.Code, h.writer.created)
	}
}
