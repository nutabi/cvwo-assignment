package repository

import (
	"context"
	"testing"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"gorm.io/driver/sqlite"
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

func initRepo() (Repository, error) {
	repo, err := ConnectSQL(sqlite.Open(":memory:"))
	if err != nil {
		return nil, err
	}
	if err := repo.Migrate(); err != nil {
		return nil, err
	}
	return repo, nil
}

func addMockUser(repo Repository, username, email, phc string) error {
	user := &model.User{
		Username: username,
		Email:    email,
		PHC:      phc,
	}
	return repo.CreateUser(context.Background(), user)
}

func initRepoWithMockUsers() (Repository, error) {
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
			if err := repo.CreateUser(t.Context(), &model.User{
				Username: mu.username,
				Email:    mu.email,
				PHC:      mu.phc,
			}); err != nil {
				t.Errorf("Failed to create user %s: %v", mu.username, err)
			}
		})
	}
}
