package primary

import (
	"context"
	"errors"
	"log/slog"

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
		slog.Error("failed to check username existence", "username", username, "error", err)
		return nil, errors.Join(service.ErrDatabaseError, err)
	}
	if usernameExists {
		slog.Warn("registration failed: username taken", "username", username)
		return nil, service.ErrUsernameTaken
	}
	emailExists, err := s.repo.CheckEmailExists(ctx, email)
	if err != nil {
		slog.Error("failed to check email existence", "email", email, "error", err)
		return nil, errors.Join(service.ErrDatabaseError, err)
	}
	if emailExists {
		slog.Warn("registration failed: email in use", "email", email)
		return nil, service.ErrEmailInUse
	}

	// Hash password
	phc, err := utility.ComputePHC(password)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
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
		slog.Error("failed to create user", "username", username, "error", err)
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	slog.Info("user registered", "user_id", newUser.ID, "username", username)

	// Return the newly created user's profile
	profile := service.ProfileFromUser(&newUser, true)
	return profile, nil
}

func (s *primaryService) FetchUserByID(ctx context.Context, id uint) (*service.UserProfile, error) {
	// Fetch user from repository
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Warn("user not found", "user_id", id)
			return nil, service.ErrUserNotFound
		}
		slog.Error("failed to fetch user", "user_id", id, "error", err)
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	slog.Debug("fetched user", "user_id", id)

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

	// Save changes via repository
	if err := s.repo.UpdateUser(
		ctx,
		user.ID,
		user.AvatarURL,
		user.Bio,
	); err != nil {
		slog.Error("failed to update user", "user_id", user.ID, "error", err)
		return errors.Join(service.ErrDatabaseError, err)
	}

	slog.Info("user updated", "user_id", user.ID)

	return nil
}
