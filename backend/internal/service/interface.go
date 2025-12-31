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

	CreateTopic(ctx context.Context, userID uint, title string, description *string) (*TopicInfo, error)
	FetchTopics(ctx context.Context, limit, offset int, userID uint, withPosts bool) ([]TopicInfo, error)
	FetchTopicByID(ctx context.Context, topicID uint) (*TopicInfo, error)
	UpdateTopic(ctx context.Context, topicID, userID uint, title, description *string) error
	DeleteTopic(ctx context.Context, topicID, userID uint) error

	// Post-related services

	CreatePost(ctx context.Context, userID uint, topicID uint, title, content string) (*PostInfo, error)
	FetchPosts(ctx context.Context, limit, offset int, postID, userID uint, withComments bool) ([]PostInfo, error)
	FetchPostByID(ctx context.Context, postID uint) (*PostInfo, error)
	UpdatePost(ctx context.Context, postID, userID uint, title, content *string) error
	DeletePost(ctx context.Context, postID, userID uint) error
}
