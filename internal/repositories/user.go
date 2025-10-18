package repositories

import (
	"context"

	"github.com/EduBarreira1212/vehicle-details-api/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		History:  []models.History{},
	}

	if err := repository.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (repository *User) GetById(ctx context.Context, id uint64) (*models.PublicUser, error) {
	var user models.PublicUser

	if err := repository.db.WithContext(ctx).
		Model(&models.User{}).
		Select("id", "name", "email").
		First(&user, id).Error; err != nil {
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

func (repository *User) GetUserHistoryById(ctx context.Context, id uint64) ([]models.History, error) {
	var history []models.History

	if err := repository.db.WithContext(ctx).
		Where("user_id = ?", id).
		Find(&history).Error; err != nil {
		return nil, err
	}

	if len(history) == 0 {
		return []models.History{}, nil
	}

	return history, nil
}

func (repository *User) UpdateUserHistory(
	ctx context.Context,
	id uint64,
	newHistory models.History,
) error {
	newHistory.UserID = id

	return repository.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_id"},
				{Name: "plate"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"model":      gorm.Expr("EXCLUDED.model"),
				"updated_at": gorm.Expr("NOW()"),
			}),
		}).
		Create(&newHistory).Error
}
