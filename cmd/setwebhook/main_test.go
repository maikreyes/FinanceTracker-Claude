package main

import "testing"

func TestValidateWebhookURL(t *testing.T) {
	cases := map[string]bool{
		"https://mi-proyecto.vercel.app/api/webhook": true,
		"http://mi-proyecto.vercel.app/api/webhook":  false, // Telegram solo entrega por HTTPS
		"mi-proyecto.vercel.app/api/webhook":         false,
		"https://":                                   false,
		"":                                           false,
	}
	for raw, wantOK := range cases {
		if err := validateWebhookURL(raw); (err == nil) != wantOK {
			t.Errorf("validateWebhookURL(%q) = %v, want ok=%v", raw, err, wantOK)
		}
	}
}
