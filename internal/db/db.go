package db

import (
	"fmt"

	"github.com/qycxf/deploy-agent/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func New() (*DB, error) {
	dsn := config.GetDatabaseURL()
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is empty")
	}

	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	return &DB{DB: gdb}, nil
}
