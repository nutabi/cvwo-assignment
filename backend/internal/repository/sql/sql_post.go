package sql

import (
	"context"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"gorm.io/gorm"
)

func (r *sqlRepository) CreatePost(ctx context.Context, post *model.Post) error {
	return gorm.
		G[model.Post](&r.db).
		Create(ctx, post)
}

func (r *sqlRepository) GetPostByID(ctx context.Context, id uint) (model.Post, error) {
	return gorm.
		G[model.Post](&r.db).
		Preload("Author", nil).
		Preload("Comments", nil).
		Where("id = ?", id).
		First(ctx)
}

func (r *sqlRepository) GetPosts(
	ctx context.Context,
	limit, offset int,
	topicID *uint,
	userID *uint,
	withComments bool,
) ([]model.Post, error) {
	// Base query
	query := gorm.
		G[model.Post](&r.db).
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

	// Preload comments if requested
	if withComments {
		query = query.Preload("Comments", nil)
	}

	// Execute query
	posts, err := query.Find(ctx)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *sqlRepository) UpdatePost(ctx context.Context, post *model.Post) error {
	_, err := gorm.
		G[model.Post](&r.db).
		Where("id = ?", post.ID).
		Updates(ctx, *post)
	return err
}

func (r *sqlRepository) DeletePost(ctx context.Context, post *model.Post) error {
	_, err := gorm.
		G[model.Post](&r.db).
		Where("id = ?", post.ID).
		Delete(ctx)
	return err
}
