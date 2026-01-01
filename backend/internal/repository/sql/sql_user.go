package sql

import (
	"context"
	"log/slog"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"gorm.io/gorm"
)

func (r *sqlRepository) CreateUser(
	ctx context.Context,
	user *model.User,
) error {
	slog.Debug("creating user", "username", user.Username)
	return gorm.
		G[model.User](r.db).
		Create(ctx, user)
}

func (r *sqlRepository) GetUserByID(
	ctx context.Context,
	id uint,
) (*model.User, error) {
	slog.Debug("fetching user by ID", "user_id", id)
	user, err := gorm.
		G[model.User](r.db).
		Where("id = ?", id).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *sqlRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (*model.User, error) {
	slog.Debug("fetching user by username", "username", username)
	user, err := gorm.
		G[model.User](r.db).
		Where("username = ?", username).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *sqlRepository) CheckUsernameExists(
	ctx context.Context,
	username string,
) (bool, error) {
	slog.Debug("checking username existence", "username", username)
	var count int64
	err := r.db.Model(&model.User{}).
		Where("username = ?", username).
		Count(&count).
		Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *sqlRepository) CheckEmailExists(
	ctx context.Context,
	email string,
) (bool, error) {
	slog.Debug("checking email existence", "email", email)
	var count int64
	err := r.db.Model(&model.User{}).
		Where("email = ?", email).
		Count(&count).
		Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *sqlRepository) UpdateUser(
	ctx context.Context,
	userID uint,
	avatarURL *string,
	bio *string,
) error {
	slog.Debug("updating user", "user_id", userID)
	_, err := gorm.
		G[model.User](r.db).
		Where("id = ?", userID).
		Updates(ctx, model.User{AvatarURL: avatarURL, Bio: bio})
	return err
}
