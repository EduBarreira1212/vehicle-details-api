package models

import (
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id,omitempty"`
	Name      string         `json:"name,omitempty"`
	Email     string         `gorm:"uniqueIndex;size:180;not null" json:"email,omitempty"`
	Password  string         `json:"password,omitempty"`
	History   datatypes.JSON `json:"history,omitempty"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (user *User) Prepare(step string) error {
	if err := user.validate(step); err != nil {
		return err
	}

	if err := user.format(); err != nil {
		return err
	}

	return nil
}

func (user *User) validate(step string) error {
	if user.Name == "" {
		return errors.New("name is mandatory")
	}

	if user.Email == "" {
		return errors.New("e-mail is mandatory")
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return errors.New("e-mail invalid")
	}

	if step == "register" && user.Password == "" {
		return errors.New("password is mandatory")
	}

	return nil
}

func (user *User) format() error {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	return nil
}
