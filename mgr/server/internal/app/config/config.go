package config

import "os"

type Config struct {
	HTTPAddr string
	WebRoot  string
}

func Load() Config {
	return Config{
		HTTPAddr: envOrDefault("CLEAR_BILL_HTTP_ADDR", ":8080"),
		WebRoot:  envOrDefault("CLEAR_BILL_WEB_ROOT", "./website"),
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
