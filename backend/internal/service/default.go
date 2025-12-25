package service

import (
	"context"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"github.com/nutabi/cvwo-assignment/backend/internal/repository"
	"gorm.io/gorm"
)

type defaultService struct {
	repo repository.Repository
}

func Default(repo repository.Repository) Service {
	return &defaultService{repo: repo}
}

func (s *defaultService) Repo() repository.Repository {
	return s.repo
}

func (s *defaultService) GetUserProfileByID(ctx context.Context, id uint) (*UserProfile, error) {
	// Fetch user from repository
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound(id)
		} else {
			return nil, ErrDatabaseUnknown(err)
		}
	}

	// Map repository user to service user profile
	userProfile := &UserProfile{
		CreatedAt: user.CreatedAt,
		Username:  user.Username,
		AvatarUrl: user.AvatarURL,
		Bio:       user.Bio,
	}

	return userProfile, nil
}

func (s *defaultService) GetCurrentUserProfile(ctx context.Context, user *model.User) (*UserProfile, error) {
	// Map repository user to service user profile
	userProfile := &UserProfile{
		CreatedAt: user.CreatedAt,
		UpdatedAt: &user.UpdatedAt,
		Username:  user.Username,
		Email:     &user.Email,
		AvatarUrl: user.AvatarURL,
		Bio:       user.Bio,
	}

	return userProfile, nil
}

func (s *defaultService) UpdateCurrentUserProfile(
	ctx context.Context,
	user *model.User,
	newAvatarUrl,
	newBio *string,
) error {
	// Update fields if provided
	if newAvatarUrl != nil {
		user.AvatarURL = newAvatarUrl
	}
	if newBio != nil {
		user.Bio = newBio
	}

	// Save changes via repository
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return ErrDatabaseUnknown(err)
	}

	return nil
}
