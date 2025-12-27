package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"github.com/nutabi/cvwo-assignment/backend/internal/repository"
	"gorm.io/gorm"
)

// Helper function to create a test user
func createTestUser(id uint, username, email string) *model.User {
	now := time.Now()
	avatarURL := "https://example.com/avatar.jpg"
	bio := "Test bio"

	return &model.User{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Username:  username,
		Email:     email,
		PHC:       "$argon2id$v=19$m=65536,t=3,p=2$testhashedpassword",
		AvatarURL: &avatarURL,
		Bio:       &bio,
	}
}

// Test GetUserProfileByID
func TestGetUserProfileByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockRepo := repository.NewMockRepository().(*repository.MockRepository)
		svc := Default(mockRepo)
		ctx := context.Background()

		testUser := createTestUser(1, "testuser", "test@example.com")
		mockRepo.AddUser(testUser)

		// Execute
		profile, err := svc.GetUserProfileByID(ctx, 1)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if profile == nil {
			t.Fatal("Expected profile, got nil")
		}
		if profile.Username != "testuser" {
			t.Errorf("Expected username 'testuser', got '%s'", profile.Username)
		}
		if profile.Email != nil {
			t.Errorf("Expected email to be nil for public profile, got %v", *profile.Email)
		}
		if profile.AvatarUrl == nil || *profile.AvatarUrl != "https://example.com/avatar.jpg" {
			t.Errorf("Expected avatar URL, got %v", profile.AvatarUrl)
		}
		if profile.Bio == nil || *profile.Bio != "Test bio" {
			t.Errorf("Expected bio, got %v", profile.Bio)
		}
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// Setup
		mockRepo := repository.NewMockRepository().(*repository.MockRepository)
		svc := Default(mockRepo)
		ctx := context.Background()

		// Execute
		profile, err := svc.GetUserProfileByID(ctx, 999)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, ErrUserNotFound) {
			t.Errorf("Expected ErrUserNotFound, got %v", err)
		}
		if profile != nil {
			t.Errorf("Expected nil profile, got %v", profile)
		}
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Setup
		mockRepo := repository.NewMockRepository().(*repository.MockRepository)
		mockRepo.GetUserByIDError = errors.New("database connection failed")
		svc := Default(mockRepo)
		ctx := context.Background()

		// Execute
		profile, err := svc.GetUserProfileByID(ctx, 1)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, ErrDatabaseUnknown) {
			t.Errorf("Expected ErrDatabaseUnknown, got %v", err)
		}
		if profile != nil {
			t.Errorf("Expected nil profile, got %v", profile)
		}
	})
}

// Test GetCurrentUserProfile
func TestGetCurrentUserProfile(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockRepo := repository.NewMockRepository().(*repository.MockRepository)
		svc := Default(mockRepo)
		ctx := context.Background()

		testUser := createTestUser(1, "currentuser", "current@example.com")

		// Execute
		profile, err := svc.GetCurrentUserProfile(ctx, testUser)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if profile == nil {
			t.Fatal("Expected profile, got nil")
		}
		if profile.Username != "currentuser" {
			t.Errorf("Expected username 'currentuser', got '%s'", profile.Username)
		}
		if profile.Email == nil || *profile.Email != "current@example.com" {
			t.Errorf("Expected email 'current@example.com', got %v", profile.Email)
		}
		if profile.UpdatedAt == nil {
			t.Error("Expected UpdatedAt to be set for current user profile")
		}
		if profile.AvatarUrl == nil || *profile.AvatarUrl != "https://example.com/avatar.jpg" {
			t.Errorf("Expected avatar URL, got %v", profile.AvatarUrl)
		}
		if profile.Bio == nil || *profile.Bio != "Test bio" {
			t.Errorf("Expected bio, got %v", profile.Bio)
		}
	})
}

// Test UpdateCurrentUserProfile
func TestUpdateCurrentUserProfile(t *testing.T) {
	t.Run("UpdateBothFields", func(t *testing.T) {
		// Setup
		mockRepo := repository.NewMockRepository().(*repository.MockRepository)
		svc := Default(mockRepo)
		ctx := context.Background()

		testUser := createTestUser(1, "testuser", "test@example.com")
		mockRepo.AddUser(testUser)

		newAvatarURL := "https://example.com/new-avatar.jpg"
		newBio := "Updated bio"

		// Execute
		err := svc.UpdateCurrentUserProfile(ctx, testUser, &newAvatarURL, &newBio)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if testUser.AvatarURL == nil || *testUser.AvatarURL != newAvatarURL {
			t.Errorf("Expected avatar URL '%s', got %v", newAvatarURL, testUser.AvatarURL)
		}
		if testUser.Bio == nil || *testUser.Bio != newBio {
			t.Errorf("Expected bio '%s', got %v", newBio, testUser.Bio)
		}
	})

	t.Run("UpdateAvatarOnly", func(t *testing.T) {
		// Setup
		mockRepo := repository.NewMockRepository().(*repository.MockRepository)
		svc := Default(mockRepo)
		ctx := context.Background()

		testUser := createTestUser(1, "testuser", "test@example.com")
		originalBio := *testUser.Bio
		mockRepo.AddUser(testUser)

		newAvatarURL := "https://example.com/new-avatar.jpg"

		// Execute
		err := svc.UpdateCurrentUserProfile(ctx, testUser, &newAvatarURL, nil)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if testUser.AvatarURL == nil || *testUser.AvatarURL != newAvatarURL {
			t.Errorf("Expected avatar URL '%s', got %v", newAvatarURL, testUser.AvatarURL)
		}
		if testUser.Bio == nil || *testUser.Bio != originalBio {
			t.Errorf("Expected bio to remain '%s', got %v", originalBio, testUser.Bio)
		}
	})

	t.Run("UpdateBioOnly", func(t *testing.T) {
		// Setup
		mockRepo := repository.NewMockRepository().(*repository.MockRepository)
		svc := Default(mockRepo)
		ctx := context.Background()

		testUser := createTestUser(1, "testuser", "test@example.com")
		originalAvatarURL := *testUser.AvatarURL
		mockRepo.AddUser(testUser)

		newBio := "Updated bio only"

		// Execute
		err := svc.UpdateCurrentUserProfile(ctx, testUser, nil, &newBio)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if testUser.AvatarURL == nil || *testUser.AvatarURL != originalAvatarURL {
			t.Errorf("Expected avatar URL to remain '%s', got %v", originalAvatarURL, testUser.AvatarURL)
		}
		if testUser.Bio == nil || *testUser.Bio != newBio {
			t.Errorf("Expected bio '%s', got %v", newBio, testUser.Bio)
		}
	})

	t.Run("NoUpdate", func(t *testing.T) {
		// Setup
		mockRepo := repository.NewMockRepository().(*repository.MockRepository)
		svc := Default(mockRepo)
		ctx := context.Background()

		testUser := createTestUser(1, "testuser", "test@example.com")
		originalAvatarURL := *testUser.AvatarURL
		originalBio := *testUser.Bio
		mockRepo.AddUser(testUser)

		// Execute
		err := svc.UpdateCurrentUserProfile(ctx, testUser, nil, nil)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if testUser.AvatarURL == nil || *testUser.AvatarURL != originalAvatarURL {
			t.Errorf("Expected avatar URL to remain '%s', got %v", originalAvatarURL, testUser.AvatarURL)
		}
		if testUser.Bio == nil || *testUser.Bio != originalBio {
			t.Errorf("Expected bio to remain '%s', got %v", originalBio, testUser.Bio)
		}
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Setup
		mockRepo := repository.NewMockRepository().(*repository.MockRepository)
		mockRepo.UpdateUserError = errors.New("database error")
		svc := Default(mockRepo)
		ctx := context.Background()

		testUser := createTestUser(1, "testuser", "test@example.com")
		newBio := "New bio"

		// Execute
		err := svc.UpdateCurrentUserProfile(ctx, testUser, nil, &newBio)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, ErrDatabaseUnknown) {
			t.Errorf("Expected ErrDatabaseUnknown, got %v", err)
		}
	})
}

// Test RegisterUser
func TestRegisterUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockRepo := repository.NewMockRepository().(*repository.MockRepository)
		svc := Default(mockRepo)
		ctx := context.Background()

		// Execute
		profile, err := svc.RegisterUser(ctx, "newuser", "newuser@example.com", "SecurePass123!")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if profile == nil {
			t.Fatal("Expected profile, got nil")
		}
		if profile.Username != "newuser" {
			t.Errorf("Expected username 'newuser', got '%s'", profile.Username)
		}
		if profile.Email == nil || *profile.Email != "newuser@example.com" {
			t.Errorf("Expected email 'newuser@example.com', got %v", profile.Email)
		}
		if profile.UpdatedAt == nil {
			t.Error("Expected UpdatedAt to be set")
		}

		// Verify user was created in repository
		user, err := mockRepo.GetUserByUsername(ctx, "newuser")
		if err != nil {
			t.Fatalf("Expected user to be created in repository, got error: %v", err)
		}
		if user.Username != "newuser" {
			t.Errorf("Expected username 'newuser', got '%s'", user.Username)
		}
		if user.PHC == "" {
			t.Error("Expected password hash to be set")
		}
	})

	t.Run("UsernameTaken", func(t *testing.T) {
		// Setup
		mockRepo := repository.NewMockRepository().(*repository.MockRepository)
		svc := Default(mockRepo)
		ctx := context.Background()

		existingUser := createTestUser(1, "existinguser", "existing@example.com")
		mockRepo.AddUser(existingUser)

		// Execute
		profile, err := svc.RegisterUser(ctx, "existinguser", "newemail@example.com", "SecurePass123!")

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, ErrUsernameTaken) {
			t.Errorf("Expected ErrUsernameTaken, got %v", err)
		}
		if profile != nil {
			t.Errorf("Expected nil profile, got %v", profile)
		}
	})

	t.Run("EmailInUse", func(t *testing.T) {
		// Setup
		mockRepo := repository.NewMockRepository().(*repository.MockRepository)
		svc := Default(mockRepo)
		ctx := context.Background()

		existingUser := createTestUser(1, "existinguser", "existing@example.com")
		mockRepo.AddUser(existingUser)

		// Execute
		profile, err := svc.RegisterUser(ctx, "newuser", "existing@example.com", "SecurePass123!")

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, ErrEmailInUse) {
			t.Errorf("Expected ErrEmailInUse, got %v", err)
		}
		if profile != nil {
			t.Errorf("Expected nil profile, got %v", profile)
		}
	})

	t.Run("UsernameCheckDatabaseError", func(t *testing.T) {
		// Setup
		mockRepo := repository.NewMockRepository().(*repository.MockRepository)
		mockRepo.CheckUsernameExistsError = errors.New("database error")
		svc := Default(mockRepo)
		ctx := context.Background()

		// Execute
		profile, err := svc.RegisterUser(ctx, "newuser", "newuser@example.com", "SecurePass123!")

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, ErrDatabaseUnknown) {
			t.Errorf("Expected ErrDatabaseUnknown, got %v", err)
		}
		if profile != nil {
			t.Errorf("Expected nil profile, got %v", profile)
		}
	})

	t.Run("EmailCheckDatabaseError", func(t *testing.T) {
		// Setup
		mockRepo := repository.NewMockRepository().(*repository.MockRepository)
		mockRepo.CheckEmailExistsError = errors.New("database error")
		svc := Default(mockRepo)
		ctx := context.Background()

		// Execute
		profile, err := svc.RegisterUser(ctx, "newuser", "newuser@example.com", "SecurePass123!")

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, ErrDatabaseUnknown) {
			t.Errorf("Expected ErrDatabaseUnknown, got %v", err)
		}
		if profile != nil {
			t.Errorf("Expected nil profile, got %v", profile)
		}
	})

	t.Run("CreateUserDatabaseError", func(t *testing.T) {
		// Setup
		mockRepo := repository.NewMockRepository().(*repository.MockRepository)
		mockRepo.CreateUserError = errors.New("database error")
		svc := Default(mockRepo)
		ctx := context.Background()

		// Execute
		profile, err := svc.RegisterUser(ctx, "newuser", "newuser@example.com", "SecurePass123!")

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, ErrDatabaseUnknown) {
			t.Errorf("Expected ErrDatabaseUnknown, got %v", err)
		}
		if profile != nil {
			t.Errorf("Expected nil profile, got %v", profile)
		}
	})
}

// Test Repo accessor
func TestRepo(t *testing.T) {
	t.Run("ReturnsRepository", func(t *testing.T) {
		// Setup
		mockRepo := repository.NewMockRepository()
		svc := Default(mockRepo)

		// Execute
		repo := svc.Repo()

		// Assert
		if repo == nil {
			t.Fatal("Expected repository, got nil")
		}
		if repo != mockRepo {
			t.Error("Expected same repository instance")
		}
	})
}
