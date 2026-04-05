package config

import "os"

// Config holds all application configuration loaded from the environment.
type Config struct {
	Port         string
	AllowOrigins string
}

// Load reads configuration from environment variables, falling back to defaults.
func Load() Config {
	return Config{
		Port:         getEnv("PORT", "3000"),
		AllowOrigins: getEnv("ALLOW_ORIGINS", "*"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
