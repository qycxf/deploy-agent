package models

import "time"

type User struct {
	ID uint `gorm:"primaryKey"`

	Username     string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

// TableName keeps the table name explicit/stable.
func (User) TableName() string {
	return "users"
}
