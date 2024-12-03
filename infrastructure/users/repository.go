package users

import (
	"context"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/vorotilkin/twitter-users/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"twitter-bff/domain/models"
	"twitter-bff/infrastructure/users/hydrators"
	"twitter-bff/pkg/grpc"
)

type Repository struct {
	client *grpc.Client
}

func (r *Repository) FetchUsersByIDs(ctx context.Context, ids []int32) (map[int32]models.User, error) {
	client := proto.NewUsersClient(r.client.Connection())

	req := proto.UsersByIDsRequest{Ids: ids}

	response, err := client.UsersByIDs(ctx, &req)
	if status.Code(err) == codes.NotFound {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "FetchUsersByIDs")
	}

	return lo.SliceToMap(response.GetUsers(), func(user *proto.User) (int32, models.User) {
		u := hydrators.DomainUser(user)
		return u.ID, u
	}), nil
}

func (r *Repository) FetchUserByEmail(ctx context.Context, email string) (models.User, error) {
	client := proto.NewUsersClient(r.client.Connection())

	req := proto.UserByEmailRequest{Email: email}

	response, err := client.UserByEmail(ctx, &req)
	if err != nil {
		return models.User{}, errors.Wrap(err, "FetchUserByEmail")
	}

	return hydrators.DomainUser(response.GetUser()), nil
}

func (r *Repository) FetchPasswordHashByEmail(ctx context.Context, email string) (string, error) {
	client := proto.NewUsersClient(r.client.Connection())

	req := proto.PasswordHashByEmailRequest{Email: email}

	response, err := client.PasswordHashByEmail(ctx, &req)
	if err != nil {
		return "", errors.Wrap(err, "FetchPasswordHashByEmail")
	}

	return response.GetPasswordHash(), nil
}

func (r *Repository) Create(ctx context.Context, name, passwordHash, username, email string) (models.User, error) {
	client := proto.NewUsersClient(r.client.Connection())

	req := proto.CreateRequest{
		Name:         name,
		PasswordHash: passwordHash,
		Username:     username,
		Email:        email,
	}

	protoUser, err := client.Create(ctx, &req)
	if err != nil {
		return models.User{}, err
	}

	return hydrators.DomainUser(protoUser.GetUser()), nil
}

func (r *Repository) UpdateUserByID(ctx context.Context, userToUpdate models.UserOption) (models.User, error) {
	client := proto.NewUsersClient(r.client.Connection())

	req := proto.UpdateByIDRequest{
		Id:           userToUpdate.ID,
		Name:         userToUpdate.Name.ToPointer(),
		Username:     userToUpdate.Username.ToPointer(),
		Bio:          userToUpdate.Bio.ToPointer(),
		ProfileImage: userToUpdate.ProfileImage.ToPointer(),
		CoverImage:   userToUpdate.CoverImage.ToPointer(),
	}

	response, err := client.UpdateByID(ctx, &req)
	if status.Code(err) == codes.NotFound {
		return models.User{}, models.ErrNotFound
	}
	if err != nil {
		return models.User{}, errors.Wrap(err, "UpdateUserByID")
	}

	return hydrators.DomainUser(response.GetUser()), nil
}

func NewRepository(client *grpc.Client) *Repository {
	return &Repository{
		client: client,
	}
}
