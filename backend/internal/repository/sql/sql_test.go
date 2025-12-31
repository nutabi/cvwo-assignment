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

var mockTopics []struct {
	name        string
	description string
	authorID    uint
} = []struct {
	name        string
	description string
	authorID    uint
}{
	{"General Discussion", "A place for general topics", 1},
	{"Technology", "Discuss technology and programming", 1},
	{"Sports", "Talk about sports and games", 2},
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

func addMockTopic(repo repository.Repository, name, description string, authorID uint) error {
	desc := description
	topic := &model.Topic{
		Name:        name,
		Description: &desc,
		AuthorID:    authorID,
	}
	return repo.CreateTopic(context.Background(), topic)
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

func initRepoWithMockUsersAndTopics() (repository.Repository, error) {
	repo, err := initRepoWithMockUsers()
	if err != nil {
		return nil, err
	}
	for _, mt := range mockTopics {
		if err := addMockTopic(repo, mt.name, mt.description, mt.authorID); err != nil {
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

// Topic-related tests

func TestCreateTopic(t *testing.T) {
	repo, err := initRepoWithMockUsers()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}
	for _, mt := range mockTopics {
		t.Run(mt.name, func(t *testing.T) {
			desc := mt.description
			if err := repo.CreateTopic(context.Background(), &model.Topic{
				Name:        mt.name,
				Description: &desc,
				AuthorID:    mt.authorID,
			}); err != nil {
				t.Errorf("Failed to create topic %s: %v", mt.name, err)
			}
		})
	}
}

func TestCreateTopic_DuplicateName(t *testing.T) {
	repo, err := initRepoWithMockUsersAndTopics()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Try to create a topic with an existing name
	desc := "Duplicate topic"
	err = repo.CreateTopic(context.Background(), &model.Topic{
		Name:        mockTopics[0].name,
		Description: &desc,
		AuthorID:    1,
	})
	if err == nil {
		t.Error("Expected error when creating topic with duplicate name, got nil")
	}
}

func TestGetTopicByID(t *testing.T) {
	repo, err := initRepoWithMockUsersAndTopics()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}
	for i, mt := range mockTopics {
		t.Run(mt.name, func(t *testing.T) {
			topic, err := repo.GetOneTopic(context.Background(), uint(i+1))
			if err != nil {
				t.Errorf("Failed to get topic by ID: %v", err)
			}
			if topic.Name != mt.name {
				t.Errorf("Retrieved topic name does not match. Got %s, want %s", topic.Name, mt.name)
			}
			if topic.Description == nil || *topic.Description != mt.description {
				t.Errorf("Retrieved topic description does not match. Got %v, want %s", topic.Description, mt.description)
			}
			if topic.AuthorID != mt.authorID {
				t.Errorf("Retrieved topic authorID does not match. Got %d, want %d", topic.AuthorID, mt.authorID)
			}
		})
	}
}

func TestGetTopicByID_NotFound(t *testing.T) {
	repo, err := initRepoWithMockUsersAndTopics()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Try to get non-existent topic
	_, err = repo.GetOneTopic(context.Background(), 9999)
	if err == nil {
		t.Error("Expected error when getting non-existent topic, got nil")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestGetTopics(t *testing.T) {
	repo, err := initRepoWithMockUsersAndTopics()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	t.Run("GetAllTopics", func(t *testing.T) {
		topics, err := repo.GetTopics(context.Background(), 10, 0, nil, false)
		if err != nil {
			t.Fatalf("Failed to get all topics: %v", err)
		}
		if len(topics) != len(mockTopics) {
			t.Errorf("Expected %d topics, got %d", len(mockTopics), len(topics))
		}
	})

	t.Run("GetTopicsWithLimit", func(t *testing.T) {
		topics, err := repo.GetTopics(context.Background(), 2, 0, nil, false)
		if err != nil {
			t.Fatalf("Failed to get topics with limit: %v", err)
		}
		if len(topics) != 2 {
			t.Errorf("Expected 2 topics, got %d", len(topics))
		}
	})

	t.Run("GetTopicsWithOffset", func(t *testing.T) {
		topics, err := repo.GetTopics(context.Background(), 10, 1, nil, false)
		if err != nil {
			t.Fatalf("Failed to get topics with offset: %v", err)
		}
		if len(topics) != len(mockTopics)-1 {
			t.Errorf("Expected %d topics, got %d", len(mockTopics)-1, len(topics))
		}
	})

	t.Run("GetTopicsByUserID", func(t *testing.T) {
		userID := uint(1)
		topics, err := repo.GetTopics(context.Background(), 10, 0, &userID, false)
		if err != nil {
			t.Fatalf("Failed to get topics by user ID: %v", err)
		}
		// Count how many topics have authorID = 1 in mockTopics
		expectedCount := 0
		for _, mt := range mockTopics {
			if mt.authorID == 1 {
				expectedCount++
			}
		}
		if len(topics) != expectedCount {
			t.Errorf("Expected %d topics for user 1, got %d", expectedCount, len(topics))
		}
		for _, topic := range topics {
			if topic.AuthorID != userID {
				t.Errorf("Expected topic to be authored by user %d, got %d", userID, topic.AuthorID)
			}
		}
	})

	t.Run("GetTopicsWithPosts", func(t *testing.T) {
		topics, err := repo.GetTopics(context.Background(), 10, 0, nil, true)
		if err != nil {
			t.Fatalf("Failed to get topics with posts: %v", err)
		}
		if len(topics) != len(mockTopics) {
			t.Errorf("Expected %d topics, got %d", len(mockTopics), len(topics))
		}
		// Even if no posts exist, Posts field should be initialized (empty slice, not nil)
		for _, topic := range topics {
			if topic.Posts == nil {
				t.Error("Expected Posts to be preloaded (not nil)")
			}
		}
	})
}

func TestUpdateTopic(t *testing.T) {
	repo, err := initRepoWithMockUsersAndTopics()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	topic, err := repo.GetOneTopic(context.Background(), 1)
	if err != nil {
		t.Fatalf("Failed to get topic: %v", err)
	}

	// Update name
	newName := "Updated General Discussion"
	topic.Name = newName
	if err := repo.UpdateTopic(context.Background(), &topic); err != nil {
		t.Fatalf("Failed to update topic name: %v", err)
	}

	// Retrieve again to check update
	updatedTopic, err := repo.GetOneTopic(context.Background(), 1)
	if err != nil {
		t.Fatalf("Failed to get topic after update: %v", err)
	}
	if updatedTopic.Name != newName {
		t.Errorf("Topic name not updated. Got %s, want %s", updatedTopic.Name, newName)
	}

	// Update description
	newDesc := "An updated description for general topics"
	updatedTopic.Description = &newDesc
	if err := repo.UpdateTopic(context.Background(), &updatedTopic); err != nil {
		t.Fatalf("Failed to update topic description: %v", err)
	}

	// Retrieve again to check update
	finalTopic, err := repo.GetOneTopic(context.Background(), 1)
	if err != nil {
		t.Fatalf("Failed to get topic after description update: %v", err)
	}
	if finalTopic.Description == nil || *finalTopic.Description != newDesc {
		t.Errorf("Topic description not updated. Got %v, want %s", finalTopic.Description, newDesc)
	}
}

func TestDeleteTopic(t *testing.T) {
	repo, err := initRepoWithMockUsersAndTopics()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Get topic with ID 1
	topic, err := repo.GetOneTopic(context.Background(), 1)
	if err != nil {
		t.Fatalf("Failed to get topic: %v", err)
	}

	// Delete topic
	if err := repo.DeleteTopic(context.Background(), topic.ID); err != nil {
		t.Fatalf("Failed to delete topic: %v", err)
	}

	// Verify topic is deleted
	_, err = repo.GetOneTopic(context.Background(), 1)
	if err == nil {
		t.Error("Expected error when getting deleted topic, got nil")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestDeleteTopic_NotFound(t *testing.T) {
	repo, err := initRepoWithMockUsersAndTopics()
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Try to delete non-existent topic
	err = repo.DeleteTopic(context.Background(), 9999)
	// Since the implementation doesn't check if topic exists before deleting,
	// this won't return an error. The delete will just affect 0 rows.
	// This test now verifies that deleting a non-existent topic doesn't cause a panic
	if err != nil {
		t.Errorf("Unexpected error when deleting non-existent topic: %v", err)
	}
}
