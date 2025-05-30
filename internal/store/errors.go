package store

import (
	"errors"
)

const (
	UniqueViolationErr = "23505"
)

var (
	ErrRecordNotFound        = errors.New("record not found")
	ErrDuplicateResourceType = errors.New("resource type already exists")
	ErrDuplicateEmail        = errors.New("duplicate email")
	ErrDuplicateUserName     = errors.New("duplicate username")
)
