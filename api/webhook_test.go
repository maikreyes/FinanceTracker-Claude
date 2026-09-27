package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_ConfiguracionInvalidaResponde500SinFiltrarDetalle(t *testing.T) {
	for _, name := range []string{"TELEGRAM_BOT_TOKEN", "TELEGRAM_WEBHOOK_SECRET", "UPSTASH_REDIS_REST_URL", "UPSTASH_REDIS_REST_TOKEN"} {
		t.Setenv(name, "")
	}

	rec := httptest.NewRecorder()
	Handler(rec, httptest.NewRequest(http.MethodPost, "/api/webhook", strings.NewReader("{}")))

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "TELEGRAM") || strings.Contains(rec.Body.String(), "UPSTASH") {
		t.Errorf("body = %q, el endpoint es público y no debe nombrar variables de entorno", rec.Body.String())
	}
}
