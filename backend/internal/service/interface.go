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

	RegisterUser(ctx context.Context, username, email, password string) (*UserProfile, error)
	FetchUserByID(ctx context.Context, id uint) (*UserProfile, error)
	FetchCurrentUser(ctx context.Context, user *model.User) (*UserProfile, error)
	UpdateCurrentUser(ctx context.Context, user *model.User, newAvatarUrl, newBio *string) error

	// Topic-related services

	CreateTopic(ctx context.Context, user *model.User, title string, description *string) (*TopicInfo, error)
	FetchTopicByID(ctx context.Context, topicID uint) (*TopicInfo, error)
	FetchTopics(ctx context.Context, limit, offset int, userID uint, withPosts bool) ([]TopicInfo, error)
	UpdateTopic(ctx context.Context, topicID, userID uint, title, description *string) error
	DeleteTopic(ctx context.Context, topicID, userID uint) error
}
