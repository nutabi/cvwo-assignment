package service

import (
	"context"
	"errors"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"github.com/nutabi/cvwo-assignment/backend/internal/repository"
	"github.com/nutabi/cvwo-assignment/backend/internal/utility"
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
			return nil, ErrUserNotFound
		} else {
			return nil, errors.Join(ErrDatabaseUnknown, err)
		}
	}

	// Map repository user to service user profile
	profile := profileFromUser(&user, false)
	return &profile, nil
}

func (s *defaultService) GetCurrentUserProfile(ctx context.Context, user *model.User) (*UserProfile, error) {
	// Map repository user to service user profile
	profile := profileFromUser(user, true)
	return &profile, nil
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
		return errors.Join(ErrDatabaseUnknown, err)
	}

	return nil
}

func (s *defaultService) RegisterUser(ctx context.Context, username, email, password string) (*UserProfile, error) {
	// Check for existing username or email
	usernameExists, err := s.repo.CheckUsernameExists(ctx, username)
	if err != nil {
		return nil, errors.Join(ErrDatabaseUnknown, err)
	}
	if usernameExists {
		return nil, ErrUsernameTaken
	}
	emailExists, err := s.repo.CheckEmailExists(ctx, email)
	if err != nil {
		return nil, errors.Join(ErrDatabaseUnknown, err)
	}
	if emailExists {
		return nil, ErrEmailInUse
	}

	// Hash password
	phc, err := utility.ComputePHC(password)
	if err != nil {
		return nil, errors.Join(ErrDatabaseUnknown, err)
	}

	// Create new user model
	newUser := model.User{
		Username: username,
		Email:    email,
		PHC:      phc,
	}

	// Save new user via repository
	if err := s.repo.CreateUser(ctx, &newUser); err != nil {
		return nil, errors.Join(ErrDatabaseUnknown, err)
	}

	// Return the newly created user's profile
	profile := profileFromUser(&newUser, true)
	return &profile, nil
}

func (s *defaultService) ListTopics(
	ctx context.Context,
	limit,
	offset int,
	userID uint,
	withPosts bool,
) ([]TopicInfo, error) {
	// Initialize userID pointer
	var userIDPtr *uint
	if userID != 0 {
		userIDPtr = &userID
	}

	// Fetch topics from repository
	topics, err := s.repo.GetTopics(ctx, limit, offset, userIDPtr, withPosts)
	if err != nil {
		return nil, errors.Join(ErrDatabaseUnknown, err)
	}

	// Convert to DTO
	infos := make([]TopicInfo, 0, len(topics))
	for _, topic := range topics {
		infos = append(infos, infoFromTopic(&topic, withPosts))
	}
	return userProfile, nil
}
