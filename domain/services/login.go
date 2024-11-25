package services

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"time"
	"twitter-bff/domain/models"
	"twitter-bff/helpers"
)

var ErrUserNotFound = errors.New("user not found")

type LoginRepository interface {
	ExistByUserAndPasswordHash(ctx context.Context, username string, passwordHash string) (bool, error)
}

type Config struct {
	SecretKey string
}

type LoginService struct {
	repo   LoginRepository
	config Config
}

func (s *LoginService) Login(ctx context.Context, email, password string) (models.JWTToken, error) {
	hash, err := helpers.GenerateHash(password)
	if err != nil {
		return models.JWTToken{}, err
	}

	exist, err := s.repo.ExistByUserAndPasswordHash(ctx, email, hash)
	if err != nil {
		return models.JWTToken{}, errors.Wrap(err, "failed to check existence of user")
	}
	if !exist {
		return models.JWTToken{}, ErrUserNotFound
	}

	// Генерируем полезные данные, которые будут храниться в токене
	payload := jwt.MapClaims{
		"sub": email,
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
