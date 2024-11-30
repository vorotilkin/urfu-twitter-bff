package hydrators

import (
	"github.com/vorotilkin/twitter-users/proto"
	"twitter-bff/domain/models"
)

func DomainUser(user *proto.User) models.User {
	if user == nil {
		return models.User{}
	}

	return models.User{
		ID:           user.GetId(),
		Name:         user.GetName(),
		PasswordHash: user.GetPasswordHash(),
		Username:     user.GetUsername(),
		Email:        user.GetEmail(),
	}
}
