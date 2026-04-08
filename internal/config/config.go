package config

import (
	"os"
)

type Config struct {
	Host               string
	Port               string
	DBPath             string
	AuthorizedKeysPath string
}

func Load() (Config, error) {

	cfg := Config{
		Host:               envOrDefault("HOST", "0.0.0.0"),
		Port:               envOrDefault("PORT", "2222"),
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
