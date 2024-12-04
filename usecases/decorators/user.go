package decorators

import (
	"fmt"
	"github.com/oapi-codegen/runtime/types"
	"github.com/samber/lo"
	"twitter-bff/domain/models"
	"twitter-bff/openapigen"
)

func EchoUsers(users []models.User) []*openapigen.User {
	return lo.Map(users, func(user models.User, _ int) *openapigen.User {
		return EchoUser(user)
	})
}

func EchoUser(user models.User) *openapigen.User {
	followingIDs := lo.Map(user.FollowingUserIds, func(id int32, _ int) string {
		return fmt.Sprint(id)
	})

	return &openapigen.User{
		Email:          lo.Ternary(len(user.Email) > 0, lo.ToPtr(types.Email(user.Email)), nil),
		Id:             lo.ToPtr(user.ID),
		Name:           lo.ToPtr(user.Name),
		Username:       lo.ToPtr(user.Username),
		Bio:            lo.ToPtr(user.Bio),
		ProfileImage:   lo.ToPtr(user.ProfileImage),
		CoverImage:     lo.ToPtr(user.CoverImage),
		FollowingIds:   lo.Ternary(len(user.FollowingUserIds) != 0, &followingIDs, nil),
		FollowersCount: lo.ToPtr(len(user.FollowerUserIds)),
	}
}
