package primary

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
	"gorm.io/gorm"
)

func (s *primaryService) CreateComment(
	ctx context.Context,
	userID uint,
	postID uint,
	content string,
) (*service.CommentInfo, error) {
	// Create new comment model
	newComment := model.Comment{
		Content:  content,
		AuthorID: userID,
		PostID:   postID,
	}

	// Save new comment via repository
	if err := s.repo.CreateComment(ctx, &newComment); err != nil {
		slog.Error("failed to create comment", "post_id", postID, "user_id", userID, "error", err)
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	// Fetch the created comment with Author preloaded
	comment, err := s.repo.GetOneComment(ctx, newComment.ID)
	if err != nil {
		slog.Error("failed to fetch created comment", "comment_id", newComment.ID, "error", err)
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	slog.Info("comment created", "comment_id", comment.ID, "post_id", postID, "user_id", userID)

	// Return the newly created comment's info
	info := service.InfoFromComment(comment)
	return info, nil
}

func (s *primaryService) FetchComments(
	ctx context.Context,
	limit,
	offset int,
	postID,
	userID uint,
) ([]*service.CommentInfo, error) {
	// Initialise userID pointer
	var userIDPtr *uint
	if userID != 0 {
		userIDPtr = &userID
	}

	// Initialise postID pointer
	var postIDPtr *uint
	if postID != 0 {
		postIDPtr = &postID
	}

	// Fetch comments from repository
	comments, err := s.repo.GetComments(ctx, limit, offset, postIDPtr, userIDPtr)
	if err != nil {
		slog.Error("failed to fetch comments", "limit", limit, "offset", offset, "error", err)
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	slog.Debug("fetched comments", "count", len(comments), "limit", limit, "offset", offset)

	// Convert to CommentInfo slice
	var commentInfos []*service.CommentInfo
	for _, comment := range comments {
		commentInfos = append(commentInfos, service.InfoFromComment(&comment))
	}

	return commentInfos, nil
}

func (s *primaryService) FetchCommentByID(
	ctx context.Context,
	commentID uint,
) (*service.CommentInfo, error) {
	// Fetch comment from repository
	comment, err := s.repo.GetOneComment(ctx, commentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Warn("comment not found", "comment_id", commentID)
			return nil, service.ErrCommentNotFound
		}
		slog.Error("failed to fetch comment", "comment_id", commentID, "error", err)
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	slog.Debug("fetched comment", "comment_id", commentID)

	// Convert to CommentInfo
	info := service.InfoFromComment(comment)
	return info, nil
}

func (s *primaryService) UpdateComment(
	ctx context.Context,
	commentID,
	userID uint,
	content string,
) error {
	// Fetch comment from repository
	comment, err := s.repo.GetOneComment(ctx, commentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Warn("comment not found for update", "comment_id", commentID)
			return service.ErrCommentNotFound
		}
		slog.Error("failed to fetch comment for update", "comment_id", commentID, "error", err)
		return errors.Join(service.ErrDatabaseError, err)
	}

	// Check if the user is the author of the comment
	if comment.AuthorID != userID {
		slog.Warn("forbidden comment update attempt", "comment_id", commentID, "user_id", userID, "author_id", comment.AuthorID)
		return service.ErrForbidden
	}

	// Save updated comment via repository
	if err := s.repo.UpdateComment(ctx, commentID, content); err != nil {
		slog.Error("failed to update comment", "comment_id", commentID, "error", err)
		return errors.Join(service.ErrDatabaseError, err)
	}

	slog.Info("comment updated", "comment_id", commentID, "user_id", userID)

	return nil
}

func (s *primaryService) DeleteComment(
	ctx context.Context,
	commentID,
	userID uint,
) error {
	// Fetch comment from repository
	comment, err := s.repo.GetOneComment(ctx, commentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Warn("comment not found for delete", "comment_id", commentID)
			return service.ErrCommentNotFound
		}
		slog.Error("failed to fetch comment for delete", "comment_id", commentID, "error", err)
		return errors.Join(service.ErrDatabaseError, err)
	}

	// Check if the user is the author of the comment
	if comment.AuthorID != userID {
		slog.Warn("forbidden comment delete attempt", "comment_id", commentID, "user_id", userID, "author_id", comment.AuthorID)
		return service.ErrForbidden
	}

	// Delete comment via repository
	if err := s.repo.DeleteComment(ctx, commentID); err != nil {
		slog.Error("failed to delete comment", "comment_id", commentID, "error", err)
		return errors.Join(service.ErrDatabaseError, err)
	}

	slog.Info("comment deleted", "comment_id", commentID, "user_id", userID)

	return nil
}
