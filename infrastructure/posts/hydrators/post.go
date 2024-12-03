package hydrators

import (
	"github.com/samber/lo"
	"github.com/vorotilkin/twitter-posts/proto"
	"twitter-bff/domain/models"
)

func DomainPosts(posts []*proto.Post) []models.Post {
	return lo.Map(posts, func(post *proto.Post, _ int) models.Post {
		return DomainPost(post)
	})
}

func DomainPost(post *proto.Post) models.Post {
	if post == nil {
		return models.Post{}
	}

	return models.Post{
		ID:        post.GetId(),
		Body:      post.GetBody(),
		CreatedAt: post.GetCreatedAt().AsTime(),
		UpdatedAt: post.GetUpdatedAt().AsTime(),
		UserID:    post.GetUserId(),
		LikeCount: post.GetLikeCounter(),
		Comments:  DomainComments(post.GetComments()),
	}
}
