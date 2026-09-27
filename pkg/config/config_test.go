package config

import (
	"reflect"
	"testing"
)

func TestParseAllowedChatIDs(t *testing.T) {
	got := ParseAllowedChatIDs(" 111, 222 ,,abc,333")
	want := []int64{111, 222, 333}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseAllowedChatIDs = %v, want %v", got, want)
	}
}

func TestParseAllowedChatIDs_Empty(t *testing.T) {
	if got := ParseAllowedChatIDs(""); len(got) != 0 {
		t.Errorf("ParseAllowedChatIDs(\"\") = %v, want vacío", got)
	}
}

func TestCheckWebhook(t *testing.T) {
	full := Config{TelegramToken: "t", WebhookSecret: "s", UpstashURL: "u", UpstashToken: "k"}
	if err := full.CheckWebhook(); err != nil {
		t.Fatalf("CheckWebhook() con todo seteado = %v, want nil", err)
	}

	for name, mutate := range map[string]func(*Config){
		"sin token":            func(c *Config) { c.TelegramToken = "" },
		"sin secreto":          func(c *Config) { c.WebhookSecret = "" },
		"sin url":              func(c *Config) { c.UpstashURL = "" },
		"sin token de upstash": func(c *Config) { c.UpstashToken = "" },
	} {
		t.Run(name, func(t *testing.T) {
			c := full
			mutate(&c)
			if err := c.CheckWebhook(); err == nil {
				t.Error("CheckWebhook() = nil, want error")
			}
		})
	}
}
