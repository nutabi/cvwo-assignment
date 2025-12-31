package sql

import (
	"context"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"gorm.io/gorm"
)

func (r *sqlRepository) CreatePost(ctx context.Context, post *model.Post) error {
	return gorm.
		G[model.Post](r.db).
		Create(ctx, post)
}

func (r *sqlRepository) GetOnePost(ctx context.Context, id uint) (*model.Post, error) {
	post, err := gorm.
		G[model.Post](r.db).
		Preload("Author", nil).
		Preload("Comments", nil).
		Where("id = ?", id).
		First(ctx)
	return &post, err
}

func (r *sqlRepository) GetPosts(
	ctx context.Context,
	limit, offset int,
	topicID *uint,
	userID *uint,
) ([]model.Post, error) {
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

func (r *sqlRepository) UpdatePost(
	ctx context.Context,
	postID uint,
	title,
	content *string,
) error {
	_, err := gorm.
		G[model.Post](r.db).
		Where("id = ?", postID).
		Updates(ctx, model.Post{Title: *title, Content: content})
	return err
}

func (r *sqlRepository) DeletePost(ctx context.Context, postID uint) error {
	_, err := gorm.
		G[model.Post](r.db).
		Where("id = ?", postID).
		Delete(ctx)
	return err
}
