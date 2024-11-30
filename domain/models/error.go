package models

import "github.com/pkg/errors"

var (
	ErrNotFound = errors.New("not found")
	ErrInternal = errors.New("internal error")
)
