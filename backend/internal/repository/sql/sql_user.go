package sql

import (
	"context"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"gorm.io/gorm"
)

func (r *sqlRepository) CreateUser(
	ctx context.Context,
	user *model.User,
) error {
	return gorm.
		G[model.User](r.db).
		Create(ctx, user)
}

func (r *sqlRepository) GetUserByID(
	ctx context.Context,
	id uint,
) (*model.User, error) {
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
	_, err := gorm.
		G[model.User](r.db).
		Where("id = ?", userID).
		Updates(ctx, model.User{AvatarURL: avatarURL, Bio: bio})
	return err
}
