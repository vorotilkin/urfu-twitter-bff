package services

import (
	"context"
	"github.com/pkg/errors"
	"twitter-bff/domain/models"
)

var ErrLikeUnknown = errors.New("cant Like")

type LikeRepository interface {
	Like(ctx context.Context, userID, postID int32, operationType models.LikeType) (bool, error)
}

type LikeService struct {
	repo LikeRepository
}

func (s *LikeService) Like(ctx context.Context, userID, postID int32, operationType models.LikeType) (bool, error) {
	if userID == 0 || postID == 0 {
		return false, errors.Wrap(models.ErrInvalidArgument, "zero id")
	}

	ok, err := s.repo.Like(ctx, userID, postID, operationType)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, ErrLikeUnknown
	}

	return ok, nil
}

func NewLikeService(repo LikeRepository) *LikeService {
	return &LikeService{
		repo: repo,
	}
}
