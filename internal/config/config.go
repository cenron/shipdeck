package config

import (
	"os"
)

type Config struct {
	DBPath             string
	AuthorizedKeysPath string
}

func Load() (Config, error) {

	cfg := Config{
		DBPath:             envOrDefault("DB_PATH", "./data/shipdeck.sqlite"),
		AuthorizedKeysPath: envOrDefault("AUTHORIZED_KEYS_PATH", "./data/authorized_keys"),
	}
	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
