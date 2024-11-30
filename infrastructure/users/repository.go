package users

import (
	"context"
	"github.com/pkg/errors"
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

func (r *Repository) FetchUserByID(ctx context.Context, id int32) (models.User, error) {
	client := proto.NewUsersClient(r.client.Connection())

	req := proto.UserByIDRequest{Id: id}

	response, err := client.UserByID(ctx, &req)
	if status.Code(err) == codes.NotFound {
		return models.User{}, models.ErrNotFound
	}
	if err != nil {
		return models.User{}, errors.Wrap(err, "FetchUserByID")
	}

	return hydrators.DomainUser(response.GetUser()), nil
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

func NewRepository(client *grpc.Client) *Repository {
	return &Repository{
		client: client,
	}
}
