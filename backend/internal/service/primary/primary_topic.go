package primary

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
	"gorm.io/gorm"
)

func (s *primaryService) CreateTopic(
	ctx context.Context,
	userID uint,
	title string,
	description *string,
) (*service.TopicInfo, error) {
	// Create new topic model
	newTopic := model.Topic{
		Name:        title,
		Description: description,
		AuthorID:    userID,
	}

	// Save new topic via repository
	if err := s.repo.CreateTopic(ctx, &newTopic); err != nil {
		slog.Error("failed to create topic", "user_id", userID, "error", err)
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	// Fetch the created topic with Author preloaded
	topic, err := s.repo.GetOneTopic(ctx, newTopic.ID)
	if err != nil {
		slog.Error("failed to fetch created topic", "topic_id", newTopic.ID, "error", err)
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	slog.Info("topic created", "topic_id", topic.ID, "user_id", userID)

	// Return the newly created topic's info
	info := service.InfoFromTopic(topic)
	return info, nil
}

func (s *primaryService) FetchTopicByID(
	ctx context.Context,
	topicID uint,
) (*service.TopicInfo, error) {
	// Fetch topic from repository
	topic, err := s.repo.GetOneTopic(ctx, topicID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Warn("topic not found", "topic_id", topicID)
			return nil, service.ErrTopicNotFound
		}
		slog.Error("failed to fetch topic", "topic_id", topicID, "error", err)
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	slog.Debug("fetched topic", "topic_id", topicID)

	// Convert to DTO
	info := service.InfoFromTopic(topic)
	return info, nil
}

func (s *primaryService) FetchTopics(
	ctx context.Context,
	limit,
	offset int,
	userID uint,
) ([]*service.TopicInfo, error) {
	// Initialise userID pointer
	var userIDPtr *uint
	if userID != 0 {
		userIDPtr = &userID
	}

	// Fetch topics from repository
	topics, err := s.repo.GetTopics(ctx, limit, offset, userIDPtr)
	if err != nil {
		slog.Error("failed to fetch topics", "limit", limit, "offset", offset, "error", err)
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	slog.Debug("fetched topics", "count", len(topics), "limit", limit, "offset", offset)

	// Convert to DTO
	infos := make([]*service.TopicInfo, len(topics))
	for i, topic := range topics {
		infos[i] = service.InfoFromTopic(&topic)
	}
	return infos, nil
}

func (s *primaryService) UpdateTopic(
	ctx context.Context,
	topicID,
	userID uint,
	title,
	description *string,
) error {
	// Fetch topic from repository
	topic, err := s.repo.GetOneTopic(ctx, topicID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Warn("topic not found for update", "topic_id", topicID)
			return service.ErrTopicNotFound
		}
		slog.Error("failed to fetch topic for update", "topic_id", topicID, "error", err)
		return errors.Join(service.ErrDatabaseError, err)
	}

	// Check if the user is the author
	if topic.AuthorID != userID {
		slog.Warn("forbidden topic update attempt", "topic_id", topicID, "user_id", userID, "author_id", topic.AuthorID)
		return service.ErrForbidden
	}

	// Update fields if provided
	if title != nil {
		topic.Name = *title
	}
	if description != nil {
		topic.Description = description
	}

	// Save updated topic via repository
	if err := s.repo.UpdateTopic(
		ctx,
		topic.ID,
		topic.Name,
		topic.Description,
	); err != nil {
		slog.Error("failed to update topic", "topic_id", topicID, "error", err)
		return errors.Join(service.ErrDatabaseError, err)
	}

	slog.Info("topic updated", "topic_id", topicID, "user_id", userID)

	return nil
}

func (s *primaryService) DeleteTopic(
	ctx context.Context,
	topicID,
	userID uint,
) error {
	// Fetch topic from repository
	topic, err := s.repo.GetOneTopic(ctx, topicID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Warn("topic not found for delete", "topic_id", topicID)
			return service.ErrTopicNotFound
		}
		slog.Error("failed to fetch topic for delete", "topic_id", topicID, "error", err)
		return errors.Join(service.ErrDatabaseError, err)
	}

	// Check if the user is the author
	if topic.AuthorID != userID {
		slog.Warn("forbidden topic delete attempt", "topic_id", topicID, "user_id", userID, "author_id", topic.AuthorID)
		return service.ErrForbidden
	}

	// Delete topic via repository
	if err := s.repo.DeleteTopic(ctx, topicID); err != nil {
		slog.Error("failed to delete topic", "topic_id", topicID, "error", err)
		return errors.Join(service.ErrDatabaseError, err)
	}

	slog.Info("topic deleted", "topic_id", topicID, "user_id", userID)

	return nil
}
