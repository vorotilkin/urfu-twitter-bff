package hydrators

import (
	"github.com/samber/lo"
	"github.com/vorotilkin/twitter-users/proto"
	"twitter-bff/domain/models"
)

func DomainUsers(users []*proto.User) []models.User {
	return lo.Map(users, func(user *proto.User, _ int) models.User {
		return DomainUser(user)
	})
}

func DomainUser(user *proto.User) models.User {
	if user == nil {
		return models.User{}
	}

	return models.User{
		ID:               user.GetId(),
		Name:             user.GetName(),
		PasswordHash:     user.GetPasswordHash(),
		Username:         user.GetUsername(),
		Email:            user.GetEmail(),
		Bio:              user.GetBio(),
		ProfileImage:     user.GetProfileImage(),
		CoverImage:       user.GetCoverImage(),
		FollowingUserIds: user.GetFollowingUserIds(),
		FollowerUserIds:  user.GetFollowerUserIds(),
	}
}
