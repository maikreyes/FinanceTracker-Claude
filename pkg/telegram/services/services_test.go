package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const secretToken = "123456:SECRET-TOKEN"

// closedServerURL devuelve la URL de un servidor ya cerrado, para forzar
// un error de conexión.
func closedServerURL(t *testing.T) string {
	t.Helper()
	server := httptest.NewServer(http.NotFoundHandler())
	url := server.URL
	server.Close()
	return url
}

func TestCall_ErrorDeRedNoFiltraElToken(t *testing.T) {
	b := newServices(secretToken, closedServerURL(t))

	err := b.SendMessage(context.Background(), 1, "hola")
	if err == nil {
		t.Fatal("SendMessage contra un servidor caído debería fallar")
	}
	if strings.Contains(err.Error(), "SECRET-TOKEN") {
		t.Errorf("el error filtra el token: %v", err)
	}
}

func TestDownloadFile_ErrorDeRedNoFiltraElToken(t *testing.T) {
	b := newServices(secretToken, closedServerURL(t))

	_, err := b.DownloadFile(context.Background(), "photos/file_1.jpg")
	if err == nil {
		t.Fatal("DownloadFile contra un servidor caído debería fallar")
	}
	if strings.Contains(err.Error(), "SECRET-TOKEN") {
		t.Errorf("el error filtra el token: %v", err)
	}
}

func TestDownloadPhoto_ResuelveFileIDYDescarga(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/getFile"):
			w.Write([]byte(`{"ok":true,"result":{"file_path":"photos/file_1.jpg"}}`))
		case strings.Contains(r.URL.Path, "/file/bot"):
			w.Write([]byte("fake-image-bytes"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	got, err := newServices("tok", server.URL).DownloadPhoto(context.Background(), "file-1")
	if err != nil {
		t.Fatalf("DownloadPhoto error inesperado: %v", err)
	}
	if string(got) != "fake-image-bytes" {
		t.Errorf("DownloadPhoto = %q, want %q", got, "fake-image-bytes")
	}
}

func TestSetWebhook_EnviaUrlYSecreto(t *testing.T) {
	var got url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/setWebhook") {
			http.NotFound(w, r)
			return
		}
		got = r.URL.Query()
		w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	err := newServices("tok", server.URL).SetWebhook(context.Background(), "https://example.com/api/webhook", "s3cret")
	if err != nil {
		t.Fatalf("SetWebhook error inesperado: %v", err)
	}
	if got.Get("url") != "https://example.com/api/webhook" || got.Get("secret_token") != "s3cret" {
		t.Errorf("query = %v, want url y secret_token", got)
	}
}

func TestDeleteWebhook_DescartaPendientesSoloSiSePide(t *testing.T) {
	var got url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.URL.Query()
		w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()
	s := newServices("tok", server.URL)

	if err := s.DeleteWebhook(context.Background(), false); err != nil {
		t.Fatalf("DeleteWebhook error inesperado: %v", err)
	}
	if got.Get("drop_pending_updates") != "" {
		t.Errorf("query = %v, no debería descartar pendientes", got)
	}

	if err := s.DeleteWebhook(context.Background(), true); err != nil {
		t.Fatalf("DeleteWebhook error inesperado: %v", err)
	}
	if got.Get("drop_pending_updates") != "true" {
		t.Errorf("query = %v, want drop_pending_updates=true", got)
	}
}
