package sql

import (
	"context"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"gorm.io/gorm"
)

func (r *sqlRepository) GetTopics(
	ctx context.Context,
	limit, offset int,
	userID *uint,
	withPosts bool,
) ([]model.Topic, error) {
	// Base query
	query := gorm.
		G[model.Topic](&r.db).
		Preload("Author", nil).
		Offset(offset).
		Limit(limit)

	// Filter by userID if provided
	if userID != nil {
		query = query.Where("author_id = ?", *userID)
	}

	// Preload posts if requested
	if withPosts {
		query = query.Preload("Posts", nil)
	}

	// Execute query
	topics, err := query.Find(ctx)
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (r *sqlRepository) CreateTopic(
	ctx context.Context,
	topic *model.Topic,
) error {
	err := gorm.
		G[model.Topic](&r.db).
		Create(ctx, topic)
	if err != nil {
		return err
	}
	return nil
}

func (r *sqlRepository) GetTopicByID(
	ctx context.Context,
	id uint,
) (model.Topic, error) {
	topic, err := gorm.
		G[model.Topic](&r.db).
		Preload("Author", nil).
		Where("id = ?", id).
		First(ctx)
	return topic, err
}

func (r *sqlRepository) UpdateTopic(
	ctx context.Context,
	topic *model.Topic,
) error {
	_, err := gorm.
		G[model.Topic](&r.db).
		Where("id = ?", topic.ID).
		Updates(ctx, model.Topic{Name: topic.Name, Description: topic.Description})
	return err
}

func (r *sqlRepository) DeleteTopic(
	ctx context.Context,
	topic *model.Topic,
) error {
	_, err := gorm.
		G[model.Topic](&r.db).
		Where("id = ?", topic.ID).
		Delete(ctx)
	return err
}
