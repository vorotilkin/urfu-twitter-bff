package decorators

import (
	"fmt"
	"github.com/oapi-codegen/runtime/types"
	"github.com/samber/lo"
	"twitter-bff/domain/models"
	"twitter-bff/openapigen"
)

func EchoComments(comments []models.Comment) []openapigen.Comment {
	return lo.Map(comments, func(comment models.Comment, _ int) openapigen.Comment {
		return EchoComment(comment)
	})
}

func EchoComment(comment models.Comment) openapigen.Comment {
	return openapigen.Comment{
		Body:      comment.Body,
		CreatedAt: types.Date{Time: comment.CreatedAt},
		Id:        comment.ID,
		UpdatedAt: types.Date{Time: comment.UpdatedAt},
		UserId:    fmt.Sprint(comment.UserID),
		User:      EchoUser(comment.User),
		PostId:    fmt.Sprint(comment.PostID),
	}
}
