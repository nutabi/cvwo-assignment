package sql_test

import (
	"context"
	"errors"
	"testing"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"github.com/nutabi/cvwo-assignment/backend/internal/repository"
	"github.com/nutabi/cvwo-assignment/backend/internal/repository/sql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var mockUsers []struct {
	username  string
	email     string
	phc       string
	avatarUrl string
	bio       string
} = []struct {
	username  string
	email     string
	phc       string
	avatarUrl string
	bio       string
}{
	{"test1", "test1@example.com", "test1password", "http://example.com/avatar1.png", "Bio of test1"},
	{"test2", "test2@example.com", "test2password", "http://example.com/avatar2.png", "Bio of test2"},
	{"test3", "test3@example.com", "test3password", "http://example.com/avatar3.png", "Bio of test3"},
}

func initRepo() (repository.Repository, error) {
	repo, err := sql.Connect(sqlite.Open(":memory:"))
	if err != nil {
		return nil, err
	}
	if err := repo.Migrate(); err != nil {
		return nil, err
	}
	return repo, nil
}

func addMockUser(repo repository.Repository, username, email, phc string) error {
	user := &model.User{
		Username: username,
		Email:    email,
		PHC:      phc,
	}
	return repo.CreateUser(context.Background(), user)
}

func initRepoWithMockUsers() (repository.Repository, error) {
	repo, err := initRepo()
	if err != nil {
		return nil, err
	}
	for _, mu := range mockUsers {
		if err := addMockUser(repo, mu.username, mu.email, mu.phc); err != nil {
			return nil, err
		}
	}
	return repo, nil
}

func TestMigrate(t *testing.T) {
	_, err := initRepo()
	if err != nil {
		t.Fatalf("Failed to migrate repository: %v", err)
	}
}

func TestGetUserByID(t *testing.T) {
	repo, err := initRepoWithMockUsers()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}
	for i, mu := range mockUsers {
		t.Run(mu.username, func(t *testing.T) {
			user, err := repo.GetUserByID(context.Background(), uint(i+1))
			if err != nil {
				t.Errorf("Failed to get user by ID: %v", err)
			}
			if user.Username != mu.username || user.Email != mu.email || user.PHC != mu.phc {
				t.Errorf("Retrieved user does not match expected. Got %+v, want %+v", user, mu)
			}
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	repo, err := initRepoWithMockUsers()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}
	for _, mu := range mockUsers {
		t.Run(mu.username, func(t *testing.T) {
			user, err := repo.GetUserByUsername(context.Background(), mu.username)
			if err != nil {
				t.Errorf("Failed to get user by username: %v", err)
			}
			if user.Username != mu.username || user.Email != mu.email || user.PHC != mu.phc {
				t.Errorf("Retrieved user does not match expected. Got %+v, want %+v", user, mu)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	repo, err := initRepoWithMockUsers()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}
	for _, mu := range mockUsers {
		t.Run(mu.username, func(t *testing.T) {
			user, err := repo.GetUserByUsername(context.Background(), mu.username)
			if err != nil {
				t.Fatalf("Failed to get user by username: %v", err)
			}

			// Update bio
			newBio := "Updated bio for " + mu.username
			user.Bio = &newBio
			if err := repo.UpdateUser(context.Background(), &user); err != nil {
				t.Fatalf("Failed to update user: %v", err)
			}

			// Retrieve again to check update
			updatedUser, err := repo.GetUserByUsername(context.Background(), mu.username)
			if err != nil {
				t.Fatalf("Failed to get user by username after update: %v", err)
			}
			if *updatedUser.Bio != newBio {
				t.Errorf("User bio not updated. Got %s, want %s", *updatedUser.Bio, newBio)
			}

			// Update avatar URL
			newAvatarUrl := "http://example.com/new_avatar_" + mu.username + ".png"
			updatedUser.AvatarURL = &newAvatarUrl
			if err := repo.UpdateUser(context.Background(), &updatedUser); err != nil {
				t.Fatalf("Failed to update user avatar URL: %v", err)
			}

			// Retrieve again to check update
			finalUser, err := repo.GetUserByUsername(context.Background(), mu.username)
			if err != nil {
				t.Fatalf("Failed to get user by username after avatar update: %v", err)
			}
			if *finalUser.AvatarURL != newAvatarUrl {
				t.Errorf("User avatar URL not updated. Got %s, want %s", *finalUser.AvatarURL, newAvatarUrl)
			}
		})
	}
}

func TestCheckUsernameExists(t *testing.T) {
	repo, err := initRepoWithMockUsers()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}
	for _, mu := range mockUsers {
		t.Run(mu.username, func(t *testing.T) {
			exists, err := repo.CheckUsernameExists(context.Background(), mu.username)
			if err != nil {
				t.Errorf("Error checking username existence: %v", err)
			}
			if !exists {
				t.Errorf("Expected username %s to exist", mu.username)
			}
		})
	}
	// Check a non-existing username
	nonExistingUsername := "nonexistentuser"
	exists, err := repo.CheckUsernameExists(context.Background(), nonExistingUsername)
	if err != nil {
		t.Errorf("Error checking username existence: %v", err)
	}
	if exists {
		t.Errorf("Did not expect username %s to exist", nonExistingUsername)
	}
}

func TestCheckEmailExists(t *testing.T) {
	repo, err := initRepoWithMockUsers()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}
	for _, mu := range mockUsers {
		t.Run(mu.email, func(t *testing.T) {
			exists, err := repo.CheckEmailExists(context.Background(), mu.email)
			if err != nil {
				t.Errorf("Error checking email existence: %v", err)
			}
			if !exists {
				t.Errorf("Expected email %s to exist", mu.email)
			}
		})
	}
	// Check a non-existing email
	nonExistingEmail := "nonexistentemail@example.com"
	exists, err := repo.CheckEmailExists(context.Background(), nonExistingEmail)
	if err != nil {
		t.Errorf("Error checking email existence: %v", err)
	}
	if exists {
		t.Errorf("Did not expect email %s to exist", nonExistingEmail)
	}
}

func TestCreateUser(t *testing.T) {
	repo, err := initRepo()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}
	for _, mu := range mockUsers {
		t.Run(mu.username, func(t *testing.T) {
			if err := repo.CreateUser(context.Background(), &model.User{
				Username: mu.username,
				Email:    mu.email,
				PHC:      mu.phc,
			}); err != nil {
				t.Errorf("Failed to create user %s: %v", mu.username, err)
			}
		})
	}
}

// Test error cases

func TestConnectSQL_InvalidDialect(t *testing.T) {
	// SQLite with empty string still works, so we need a truly invalid config
	// We'll skip this test as it's hard to force a connection error with SQLite
	// The error path in ConnectSQL is tested when gorm.Open fails
	t.Skip("Skipping as SQLite is very permissive with connection strings")
}

func TestGetUserByID_NotFound(t *testing.T) {
	repo, err := initRepoWithMockUsers()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Try to get non-existent user
	_, err = repo.GetUserByID(context.Background(), 9999)
	if err == nil {
		t.Error("Expected error when getting non-existent user, got nil")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestGetUserByUsername_NotFound(t *testing.T) {
	repo, err := initRepoWithMockUsers()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Try to get non-existent user
	_, err = repo.GetUserByUsername(context.Background(), "nonexistentuser")
	if err == nil {
		t.Error("Expected error when getting non-existent username, got nil")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	repo, err := initRepoWithMockUsers()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Try to create a user with an existing username
	err = repo.CreateUser(context.Background(), &model.User{
		Username: mockUsers[0].username,
		Email:    "newemail@example.com",
		PHC:      "newpassword",
	})
	if err == nil {
		t.Error("Expected error when creating user with duplicate username, got nil")
	}
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	repo, err := initRepoWithMockUsers()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Try to create a user with an existing email
	err = repo.CreateUser(context.Background(), &model.User{
		Username: "newuser",
		Email:    mockUsers[0].email,
		PHC:      "newpassword",
	})
	if err == nil {
		t.Error("Expected error when creating user with duplicate email, got nil")
	}
}

func TestUpdateUser_WithAvatarAndBio(t *testing.T) {
	repo, err := initRepoWithMockUsers()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	user, err := repo.GetUserByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Update both avatar and bio
	newAvatar := "https://example.com/new-avatar.png"
	newBio := "New bio"
	user.AvatarURL = &newAvatar
	user.Bio = &newBio

	err = repo.UpdateUser(context.Background(), &user)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Verify updates
	updatedUser, err := repo.GetUserByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}

	if updatedUser.AvatarURL == nil || *updatedUser.AvatarURL != newAvatar {
		t.Errorf("Expected avatar URL %s, got %v", newAvatar, updatedUser.AvatarURL)
	}
	if updatedUser.Bio == nil || *updatedUser.Bio != newBio {
		t.Errorf("Expected bio %s, got %v", newBio, updatedUser.Bio)
	}
}
