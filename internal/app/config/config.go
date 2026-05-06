package config

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	SyncURL string
	APIKey  string
	Timeout time.Duration
	Verbose bool
}

func Load() Config {
	return Config{
		SyncURL: getenv("SYNC_URL", ""),
		APIKey:  getenv("API_KEY", ""),
		Timeout: durationEnv("TIMEOUT", 30*time.Second),
		Verbose: boolEnv("VERBOSE", false),
	}
}

func getenv(k, def string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return def
}

func durationEnv(k string, def time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}

func boolEnv(k string, def bool) bool {
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return def
	}
	return v == "true" || v == "1" || v == "yes"
}
