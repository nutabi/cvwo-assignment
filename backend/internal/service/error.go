package service

import (
	"errors"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUsernameTaken   = errors.New("username taken")
	ErrEmailInUse      = errors.New("email in use")
	ErrDatabaseUnknown = errors.New("unknown database error")
)
