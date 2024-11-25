package services

import (
	"context"
	"twitter-bff/domain/models"
	"twitter-bff/helpers"
)

type CreateRepository interface {
	Create(ctx context.Context, name string, passwordHash string, username string, email string) (models.User, error)
}

type CreateUserService struct {
	repo CreateRepository
}

func (s *CreateUserService) Create(ctx context.Context, name, password, username, email string) (models.User, error) {
	hash, err := helpers.GenerateHash(password)
	if err != nil {
		return models.User{}, err
	}

	return s.repo.Create(ctx, name, hash, username, email)
}

func NewCreateUserService(repo CreateRepository) *CreateUserService {
	return &CreateUserService{repo: repo}
}
