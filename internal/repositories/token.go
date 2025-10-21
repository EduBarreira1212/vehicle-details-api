package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/EduBarreira1212/vehicle-details-api/internal/models"
	"gorm.io/gorm"
)

type Token struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *Token {
	return &Token{db: db}
}

func (repository *Token) CreatePasswordResetToken(ctx context.Context, reset *models.PasswordResetToken) error {
	return repository.db.Create(reset).Error
}

func (repository *Token) ExpireActiveResetTokens(ctx context.Context, userID uint64) error {
	return repository.db.Model(&models.PasswordResetToken{}).
		Where("user_id = ? AND used_at IS NULL AND expires_at > ?", userID, time.Now()).
		Update("expires_at", time.Now().Add(-1*time.Minute)).
		Error
}

func (repository *Token) FindPasswordResetTokenByHash(ctx context.Context, tokenHash string) (*models.PasswordResetToken, error) {
	var prt models.PasswordResetToken

	if err := repository.db.Where("token_hash = ?", tokenHash).First(&prt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &prt, nil
}

func (repository *Token) ResetPasswordWithToken(ctx context.Context, userID, tokenID uint64, newHashedPassword string) error {
	return repository.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).
			Where("id = ?", userID).
			Update("password", newHashedPassword).Error; err != nil {
			return err
		}

		now := time.Now()
		if err := tx.Model(&models.PasswordResetToken{}).
			Where("id = ?", tokenID).
			Updates(map[string]any{
				"used_at": &now,
			}).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.PasswordResetToken{}).
			Where("user_id = ? AND used_at IS NULL AND expires_at > ?", userID, now).
			Update("expires_at", now.Add(-1*time.Minute)).Error; err != nil {
			return err
		}

		return nil
	})
}
