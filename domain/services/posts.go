package services

import (
	"context"
	"github.com/pkg/errors"
	"twitter-bff/domain/models"
)

const defaultPostLimit = 100

type PostsRepository interface {
	Create(ctx context.Context, userID int32, body string) (models.Post, error)
	PostsByUserID(ctx context.Context, userID int32) ([]models.Post, error)
	LatestPosts(ctx context.Context, limit int32) ([]models.Post, error)
	PostByID(ctx context.Context, postID int32) (models.Post, error)
	CommentsByPostID(ctx context.Context, postID int32) ([]models.Comment, error)
}

type PostsService struct {
	repo PostsRepository
}

func (s *PostsService) Create(ctx context.Context, userID int32, body string) (models.Post, error) {
	if userID == 0 {
		return models.Post{}, errors.Wrap(models.ErrInvalidArgument, "invalid user id")
	}

	if len(body) == 0 {
		return models.Post{}, errors.Wrap(models.ErrInvalidArgument, "invalid post body")
	}

	post, err := s.repo.Create(ctx, userID, body)
	if err != nil {
		return models.Post{}, errors.Wrap(err, "create repo err")
	}

	return post, nil
}

func (s *PostsService) Posts(ctx context.Context, userID int32) ([]models.Post, error) {
	if userID == 0 {
		posts, err := s.repo.LatestPosts(ctx, defaultPostLimit)
		if err != nil {
			return nil, errors.Wrap(err, "get posts by user err")
		}
		return posts, nil
	}

	posts, err := s.repo.PostsByUserID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "get posts err")
	}

	return posts, nil
}

func (s *PostsService) PostByID(ctx context.Context, postID int32) (models.Post, error) {
	if postID == 0 {
		return models.Post{}, errors.Wrap(models.ErrInvalidArgument, "invalid post id")
	}

	post, err := s.repo.PostByID(ctx, postID)
	if err != nil {
		return models.Post{}, errors.Wrap(err, "posts repo err")
	}

	return post, nil
}

func (s *PostsService) CommentsByPostID(ctx context.Context, postID int32) ([]models.Comment, error) {
	if postID == 0 {
		return nil, errors.Wrap(models.ErrInvalidArgument, "invalid post id")
	}

	comments, err := s.repo.CommentsByPostID(ctx, postID)
	if err != nil {
		return nil, errors.Wrap(err, "posts repo err")
	}

	return comments, nil
}

func NewPostsService(repo PostsRepository) *PostsService {
	return &PostsService{repo: repo}
}
