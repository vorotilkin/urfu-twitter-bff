package services

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"time"
	"twitter-bff/domain/models"
	"twitter-bff/helpers"
)

type LoginRepository interface {
	FetchUserByEmail(ctx context.Context, email string) (models.User, error)
}

type Config struct {
	SecretKey string
}

type LoginService struct {
	repo   LoginRepository
	config Config
}

func (s *LoginService) Login(ctx context.Context, email, password string) (models.JWTToken, error) {
	user, err := s.repo.FetchUserByEmail(ctx, email)
	if err != nil {
		return models.JWTToken{}, errors.Wrap(err, "failed to fetch user")
	}

	err = helpers.CompareHashAndPassword(user.PasswordHash, password)
	if err != nil {
		return models.JWTToken{}, errors.Wrap(err, "failed to compare password hash")
	}

	// Генерируем полезные данные, которые будут храниться в токене
	payload := jwt.MapClaims{
		"sub": fmt.Sprint(user.ID),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	// Создаем новый JWT-токен и подписываем его по алгоритму HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signedString, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return models.JWTToken{}, errors.Wrap(err, "failed to sign JWT")
	}

	return models.JWTToken{Token: signedString}, nil
}

func NewLoginService(repo LoginRepository, c Config) *LoginService {
	return &LoginService{
		repo:   repo,
		config: c,
	}
}
