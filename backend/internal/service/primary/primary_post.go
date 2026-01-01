package primary

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
	"gorm.io/gorm"
)

func (s *primaryService) CreatePost(
	ctx context.Context,
	userID uint,
	topicID uint,
	title string,
	content *string,
) (*service.PostInfo, error) {
	// Check if topic exists
	topicExists, err := s.repo.CheckTopicExists(ctx, topicID)
	if err != nil {
		slog.Error("failed to check topic existence", "topic_id", topicID, "error", err)
		return nil, errors.Join(service.ErrDatabaseError, err)
	}
	if !topicExists {
		slog.Warn("topic not found for creating post", "topic_id", topicID)
		return nil, service.ErrTopicNotFound
	}
	// Create new post model
	newPost := model.Post{
		Title:    title,
		Content:  content,
		AuthorID: userID,
		TopicID:  topicID,
	}

	// Save new post via repository
	if err := s.repo.CreatePost(ctx, &newPost); err != nil {
		slog.Error("failed to create post", "topic_id", topicID, "user_id", userID, "error", err)
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	// Fetch the created post with Author preloaded
	post, err := s.repo.GetOnePost(ctx, newPost.ID)
	if err != nil {
		slog.Error("failed to fetch created post", "post_id", newPost.ID, "error", err)
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	slog.Info("post created", "post_id", post.ID, "topic_id", topicID, "user_id", userID)

	// Return the newly created post's info
	info := service.InfoFromPost(post)
	return info, nil
}

func (s *primaryService) FetchPosts(
	ctx context.Context,
	limit,
	offset int,
	topicID,
	userID uint,
) ([]*service.PostInfo, error) {
	// Initialise userID pointer
	var userIDPtr *uint
	if userID != 0 {
		userIDPtr = &userID
	}

	// Initialise topicID pointer
	var topicIDPtr *uint
	if topicID != 0 {
		topicIDPtr = &topicID
	}

	// Fetch posts from repository
	posts, err := s.repo.GetPosts(ctx, limit, offset, topicIDPtr, userIDPtr)
	if err != nil {
		slog.Error("failed to fetch posts", "limit", limit, "offset", offset, "error", err)
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	slog.Debug("fetched posts", "count", len(posts), "limit", limit, "offset", offset)

	// Convert to DTOs
	postInfos := make([]*service.PostInfo, len(posts))
	for i, post := range posts {
		info := service.InfoFromPost(&post)
		postInfos[i] = info
	}

	return postInfos, nil
}

func (s *primaryService) FetchPostByID(
	ctx context.Context,
	postID uint,
) (*service.PostInfo, error) {
	// Fetch post from repository
	post, err := s.repo.GetOnePost(ctx, postID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Warn("post not found", "post_id", postID)
			return nil, service.ErrPostNotFound
		}
		slog.Error("failed to fetch post", "post_id", postID, "error", err)
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	slog.Debug("fetched post", "post_id", postID)

	// Convert to DTO
	info := service.InfoFromPost(post)
	return info, nil
}

func (s *primaryService) UpdatePost(
	ctx context.Context,
	postID,
	userID uint,
	title,
	content *string,
) error {
	// Fetch existing post
	post, err := s.repo.GetOnePost(ctx, postID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Warn("post not found for update", "post_id", postID)
			return service.ErrPostNotFound
		}
		slog.Error("failed to fetch post for update", "post_id", postID, "error", err)
		return errors.Join(service.ErrDatabaseError, err)
	}

	// Check if the user is the author
	if post.AuthorID != userID {
		slog.Warn("forbidden post update attempt", "post_id", postID, "user_id", userID, "author_id", post.AuthorID)
		return service.ErrForbidden
	}

	// Update fields if provided
	if title != nil {
		post.Title = *title
	}
	if content != nil {
		post.Content = content
	}

	// Save updated post via repository
	if err := s.repo.UpdatePost(
		ctx,
		post.ID,
		post.Title,
		post.Content,
	); err != nil {
		slog.Error("failed to update post", "post_id", postID, "error", err)
		return errors.Join(service.ErrDatabaseError, err)
	}

	slog.Info("post updated", "post_id", postID, "user_id", userID)

	return nil
}

func (s *primaryService) DeletePost(
	ctx context.Context,
	postID,
	userID uint,
) error {
	// Fetch existing post
	post, err := s.repo.GetOnePost(ctx, postID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Warn("post not found for delete", "post_id", postID)
			return service.ErrPostNotFound
		}
		slog.Error("failed to fetch post for delete", "post_id", postID, "error", err)
		return errors.Join(service.ErrDatabaseError, err)
	}

	// Check if the user is the author
	if post.AuthorID != userID {
		slog.Warn("forbidden post delete attempt", "post_id", postID, "user_id", userID, "author_id", post.AuthorID)
		return service.ErrForbidden
	}

	// Delete post via repository
	if err := s.repo.DeletePost(ctx, postID); err != nil {
		slog.Error("failed to delete post", "post_id", postID, "error", err)
		return errors.Join(service.ErrDatabaseError, err)
	}

	slog.Info("post deleted", "post_id", postID, "user_id", userID)

	return nil
}
