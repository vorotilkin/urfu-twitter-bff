package users

import (
	"context"
	"github.com/vorotilkin/twitter-users/proto"
	"twitter-bff/domain/models"
	"twitter-bff/pkg/grpc"
)

type Repository struct {
	client *grpc.Client
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

	return models.User{
		ID:           protoUser.GetId(),
		Name:         protoUser.GetName(),
		PasswordHash: protoUser.GetPasswordHash(),
		Username:     protoUser.GetUsername(),
		Email:        protoUser.GetEmail(),
	}, err
}

func NewRepository(client *grpc.Client) *Repository {
	return &Repository{
		client: client,
	}
}
