package service

import (
	"context"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"github.com/nutabi/cvwo-assignment/backend/internal/repository"
)

type Service interface {
	// Get repository layer directly
	Repo() repository.Repository

	// User-related services
	GetUserProfileByID(ctx context.Context, id uint) (*UserProfile, error)
	GetCurrentUserProfile(ctx context.Context, user *model.User) (*UserProfile, error)
	UpdateCurrentUserProfile(ctx context.Context, user *model.User, newAvatarUrl, newBio *string) error
	RegisterUser(ctx context.Context, username, email, password string) (*UserProfile, error)
}
