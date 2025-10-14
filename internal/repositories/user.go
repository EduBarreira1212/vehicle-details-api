package repositories

import (
	"context"

	"github.com/EduBarreira1212/vehicle-details-api/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *User {
	return &User{db: db}
}

func (repository *User) Create(ctx context.Context, name, email, password string) (*models.User, error) {
	user := &models.User{
		Name:     name,
		Email:    email,
		Password: password,
		History:  datatypes.JSON([]byte("[]")),
	}

	if err := repository.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (repository *User) GetById(ctx context.Context, id uint64) (*models.User, error) {
	var user models.User

	if err := repository.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (repository *User) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	if err := repository.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (repository *User) Update(ctx context.Context, id uint64, name, email string) error {
	var user models.User

	if err := repository.db.WithContext(ctx).
		Model(&user).
		Where("id = ?", id).
		Updates(models.User{
			Name:  name,
			Email: email,
		}).Error; err != nil {
		return err
	}

	return nil
}

func (repository *User) GetPassword(ctx context.Context, id uint64) (string, error) {
	var user models.User

	if err := repository.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return "", err
	}

	return user.Password, nil
}

func (repository *User) UpdatePassword(ctx context.Context, id uint64, newPassword string) error {
	var user models.User

	if err := repository.db.WithContext(ctx).
		Model(&user).
		Where("id = ?", id).
		Updates(models.User{
			Password: newPassword,
		}).Error; err != nil {
		return err
	}

	return nil
}

func (repository *User) Delete(ctx context.Context, id uint64) error {
	if err := repository.db.WithContext(ctx).Delete(&models.User{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repository *User) GetUserHistoryById(ctx context.Context, id uint64) (datatypes.JSON, error) {
	var user models.User

	if err := repository.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}

	if len(user.History) == 0 {
		return datatypes.JSON([]byte("[]")), nil
	}

	return user.History, nil
}

func (repository *User) UpdateUserHistory(ctx context.Context, id uint64, newPlateSearched string) error {
	return repository.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("history", gorm.Expr(
			`COALESCE(
            CASE WHEN jsonb_typeof(history) = 'array' THEN history ELSE '[]'::jsonb END,
            '[]'::jsonb
        ) || to_jsonb(?::text)`,
			newPlateSearched,
		)).
		Error
}
