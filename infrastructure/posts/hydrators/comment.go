package hydrators

import (
	"github.com/samber/lo"
	"github.com/vorotilkin/twitter-posts/proto"
	"twitter-bff/domain/models"
)

func DomainComments(comments []*proto.Comment) []models.Comment {
	return lo.Map(comments, func(comment *proto.Comment, _ int) models.Comment {
		return DomainComment(comment)
	})
}

func DomainComment(comment *proto.Comment) models.Comment {
	if comment == nil {
		return models.Comment{}
	}

	return models.Comment{
		ID:        comment.GetId(),
		Body:      comment.GetBody(),
		CreatedAt: comment.GetCreatedAt().AsTime(),
		UpdatedAt: comment.GetUpdatedAt().AsTime(),
		UserID:    comment.GetUserId(),
		PostID:    comment.GetPostId(),
	}
}
