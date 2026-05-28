package config

import "os"

// GetPort returns the service port from environment or default 8080.
func GetPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

// GetDatabaseURL returns DATABASE_URL from environment.
// Example: postgres://user:pass@host:5432/db?sslmode=disable
func GetDatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}
