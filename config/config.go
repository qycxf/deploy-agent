package config

import "os"

// GetPort returns the service port from environment or default 8080.
func GetPort() string {
    if p := os.Getenv("PORT"); p != "" {
        return p
    }
    return "8080"
}
