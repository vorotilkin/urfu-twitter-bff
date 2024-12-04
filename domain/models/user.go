package models

import "github.com/samber/mo"

type User struct {
	ID               int32
	Name             string
	PasswordHash     string
	Username         string
	Email            string
	Bio              string
	ProfileImage     string
	CoverImage       string
	FollowingUserIds []int32
	FollowerUserIds  []int32
}

type UserOption struct {
	ID           int32
	Name         mo.Option[string]
	Username     mo.Option[string]
	Bio          mo.Option[string]
	ProfileImage mo.Option[string]
	CoverImage   mo.Option[string]
}
