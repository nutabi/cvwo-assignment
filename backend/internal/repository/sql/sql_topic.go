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
) ([]model.Topic, error) {
	// Base query
	query := gorm.
		G[model.Topic](r.db).
		Preload("Author", nil).
		Offset(offset).
		Limit(limit)

	// Filter by userID if provided
	if userID != nil {
		query = query.Where("author_id = ?", *userID)
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
	return gorm.
		G[model.Topic](r.db).
		Create(ctx, topic)
}

func (r *sqlRepository) GetOneTopic(
	ctx context.Context,
	id uint,
) (*model.Topic, error) {
	topic, err := gorm.
		G[model.Topic](r.db).
		Preload("Author", nil).
		Preload("Posts", nil).
		Where("id = ?", id).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return &topic, err
}

func (r *sqlRepository) UpdateTopic(
	ctx context.Context,
	topicID uint,
	name string,
	description *string,
) error {
	_, err := gorm.
		G[model.Topic](r.db).
		Where("id = ?", topicID).
		Updates(ctx, model.Topic{Name: name, Description: description})
	return err
}

func (r *sqlRepository) DeleteTopic(
	ctx context.Context,
	topicID uint,
) error {
	_, err := gorm.
		G[model.Topic](r.db).
		Where("id = ?", topicID).
		Delete(ctx)
	return err
}
