package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("HOST", "")
	t.Setenv("PORT", "")
	t.Setenv("DB_PATH", "")
	t.Setenv("AUTHORIZED_KEYS_PATH", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Host != "0.0.0.0" {
		t.Fatalf("expected default host, got %q", cfg.Host)
	}
	if cfg.Port != "2222" {
		t.Fatalf("expected default port, got %q", cfg.Port)
	}
	if cfg.DBPath != "./data/shipdeck.sqlite" {
		t.Fatalf("expected default db path, got %q", cfg.DBPath)
	}
	if cfg.AuthorizedKeysPath != "./data/authorized_keys" {
		t.Fatalf("expected default authorized keys path, got %q", cfg.AuthorizedKeysPath)
	}
}

func TestLoadEnvOverrides(t *testing.T) {
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("PORT", "2323")
	t.Setenv("DB_PATH", "/tmp/test.sqlite")
	t.Setenv("AUTHORIZED_KEYS_PATH", "/tmp/authorized_keys")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Host != "127.0.0.1" || cfg.Port != "2323" || cfg.DBPath != "/tmp/test.sqlite" || cfg.AuthorizedKeysPath != "/tmp/authorized_keys" {
		t.Fatalf("unexpected config values: %+v", cfg)
	}
}

func TestEnvOrDefault(t *testing.T) {
	t.Setenv("FOO_TEST_VALUE", "")
	if got := envOrDefault("FOO_TEST_VALUE", "fallback"); got != "fallback" {
		t.Fatalf("expected fallback, got %q", got)
	}

	t.Setenv("FOO_TEST_VALUE", "set")
	if got := envOrDefault("FOO_TEST_VALUE", "fallback"); got != "set" {
		t.Fatalf("expected env value, got %q", got)
	}
}
