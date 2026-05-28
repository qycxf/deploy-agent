package repository

import (
	"errors"

	"github.com/qycxf/deploy-agent/internal/db/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	gdb *gorm.DB
}

func NewUserRepository(gdb *gorm.DB) *UserRepository {
	return &UserRepository{gdb: gdb}
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var u models.User
	err := r.gdb.Where("username = ?", username).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// EnsureDemoUser creates a demo user if it doesn't exist yet.
// username: demo, password: secret
func (r *UserRepository) EnsureDemoUser() error {
	const username = "demo"
	const password = "secret"

	u, err := r.FindByUsername(username)
	if err != nil {
		return err
	}
	if u != nil {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Username:     username,
		PasswordHash: string(hash),
	}

	return r.gdb.Create(&user).Error
}

// VerifyPassword compares plaintext password with stored bcrypt hash.
func VerifyPassword(password string, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}
