package services

import (
	"context"
	"twitter-bff/domain/models"
)

type UpdateUserByIDRepository interface {
	UpdateUserByID(ctx context.Context, userToUpdate models.UserOption) (models.User, error)
}

type UpdateUserByIDService struct {
	repo UpdateUserByIDRepository
}

func (s *UpdateUserByIDService) UpdateUserByID(ctx context.Context, userToUpdate models.UserOption) (models.User, error) {
	return s.repo.UpdateUserByID(ctx, userToUpdate)
}

func NewUpdateUserByIDService(repo UpdateUserByIDRepository) *UpdateUserByIDService {
	return &UpdateUserByIDService{repo: repo}
}
