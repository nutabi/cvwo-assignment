package sql

import (
	"context"
	"log/slog"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"gorm.io/gorm"
)

func (r *sqlRepository) CreateComment(ctx context.Context, comment *model.Comment) error {
	slog.Debug("creating comment", "post_id", comment.PostID, "author_id", comment.AuthorID)
	return gorm.
		G[model.Comment](r.db).
		Create(ctx, comment)
}

func (r *sqlRepository) GetOneComment(ctx context.Context, id uint) (*model.Comment, error) {
	slog.Debug("fetching comment by ID", "comment_id", id)
	comment, err := gorm.
		G[model.Comment](r.db).
		Preload("Author", nil).
		Where("id = ?", id).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *sqlRepository) GetComments(
	ctx context.Context,
	limit, offset int,
	postID *uint,
	userID *uint,
) ([]model.Comment, error) {
	slog.Debug("fetching comments", "limit", limit, "offset", offset)

	// Base query
	query := gorm.
		G[model.Comment](r.db).
		Preload("Author", nil).
		Offset(offset).
		Limit(limit)

	// Filter by postID if provided
	if postID != nil {
		query = query.Where("post_id = ?", *postID)
	}

	// Filter by userID if provided
	if userID != nil {
		query = query.Where("author_id = ?", *userID)
	}

	// Execute query
	comments, err := query.Find(ctx)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *sqlRepository) UpdateComment(ctx context.Context, commentID uint, content string) error {
	slog.Debug("updating comment", "comment_id", commentID)
	_, err := gorm.
		G[model.Comment](r.db).
		Where("id = ?", commentID).
		Updates(ctx, model.Comment{Content: content})
	return err
}

func (r *sqlRepository) DeleteComment(ctx context.Context, commentID uint) error {
	slog.Debug("deleting comment", "comment_id", commentID)
	_, err := gorm.
		G[model.Comment](r.db).
		Where("id = ?", commentID).
		Delete(ctx)
	return err
}
