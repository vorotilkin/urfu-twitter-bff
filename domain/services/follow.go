package services

import (
	"context"
	"github.com/pkg/errors"
	"twitter-bff/domain/models"
)

var ErrFollowUnknown = errors.New("cant follow")

type FollowRepository interface {
	Follow(ctx context.Context, userID, targetUserID int32) (bool, error)
	Unfollow(ctx context.Context, userID, targetUserID int32) (bool, error)
	FetchUsersByIDs(ctx context.Context, ids []int32) (map[int32]models.User, error)
}

type FollowService struct {
	repo FollowRepository
}

func (s *FollowService) Follow(ctx context.Context, userID, targetUserID int32) (models.User, error) {
	if err := checkUserIDs(userID, targetUserID); err != nil {
		return models.User{}, err
	}

	ok, err := s.repo.Follow(ctx, userID, targetUserID)
	if err != nil {
		return models.User{}, errors.Wrap(err, "follow repo err")
	}

	if !ok {
		return models.User{}, ErrFollowUnknown
	}

	usersByIDs, err := s.repo.FetchUsersByIDs(ctx, []int32{userID})
	if err != nil {
		return models.User{}, errors.Wrap(err, "user repo err")
	}

	return usersByIDs[userID], nil
}

func (s *FollowService) Unfollow(ctx context.Context, userID, targetUserID int32) (models.User, error) {
	if err := checkUserIDs(userID, targetUserID); err != nil {
		return models.User{}, err
	}

	ok, err := s.repo.Unfollow(ctx, userID, targetUserID)
	if err != nil {
		return models.User{}, errors.Wrap(err, "follow repo err")
	}

	if !ok {
		return models.User{}, ErrFollowUnknown
	}

	usersByIDs, err := s.repo.FetchUsersByIDs(ctx, []int32{userID})
	if err != nil {
		return models.User{}, errors.Wrap(err, "user repo err")
	}

	return usersByIDs[userID], nil
}

func checkUserIDs(userID int32, targetUserID int32) error {
	if userID == 0 || targetUserID == 0 {
		return errors.Wrap(models.ErrInvalidArgument, "zero user id")
	}

	if userID == targetUserID {
		return errors.Wrap(models.ErrInvalidArgument, "equals user ids")
	}

	return nil
}

func NewFollowService(repo FollowRepository) *FollowService {
	return &FollowService{
		repo: repo,
	}
}
