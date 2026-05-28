package db

import (
	"fmt"

	"github.com/qycxf/deploy-agent/config"
	v1 "github.com/qycxf/deploy-agent/internal/db/automigrate/v1"
	"github.com/qycxf/deploy-agent/internal/db/models"
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

	// 按数据库版本迭代进行 AutoMigrate（先实现 v1）
	if err := v1.Migrate(gdb, &models.User{}); err != nil {
		return nil, fmt.Errorf("auto-migrate v1: %w", err)
	}

	return &DB{DB: gdb}, nil
}
