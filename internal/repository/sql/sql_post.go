package sql

import (
	"context"
	"log/slog"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"gorm.io/gorm"
)

func (r *sqlRepository) CreatePost(ctx context.Context, post *model.Post) error {
	slog.Debug("creating post", "topic_id", post.TopicID, "author_id", post.AuthorID)
	return gorm.
		G[model.Post](r.db).
		Create(ctx, post)
}

func (r *sqlRepository) GetOnePost(ctx context.Context, id uint) (*model.Post, error) {
	slog.Debug("fetching post by ID", "post_id", id)
	post, err := gorm.
		G[model.Post](r.db).
		Preload("Author", nil).
		Preload("Comments", nil).
		Where("id = ?", id).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *sqlRepository) GetPosts(
	ctx context.Context,
	limit, offset int,
	topicID *uint,
	userID *uint,
) ([]model.Post, error) {
	slog.Debug("fetching posts", "limit", limit, "offset", offset)

	// Base query
	query := gorm.
		G[model.Post](r.db).
		Preload("Author", nil).
		Offset(offset).
		Limit(limit)

	// Filter by topicID if provided
	if topicID != nil {
		query = query.Where("topic_id = ?", *topicID)
	}

	// Filter by userID if provided
	if userID != nil {
		query = query.Where("author_id = ?", *userID)
	}

	// Execute query
	posts, err := query.Find(ctx)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *sqlRepository) CheckPostExists(
	ctx context.Context,
	id uint,
) (bool, error) {
	slog.Debug("checking post existence", "post_id", id)
	count, err := gorm.
		G[model.Post](r.db).
		Where("id = ?", id).
		Count(ctx, "id")
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *sqlRepository) UpdatePost(
	ctx context.Context,
	postID uint,
	title string,
	content *string,
) error {
	slog.Debug("updating post", "post_id", postID)
	_, err := gorm.
		G[model.Post](r.db).
		Where("id = ?", postID).
		Updates(ctx, model.Post{Title: title, Content: content})
	return err
}

func (r *sqlRepository) DeletePost(ctx context.Context, postID uint) error {
	slog.Debug("deleting post (hard delete)", "post_id", postID)
	_, err := gorm.
		G[model.Post](r.db.Unscoped()).
		Where("id = ?", postID).
		Delete(ctx)
	return err
}
