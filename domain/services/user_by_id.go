package services

import (
	"context"
	"twitter-bff/domain/models"
)

type UserByIDRepository interface {
	FetchUserByID(ctx context.Context, id int32) (models.User, error)
}

type UserByIDService struct {
	repo UserByIDRepository
}

func (s *UserByIDService) UserByID(ctx context.Context, id int32) (models.User, error) {
	return s.repo.FetchUserByID(ctx, id)
}

func NewUserByIDService(repo UserByIDRepository) *UserByIDService {
	return &UserByIDService{repo: repo}
}
