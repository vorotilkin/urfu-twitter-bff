package services

import (
	"context"
	"twitter-bff/domain/models"
	"twitter-bff/helpers"
	"twitter-bff/infrastructure/users"
)

type CreateUserService struct {
	repo *users.Repository
}

func (s *CreateUserService) Create(ctx context.Context, name, password, username, email string) (models.User, error) {
	hash, err := helpers.GenerateHash(password)
	if err != nil {
		return models.User{}, err
	}

	return s.repo.Create(ctx, name, hash, username, email)
}

func NewCreateUserService(repo *users.Repository) *CreateUserService {
	return &CreateUserService{repo: repo}
}
