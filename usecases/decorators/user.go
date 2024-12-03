package decorators

import (
	"github.com/oapi-codegen/runtime/types"
	"github.com/samber/lo"
	"twitter-bff/domain/models"
	"twitter-bff/openapigen"
)

func EchoUser(user models.User) *openapigen.UserResponse {
	return &openapigen.UserResponse{
		Email:        lo.ToPtr(types.Email(user.Email)),
		Id:           lo.ToPtr(user.ID),
		Name:         lo.ToPtr(user.Name),
		Username:     lo.ToPtr(user.Username),
		Bio:          lo.ToPtr(user.Bio),
		ProfileImage: lo.ToPtr(user.ProfileImage),
		CoverImage:   lo.ToPtr(user.CoverImage),
	}
}
