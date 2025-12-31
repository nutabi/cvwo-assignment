package primary

import (
	"context"
	"errors"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
	"github.com/nutabi/cvwo-assignment/backend/internal/utility"
	"gorm.io/gorm"
)

func (s *primaryService) RegisterUser(
	ctx context.Context,
	username, email, password string,
) (*service.UserProfile, error) {
	// Check for existing username or email
	usernameExists, err := s.repo.CheckUsernameExists(ctx, username)
	if err != nil {
		return nil, errors.Join(service.ErrDatabaseError, err)
	}
	if usernameExists {
		return nil, service.ErrUsernameTaken
	}
	emailExists, err := s.repo.CheckEmailExists(ctx, email)
	if err != nil {
		return nil, errors.Join(service.ErrDatabaseError, err)
	}
	if emailExists {
		return nil, service.ErrEmailInUse
	}

	// Hash password
	phc, err := utility.ComputePHC(password)
	if err != nil {
		return nil, errors.Join(service.ErrCryptoError, err)
	}

	// Create new user model
	newUser := model.User{
		Username: username,
		Email:    email,
		PHC:      phc,
	}

	// Save new user via repository
	if err := s.repo.CreateUser(ctx, &newUser); err != nil {
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	// Return the newly created user's profile
	profile := service.ProfileFromUser(&newUser, true)
	return profile, nil
}

func (s *primaryService) FetchUserByID(ctx context.Context, id uint) (*service.UserProfile, error) {
	// Fetch user from repository
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, service.ErrUserNotFound
		} else {
			return nil, errors.Join(service.ErrDatabaseError, err)
		}
	}

	// Map repository user to service user profile
	profile := service.ProfileFromUser(user, false)
	return profile, nil
}

func (s *primaryService) FetchCurrentUser(ctx context.Context, user *model.User) (*service.UserProfile, error) {
	// Map repository user to service user profile
	profile := service.ProfileFromUser(user, true)
	return profile, nil
}

func (s *primaryService) UpdateCurrentUser(
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

	// Make sure at least one field is being updated
	if newAvatarUrl == nil && newBio == nil {
		return service.ErrNoUpdateFields
	}

	// Save changes via repository
	if err := s.repo.UpdateUser(
		ctx,
		user.ID,
		user.AvatarURL,
		user.Bio,
	); err != nil {
		return errors.Join(service.ErrDatabaseError, err)
	}

	return nil
}
