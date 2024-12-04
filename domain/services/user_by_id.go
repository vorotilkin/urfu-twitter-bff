package services

import (
	"context"
	"twitter-bff/domain/models"
)

const defaultUsersLimit = 100

type UserByIDRepository interface {
	FetchUsersByIDs(ctx context.Context, ids []int32) (map[int32]models.User, error)
	NewUsers(ctx context.Context, limit int32) ([]models.User, error)
}

type UserByIDService struct {
	repo UserByIDRepository
}

func (s *UserByIDService) UserByID(ctx context.Context, id int32) (models.User, error) {
	usersByID, err := s.repo.FetchUsersByIDs(ctx, []int32{id})
	if err != nil {
		return models.User{}, err
	}

	user, ok := usersByID[id]
	if !ok {
		return models.User{}, models.ErrNotFound
	}

	return user, nil
}

func (s *UserByIDService) NewUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.repo.NewUsers(ctx, defaultUsersLimit)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func NewUserByIDService(repo UserByIDRepository) *UserByIDService {
	return &UserByIDService{repo: repo}
}
