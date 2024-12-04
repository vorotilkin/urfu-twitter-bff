package hydrators

import (
	"github.com/vorotilkin/twitter-posts/proto"
	"twitter-bff/domain/models"
)

func ProtoLikeOperationType(likeType models.LikeType) proto.LikeRequest_OperationType {
	switch likeType {
	default:
		return proto.LikeRequest_OPERATION_TYPE_LIKE_UNSPECIFIED
	case models.Dislike:
		return proto.LikeRequest_OPERATION_TYPE_DISLIKE
	}
}
