package primary

import (
	"context"
	"errors"

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
	// Create new post model
	newPost := model.Post{
		Title:    title,
		Content:  content,
		AuthorID: userID,
		TopicID:  topicID,
	}

	// Save new post via repository
	if err := s.repo.CreatePost(ctx, &newPost); err != nil {
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	// Fetch the created post with Author preloaded
	post, err := s.repo.GetOnePost(ctx, newPost.ID)
	if err != nil {
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

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
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

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
			return nil, service.ErrPostNotFound
		}
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

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
			return service.ErrPostNotFound
		}
		return errors.Join(service.ErrDatabaseError, err)
	}

	// Check if the user is the author
	if post.AuthorID != userID {
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
		return errors.Join(service.ErrDatabaseError, err)
	}

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
			return service.ErrPostNotFound
		}
		return errors.Join(service.ErrDatabaseError, err)
	}

	// Check if the user is the author
	if post.AuthorID != userID {
		return service.ErrForbidden
	}

	// Delete post via repository
	if err := s.repo.DeletePost(ctx, postID); err != nil {
		return errors.Join(service.ErrDatabaseError, err)
	}

	return nil
}
