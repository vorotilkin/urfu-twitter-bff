package models

import (
	"github.com/pkg/errors"
	"time"
)

var (
	ErrInvalidUserID = errors.New("invalid user id")
	ErrTokenExpired  = errors.New("token expired")
)

const JWTCookieName = "user-jwt"

type JWTToken struct {
	Token string
}

type JWTUser struct {
	UserID    int32
	ExpiredAt time.Time
}

func (u JWTUser) IsOK() error {
	if u.UserID <= 0 {
		return ErrInvalidUserID
	}

	if time.Now().After(u.ExpiredAt) {
		return ErrTokenExpired
	}

	return nil
}
