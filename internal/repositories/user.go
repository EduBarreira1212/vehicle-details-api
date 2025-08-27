package repositories

import (
	"context"

	"github.com/EduBarreira1212/vehicle-details-api/internal/models"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *User {
	return &User{db: db}
}

func (repository User) Create(ctx context.Context, name, email, password string) (*models.User, error) {
	user := &models.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	if err := repository.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
