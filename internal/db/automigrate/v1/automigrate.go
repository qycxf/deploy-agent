package v1

import "gorm.io/gorm"

// Migrate performs schema migration for DB schema version v1.
// This implementation uses GORM AutoMigrate so we don't maintain SQL migration files.
func Migrate(gdb *gorm.DB, models ...interface{}) error {
	return gdb.AutoMigrate(models...)
}
