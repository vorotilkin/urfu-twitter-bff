package models

import "time"

type Post struct {
	ID        int32
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    int32
	LikeCount int32
	Comments  []Comment
}
