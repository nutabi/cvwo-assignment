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
	FetchTopics(ctx context.Context, limit, offset int, userID uint) ([]*TopicInfo, error)
	FetchTopicByID(ctx context.Context, topicID uint) (*TopicInfo, error)
	UpdateTopic(ctx context.Context, topicID, userID uint, title, description *string) error
	DeleteTopic(ctx context.Context, topicID, userID uint) error

	// Post-related services

	CreatePost(ctx context.Context, userID uint, topicID uint, title string, content *string) (*PostInfo, error)
	FetchPosts(ctx context.Context, limit, offset int, topicID, userID uint) ([]*PostInfo, error)
	FetchPostByID(ctx context.Context, postID uint) (*PostInfo, error)
	UpdatePost(ctx context.Context, postID, userID uint, title, content *string) error
	DeletePost(ctx context.Context, postID, userID uint) error

	// Comment-related services

	CreateComment(ctx context.Context, userID uint, postID uint, content string) (*CommentInfo, error)
	FetchComments(ctx context.Context, limit, offset int, postID, userID uint) ([]*CommentInfo, error)
	FetchCommentByID(ctx context.Context, commentID uint) (*CommentInfo, error)
	UpdateComment(ctx context.Context, commentID, userID uint, content string) error
	DeleteComment(ctx context.Context, commentID, userID uint) error
}
