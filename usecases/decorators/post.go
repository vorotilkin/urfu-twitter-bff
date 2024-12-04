package decorators

import (
	"fmt"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/samber/lo"
	"twitter-bff/domain/models"
	"twitter-bff/openapigen"
)

func EchoPosts(posts []models.Post) []openapigen.Post {
	return lo.Map(posts, func(post models.Post, _ int) openapigen.Post {
		return EchoPost(post)
	})
}

func EchoPost(post models.Post) openapigen.Post {
	return openapigen.Post{
		Body:              post.Body,
		Comments:          EchoComments(post.Comments),
		CreatedAt:         openapi_types.Date{Time: post.CreatedAt},
		Id:                post.ID,
		LikeCount:         post.LikeCount,
		IsCurrentUserLike: lo.ToPtr(post.IsCurrentUserLike),
		UpdatedAt:         openapi_types.Date{Time: post.UpdatedAt},
		UserId:            fmt.Sprint(post.UserID),
		User:              EchoUser(post.User),
	}
}
