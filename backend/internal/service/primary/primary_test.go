package primary_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"github.com/nutabi/cvwo-assignment/backend/internal/repository"
	"github.com/nutabi/cvwo-assignment/backend/internal/repository/sql"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
	"github.com/nutabi/cvwo-assignment/backend/internal/service/primary"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func createMockRepository() (repository.Repository, error) {
	repo, err := sql.Connect(sqlite.Open(":memory:"))
	if err != nil {
		return nil, err
	}
	if err := repo.Migrate(); err != nil {
		return nil, err
	}
	return repo, nil
}

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
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		testUser := createTestUser(1, "testuser", "test@example.com")
		mockRepo.CreateUser(ctx, testUser)

		// Execute
		profile, err := svc.FetchUserByID(ctx, 1)

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
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Execute
		profile, err := svc.FetchUserByID(ctx, 999)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrUserNotFound) {
			t.Errorf("Expected ErrUserNotFound, got %v", err)
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
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		testUser := createTestUser(1, "currentuser", "current@example.com")

		// Execute
		profile, err := svc.FetchCurrentUser(ctx, testUser)

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
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		testUser := createTestUser(1, "testuser", "test@example.com")
		mockRepo.CreateUser(ctx, testUser)

		newAvatarURL := "https://example.com/new-avatar.jpg"
		newBio := "Updated bio"

		// Execute
		err = svc.UpdateCurrentUser(ctx, testUser, &newAvatarURL, &newBio)

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
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		testUser := createTestUser(1, "testuser", "test@example.com")
		originalBio := *testUser.Bio
		mockRepo.CreateUser(ctx, testUser)
		newAvatarURL := "https://example.com/new-avatar.jpg"

		// Execute
		err = svc.UpdateCurrentUser(ctx, testUser, &newAvatarURL, nil)

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
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		testUser := createTestUser(1, "testuser", "test@example.com")
		originalAvatarURL := *testUser.AvatarURL
		mockRepo.CreateUser(ctx, testUser)

		newBio := "Updated bio only"

		// Execute
		err = svc.UpdateCurrentUser(ctx, testUser, nil, &newBio)

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
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		testUser := createTestUser(1, "testuser", "test@example.com")
		mockRepo.CreateUser(ctx, testUser)

		// Execute
		err = svc.UpdateCurrentUser(ctx, testUser, nil, nil)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrNoUpdateFields) {
			t.Errorf("Expected ErrNoUpdateFields, got %v", err)
		}
	})
}

// Test RegisterUser
func TestRegisterUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
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
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		existingUser := createTestUser(1, "existinguser", "existing@example.com")
		mockRepo.CreateUser(ctx, existingUser)

		// Execute
		profile, err := svc.RegisterUser(ctx, "existinguser", "newemail@example.com", "SecurePass123!")

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrUsernameTaken) {
			t.Errorf("Expected ErrUsernameTaken, got %v", err)
		}
		if profile != nil {
			t.Errorf("Expected nil profile, got %v", profile)
		}
	})

	t.Run("EmailInUse", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		existingUser := createTestUser(1, "existinguser", "existing@example.com")
		mockRepo.CreateUser(ctx, existingUser)

		// Execute
		profile, err := svc.RegisterUser(ctx, "newuser", "existing@example.com", "SecurePass123!")

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrEmailInUse) {
			t.Errorf("Expected ErrEmailInUse, got %v", err)
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
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)

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

// Helper function to create a test topic
func createTestTopic(id uint, name string, description *string, authorID uint) *model.Topic {
	now := time.Now()
	return &model.Topic{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:        name,
		Description: description,
		AuthorID:    authorID,
	}
}

// Helper function to create a test post
func createTestPost(id uint, title string, content *string, authorID, topicID uint) *model.Post {
	now := time.Now()
	return &model.Post{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Title:    title,
		Content:  content,
		AuthorID: authorID,
		TopicID:  topicID,
	}
}

// Helper function to create a test comment
func createTestComment(id uint, content string, authorID, postID uint) *model.Comment {
	now := time.Now()
	return &model.Comment{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Content:  content,
		AuthorID: authorID,
		PostID:   postID,
	}
}

// Test CreateTopic
func TestCreateTopic(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author user
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		description := "Test topic description"

		// Execute
		topicInfo, err := svc.CreateTopic(ctx, 1, "Test Topic", &description)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if topicInfo == nil {
			t.Fatal("Expected topic info, got nil")
		}
		if topicInfo.Name != "Test Topic" {
			t.Errorf("Expected name 'Test Topic', got '%s'", topicInfo.Name)
		}
		if topicInfo.Description == nil || *topicInfo.Description != description {
			t.Errorf("Expected description '%s', got %v", description, topicInfo.Description)
		}
		if topicInfo.Author == nil {
			t.Fatal("Expected author, got nil")
		}
		if topicInfo.Author.Username != "author" {
			t.Errorf("Expected author username 'author', got '%s'", topicInfo.Author.Username)
		}
	})

	t.Run("SuccessWithNullDescription", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author user
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		// Execute
		topicInfo, err := svc.CreateTopic(ctx, 1, "Test Topic", nil)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if topicInfo == nil {
			t.Fatal("Expected topic info, got nil")
		}
		if topicInfo.Description != nil {
			t.Errorf("Expected nil description, got %v", topicInfo.Description)
		}
	})
}

// Test FetchTopicByID
func TestFetchTopicByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author and topic
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		description := "Test description"
		topic := createTestTopic(1, "Test Topic", &description, 1)
		mockRepo.CreateTopic(ctx, topic)

		// Execute
		topicInfo, err := svc.FetchTopicByID(ctx, 1)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if topicInfo == nil {
			t.Fatal("Expected topic info, got nil")
		}
		if topicInfo.Name != "Test Topic" {
			t.Errorf("Expected name 'Test Topic', got '%s'", topicInfo.Name)
		}
	})

	t.Run("TopicNotFound", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Execute
		topicInfo, err := svc.FetchTopicByID(ctx, 999)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrTopicNotFound) {
			t.Errorf("Expected ErrTopicNotFound, got %v", err)
		}
		if topicInfo != nil {
			t.Errorf("Expected nil topic info, got %v", topicInfo)
		}
	})
}

// Test FetchTopics
func TestFetchTopics(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create authors and topics
		author1 := createTestUser(1, "author1", "author1@example.com")
		author2 := createTestUser(2, "author2", "author2@example.com")
		mockRepo.CreateUser(ctx, author1)
		mockRepo.CreateUser(ctx, author2)

		desc1 := "Description 1"
		desc2 := "Description 2"
		topic1 := createTestTopic(1, "Topic 1", &desc1, 1)
		topic2 := createTestTopic(2, "Topic 2", &desc2, 2)
		mockRepo.CreateTopic(ctx, topic1)
		mockRepo.CreateTopic(ctx, topic2)

		// Execute
		topics, err := svc.FetchTopics(ctx, 10, 0, 0)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(topics) != 2 {
			t.Errorf("Expected 2 topics, got %d", len(topics))
		}
	})

	t.Run("FilterByUserID", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create authors and topics
		author1 := createTestUser(1, "author1", "author1@example.com")
		author2 := createTestUser(2, "author2", "author2@example.com")
		mockRepo.CreateUser(ctx, author1)
		mockRepo.CreateUser(ctx, author2)

		desc1 := "Description 1"
		desc2 := "Description 2"
		topic1 := createTestTopic(1, "Topic 1", &desc1, 1)
		topic2 := createTestTopic(2, "Topic 2", &desc2, 2)
		mockRepo.CreateTopic(ctx, topic1)
		mockRepo.CreateTopic(ctx, topic2)

		// Execute
		topics, err := svc.FetchTopics(ctx, 10, 0, 1)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(topics) != 1 {
			t.Errorf("Expected 1 topic, got %d", len(topics))
		}
		if len(topics) > 0 && topics[0].Author.Username != "author1" {
			t.Errorf("Expected author 'author1', got '%s'", topics[0].Author.Username)
		}
	})
}

// Test UpdateTopic
func TestUpdateTopic(t *testing.T) {
	t.Run("UpdateBothFields", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author and topic
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc := "Original description"
		topic := createTestTopic(1, "Original Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		newTitle := "Updated Topic"
		newDesc := "Updated description"

		// Execute
		err = svc.UpdateTopic(ctx, 1, 1, &newTitle, &newDesc)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify update
		updated, err := mockRepo.GetOneTopic(ctx, 1)
		if err != nil {
			t.Fatalf("Failed to fetch updated topic: %v", err)
		}
		if updated.Name != newTitle {
			t.Errorf("Expected name '%s', got '%s'", newTitle, updated.Name)
		}
		if updated.Description == nil || *updated.Description != newDesc {
			t.Errorf("Expected description '%s', got %v", newDesc, updated.Description)
		}
	})

	t.Run("UpdateTitleOnly", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author and topic
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc := "Original description"
		topic := createTestTopic(1, "Original Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		newTitle := "Updated Topic"

		// Execute
		err = svc.UpdateTopic(ctx, 1, 1, &newTitle, nil)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify update
		updated, err := mockRepo.GetOneTopic(ctx, 1)
		if err != nil {
			t.Fatalf("Failed to fetch updated topic: %v", err)
		}
		if updated.Name != newTitle {
			t.Errorf("Expected name '%s', got '%s'", newTitle, updated.Name)
		}
	})

	t.Run("TopicNotFound", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		newTitle := "Updated Topic"

		// Execute
		err = svc.UpdateTopic(ctx, 999, 1, &newTitle, nil)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrTopicNotFound) {
			t.Errorf("Expected ErrTopicNotFound, got %v", err)
		}
	})

	t.Run("Forbidden", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create authors and topic
		author1 := createTestUser(1, "author1", "author1@example.com")
		author2 := createTestUser(2, "author2", "author2@example.com")
		mockRepo.CreateUser(ctx, author1)
		mockRepo.CreateUser(ctx, author2)

		desc := "Original description"
		topic := createTestTopic(1, "Original Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		newTitle := "Updated Topic"

		// Execute (author2 trying to update author1's topic)
		err = svc.UpdateTopic(ctx, 1, 2, &newTitle, nil)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrForbidden) {
			t.Errorf("Expected ErrForbidden, got %v", err)
		}
	})

	t.Run("NoUpdateFields", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author and topic
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc := "Original description"
		topic := createTestTopic(1, "Original Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		// Execute
		err = svc.UpdateTopic(ctx, 1, 1, nil, nil)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrNoUpdateFields) {
			t.Errorf("Expected ErrNoUpdateFields, got %v", err)
		}
	})
}

// Test DeleteTopic
func TestDeleteTopic(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author and topic
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc := "Test description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		// Execute
		err = svc.DeleteTopic(ctx, 1, 1)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify deletion
		_, err = mockRepo.GetOneTopic(ctx, 1)
		if err == nil {
			t.Error("Expected topic to be deleted")
		}
	})

	t.Run("TopicNotFound", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Execute
		err = svc.DeleteTopic(ctx, 999, 1)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrTopicNotFound) {
			t.Errorf("Expected ErrTopicNotFound, got %v", err)
		}
	})

	t.Run("Forbidden", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create authors and topic
		author1 := createTestUser(1, "author1", "author1@example.com")
		author2 := createTestUser(2, "author2", "author2@example.com")
		mockRepo.CreateUser(ctx, author1)
		mockRepo.CreateUser(ctx, author2)

		desc := "Test description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		// Execute (author2 trying to delete author1's topic)
		err = svc.DeleteTopic(ctx, 1, 2)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrForbidden) {
			t.Errorf("Expected ErrForbidden, got %v", err)
		}
	})
}

// Test CreatePost
func TestCreatePost(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author and topic
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content := "Test post content"

		// Execute
		postInfo, err := svc.CreatePost(ctx, 1, 1, "Test Post", &content)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if postInfo == nil {
			t.Fatal("Expected post info, got nil")
		}
		if postInfo.Title != "Test Post" {
			t.Errorf("Expected title 'Test Post', got '%s'", postInfo.Title)
		}
		if postInfo.Content == nil || *postInfo.Content != content {
			t.Errorf("Expected content '%s', got %v", content, postInfo.Content)
		}
		if postInfo.TopicID != 1 {
			t.Errorf("Expected topic ID 1, got %d", postInfo.TopicID)
		}
	})

	t.Run("SuccessWithNullContent", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author and topic
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		// Execute
		postInfo, err := svc.CreatePost(ctx, 1, 1, "Test Post", nil)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if postInfo == nil {
			t.Fatal("Expected post info, got nil")
		}
		if postInfo.Content != nil {
			t.Errorf("Expected nil content, got %v", postInfo.Content)
		}
	})
}

// Test FetchPostByID
func TestFetchPostByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author, topic, and post
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content := "Test content"
		post := createTestPost(1, "Test Post", &content, 1, 1)
		mockRepo.CreatePost(ctx, post)

		// Execute
		postInfo, err := svc.FetchPostByID(ctx, 1)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if postInfo == nil {
			t.Fatal("Expected post info, got nil")
		}
		if postInfo.Title != "Test Post" {
			t.Errorf("Expected title 'Test Post', got '%s'", postInfo.Title)
		}
	})

	t.Run("PostNotFound", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Execute
		postInfo, err := svc.FetchPostByID(ctx, 999)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrPostNotFound) {
			t.Errorf("Expected ErrPostNotFound, got %v", err)
		}
		if postInfo != nil {
			t.Errorf("Expected nil post info, got %v", postInfo)
		}
	})
}

// Test FetchPosts
func TestFetchPosts(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create authors, topics, and posts
		author1 := createTestUser(1, "author1", "author1@example.com")
		author2 := createTestUser(2, "author2", "author2@example.com")
		mockRepo.CreateUser(ctx, author1)
		mockRepo.CreateUser(ctx, author2)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content1 := "Content 1"
		content2 := "Content 2"
		post1 := createTestPost(1, "Post 1", &content1, 1, 1)
		post2 := createTestPost(2, "Post 2", &content2, 2, 1)
		mockRepo.CreatePost(ctx, post1)
		mockRepo.CreatePost(ctx, post2)

		// Execute
		posts, err := svc.FetchPosts(ctx, 10, 0, 0, 0)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(posts) != 2 {
			t.Errorf("Expected 2 posts, got %d", len(posts))
		}
	})

	t.Run("FilterByTopicID", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author and topics
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc1 := "Topic 1"
		desc2 := "Topic 2"
		topic1 := createTestTopic(1, "Test Topic 1", &desc1, 1)
		topic2 := createTestTopic(2, "Test Topic 2", &desc2, 1)
		mockRepo.CreateTopic(ctx, topic1)
		mockRepo.CreateTopic(ctx, topic2)

		content1 := "Content 1"
		content2 := "Content 2"
		post1 := createTestPost(1, "Post 1", &content1, 1, 1)
		post2 := createTestPost(2, "Post 2", &content2, 1, 2)
		mockRepo.CreatePost(ctx, post1)
		mockRepo.CreatePost(ctx, post2)

		// Execute
		posts, err := svc.FetchPosts(ctx, 10, 0, 1, 0)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(posts) != 1 {
			t.Errorf("Expected 1 post, got %d", len(posts))
		}
		if len(posts) > 0 && posts[0].TopicID != 1 {
			t.Errorf("Expected topic ID 1, got %d", posts[0].TopicID)
		}
	})

	t.Run("FilterByUserID", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create authors and topic
		author1 := createTestUser(1, "author1", "author1@example.com")
		author2 := createTestUser(2, "author2", "author2@example.com")
		mockRepo.CreateUser(ctx, author1)
		mockRepo.CreateUser(ctx, author2)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content1 := "Content 1"
		content2 := "Content 2"
		post1 := createTestPost(1, "Post 1", &content1, 1, 1)
		post2 := createTestPost(2, "Post 2", &content2, 2, 1)
		mockRepo.CreatePost(ctx, post1)
		mockRepo.CreatePost(ctx, post2)

		// Execute
		posts, err := svc.FetchPosts(ctx, 10, 0, 0, 1)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(posts) != 1 {
			t.Errorf("Expected 1 post, got %d", len(posts))
		}
		if len(posts) > 0 && posts[0].Author.Username != "author1" {
			t.Errorf("Expected author 'author1', got '%s'", posts[0].Author.Username)
		}
	})
}

// Test UpdatePost
func TestUpdatePost(t *testing.T) {
	t.Run("UpdateBothFields", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author, topic, and post
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content := "Original content"
		post := createTestPost(1, "Original Post", &content, 1, 1)
		mockRepo.CreatePost(ctx, post)

		newTitle := "Updated Post"
		newContent := "Updated content"

		// Execute
		err = svc.UpdatePost(ctx, 1, 1, &newTitle, &newContent)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify update
		updated, err := mockRepo.GetOnePost(ctx, 1)
		if err != nil {
			t.Fatalf("Failed to fetch updated post: %v", err)
		}
		if updated.Title != newTitle {
			t.Errorf("Expected title '%s', got '%s'", newTitle, updated.Title)
		}
		if updated.Content == nil || *updated.Content != newContent {
			t.Errorf("Expected content '%s', got %v", newContent, updated.Content)
		}
	})

	t.Run("UpdateTitleOnly", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author, topic, and post
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content := "Original content"
		post := createTestPost(1, "Original Post", &content, 1, 1)
		mockRepo.CreatePost(ctx, post)

		newTitle := "Updated Post"

		// Execute
		err = svc.UpdatePost(ctx, 1, 1, &newTitle, nil)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify update
		updated, err := mockRepo.GetOnePost(ctx, 1)
		if err != nil {
			t.Fatalf("Failed to fetch updated post: %v", err)
		}
		if updated.Title != newTitle {
			t.Errorf("Expected title '%s', got '%s'", newTitle, updated.Title)
		}
	})

	t.Run("PostNotFound", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		newTitle := "Updated Post"

		// Execute
		err = svc.UpdatePost(ctx, 999, 1, &newTitle, nil)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrPostNotFound) {
			t.Errorf("Expected ErrPostNotFound, got %v", err)
		}
	})

	t.Run("Forbidden", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create authors, topic, and post
		author1 := createTestUser(1, "author1", "author1@example.com")
		author2 := createTestUser(2, "author2", "author2@example.com")
		mockRepo.CreateUser(ctx, author1)
		mockRepo.CreateUser(ctx, author2)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content := "Original content"
		post := createTestPost(1, "Original Post", &content, 1, 1)
		mockRepo.CreatePost(ctx, post)

		newTitle := "Updated Post"

		// Execute (author2 trying to update author1's post)
		err = svc.UpdatePost(ctx, 1, 2, &newTitle, nil)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrForbidden) {
			t.Errorf("Expected ErrForbidden, got %v", err)
		}
	})

	t.Run("NoUpdateFields", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author, topic, and post
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content := "Original content"
		post := createTestPost(1, "Original Post", &content, 1, 1)
		mockRepo.CreatePost(ctx, post)

		// Execute
		err = svc.UpdatePost(ctx, 1, 1, nil, nil)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrNoUpdateFields) {
			t.Errorf("Expected ErrNoUpdateFields, got %v", err)
		}
	})
}

// Test DeletePost
func TestDeletePost(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author, topic, and post
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content := "Test content"
		post := createTestPost(1, "Test Post", &content, 1, 1)
		mockRepo.CreatePost(ctx, post)

		// Execute
		err = svc.DeletePost(ctx, 1, 1)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify deletion
		_, err = mockRepo.GetOnePost(ctx, 1)
		if err == nil {
			t.Error("Expected post to be deleted")
		}
	})

	t.Run("PostNotFound", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Execute
		err = svc.DeletePost(ctx, 999, 1)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrPostNotFound) {
			t.Errorf("Expected ErrPostNotFound, got %v", err)
		}
	})

	t.Run("Forbidden", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create authors, topic, and post
		author1 := createTestUser(1, "author1", "author1@example.com")
		author2 := createTestUser(2, "author2", "author2@example.com")
		mockRepo.CreateUser(ctx, author1)
		mockRepo.CreateUser(ctx, author2)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content := "Test content"
		post := createTestPost(1, "Test Post", &content, 1, 1)
		mockRepo.CreatePost(ctx, post)

		// Execute (author2 trying to delete author1's post)
		err = svc.DeletePost(ctx, 1, 2)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrForbidden) {
			t.Errorf("Expected ErrForbidden, got %v", err)
		}
	})
}

// Test CreateComment
func TestCreateComment(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author, topic, and post
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content := "Test content"
		post := createTestPost(1, "Test Post", &content, 1, 1)
		mockRepo.CreatePost(ctx, post)

		// Execute
		commentInfo, err := svc.CreateComment(ctx, 1, 1, "This is a test comment")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if commentInfo == nil {
			t.Fatal("Expected comment info, got nil")
		}
		if commentInfo.Content != "This is a test comment" {
			t.Errorf("Expected content 'This is a test comment', got '%s'", commentInfo.Content)
		}
		if commentInfo.Author == nil {
			t.Fatal("Expected author, got nil")
		}
		if commentInfo.Author.Username != "author" {
			t.Errorf("Expected author username 'author', got '%s'", commentInfo.Author.Username)
		}
	})
}

// Test FetchCommentByID
func TestFetchCommentByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author, topic, post, and comment
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content := "Test content"
		post := createTestPost(1, "Test Post", &content, 1, 1)
		mockRepo.CreatePost(ctx, post)

		comment := createTestComment(1, "Test comment", 1, 1)
		mockRepo.CreateComment(ctx, comment)

		// Execute
		commentInfo, err := svc.FetchCommentByID(ctx, 1)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if commentInfo == nil {
			t.Fatal("Expected comment info, got nil")
		}
		if commentInfo.Content != "Test comment" {
			t.Errorf("Expected content 'Test comment', got '%s'", commentInfo.Content)
		}
	})

	t.Run("CommentNotFound", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Execute
		commentInfo, err := svc.FetchCommentByID(ctx, 999)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrCommentNotFound) {
			t.Errorf("Expected ErrCommentNotFound, got %v", err)
		}
		if commentInfo != nil {
			t.Errorf("Expected nil comment info, got %v", commentInfo)
		}
	})
}

// Test FetchComments
func TestFetchComments(t *testing.T) {
	t.Run("FetchAllComments", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create authors, topic, post, and comments
		author1 := createTestUser(1, "author1", "author1@example.com")
		author2 := createTestUser(2, "author2", "author2@example.com")
		mockRepo.CreateUser(ctx, author1)
		mockRepo.CreateUser(ctx, author2)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content := "Test content"
		post := createTestPost(1, "Test Post", &content, 1, 1)
		mockRepo.CreatePost(ctx, post)

		comment1 := createTestComment(1, "Comment 1", 1, 1)
		comment2 := createTestComment(2, "Comment 2", 2, 1)
		mockRepo.CreateComment(ctx, comment1)
		mockRepo.CreateComment(ctx, comment2)

		// Execute
		comments, err := svc.FetchComments(ctx, 10, 0, 0, 0)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(comments) != 2 {
			t.Errorf("Expected 2 comments, got %d", len(comments))
		}
	})

	t.Run("FilterByPostID", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create authors, topic, posts, and comments
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content1 := "Test content 1"
		content2 := "Test content 2"
		post1 := createTestPost(1, "Test Post 1", &content1, 1, 1)
		post2 := createTestPost(2, "Test Post 2", &content2, 1, 1)
		mockRepo.CreatePost(ctx, post1)
		mockRepo.CreatePost(ctx, post2)

		comment1 := createTestComment(1, "Comment for post 1", 1, 1)
		comment2 := createTestComment(2, "Another comment for post 1", 1, 1)
		comment3 := createTestComment(3, "Comment for post 2", 1, 2)
		mockRepo.CreateComment(ctx, comment1)
		mockRepo.CreateComment(ctx, comment2)
		mockRepo.CreateComment(ctx, comment3)

		// Execute
		comments, err := svc.FetchComments(ctx, 10, 0, 1, 0)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(comments) != 2 {
			t.Errorf("Expected 2 comments for post 1, got %d", len(comments))
		}
	})

	t.Run("FilterByUserID", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create authors, topic, post, and comments
		author1 := createTestUser(1, "author1", "author1@example.com")
		author2 := createTestUser(2, "author2", "author2@example.com")
		mockRepo.CreateUser(ctx, author1)
		mockRepo.CreateUser(ctx, author2)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content := "Test content"
		post := createTestPost(1, "Test Post", &content, 1, 1)
		mockRepo.CreatePost(ctx, post)

		comment1 := createTestComment(1, "Comment by author 1", 1, 1)
		comment2 := createTestComment(2, "Another comment by author 1", 1, 1)
		comment3 := createTestComment(3, "Comment by author 2", 2, 1)
		mockRepo.CreateComment(ctx, comment1)
		mockRepo.CreateComment(ctx, comment2)
		mockRepo.CreateComment(ctx, comment3)

		// Execute
		comments, err := svc.FetchComments(ctx, 10, 0, 0, 1)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(comments) != 2 {
			t.Errorf("Expected 2 comments by author 1, got %d", len(comments))
		}
	})

	t.Run("WithPagination", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author, topic, post, and multiple comments
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content := "Test content"
		post := createTestPost(1, "Test Post", &content, 1, 1)
		mockRepo.CreatePost(ctx, post)

		for i := 1; i <= 5; i++ {
			comment := createTestComment(uint(i), "Comment "+string(rune(i)), 1, 1)
			mockRepo.CreateComment(ctx, comment)
		}

		// Execute - get first page
		commentsPage1, err := svc.FetchComments(ctx, 2, 0, 0, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Execute - get second page
		commentsPage2, err := svc.FetchComments(ctx, 2, 2, 0, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Assert
		if len(commentsPage1) != 2 {
			t.Errorf("Expected 2 comments on page 1, got %d", len(commentsPage1))
		}
		if len(commentsPage2) != 2 {
			t.Errorf("Expected 2 comments on page 2, got %d", len(commentsPage2))
		}
		if commentsPage1[0].CommentID == commentsPage2[0].CommentID {
			t.Error("Expected different comments on different pages")
		}
	})
}

// Test UpdateComment
func TestUpdateComment(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author, topic, post, and comment
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content := "Test content"
		post := createTestPost(1, "Test Post", &content, 1, 1)
		mockRepo.CreatePost(ctx, post)

		comment := createTestComment(1, "Original comment", 1, 1)
		mockRepo.CreateComment(ctx, comment)

		// Execute
		err = svc.UpdateComment(ctx, 1, 1, "Updated comment")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify update
		updatedComment, err := svc.FetchCommentByID(ctx, 1)
		if err != nil {
			t.Fatalf("Failed to fetch updated comment: %v", err)
		}
		if updatedComment.Content != "Updated comment" {
			t.Errorf("Expected content 'Updated comment', got '%s'", updatedComment.Content)
		}
	})

	t.Run("CommentNotFound", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Execute
		err = svc.UpdateComment(ctx, 999, 1, "Updated comment")

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrCommentNotFound) {
			t.Errorf("Expected ErrCommentNotFound, got %v", err)
		}
	})

	t.Run("Forbidden", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create authors, topic, post, and comment
		author1 := createTestUser(1, "author1", "author1@example.com")
		author2 := createTestUser(2, "author2", "author2@example.com")
		mockRepo.CreateUser(ctx, author1)
		mockRepo.CreateUser(ctx, author2)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content := "Test content"
		post := createTestPost(1, "Test Post", &content, 1, 1)
		mockRepo.CreatePost(ctx, post)

		comment := createTestComment(1, "Comment by author1", 1, 1)
		mockRepo.CreateComment(ctx, comment)

		// Execute (author2 trying to update author1's comment)
		err = svc.UpdateComment(ctx, 1, 2, "Updated comment")

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrForbidden) {
			t.Errorf("Expected ErrForbidden, got %v", err)
		}
	})
}

// Test DeleteComment
func TestDeleteComment(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create author, topic, post, and comment
		author := createTestUser(1, "author", "author@example.com")
		mockRepo.CreateUser(ctx, author)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content := "Test content"
		post := createTestPost(1, "Test Post", &content, 1, 1)
		mockRepo.CreatePost(ctx, post)

		comment := createTestComment(1, "Test comment", 1, 1)
		mockRepo.CreateComment(ctx, comment)

		// Execute
		err = svc.DeleteComment(ctx, 1, 1)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify deletion
		_, err = svc.FetchCommentByID(ctx, 1)
		if err == nil {
			t.Error("Expected error when fetching deleted comment, got nil")
		}
		if !errors.Is(err, service.ErrCommentNotFound) {
			t.Errorf("Expected ErrCommentNotFound, got %v", err)
		}
	})

	t.Run("CommentNotFound", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Execute
		err = svc.DeleteComment(ctx, 999, 1)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrCommentNotFound) {
			t.Errorf("Expected ErrCommentNotFound, got %v", err)
		}
	})

	t.Run("Forbidden", func(t *testing.T) {
		// Setup
		mockRepo, err := createMockRepository()
		if err != nil {
			t.Fatalf("Failed to create mock repository: %v", err)
		}
		svc := primary.NewService(mockRepo)
		ctx := context.Background()

		// Create authors, topic, post, and comment
		author1 := createTestUser(1, "author1", "author1@example.com")
		author2 := createTestUser(2, "author2", "author2@example.com")
		mockRepo.CreateUser(ctx, author1)
		mockRepo.CreateUser(ctx, author2)

		desc := "Topic description"
		topic := createTestTopic(1, "Test Topic", &desc, 1)
		mockRepo.CreateTopic(ctx, topic)

		content := "Test content"
		post := createTestPost(1, "Test Post", &content, 1, 1)
		mockRepo.CreatePost(ctx, post)

		comment := createTestComment(1, "Comment by author1", 1, 1)
		mockRepo.CreateComment(ctx, comment)

		// Execute (author2 trying to delete author1's comment)
		err = svc.DeleteComment(ctx, 1, 2)

		// Assert
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, service.ErrForbidden) {
			t.Errorf("Expected ErrForbidden, got %v", err)
		}
	})
}
