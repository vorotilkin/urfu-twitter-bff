package models

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int32
	Name         string
	PasswordHash string
	Username     string
	Email        string
}

func (u *User) SetPasswordHash(password string) error {
	// bcrypt.GenerateFromPassword возвращает хэш пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	u.PasswordHash = string(hash)

	return nil
}
