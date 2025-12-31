package primary

import (
	"context"
	"errors"

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
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	// Return the newly created topic's info
	info := service.InfoFromTopic(&newTopic)
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
			return nil, service.ErrTopicNotFound
		}
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

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
		return nil, errors.Join(service.ErrDatabaseError, err)
	}

	// Convert to DTO
	infos := make([]*service.TopicInfo, 0, len(topics))
	for _, topic := range topics {
		infos = append(infos, service.InfoFromTopic(&topic))
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
			return service.ErrTopicNotFound
		}
		return errors.Join(service.ErrDatabaseError, err)
	}

	// Check if the user is the author
	if topic.AuthorID != userID {
		return service.ErrForbidden
	}

	// Make sure at least one field is being updated
	if title == nil && description == nil {
		return service.ErrNoUpdateFields
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
		return errors.Join(service.ErrDatabaseError, err)
	}

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
			return service.ErrTopicNotFound
		}
		return errors.Join(service.ErrDatabaseError, err)
	}

	// Check if the user is the author
	if topic.AuthorID != userID {
		return service.ErrForbidden
	}

	// Delete topic via repository
	if err := s.repo.DeleteTopic(ctx, topicID); err != nil {
		return errors.Join(service.ErrDatabaseError, err)
	}

	return nil
}
