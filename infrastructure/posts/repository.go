package posts

import (
	"context"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/vorotilkin/twitter-posts/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"twitter-bff/domain/models"
	"twitter-bff/infrastructure/posts/hydrators"
	"twitter-bff/pkg/grpc"
)

type Repository struct {
	client *grpc.Client
}

func (r *Repository) Create(ctx context.Context, userID int32, body string) (models.Post, error) {
	client := proto.NewPostsClient(r.client.Connection())

	req := proto.CreateRequest{
		UserId: userID,
		Body:   body,
	}

	response, err := client.Create(ctx, &req)
	if err != nil {
		return models.Post{}, err
	}

	return hydrators.DomainPost(response.GetPost()), nil
}

func (r *Repository) PostsByUserID(ctx context.Context, userID int32) ([]models.Post, error) {
	if userID == 0 {
		return nil, errors.Wrap(models.ErrInvalidArgument, "user id = 0")
	}

	client := proto.NewPostsClient(r.client.Connection())

	req := proto.PostsRequest{
		Filters: &proto.PostsRequest_Filters{
			Sort: &proto.FilterByOrder{
				Sort: proto.FilterByOrder_SORT_DESC,
			},
			FilterUsers: &proto.FilterByUserIDs{
				UserIds: []int32{userID},
			},
		},
	}

	response, err := client.Posts(ctx, &req)
	if err != nil {
		return nil, err
	}

	return hydrators.DomainPosts(response.GetPosts()), nil
}

func (r *Repository) LatestPosts(ctx context.Context, userIDs []int32, currentUserId, limit int32) ([]models.Post, error) {
	client := proto.NewPostsClient(r.client.Connection())

	req := proto.PostsRequest{
		Filters: &proto.PostsRequest_Filters{
			FilterUsers: &proto.FilterByUserIDs{
				UserIds: userIDs,
			},
			Pagination: &proto.FilterByPagination{
				PerPage: lo.Ternary(limit != 0, limit, 1000),
			},
			Sort: &proto.FilterByOrder{
				Sort: proto.FilterByOrder_SORT_DESC,
			},
		},
		CurrentUserId: currentUserId,
	}

	response, err := client.Posts(ctx, &req)
	if err != nil {
		return nil, err
	}

	return hydrators.DomainPosts(response.GetPosts()), nil
}

func (r *Repository) PostByID(ctx context.Context, postID int32, userID int32) (models.Post, error) {
	if postID == 0 {
		return models.Post{}, errors.Wrap(models.ErrInvalidArgument, "post id = 0")
	}

	client := proto.NewPostsClient(r.client.Connection())

	req := proto.PostByIDRequest{Id: postID, UserId: userID}

	response, err := client.PostByID(ctx, &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return models.Post{}, models.ErrNotFound
		}
		return models.Post{}, err
	}

	return hydrators.DomainPost(response.GetPost()), nil
}

func (r *Repository) CommentsByPostID(ctx context.Context, postID int32) ([]models.Comment, error) {
	client := proto.NewPostsClient(r.client.Connection())

	req := proto.CommentsByPostIDRequest{PostId: postID}

	response, err := client.CommentsByPostID(ctx, &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return hydrators.DomainComments(response.GetComments()), nil
}

func (r *Repository) Like(ctx context.Context, userID, postID int32, operationType models.LikeType) (bool, error) {
	client := proto.NewPostsClient(r.client.Connection())

	req := proto.LikeRequest{
		PostId:        postID,
		UserId:        userID,
		OperationType: hydrators.ProtoLikeOperationType(operationType),
	}

	response, err := client.Like(ctx, &req)
	if err != nil {
		return false, err
	}

	return response.GetOk(), nil
}

func NewRepository(client *grpc.Client) *Repository {
	return &Repository{
		client: client,
	}
}
