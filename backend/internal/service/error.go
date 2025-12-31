package service

import (
	"errors"
)

var (
	// Cryptographic errors
	ErrCryptoError = errors.New("cryptographic error")

	// Database errors
	ErrDatabaseError = errors.New("unknown database error")

	// Email in use error
	ErrEmailInUse = errors.New("email in use")

	// No updated fields error
	ErrNoUpdateFields = errors.New("no fields to update")

	// Topic not found error
	ErrTopicNotFound = errors.New("topic not found")

	// Post not found error
	ErrPostNotFound = errors.New("post not found")

	// Forbidden access error
	ErrForbidden = errors.New("forbidden")

	// User not found error
	ErrUserNotFound = errors.New("user not found")

	// Username taken error
	ErrUsernameTaken = errors.New("username taken")
)
