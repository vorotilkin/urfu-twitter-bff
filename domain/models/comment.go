package models

import "time"

type Comment struct {
	ID        int32
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    int32
	PostID    int32
}
