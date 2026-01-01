package repository

import (
	"context"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
)

type Repository interface {
	Migrate() error

	// User-related repository methods

	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	UpdateUser(ctx context.Context, userID uint, avatarURL, bio *string) error

	// Topic-related repository methods

	CreateTopic(ctx context.Context, topic *model.Topic) error
	GetOneTopic(ctx context.Context, id uint) (*model.Topic, error)
	GetTopics(ctx context.Context, limit, offset int, userID *uint) ([]model.Topic, error)
	CheckTopicExists(ctx context.Context, id uint) (bool, error)
	UpdateTopic(ctx context.Context, topicID uint, name string, description *string) error
	DeleteTopic(ctx context.Context, topicID uint) error

	// Post-related repository methods

	CreatePost(ctx context.Context, post *model.Post) error
	GetOnePost(ctx context.Context, id uint) (*model.Post, error)
	GetPosts(ctx context.Context, limit, offset int, topicID *uint, userID *uint) ([]model.Post, error)
	CheckPostExists(ctx context.Context, id uint) (bool, error)
	UpdatePost(ctx context.Context, postID uint, title string, content *string) error
	DeletePost(ctx context.Context, postID uint) error

	// Comment-related repository methods

	CreateComment(ctx context.Context, comment *model.Comment) error
	GetOneComment(ctx context.Context, id uint) (*model.Comment, error)
	GetComments(ctx context.Context, limit, offset int, postID *uint, userID *uint) ([]model.Comment, error)
	UpdateComment(ctx context.Context, commentID uint, content string) error
	DeleteComment(ctx context.Context, commentID uint) error
}
