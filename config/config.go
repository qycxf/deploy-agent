package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/goccy/go-yaml"
)

type ConfigFile struct {
	Database DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode"`
}

var (
	loadOnce  sync.Once
	cachedCfg ConfigFile
	loadErr   error
)

// getConfigPath returns the config file path.
// You can override it via CONFIG_PATH.
func getConfigPath() string {
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		return p
	}
	// Relative to current working directory.
	return filepath.Join(".", "config", "config.yaml")
}

func loadConfig() (ConfigFile, error) {
	loadOnce.Do(func() {
		b, err := os.ReadFile(getConfigPath())
		if err != nil {
			loadErr = err
			return
		}
		loadErr = yaml.Unmarshal(b, &cachedCfg)
	})
	return cachedCfg, loadErr
}

// GetPort returns the service port from environment or default 8080.
func GetPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

// GetDatabaseURL builds Postgres DSN.
// Priority:
// 1) DATABASE_URL (raw DSN string, e.g. postgres://user:pass@host:5432/db?sslmode=disable)
// 2) config/config.yaml (database.* fields)
//
// If it can't build a DSN, it returns empty string and caller can decide.
func GetDatabaseURL() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}

	cfg, err := loadConfig()
	if err != nil {
		return ""
	}

	host := cfg.Database.Host
	user := cfg.Database.User
	password := cfg.Database.Password
	dbname := cfg.Database.Dbname
	port := cfg.Database.Port
	sslmode := cfg.Database.Sslmode

	if host == "" || user == "" || password == "" || dbname == "" {
		return ""
	}
	if port == 0 {
		port = 5432
	}
	if sslmode == "" {
		sslmode = "disable"
	}

	// DSN format accepted by github.com/lib/pq and used by gorm postgres driver:
	// host=... user=... password=... dbname=... port=... sslmode=...
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)
}
