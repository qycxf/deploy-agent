package main

import (
	"log"

	"github.com/qycxf/deploy-agent/config"
	automigratev1 "github.com/qycxf/deploy-agent/internal/db/automigrate/v1"
	"github.com/qycxf/deploy-agent/internal/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := config.GetDatabaseURL()
	if dsn == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("open postgres: %v", err)
	}

	if err := automigratev1.Migrate(gdb, &models.User{}); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	log.Println("migrations applied successfully")
}
