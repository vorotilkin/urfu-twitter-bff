package services

import (
	"context"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"twitter-bff/domain/models"
)

const defaultPostLimit = 100

type PostsRepository interface {
	Create(ctx context.Context, userID int32, body string) (models.Post, error)
	PostsByUserID(ctx context.Context, userID int32) ([]models.Post, error)
	LatestPosts(ctx context.Context, userIDs []int32, currentUserId, limit int32) ([]models.Post, error)
	PostByID(ctx context.Context, postID int32, userID int32) (models.Post, error)
	CommentsByPostID(ctx context.Context, postID int32) ([]models.Comment, error)
}

type PostsUsersByIDsRepository interface {
	FetchUsersByIDs(ctx context.Context, ids []int32) (map[int32]models.User, error)
}

type PostsService struct {
	repo      PostsRepository
	usersRepo PostsUsersByIDsRepository
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

	usersByID, err := s.usersRepo.FetchUsersByIDs(ctx, []int32{userID})
	if err != nil {
		return models.Post{}, errors.Wrap(err, "get users err")
	}

	post.User = usersByID[userID]

	return post, nil
}

func (s *PostsService) PostsByUserID(ctx context.Context, userID int32) ([]models.Post, error) {
	if userID == 0 {
		return nil, errors.Wrap(models.ErrInvalidArgument, "invalid user id")
	}

	posts, err := s.repo.PostsByUserID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "get posts err")
	}

	if len(posts) == 0 {
		return []models.Post{}, nil
	}

	userIDs := lo.Uniq(lo.Map(posts, func(post models.Post, _ int) int32 {
		return post.UserID
	}))

	usersByID, err := s.usersRepo.FetchUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, errors.Wrap(err, "get users err")
	}

	for i, post := range posts {
		user := usersByID[post.UserID]
		u := models.User{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		}

		posts[i].User = u
	}

	return posts, nil
}

func (s *PostsService) FeedPosts(ctx context.Context, userID int32) ([]models.Post, error) {
	if userID == 0 {
		return []models.Post{}, nil
	}

	users, err := s.usersRepo.FetchUsersByIDs(ctx, []int32{userID})
	if err != nil {
		return nil, errors.Wrap(err, "get current user err")
	}

	currentUser, ok := users[userID]
	if !ok {
		return nil, errors.Wrap(models.ErrNotFound, "current user not found")
	}

	userIDs := make([]int32, len(currentUser.FollowingUserIds), len(currentUser.FollowingUserIds)+1)
	copy(userIDs, currentUser.FollowingUserIds)
	userIDs = append(userIDs, userID)

	posts, err := s.repo.LatestPosts(ctx, userIDs, userID, defaultPostLimit)
	if err != nil {
		return nil, errors.Wrap(err, "feed posts err")
	}

	users, err = s.usersRepo.FetchUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, errors.Wrap(err, "get current user err")
	}

	for i, post := range posts {
		user := users[post.UserID]
		u := models.User{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		}

		posts[i].User = u
	}

	return posts, nil
}

func (s *PostsService) PostByID(ctx context.Context, postID int32, userID int32) (models.Post, error) {
	if postID == 0 {
		return models.Post{}, errors.Wrap(models.ErrInvalidArgument, "invalid post id")
	}

	post, err := s.repo.PostByID(ctx, postID, userID)
	if err != nil {
		return models.Post{}, errors.Wrap(err, "posts repo err")
	}

	userIDs := make([]int32, 0, len(post.Comments)+1)
	seen := make(map[int32]struct{}, len(post.Comments)+1)

	userIDs = append(userIDs, post.UserID)
	seen[post.UserID] = struct{}{}

	for _, comment := range post.Comments {
		if _, ok := seen[comment.ID]; ok {
			continue
		}
		userIDs = append(userIDs, comment.UserID)
		seen[comment.ID] = struct{}{}
	}

	usersByID, err := s.usersRepo.FetchUsersByIDs(ctx, userIDs)
	if err != nil {
		return models.Post{}, errors.Wrap(err, "get users err")
	}

	user := usersByID[post.UserID]
	u := models.User{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
	}

	post.User = u

	for i, comment := range post.Comments {
		user := usersByID[comment.UserID]
		u := models.User{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		}

		post.Comments[i].User = u
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

func NewPostsService(repo PostsRepository, usersRepo PostsUsersByIDsRepository) *PostsService {
	return &PostsService{repo: repo, usersRepo: usersRepo}
}
