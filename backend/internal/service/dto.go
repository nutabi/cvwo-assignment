package service

import (
	"time"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
)

type UserProfile struct {
	UserID    uint       `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Username  string     `json:"username"`
	Email     *string    `json:"email,omitempty"`
	AvatarUrl *string    `json:"avatar_url"`
	Bio       *string    `json:"bio"`
}

type TopicInfo struct {
	TopicID     uint         `json:"topic_id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Name        string       `json:"name"`
	Description *string      `json:"description"`
	Author      *UserProfile `json:"author"`

	Posts []PostInfo `json:"posts,omitempty"`
}

type PostInfo struct {
	PostID    uint         `json:"post_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Title     string       `json:"title"`
	Content   *string      `json:"content"`
	Author    *UserProfile `json:"author"`
	TopicID   uint         `json:"topic_id"`

	Comments []CommentInfo `json:"comments,omitempty"`
}

type CommentInfo struct {
	CommentID uint         `json:"comment_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Content   string       `json:"content"`
	Author    *UserProfile `json:"author"`
	PostID    uint         `json:"post_id"`
}

func ProfileFromUser(user *model.User, isPrivate bool) *UserProfile {
	profile := &UserProfile{
		UserID:    user.ID,
		CreatedAt: user.CreatedAt,
		Username:  user.Username,
		AvatarUrl: user.AvatarURL,
		Bio:       user.Bio,
	}
	if isPrivate {
		profile.UpdatedAt = &user.UpdatedAt
		profile.Email = &user.Email
	}
	return profile
}

func InfoFromTopic(topic *model.Topic) *TopicInfo {
	return &TopicInfo{
		TopicID:     topic.ID,
		CreatedAt:   topic.CreatedAt,
		UpdatedAt:   topic.UpdatedAt,
		Name:        topic.Name,
		Description: topic.Description,
		Author:      ProfileFromUser(topic.Author, false),
	}
}

func InfoFromPost(post *model.Post) *PostInfo {
	return &PostInfo{
		PostID:    post.ID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Title:     post.Title,
		Content:   post.Content,
		Author:    ProfileFromUser(post.Author, false),
		TopicID:   post.TopicID,
	}
}

func InfoFromComment(comment *model.Comment) *CommentInfo {
	return &CommentInfo{
		CommentID: comment.ID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		Content:   comment.Content,
		Author:    ProfileFromUser(comment.Author, false),
		PostID:    comment.PostID,
	}
}
